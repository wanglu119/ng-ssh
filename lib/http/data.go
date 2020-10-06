package http

import (
	"net/http"
	"fmt"
	"strconv"
	
	"github.com/wanglu119/ng-ssh/lib/config"
	"github.com/wanglu119/ng-ssh/lib/common"
)

type data struct {
	server *config.Server
}

type handleFunc func(w http.ResponseWriter, r *http.Request, d *data) (int, error)

func handle(fn handleFunc, prefix string) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		status, err := fn(w, r, &data{
			server: config.GetServer(),
		})
		
		if status != 0 {
			txt := http.StatusText(status)
			http.Error(w, strconv.Itoa(status)+" "+txt, status)
		}
		
		if status >= 400 || err != nil {
			log.Error(fmt.Sprintf("%s: %v %v", r.URL.Path, status, err))
		}
	})
	
	return http.StripPrefix(prefix, handler)
}

func sshConfigToMachine(sshConfig *config.SshConfig) *common.Machine {
	return &common.Machine{
		Name: sshConfig.Name,
		Host: sshConfig.Host,
		Port: sshConfig.Port,
		User: sshConfig.Username,
		Password: config.Decode(sshConfig),
		Type: sshConfig.Type,
	}
}

