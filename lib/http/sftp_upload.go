package http

import (
	"net/http"
	"io"
	"io/ioutil"
	"strings"
	"path/filepath"
	
	"github.com/wanglu119/ng-ssh/lib/sftp"
)

var sftpResourcePostPutHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	defer func() {
		_, _ = io.Copy(ioutil.Discard, r.Body)
	}()
	
	mc, err := parseQueryMachine(r)
	if err != nil {
		return http.StatusBadRequest, err
	}
	
	// For directories, only allow POST for creation.
	if strings.HasSuffix(r.URL.Path, "/") {
		if r.Method == http.MethodPut {
			return http.StatusMethodNotAllowed, nil
		}
		
		err = sftp.SftpCreateDir(r.URL.Path, mc)
		if err != nil {
			return http.StatusInternalServerError, err
		} else {
			return 0, nil
		}
	}
	
	if r.Method == http.MethodPost && r.URL.Query().Get("override") != "true" {
		if _, err := sftp.SftpStat(r.URL.Path, mc); err == nil {
			return http.StatusConflict, nil
		}
	}
	
//	action := "upload"
//	if r.Method == http.MethodPut {
//		action = "save"
//	}
	
	dir, _ := filepath.Split(r.URL.Path)
	err = sftp.SftpCreateDir(dir, mc)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	
	err = sftp.SftpUpRequest(r.URL.Path, r, w, mc)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	
	return 0, nil
	
})



