package http

import (
	"net/http"
	"net/url"
	"fmt"
	"errors"
	"strings"
	"path/filepath"
	"path"
	
	"github.com/mholt/archiver"
	
	"github.com/wanglu119/ng-ssh/lib/config"
	"github.com/wanglu119/ng-ssh/lib/common"
	"github.com/wanglu119/ng-ssh/lib/sftp"
)

func parseQueryIsDir(r *http.Request) (bool, error) {
	switch r.URL.Query().Get("is_dir") {
		case "true":
			return true, nil
		case "false":
			return false, nil
		default:
			return false, errors.New("not found param is_dir")
	}
}

func parseQueryMachine(r *http.Request) (*common.Machine, error) {
	name := r.URL.Query().Get("name")
	if name == "" {
		return nil, errors.New("not found param name")
	}
	sshConfigs := config.GetSshConfigs()
	sshConfig, ok := sshConfigs[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("not found machine by %s", name))
	}
	
	mc := sshConfigToMachine(sshConfig)
	
	return mc, nil
}

var sftpRaw = func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	isDir, err := parseQueryIsDir(r)
	if err != nil {
		return http.StatusBadRequest,err
	}
	mc, err := parseQueryMachine(r)
	if err != nil {
		return http.StatusBadRequest, err
	}
	
	if !isDir {
		sftpRawFile(w, r, r.URL.Path, mc)
	}
	
	return sftpRawDir(w, r, r.URL.Path, mc)
}

func setContentDisposition(w http.ResponseWriter, r *http.Request, fileName string) {
	if r.URL.Query().Get("inline") == "true" {
		w.Header().Set("Content-Disposition", "inline")
	} else {
		// As per RFC6266 section 4.3
		w.Header().Set("Content-Disposition", "attachment; filename*=utf-8''"+url.PathEscape(fileName))
	}
}

func sftpRawFile(w http.ResponseWriter, r *http.Request, fullPath string, mc *common.Machine) (int, error) {
	file, fileInfo, err := sftp.SftpFetchFile(fullPath, mc)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer file.Close()
	
	setContentDisposition(w, r, fileInfo.Name())
	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
	
	return 0, nil
} 

// ================================================================================

func slashClean(name string) string {
	if name == "" || name[0] != '/' {
		name = "/" + name
	}
	return path.Clean(name)
}

func parseQueryFiles(r *http.Request, dirPath string) ([]string, error) {
	var fileSlice []string
	names := strings.Split(r.URL.Query().Get("files"), ",")
	if len(names) == 0 {
		fileSlice = append(fileSlice, dirPath)
	} else {
		for _, name := range names {
			name, err := url.QueryUnescape(strings.Replace(name, "+", "%2B", -1)) //nolint:shadow
			if err != nil {
				return nil, err
			}

			name = slashClean(name)
			fileSlice = append(fileSlice, filepath.Join(dirPath, name))
		}
	}
	
	return fileSlice, nil
}


//nolint: goconst
func parseQueryAlgorithm(r *http.Request) (string, archiver.Writer, error) {
	// TODO: use enum
	switch r.URL.Query().Get("algo") {
	case "zip", "true", "":
		return ".zip", archiver.NewZip(), nil
	case "tar":
		return ".tar", archiver.NewTar(), nil
	case "targz":
		return ".tar.gz", archiver.NewTarGz(), nil
	case "tarbz2":
		return ".tar.bz2", archiver.NewTarBz2(), nil
	case "tarxz":
		return ".tar.xz", archiver.NewTarXz(), nil
	case "tarlz4":
		return ".tar.lz4", archiver.NewTarLz4(), nil
	case "tarsz":
		return ".tar.sz", archiver.NewTarSz(), nil
	default:
		return "", nil, errors.New("format not implemented")
	}
}

func sftpRawDir(w http.ResponseWriter, r *http.Request, fullPath string, mc *common.Machine) (int, error) {
	dirPath := r.URL.Path
	
	filenames, err := parseQueryFiles(r, dirPath)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return http.StatusBadRequest, err
	}
	
	extension, ar, err := parseQueryAlgorithm(r)
	if err != nil {
		return http.StatusBadRequest, err
	}
	
	err = ar.Create(w)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, err
	}
	defer ar.Close()
	
	name, err := sftp.SftpFetchToArchive(dirPath, filenames, ar, extension, mc)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, err
	}
	
	w.Header().Set("Content-Disposition", "attachment; filename*=utf-8''"+url.PathEscape(name))
	
	return 0,nil
}







