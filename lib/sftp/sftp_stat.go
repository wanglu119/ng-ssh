package sftp

import (
	"fmt"
	
	"github.com/wanglu119/ng-ssh/lib/common"
)


func SftpStat(dirPath string, mc *common.Machine) (*Ls, error) {
	sftpClient, err := NewSftpClient(mc)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return nil, err
	}
	defer sftpClient.Close()
	
	stat, err := sftpClient.Stat(dirPath)
	if err != nil {
		return nil, err
	}
	
	ls := &Ls{
	 	Name: stat.Name(), 
	 	Size: stat.Size(), 
	 	Path: dirPath, 
	 	Time: stat.ModTime(), 
	 	Mod: stat.Mode().String(), 
	 	IsDir: stat.IsDir(),
	 }
	
	return ls, nil 
}


