package sftp

import (
	"testing"
	"fmt"
	
	"github.com/wanglu119/ng-ssh/lib/common"
)

func TestCopy(t *testing.T) {
	mc := &common.Machine {
		Ip: "10.130.17.154",
		Port: 22,
		User: "tusimple",
		Password: "tusimple2017",
		Type: "password",
	}
	srcFullPath := "/home/tusimple/test/fio"
	destFullPath := "/home/tusimple/test/tmp/t"
	
	info, err := SftpStat(srcFullPath, mc)
	if err != nil {
		t.Fatal(err)
	} 
	fmt.Println(info.Name())
	
	err = SftpCopy(srcFullPath, destFullPath, mc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestXXX(t *testing.T) {
	fmt.Println("------------------------")
}
