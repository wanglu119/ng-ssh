package sftp

import (
	"errors"
	"fmt"
	
	"github.com/wanglu119/ng-ssh/lib/common"
)

func SftpCreateFile(fullPath string, mc *common.Machine) error {
	
	sftpClient, err := NewSftpClient(mc)
	if err != nil {
		return err
	}
	defer sftpClient.Close()
	
	_, err = sftpClient.Stat(fullPath)
	if err == nil {
		return errors.New(fmt.Sprintf("%s has exsit", fullPath))
	}
	
	_, err = sftpClient.Create(fullPath)
	if err != nil {
		return err
	}
	
	return nil
}

func SftpCreateDir(fullPath string, mc *common.Machine) error {
	sftpClient, err := NewSftpClient(mc)
	if err != nil {
		return err
	}
	defer sftpClient.Close()
	
	err = sftpClient.MkdirAll(fullPath)
	if err != nil {
		return err
	}
	
	return nil
}
