package engine

import (
	"errors"
	"log"
	"os"
)

type SegFile struct {
	name string
}

// Functions need to cover:
// get(offset), get key-value pair from SegmentFile, offset is error cases
// Append: append new key-value to SegmentFile, return current SegmentFile offset
// NewSegmentFile: create a new empty file
// GetAll(): return all Key-value pair as array, because there will be duplicate Key-value
// Compress and Merge
//

func NewSegFile(filename string) (*SegFile, error) {
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
