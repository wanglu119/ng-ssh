package http

import (
	"fmt"
	"net/http"
	"encoding/base64"
	"encoding/json"
	"sync"
	"time"
	
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	
	"github.com/wanglu119/ng-ssh/lib/config"
	"github.com/wanglu119/ng-ssh/lib/ssh"
)

type initMessage struct {
	Name string `json:"name,omitempty"`
	AuthToken string `json:"auth_token,omitempty"`
}

// Protocols defines the name of this protocol,
// which is supposed to be used to the subprotocol of Websocket streams.
var Protocols = []string{"webtty"}

const (
	// Unknown message type, maybe send by a bug
	UnknownInput = '0'
	// User input typically from a keyboard
	Input = '1'
	// Ping to the server
	Ping = '2'
	// Notify that the brower size has been changed
	ResizeTerminal = '3'
)

const (
	// Unknown message type, maybe set by a bug
	UnknownOutput = '0'
	// Normal output to the terminal
	Output = '1'
	// Pong to the brower
	Pong = '2'
	// Set window title of the terminal
	SetWindowTitle = '3'
	// Set terminal preference
	SetPreferences = '4'
	// Make terminal to reconnect 
	SetReconnect = '5'
)

var (
	errWsClosed = errors.New("ws closed")
	errSshClientClosed = errors.New("ssh client closed")
)

var upGrader = websocket.Upgrader {
	ReadBufferSize: 1024, 
	WriteBufferSize: 1024*1024*10,
	Subprotocols: Protocols,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var sshWs = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	closeReason := "unkonwn reason"
	
	defer func() {
		log.Warn(fmt.Sprintf("Connection closed by %s: %s", closeReason, r.RemoteAddr))
	}()
	
	log.Info(fmt.Sprintf("New client connected: %s", r.RemoteAddr))
	
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		closeReason = err.Error()
		return http.StatusInternalServerError, err
	}
	
	defer func() {
		if err != nil {
			safeMessage := base64.StdEncoding.EncodeToString(([]byte(fmt.Sprintf("error:%v", err))))
			conn.WriteMessage(websocket.TextMessage, []byte(safeMessage))
		}
		
		conn.Close()
	}()
	
	err = processWSConn(conn)
	
	switch err {
		case errWsClosed:
			closeReason = "client"
		case errSshClientClosed:
			closeReason = "ssh server"
		default:
			closeReason = fmt.Sprintf("an error: %s", err)
	}
	
	return 0, nil
})

func processWSConn(conn *websocket.Conn) error {
	wswConn := &wsWrapper{conn}
	
	typ, initLine, err := conn.ReadMessage()
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return errors.Wrapf(err, "failed to authenticate websocket connection")
	}
	if typ != websocket.TextMessage {
		err = errors.New("failed to authenticate websocket connection: invalid message type")
		log.Error(fmt.Sprintf("%v", err))
		return err
	}
	
	var init initMessage
	err = json.Unmarshal(initLine, &init)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return errors.Wrapf(err, "failed to authenticate websocket connection")
	}
	
	log.Info(fmt.Sprintf("websocket initMessage is: %v", init))
	
	// find ssh config by name
	sshConfigs := config.GetSshConfigs()
	sshConfig, ok := sshConfigs[init.Name]
	if !ok {
		err = errors.New(fmt.Sprintf("not find ssh config by: %s", init.Name))
		log.Error(fmt.Sprintf("%v", err))
		return err
	}
	
	err = runSshClient(sshConfig, wswConn)
	
	return err
}

// =====================================================================

func runSshClient(sshConfig *config.SshConfig, wswConn *wsWrapper) error {
	
	webSshClient, err := NewWebSshClient(sshConfig, wswConn)
	if err != nil {
		return err
	}
	
	err = webSshClient.Run()
	
	return err
}

type WebSshClient struct {
	wswConn *wsWrapper
	sshSession *ssh.LogicSshSession
	
	bufferSize int
	writeMutex sync.Mutex
}

func NewWebSshClient(sshConfig *config.SshConfig, wswConn *wsWrapper) (*WebSshClient, error) {
	mc := sshConfigToMachine(sshConfig)
	
	sshSession, err := ssh.NewLogicSshSession(mc)
	if err != nil {
		return nil, err
	}
	
	return &WebSshClient{
		sshSession: sshSession,
		wswConn: wswConn,
		bufferSize: 1024,
	}, nil
}

func (wsc *WebSshClient) Run() error {
	
	var err error
	errs := make(chan error, 2)
	
	go func() {
		ticker := time.NewTicker(100*time.Millisecond)
		defer ticker.Stop()
		
		var handleErr error
		for {
			select {
				case  <- ticker.C: {
					handleErr = wsc.handleSshSessionReadEvent()
					if handleErr != nil {
						errs <- handleErr
						return
					}
				}
			}
		}
		
	}()
	
	go func() {
		errs <- func() error {
			buffer := make([]byte, wsc.bufferSize) 
			for {
				n, err := wsc.wswConn.Read(buffer)
				if err != nil {
					return errWsClosed
				}
				
				err = wsc.handleWsReadEvent(buffer[:n])
				if err != nil {
					return err
				}
			}
		}()
	}()
	
	select {
		case err = <- errs:
	}
	
	return err
}

func (wsc *WebSshClient) handleSshSessionReadEvent() error {
	data := wsc.sshSession.Read()
	if len(data) > 0 {
		if len(data)>10000 {
			data = data[len(data)-10000:]
		}
		safeMessage := base64.StdEncoding.EncodeToString(data)
		err := wsc.wsWrite(append([]byte{Output},[]byte(safeMessage)...))
		if err != nil {
			log.Error(fmt.Sprintf("%v", err))
			return errors.Wrapf(err, "failed to send message to master")
		}
	}
	
	return nil
}

func (wsc *WebSshClient) wsWrite(data []byte) error {
	wsc.writeMutex.Lock()
	defer wsc.writeMutex.Unlock()
	
	_, err := wsc.wswConn.Write(data)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return errors.Wrapf(err, "failed to write to websocket") 
	}
	return nil
}

func (wsc *WebSshClient) handleWsReadEvent(data []byte) error {
	if len(data) == 0 {
		err := errors.New("unexpected zero length read from master")
		log.Error(fmt.Sprintf("%v", err))
		return err
	}
	
	switch data[0] {
		case Input:
			if len(data) <= 1 {
				return nil
			}
			err := wsc.sshSession.Write(data[1:])
			if err != nil {
				log.Error(fmt.Sprintf("%v", err))
				return errors.Wrapf(err, "failed to write received data to slave")
			}
		case Ping:
			_,err := wsc.wswConn.Write([]byte{Pong})
			if err != nil {
				log.Error(fmt.Sprintf("%v", err))
				return errors.Wrapf(err, "failed to return Pong message to master")
			}
		case ResizeTerminal:
			if len(data) <= 1 {
				err := errors.New("received malformed remote command for terminal resize: empty payload")
				log.Error(fmt.Sprintf("%v", err))
				return err
			}
			
			var args argResizeTerminal
			err := json.Unmarshal(data[1:], &args)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err))
				return errors.Wrapf(err, "received malformed data for terminal resize")
			}
			err = wsc.sshSession.ResizeTerminal(args.Rows, args.Columns)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err))
			}
		default:
			err := errors.Errorf("unknown message type '%c'", data[0])
			log.Error(fmt.Sprintf("%v", err))
			return err
	}
	
	return nil
}

type argResizeTerminal struct {
	Columns int
	Rows int
}



