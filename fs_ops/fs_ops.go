package fs_ops

import (
	"errors"
	"io/fs"
	"os"
)

type FsOps struct {}

func (f *FsOps) ReadKey(name string) ([]byte, error) {
	// TODO: Add locking mechanism to prevent concurrent writes

	file, err := os.Open(f.fileName(name))
	if err != nil {
		return nil, err
	}
	
	defer file.Close()
	
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := stat.Size()

	data := make([]byte, fileSize)
	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (f *FsOps) WriteKey(name string, data []byte, perm fs.FileMode) error {
	_, err := os.Stat(f.fileName(name))

	// Create file if it doesn't exist
	if err != nil {
		if errors.Is(err, os.ErrNotExist){
			_, err := os.Create(f.fileName(name))

			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// Open file for writing
	file, err := os.OpenFile(f.fileName(name), os.O_RDWR, perm)
	if err != nil {
		return err
	}

	defer file.Close()


	// Write data to file
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (f *FsOps) DeleteKey(name string) error {
	err := os.Remove(f.fileName(name))
	if err != nil {
		return err
	}
	return nil
}

func CreateFsOps() *FsOps {
	return &FsOps{}
}

func (f *FsOps) fileName(name string) string {
	return "store/" + name
}