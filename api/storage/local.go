package storage

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type LocalStorage struct{}

func NewLocalStorage() Storage {
	return &LocalStorage{}
}

func (s LocalStorage) Store(filename string, contents []byte) (string, error) {
	fileRelativePath := fmt.Sprintf("./public/%s", filename)
	err := ioutil.WriteFile(fileRelativePath, contents, 0666)
	if err != nil {
		return "", err
	}

	return filepath.Abs(fileRelativePath)
}

func (s LocalStorage) Get(filename string) ([]byte, error) {
	filepath := fmt.Sprintf("./public/%s", filename)
	return ioutil.ReadFile(filepath)
}

func (s LocalStorage) Delete(filename string) error {
	return nil
}
