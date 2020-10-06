package ssh

import (
	"bytes"
	"sync"
	"io"
	"strings"
	"fmt"
	
	"golang.org/x/crypto/ssh"
	
	"github.com/wanglu119/ng-ssh/lib/common"
)

const sudoPrefix, sudoSuffix = "[sudo] password for ", ": "
const sudoPrefixLen = len(sudoPrefix)

type safeBuffer struct {
	buffer bytes.Buffer
	mu     sync.Mutex
}

func (w *safeBuffer) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Write(p)
}
func (w *safeBuffer) Bytes() []byte {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Bytes()
}
func (w *safeBuffer) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buffer.Reset()
}

type LogicSshSession struct {
	stdinPipe io.WriteCloser
	// ssh terminal output
	comboOutput *safeBuffer
	session *ssh.Session
	client *ssh.Client
	mc *common.Machine
}

func NewLogicSshSession(mc *common.Machine) (*LogicSshSession, error) {
	
	sshClient, err := NewSshClient(mc)
	if err != nil {
		return nil, err
	}
	
	sshSession, err := sshClient.NewSession()
	if err != nil {
		return nil, err
	}
	
	stdinP, err := sshSession.StdinPipe()
	if err != nil {
		return nil, err
	}
	
	comboWriter := new(safeBuffer)
	//ssh.stdout and stderr will write output into comboWriter
	sshSession.Stdout = comboWriter
	sshSession.Stderr = comboWriter
	
	modes := ssh.TerminalModes {
		ssh.ECHO:          1,     // disable echo
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	cols, rows := 80, 40
	// Request pseudo terminal
	if err := sshSession.RequestPty("xterm", rows, cols, modes); err != nil {
		return nil, err
	}
	
	// Start remote shell
	if err := sshSession.Shell(); err != nil {
		return nil, err
	}
	
	return &LogicSshSession{
		stdinPipe: stdinP,
		comboOutput: comboWriter,
		session: sshSession,
		client: sshClient,
		mc: mc,
	}, nil
}

// Close
func (ss *LogicSshSession) Close() {
	if ss.session != nil {
		ss.session.Close()
	}
	if ss.client != nil {
		ss.client.Close()
	}
	if ss.comboOutput != nil {
		ss.comboOutput = nil
	}
}

func (ss *LogicSshSession) Read() []byte {
	bs := ss.comboOutput.Bytes()
	if len(bs) > 0 {
		ss.comboOutput.Reset()
		if ss.mc.Name != "ngssh-server" {
			line := string(bs)
			lines := strings.Split(line, "\n")
			line = lines[len(lines)-1]
			
			if len(line) >= sudoPrefixLen && strings.HasPrefix(line, sudoPrefix) && strings.HasSuffix(line, sudoSuffix) && strings.Contains(line, ss.mc.User) {
				err := ss.Write([]byte(ss.mc.Password + "\n"))
				if err == nil {
					msg := fmt.Sprintf("\r\n\033[32m ng-ssh has automatically input password for\033[0m\033[33m %s \033[0m", ss.mc.User)
					bmsg := []byte(msg)
					bs = append(bs, bmsg...)
				}
			}
		}
	}
	return bs
}

func (ss *LogicSshSession) Write(data []byte) error {
	if _, err := ss.stdinPipe.Write(data); err != nil {
		return err
	}
	
	return nil
}

func (ss *LogicSshSession) ResizeTerminal(rows, cols int) error {
	return ss.session.WindowChange(rows, cols)
}



