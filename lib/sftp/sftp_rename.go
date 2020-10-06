package sftp

import (
	"github.com/wanglu119/ng-ssh/lib/common"
)

func SftpRename(oldFullPath, newFullPath string, mc *common.Machine) error {
	
	sftpClient, err := NewSftpClient(mc)
	if err != nil {
		return err
	}
	err = sftpClient.Rename(oldFullPath, newFullPath)
	if err != nil {
		return err
	}
	
	return nil
}