package utils

import (
	"os"
	"path/filepath"
)

func CreateFile(filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

/*
文件是否存在
*/
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func WiteFile(path, data string) error {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}
	fileHandler, err := os.Create(path)
	if nil != err {
		return err
	}
	defer fileHandler.Close()

	_, err = fileHandler.Write([]byte(data))
	if err != nil {
		return err
	}
	return nil
}
