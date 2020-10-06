package config

import (
	"strconv"
)

const (
	SERVER_MODE_LOCAL = "local"
	SERVER_MODE_NG = "ng"
)

type Server struct {
	Address string
	Port int
	Token string
	Mode string
	Username string
	Password string
}

func (s *Server) GetFullAddress() string {
	return s.Address+":"+strconv.Itoa(s.Port)
}






