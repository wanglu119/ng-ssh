package scp

import (
	"path/filepath"
)

func toUnixPath(path string) string {
	return filepath.Clean(path)
}


