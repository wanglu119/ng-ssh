package sftp

import (
	"path"
	"fmt"
	"io"
	
	"github.com/pkg/sftp"
	
	"github.com/wanglu119/ng-ssh/lib/common"
)

func subDirCopy(srcFullPath, destFullPath string, sftpClient *sftp.Client) error{
	
	srcStat, err := sftpClient.Stat(srcFullPath)
	if err != nil {
		log.Error(fmt.Sprintf("%v",err))
		return err
	}
	
	if srcStat.IsDir() {
		err = sftpClient.MkdirAll(destFullPath)
		if err != nil {
			return err
		}
		files, err := sftpClient.ReadDir(srcFullPath)
		if err != nil {
			return err
		}
		for _, file := range files {
			if file.IsDir() {
				err = subDirCopy(path.Join(srcFullPath, file.Name()), path.Join(destFullPath, file.Name()), sftpClient)
				if err != nil {
					return err
				}
			} else {
				err = subDirCopy(path.Join(srcFullPath, file.Name()), path.Join(destFullPath, file.Name()), sftpClient)
				if err != nil {
					return err
				}
			}
		}
	} else {
		srcFile, err := sftpClient.Open(srcFullPath)
		if err != nil {
			return err
		}
		defer srcFile.Close()
		destFile, err := sftpClient.Create(destFullPath)
		if err != nil {
			return err
		}
		defer destFile.Close()
		io.Copy(destFile, srcFile)
	}
	
	return nil
}

func SftpCopy(srcFullPath, destFullPath string, mc *common.Machine) error {
	sftpClient, err := NewSftpClient(mc)
	if err != nil {
		return err
	}
	defer sftpClient.Close()
	
	return subDirCopy(srcFullPath, destFullPath, sftpClient)
}

