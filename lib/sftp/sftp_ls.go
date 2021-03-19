package sftp

import (
	"fmt"
	"time"
	"path"
	"sort"
	
	"github.com/wanglu119/ng-ssh/lib/common"
)

type Ls struct {
	Name  string    `json:"name"`
	Path  string    `json:"path"` // including Name
	Size  int64     `json:"size"`
	Time  time.Time `json:"time"`
	Mod   string    `json:"mod"`
	IsDir bool      `json:"is_dir"`
	Index int 		`json:"index"`
}

type LsMeta struct {
	DirPath string `json:"dir_path"`
	FileCount int `json:"file_count"`
	DirCount int `json:"dir_count"`
	
	Name  string    `json:"name"`
	Size  int64     `json:"size"`
	Time  time.Time `json:"time"`
	Mod   string    `json:"mod"`
	IsDir bool      `json:"is_dir"`
}


func SftpLs(dirPath string, mc *common.Machine) ([]*Ls,*LsMeta, error) {
	sftpClient, err := NewSftpClient(mc)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return nil,nil, err
	}
	defer sftpClient.Close()

	if dirPath == "$HOME" {
		wd, err := sftpClient.Getwd()
		if err != nil {
			log.Error(fmt.Sprintf("%v", err))
			return nil,nil, err
		}
		dirPath = wd
	}
	
	currStat, err := sftpClient.Stat(dirPath)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return nil,nil,err
	}
	
	files, err := sftpClient.ReadDir(dirPath)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return nil,nil,err
	}
	// this will not be converted to null if slice is empty.
	fileList := make([]*Ls, 0) 
	dirList := make([]*Ls, 0) 
	
	for _, file := range files {
		tt := &Ls{Name: file.Name(), Size: file.Size(), Path: path.Join(dirPath, file.Name()), Time: file.ModTime(), Mod: file.Mode().String(), IsDir: file.IsDir()}
		if file.IsDir() {
			dirList = append(dirList, tt)
		} else {
			fileList = append(fileList, tt)
		}
	}
	sort.Slice(dirList, func(i,j int) bool {
		return dirList[i].Name < dirList[j].Name
	})
	sort.Slice(fileList, func(i,j int) bool {
		return fileList[i].Name < fileList[j].Name
	})
	// set index
	allFiles := make([]*Ls, 0) 
	index := 0
	for _, f := range dirList {
		f.Index = index
		index++
		allFiles = append(allFiles,f)
	}
	for _, f := range fileList {
		f.Index = index
		index++
		allFiles = append(allFiles,f)
	}
	
	
	lsMeta := &LsMeta{
		DirPath: dirPath,
		DirCount: len(dirList),
		FileCount: len(fileList),
		
		Name: currStat.Name(),
		IsDir: currStat.IsDir(),
		Mod: currStat.Mode().String(),
		Size: currStat.Size(),
		Time: currStat.ModTime(),
	}
	
	return allFiles,lsMeta, nil
}


