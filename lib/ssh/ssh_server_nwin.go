// +build !windows

package ssh

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"unsafe"
	"net"

	"github.com/gliderlabs/ssh"
	"github.com/kr/pty"
	
	"github.com/wanglu119/ng-ssh/lib/config"
	"github.com/wanglu119/ng-ssh/lib/util"
)


func setWinsize(f *os.File, w, h int) {
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
}

func RunSshServer() {
	go func() {
		ssh.Handle(func(s ssh.Session) {
			cmd := exec.Command("/bin/bash")
			ptyReq, winCh, isPty := s.Pty()
			if isPty {
				cmd.Env = append(cmd.Env, fmt.Sprintf("TERM=%s", ptyReq.Term))
				f, err := pty.Start(cmd)
				if err != nil {
					panic(err)
				}
				go func() {
					for win := range winCh {
						setWinsize(f, win.Width, win.Height)
					}
				}()
				go func() {
					io.Copy(f, s) // stdin
				}()
				io.Copy(s, f) // stdout
				cmd.Wait()
			} else {
				io.WriteString(s, "No PTY requested.\n")
				s.Exit(1)
			}
		})
		passwordOption := func(server *ssh.Server) error {
			server.PasswordHandler = func (ctx ssh.Context, password string) bool {
				user := ctx.User()
				return user == "ng-ssh"
			}
			return nil
		}
		
		port := util.GetFreePort()
		var listener net.Listener
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil{
			log.Error(fmt.Sprintf("%v", err))
			return
		}
		defer listener.Close()
		
		sshConfig := &config.SshConfig{
			Name: "ngssh-server",
			Host: "localhost",
			Port: uint(port),
			Type: config.SSH_CONFIG_TYPE_PASSWORD,
			Username: "ng-ssh",
			Password: "",
		}
		sshConfigs := config.GetSshConfigs()
		sshConfigs[sshConfig.Name] = sshConfig
		
		ssh.Serve(listener, nil, passwordOption)
	}()
}
