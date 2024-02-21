package engine

import (
	"errors"
	"fmt"
	"os"
)

// File I/O functions

func CheckExist(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		fmt.Printf("error opening file: %s", filename)
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
