package util

import (
	"net"
)

func GetFreePort() (port int) {
	l,_ := net.Listen("tcp", ":0")
	port = l.Addr().(*net.TCPAddr).Port
	
	l.Close()
	
	return 
}