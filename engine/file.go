package engine

import (
	"errors"
	"os"
)

// File I/O functions

func CheckExist(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return info.Size() >= 0
}

func CreateFile(filename string) error {
	if CheckExist(filename) {
		return errors.New("file already exist")
	}

	// call os.Create on existing file will override the existing file
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func DeleteFile(filename string) error {
	return os.Remove(filename)
}

func OpenFile(filename string) (*os.File, error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func WriteFile(filename string, content []byte) error {
	return os.WriteFile(filename, content, 0666)
}

func ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}
