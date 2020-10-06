package ssh

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	
	"github.com/wanglu119/ng-ssh/lib/common"
)

func NewSshClient(mc *common.Machine) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		Timeout:         time.Second * 5,
		User:            mc.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
		//HostKeyCallback: hostKeyCallBackFunc(mc.Host),
	}
	if mc.Type == "password" {
		config.Auth = []ssh.AuthMethod{ssh.Password(mc.Password)}
	} else {
		config.Auth = []ssh.AuthMethod{publicKeyAuthFunc(mc.Key)}
	}
	addr := fmt.Sprintf("%s:%d", mc.Host, mc.Port)
	c, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func hostKeyCallBackFunc(host string) ssh.HostKeyCallback {
	hostPath, err := homedir.Expand("~/.ssh/known_hosts")
	if err != nil {
		log.Error(fmt.Sprintf("find known_hosts's home dir failed: %v", err))
	}
	file, err := os.Open(hostPath)
	if err != nil {
		log.Error(fmt.Sprintf("can't find known_host file: %v", err))
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				log.Error(fmt.Sprintf("error parsing %q: %v", fields[2], err))
			}
			break
		}
	}
	if hostKey == nil {
		log.Error(fmt.Sprintf("no hostkey for %s,%v", host, err))
	}
	return ssh.FixedHostKey(hostKey)
}

func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
	keyPath, err := homedir.Expand(kPath)
	if err != nil {
		log.Error(fmt.Sprintf("find key's home dir failed: %v", err))
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Error(fmt.Sprintf("ssh key file read failed: %v", err))
	}
	// CreateUserOfRole the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Error(fmt.Sprintf("ssh key signer failed: %v", err))
	}
	return ssh.PublicKeys(signer)
}
func runCommand(client *ssh.Client, command string) (stdout string, err error) {
	session, err := client.NewSession()
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return
	}
	defer session.Close()

	var buf bytes.Buffer
	session.Stdout = &buf
	err = session.Run(command)
	if err != nil {
		//log.Print(err)
		return
	}
	stdout = string(buf.Bytes())

	return
}
