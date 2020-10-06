package http

import (
	"net/http"
	"encoding/json"
	"fmt"
	
	"github.com/wanglu119/ng-ssh/lib/config"
)

var getSshConfigs = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	sshConfigs := config.GetSshConfigs()
	return renderJSON(w, r, sshConfigs)
})

var getSshConfigNames = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	sshConfigs := config.GetSshConfigs()
	configNames := []string{}
	for _,conf := range sshConfigs {
		configNames = append(configNames, conf.Name)
	}
	return renderJSON(w, r, configNames)
})

var getSftpConfigNames = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	sshConfigs := config.GetSshConfigs()
	configNames := []string{}
	for _,conf := range sshConfigs {
		if conf.Name != "ngssh-server" {
			configNames = append(configNames, conf.Name)
		}
	}
	return renderJSON(w, r, configNames)
})

var addSshConfig = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	sshConfig := &config.SshConfig{}
	err := json.NewDecoder(r.Body).Decode(sshConfig)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, err
	}
	sshConfig.Type = config.SSH_CONFIG_TYPE_PASSWORD
	config.Encode(sshConfig)
	sshConfigs := config.GetSshConfigs()
	sshConfigs[sshConfig.Name] = sshConfig
	config.UpdateConfig()
	
	return 0, nil
})

var deleteSshConfig = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	sshConfig := &config.SshConfig{}
	err := json.NewDecoder(r.Body).Decode(sshConfig)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, err
	}
	
	sshConfigs := config.GetSshConfigs()
	if _, ok := sshConfigs[sshConfig.Name]; !ok {
		return http.StatusBadRequest, nil
	}
	delete(sshConfigs, sshConfig.Name)
	config.UpdateConfig()
	
	return 0, nil
})
