package sftp

import (
	"fmt"
	"mime/multipart"
	"path"
	"io"
	"os"
	
	"github.com/pkg/sftp"
	
	"github.com/wanglu119/ng-ssh/lib/common"
)

func SftpUpFileHeader(header *multipart.FileHeader, desDir string, mc *common.Machine) {
	sftpClient, err := NewSftpClient(mc)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sftpClient.Close()
	
	err = uploadFileHeader(desDir, sftpClient, header)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func uploadFileHeader(desDir string, client *sftp.Client, header *multipart.FileHeader) error {

	srcFile, err := header.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()
	
	return uploadReader(desDir, header.Filename, client, srcFile)
}

func uploadReader(desDir, filename string, client *sftp.Client, srcFile io.Reader) error {
	if desDir == "$HOME" {
		wd, err := client.Getwd()
		if err != nil {
			return err
		}
		desDir = wd
	}
	
	dstFile, err := client.Create(path.Join(desDir, filename))
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = dstFile.ReadFrom(srcFile)
	if err != nil {
		return err
	}
	return nil
}

func SftpUpFile(file *os.File, desDir string, mc *common.Machine) error {
	sftpClient, err := NewSftpClient(mc)
	if err != nil {
		return err
	}
	defer sftpClient.Close()
	
	err = uploadReader(desDir,file.Name(), sftpClient, file)
	if err != nil {
		return err
	}
	
	return nil
}





