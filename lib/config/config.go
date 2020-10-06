package config

import (
	"fmt"
	"reflect"
	"sync"
	
	"github.com/wanglu119/me-deps/viper"
)

var sshConfigsOnce sync.Once

func GetSshConfigs() map[string]*SshConfig {
	var sshConfigs map[string]*SshConfig
	
	vper := viper.GetViper()
	
	sshConfigsOnce.Do(func(){
		sub := vper.Sub("ssh_configs")
		if sub != nil {
			err := sub.Unmarshal(&sshConfigs)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err))
			}
			vper.Set("ssh_configs", sshConfigs)
		}
	})
	
	raw := vper.Get("ssh_configs")
	switch val := raw.(type) {
		case *viper.Viper:
			val.Unmarshal(&sshConfigs)
			vper.Set("ssh_configs", sshConfigs)
		case map[string]*SshConfig:
			sshConfigs = val
		case nil:
			sshConfigs = make(map[string]*SshConfig)
			vper.Set("ssh_configs", sshConfigs)
		default:
			log.Error(fmt.Sprintf("ssh_configs type error: %s", reflect.TypeOf(val).Name()))
	}

	
	return sshConfigs
}


func GetServer() *Server {
	var server *Server
	
	vper := viper.GetViper()
	
	sub := vper.Sub("server")
	if sub != nil {
		err := sub.Unmarshal(&server)
		if err != nil {
			log.Error(fmt.Sprintf("%v", err))
		}
		vper.Set("server", server)
	} else {
		raw := vper.Get("server")
		switch val := raw.(type){
			case *viper.Viper:
				val.Unmarshal(&server)
				vper.Set("server", server)
			case *Server:
				server = val
			default:
				log.Error("server type error")
		}
	}
	
	return server
}

func UpdateConfig() error {
	return viper.GetViper().WriteConfig()
}


