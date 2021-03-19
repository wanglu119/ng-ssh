package http

import (
	"net/http"
	"encoding/json"
	"fmt"
	
	"github.com/wanglu119/ng-ssh/lib/config"
	"github.com/wanglu119/ng-ssh/lib/common"
	"github.com/wanglu119/ng-ssh/lib/sftp"
)

type sftpParam struct {
	Name string `json:"name,omitempty"`
	DirPath string `json:"dir_path,omitempty"`
	OldFullPath string `json:"old_full_path,omitempty"`
	NewFullPath string `json:"new_full_path,omitempty"`
	IsDir bool `json:"is_dir,omitempty"`
}

func parseSftpParam(r *http.Request) (*sftpParam, *common.Machine, error) {
	var param sftpParam
	err := json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return nil,nil, err
	}
	sshConfigs := config.GetSshConfigs()
	sshConfig, ok := sshConfigs[param.Name]
	if !ok {
		return nil,nil, err
	}
	
	mc := sshConfigToMachine(sshConfig)
	if param.DirPath == "" {
		param.DirPath = "$HOME"
	}
	
	return &param, mc, nil
}

var sftpListFiles = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	param, mc, err := parseSftpParam(r)
	if err != nil {
		return http.StatusBadRequest, nil
	}
	
	fileList,lsMeta,err := sftp.SftpLs(param.DirPath, mc)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	
	return renderJSON(w, r, map[string]interface{}{"ls_meta":lsMeta,"file_list":fileList})
})

var sftpRname = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	param, mc, err := parseSftpParam(r)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return http.StatusBadRequest, nil
	}
	
	if param.OldFullPath == "" || param.NewFullPath == "" {
		return http.StatusBadRequest, nil
	}
	
	err = sftp.SftpRename(param.OldFullPath, param.NewFullPath, mc)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, err
	}
	
	return 0,nil
})

var sftpCreateFile = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	param, mc, err := parseSftpParam(r)
	if err != nil {
		return http.StatusBadRequest, nil
	}
	
	if param.NewFullPath == "" {
		return http.StatusBadRequest, nil
	}
	
	err = sftp.SftpCreateFile(param.NewFullPath, mc)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, err
	}
	
	return 0,nil
})

var sftpCreateDir = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	param, mc, err := parseSftpParam(r)
	if err != nil {
		return http.StatusBadRequest, nil
	}
	
	if param.NewFullPath == "" {
		return http.StatusBadRequest, nil
	}
	
	err = sftp.SftpCreateDir(param.NewFullPath, mc)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, err
	}
	
	return 0,nil
})

var sftpRemoveFile = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	param, mc, err := parseSftpParam(r)
	if err != nil {
		return http.StatusBadRequest, nil
	}
	
	if param.OldFullPath == "" {
		return http.StatusBadRequest, nil
	}
	
	err = sftp.SftpRmFile(param.OldFullPath, mc)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, err
	}
	
	return 0,nil
})

var sftpRemoveDir = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	param, mc, err := parseSftpParam(r)
	if err != nil {
		return http.StatusBadRequest, nil
	}
	
	if param.OldFullPath == "" {
		return http.StatusBadRequest, nil
	}
	
	err = sftp.SftpRmDir(param.OldFullPath, mc)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, err
	}
	
	return 0,nil
})

var sftpStat = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	param, mc, err := parseSftpParam(r)
	if err != nil {
		return http.StatusBadRequest, nil
	}
	
	if param.DirPath == "" {
		return http.StatusBadRequest, nil
	}
	
	stat, err := sftp.SftpStat(param.DirPath, mc)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, err
	}
	
	return renderJSON(w, r, stat)
})

var sftpCopy = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	param, mc, err := parseSftpParam(r)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return http.StatusBadRequest, nil
	}
	
	if param.OldFullPath == "" || param.NewFullPath == "" {
		return http.StatusBadRequest, nil
	}
	
	err = sftp.SftpCopy(param.OldFullPath, param.NewFullPath, mc)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, err
	}
	
	return 0, nil
})





