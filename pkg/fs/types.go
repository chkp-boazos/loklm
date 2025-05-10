package fs

import "os"

type DirCreator interface {
	CreateDirectory(path string, perm os.FileMode) error
}

type FileWriter interface {
	WriteFile(path string, data []byte, perm os.FileMode) error
}

type Remover interface {
	Remove(path string) error
}
