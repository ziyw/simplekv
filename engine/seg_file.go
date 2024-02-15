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

type Pair struct {
	key   string
	value string
}

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

func (s SegFile) GetAll() ([]Pair, error) {
	// read file from offset to ENTRY_SEPARATOR, from ENTRY_SEPARATOR to ENTRY_END
	info, err := os.Stat(s.name)
	if err != nil {
		return nil, err
	}
	fileSize := Offset(info.Size())

	// this actually just load the file to memeory, think about how to better manager it
	buf, err := os.ReadFile(s.name)
	if err != nil {
		log.Fatal(err)
	}

	pairs := []Pair{}

	var start Offset = 0
	for start < fileSize {
		keyEnd := start + 1
		for buf[keyEnd] != []byte(ENTRY_SEPARATOR)[0] {
			keyEnd++
		}
		key := string(buf[start:keyEnd])

		valueEnd := keyEnd + 1
		for buf[valueEnd] != []byte(ENTRY_END)[0] {
			valueEnd++
		}
		value := string(buf[keyEnd+1 : valueEnd])
		pairs = append(pairs, Pair{
			key:   key,
			value: value,
		})
		start = valueEnd + 1
	}
	return pairs, nil
}

func (s SegFile) Get(offset Offset) (Pair, error) {
	// read file from offset to ENTRY_SEPARATOR, from ENTRY_SEPARATOR to ENTRY_END
	_, err := os.Stat(s.name)
	if err != nil {
		return Pair{}, err
	}

	// this actually just load the file to memeory, think about how to better manager it
	buf, err := os.ReadFile(s.name)
	if err != nil {
		log.Fatal(err)
	}

	keyEnd := offset + 1
	for buf[keyEnd] != []byte(ENTRY_SEPARATOR)[0] {
		keyEnd++
	}
	key := string(buf[offset:keyEnd])

	valueEnd := keyEnd + 1
	for buf[valueEnd] != []byte(ENTRY_END)[0] {
		valueEnd++
	}
	value := string(buf[keyEnd+1 : valueEnd])
	return Pair{
		key:   key,
		value: value,
	}, nil
}
