package sftp

import (
	"fmt"
	"mime/multipart"
	"path"
	"io"
	"os"
	"net/http"
	
	"github.com/pkg/sftp"
	
	"github.com/wanglu119/ng-ssh/lib/common"
)

func SftpUpFileHeader(header *multipart.FileHeader, desDir string, mc *common.Machine) {
	sftpClient, err := NewSftpClient(mc)
	if err != nil {
		log.Error(fmt.Sprintf("%v",err))
		return
	}
	defer sftpClient.Close()
	
	err = uploadFileHeader(desDir, sftpClient, header)
	if err != nil {
		log.Error(fmt.Sprintf("%v",err))
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


// ==========================================================

func SftpUpRequest(desDir string, r *http.Request, w http.ResponseWriter, mc *common.Machine) error {
	sftpClient, err := NewSftpClient(mc)
	if err != nil {
		return err
	}
	defer sftpClient.Close()
	
	dstFile, err := sftpClient.Create(desDir)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	
	_, err = io.Copy(dstFile, r.Body)
	if err != nil {
		sftpClient.Remove(desDir)
		return err
	}
	
	// Gets the info about the file.
	info, err := dstFile.Stat()
	if err != nil {
		sftpClient.Remove(desDir)
		return err
	}
	
	etag := fmt.Sprintf(`"%x%x"`, info.ModTime().UnixNano(), info.Size())
	w.Header().Set("ETag", etag)
	
	return nil
}


