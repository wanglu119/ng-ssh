package sftp

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"path/filepath"
	"net/http"
	
	"github.com/pkg/sftp"
	"github.com/mholt/archiver"
	
	"github.com/wanglu119/ng-ssh/lib/common"
)

func SftpFetchFile(fullPath string, mc *common.Machine) (*sftp.File, os.FileInfo, error) {
	sftpClient, err := NewSftpClient(mc)
	if err != nil {
		return nil, nil, err
	}
	defer sftpClient.Close()
	
	fileInfo, err := sftpClient.Stat(fullPath)
	if err != nil {
		return nil, nil, err
	}
	if fileInfo.IsDir() {
		return nil, nil, fmt.Errorf("%s is not a file", fullPath)
	}
	f, err := sftpClient.Open(fullPath)
	return f, fileInfo, err
}

func setContentDisposition(w http.ResponseWriter, r *http.Request, fileName string) {
	if r.URL.Query().Get("inline") == "true" {
		w.Header().Set("Content-Disposition", "inline")
	} else {
		// As per RFC6266 section 4.3
		w.Header().Set("Content-Disposition", "attachment; filename*=utf-8''"+fileName)
	}
}

func SftpFetchFileToResponse(fullPath string, mc *common.Machine,w http.ResponseWriter, r *http.Request) error {
	sftpClient, err := NewSftpClient(mc)
	if err != nil {
		return err
	}
	defer sftpClient.Close()
	
	fileInfo, err := sftpClient.Stat(fullPath)
	if err != nil {
		return  err
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("%s is not a file", fullPath)
	}
	f, err := sftpClient.Open(fullPath)
	if err != nil {
		return err
	}
	setContentDisposition(w, r, fileInfo.Name())
	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), f)
	return  err
}

func SftpDownloadLocal(fullPath string, mc *common.Machine) {
	file, _, err := SftpFetchFile(fullPath, mc)
	defer file.Close()
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return
	}
	/*
	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, fileInfo.Name()),
	}
	c.DataFromReader(http.StatusOK, fileInfo.Size(), "application/octet-stream", file, extraHeaders)
	*/
}
func SftpCat(fullPath string, mc *common.Machine) {
	file, fileInfo, err := SftpFetchFile(fullPath, mc)
	defer file.Close()
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return
	}
	//c.String(200,"utf-8",file)
//	c.AbortWithStatusJSON(200, gin.H{"ok": true, "data": string(b), "msg": fileInfo.Name()})
	fmt.Println(fileInfo.Name())
	fmt.Println(string(b))
}

func SftpFetchToArchive(dirPath string, filenames []string, ar archiver.Writer, extension string, mc *common.Machine) error {
	sftpClient, err := NewSftpClient(mc)
	if err != nil {
		return err
	}
	defer sftpClient.Close()
	
	_,err = sftpClient.Stat(dirPath)
	if err != nil {
		return err
	}
	
	dirPath = filepath.Dir(filenames[0])
	
	for _, fname := range filenames {
		err := archiveAddFile(ar, fname,dirPath, sftpClient)
		if err != nil {
			return err
		}
	}
	
	return nil
}


func archiveAddFile(ar archiver.Writer, path,prefixPath string,  sftpClient *sftp.Client) error {
	// Checks are always done with paths with "/" as path separator.
	path = strings.Replace(path, "\\", "/", -1)
	
	info, err := sftpClient.Stat(path)
	if err != nil {
		log.Error(fmt.Sprintf("path: %s ,error: %v ",path, err))
		return nil
	}
	
	if info.IsDir() {
		subFiles, err := sftpClient.ReadDir(path)
		if err != nil {
			log.Error(fmt.Sprintf("path: %s, error: %v", subFiles, err))
			return err
		}
		
		for _,sfile := range subFiles {
			err = archiveAddFile(ar, filepath.Join(path, sfile.Name()),prefixPath,sftpClient)
			if err != nil {
				return err
			}
		} 
	} else {
		file, err := sftpClient.Open(path)
		if err != nil {
			log.Error(fmt.Sprintf("path: %s, error: %v", path, err))
			return err
		}
		defer file.Close()
		
		err = ar.Write(archiver.File{
			FileInfo: archiver.FileInfo{
				FileInfo: info,
				CustomName: strings.TrimPrefix(strings.TrimPrefix(path, prefixPath),"/"),
			},
			ReadCloser: file,
		})
		
		if err != nil {
			log.Error(fmt.Sprintf("error: %v", err))
			return err
		}
	}
	
	return nil
}



