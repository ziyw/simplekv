package engine

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	ENTRY_SEPARATOR = ","
	ENTRY_END       = "\n"
)

// Create a new string entry from key-value.
func NewEntry(key, value string) string {
	if strings.Contains(key, ENTRY_SEPARATOR) || strings.Contains(key, ENTRY_END) {
		return ""
	}
	if strings.Contains(value, ENTRY_SEPARATOR) || strings.Contains(value, ENTRY_END) {
		return ""
	}

	return fmt.Sprintf("%s%s%s%s", key, ENTRY_SEPARATOR, value, ENTRY_END)
}

// Parse a string entry back to key-value pair.
func ParseEntry(entry string) (string, string, error) {
	if len(entry) == 0 {
		return "", "", errors.New("illegal entry: can't be empty")
	}
	if !strings.Contains(entry, ENTRY_SEPARATOR) && !strings.Contains(entry, ENTRY_END) {
		return "", "", errors.New("illegal entry, should contain separators")
	}

	keyEnd := strings.Index(entry, ENTRY_SEPARATOR)
	valueEnd := strings.Index(entry, ENTRY_END)

	b := []byte(entry)
	return string(b[0:keyEnd]), string(b[keyEnd+1 : valueEnd]), nil
}

type SegFile struct {
	name string
}

func NewSegFile(filename string) (*SegFile, error) {
	// TODO: maybe when file already exist, use the existing file, no error
	// file already exist, return nil
	_, err := os.Stat(filename)
	if err == nil {
		return nil, errors.New("file already exist")
	}

	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return &SegFile{
		name: filename,
	}, nil
}

func (s SegFile) Delete() {
	if err := os.Remove(s.name); err != nil {
		log.Fatal(err)
	}
}

func (s SegFile) Append(entry string) (Offset, error) {
	// segment file should exist before append to the file
	var fileSize Offset = 0
	info, err := os.Stat(s.name)
	if err != nil {
		return 0, err
	}
	fileSize = Offset(info.Size())

	f, err := os.OpenFile(s.name, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fileSize, err
	}
	defer f.Close()

	if _, err := f.WriteString(entry); err != nil {
		return fileSize, err
	}
	return fileSize, nil
}
