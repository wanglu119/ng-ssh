package http

import (
	"net/http"
	
	"github.com/wanglu119/me-deps/mux"
)

func Setup(r *mux.Router) {
	monkey := func(fn handleFunc, prefix string) http.Handler {
		return handle(fn, prefix)
	}
	
	api := r.PathPrefix("/ngssh_api").Subrouter()
	
	sshConfig := api.PathPrefix("/ssh_config").Subrouter()
	sshConfig.Handle("/getSshConfigs", monkey(getSshConfigs,"")).Methods("OPTIONS", "GET")
	sshConfig.Handle("/addSshConfig", monkey(addSshConfig,"")).Methods("OPTIONS", "POST")
	sshConfig.Handle("/deleteSshConfig", monkey(deleteSshConfig,"")).Methods("OPTIONS", "DELETE")
	sshConfig.Handle("/getSshConfigNames", monkey(getSshConfigNames,"")).Methods("OPTIONS", "GET")
	sshConfig.Handle("/getSftpConfigNames", monkey(getSftpConfigNames,"")).Methods("OPTIONS", "GET")
	
	sshApi := api.PathPrefix("/ssh").Subrouter()
	sshApi.Handle("/ws", monkey(sshWs,"")).Methods("OPTIONS", "GET")
	
	sftpApi := api.PathPrefix("/sftp").Subrouter()
	sftpApi.Handle("/listFiles", monkey(sftpListFiles, "")).Methods("OPTIONS", "POST")
	sftpApi.Handle("/rename", monkey(sftpRname,"")).Methods("OPTIONS", "PUT")
	sftpApi.Handle("/createFile", monkey(sftpCreateFile,"")).Methods("OPTIONS", "POST")
	sftpApi.Handle("/createDir", monkey(sftpCreateDir,"")).Methods("OPTIONS", "POST")
	sftpApi.Handle("/removeFile", monkey(sftpRemoveFile,"")).Methods("OPTIONS", "DELETE")
	sftpApi.Handle("/removeDir", monkey(sftpRemoveDir,"")).Methods("OPTIONS", "DELETE")
	sftpApi.PathPrefix("/raw").Handler(monkey(sftpRaw,"/ngssh_api/sftp/raw")).Methods("OPTIONS", "GET")
	
	authApi := api.PathPrefix("/auth").Subrouter()
	authApi.Handle("/renew", monkey(renew,"")).Methods("OPTIONS", "GET")
	authApi.Handle("/login", monkey(login,"")).Methods("OPTIONS", "POST")
	
	r.PathPrefix("/ng-sftp/").Handler(monkey(indexSftp,"")).Methods("OPTIONS", "GET")
	r.PathPrefix("/").Handler(monkey(indexSsh,"")).Methods("OPTIONS", "GET")
	
}


