package fs

import "os"

type GenDirCreator struct {
	create func(path string, perm os.FileMode) error
}

func (creator GenDirCreator) CreateDirectory(path string, perm os.FileMode) error {
	return creator.create(path, perm)
}

func GetDefaultDirCreator() DirCreator {
	creator := GenDirCreator{
		create: os.MkdirAll,
	}

	return creator
}

type GenFileWriter struct {
	write func(path string, data []byte, perm os.FileMode) error
}

func (writer GenFileWriter) WriteFile(path string, data []byte, perm os.FileMode) error {
	return writer.write(path, data, perm)
}

func GetDefaultFileWriter() FileWriter {
	writer := GenFileWriter{
		write: os.WriteFile,
	}
	return writer
}

type GenRemover struct {
	remove func(path string) error
}

func (remover GenRemover) Remove(path string) error {
	return remover.remove(path)
}

func GetDefaultDirRemover() Remover {
	remover := GenRemover{
		remove: os.RemoveAll,
	}
	return remover
}
