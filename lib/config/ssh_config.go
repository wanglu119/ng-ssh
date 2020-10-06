package config

import (
	"crypto/rc4"
	"encoding/base64"
)

const (
	SSH_CONFIG_TYPE_PASSWORD = "password"
)


type SshConfig struct {
	Name string `json:"name" yaml:"name" mapstructure:"name"`
	Host string `json:"host" yaml:"host" mapstructure:"host"`
	Port uint `json:"port" yaml:"port" mapstructure:"port"`
	Username string `json:"username" yaml:"username" mapstructure:"username"`
	Password string `json:"password" yaml:"password" mapstructure:"password"`
	Type string `json:"type" yaml:"type" mapstructure:"type"`
}

func Encode(sshConfig *SshConfig) {
	key := []byte(sshConfig.Name)
	src := []byte(sshConfig.Password)
	
	cipher, _ := rc4.NewCipher(key)
	cipher.XORKeyStream(src, src)
	sshConfig.Password = base64.StdEncoding.EncodeToString(src)
}

func Decode(sshConfig *SshConfig) string {
	key := []byte(sshConfig.Name)
	src,_ := base64.StdEncoding.DecodeString(sshConfig.Password)
	
	cipher, _ := rc4.NewCipher(key)
	cipher.XORKeyStream(src, src)
	return string(src)
}


