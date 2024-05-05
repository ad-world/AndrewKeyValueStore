package fs_ops

import (
	"io/fs"
)

type FsOpsInterface interface {
	ReadKey(name string) ([]byte, error)
	WriteKey(name string, data []byte, perm fs.FileMode) error
	DeleteKey(name string) error
	fileName(name string) string
}

