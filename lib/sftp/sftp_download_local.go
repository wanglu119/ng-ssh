package sftp

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"path/filepath"
	
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
func SftpDownloadLocal(fullPath string, mc *common.Machine) {
	file, fileInfo, err := SftpFetchFile(fullPath, mc)
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
	fmt.Println(fileInfo.Name())
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

func SftpFetchToArchive(dirPath string, filenames []string, ar archiver.Writer, extension string, mc *common.Machine) (string, error) {
	sftpClient, err := NewSftpClient(mc)
	if err != nil {
		return "", err
	}
	defer sftpClient.Close()
	
	fileInfo,err := sftpClient.Stat(dirPath)
	if err != nil {
		return "", err
	}
	
	name := fileInfo.Name()
	if name == "." || name == "" {
		name = "archive"
	}
	name += extension
	
	for _, fname := range filenames {
		err := archiveAddFile(ar, fname,dirPath, sftpClient)
		if err != nil {
			return "",err
		}
	}
	
	return name, nil
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



