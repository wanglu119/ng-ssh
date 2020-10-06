package sftp

import (
	"path/filepath"

	"github.com/pkg/sftp"
	
	"github.com/wanglu119/ng-ssh/lib/common"
	"github.com/wanglu119/ng-ssh/lib/ssh"
)

const maxPacket = 1 << 15

func NewSftpClient(mc *common.Machine) (*sftp.Client, error) {
	conn, err := ssh.NewSshClient(mc)
	if err != nil {
		return nil, err
	}
	return sftp.NewClient(conn, sftp.MaxPacket(maxPacket))
}
func toUnixPath(path string) string {
	return filepath.Clean(path)
}
