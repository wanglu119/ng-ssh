package sftp

import (
	"path"
	"fmt"
	
	"github.com/pkg/sftp"
	
	"github.com/wanglu119/ng-ssh/lib/common"
)

func SftpRmFile(fullPath string, mc *common.Machine) error {
	sftpClient, err := NewSftpClient(mc)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	err = sftpClient.Remove(fullPath)
	if err != nil {
		return err
	}
	
	return nil
}

func rmR(subFullPath string, sftpClient *sftp.Client) error {
	files, err := sftpClient.ReadDir(subFullPath)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			err = rmR(path.Join(subFullPath,file.Name()), sftpClient)
			if err != nil {
				return err
			}
		} else {
			err = sftpClient.Remove(path.Join(subFullPath,file.Name()))
			if err != nil {
				return err
			}
		}
	}
	
	return sftpClient.RemoveDirectory(subFullPath)
}

func SftpRmDir(fullPath string, mc *common.Machine) error {
	sftpClient, err := NewSftpClient(mc)
	if err != nil {
		return err
	}
	defer sftpClient.Close()
	
	err = rmR(fullPath, sftpClient)
	if err != nil {
		return err
	}
	
	return nil
}

