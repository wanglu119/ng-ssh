package sftp

import (
	wlog "github.com/wanglu119/me-deps/log"
)

var log *wlog.Logger

func init() {
	log = wlog.GetLogger()
}
