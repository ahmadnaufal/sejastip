package storage

import (
	"fmt"
	"io/ioutil"
)

type LocalStorage struct{}

func NewLocalStorage() Storage {
	return &LocalStorage{}
}

func (s LocalStorage) Store(filename string, contents []byte) error {
	filepath := fmt.Sprintf("./public/%s", filename)
	err := ioutil.WriteFile(filepath, contents, 0666)
	return err
}

func (s LocalStorage) Get(filename string) ([]byte, error) {
	filepath := fmt.Sprintf("./public/%s", filename)
	return ioutil.ReadFile(filepath)
}

func (s LocalStorage) Delete(filename string) error {
	return nil
}
