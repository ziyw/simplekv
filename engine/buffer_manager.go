package engine

import (
	"errors"
	"fmt"
)

type BufferManager struct {
}

// Put: should create new segment if not existing segment
// if segment already exist, put to existing active Segment
// if exist active segment is full, create a new segment
func (b *BufferManager) Put(key, value string) error {

	sortedFileNames := ListSegmentFileNames(".", true)
	fmt.Printf("existing seg files: %v\n", sortedFileNames)

	if len(sortedFileNames) == 0 {
		activeSegment, err := NewSegment(1)
		if err != nil {
			return err
		}
		err = activeSegment.Put(key, value)
		if err != nil {
			return err
		}
		return activeSegment.hashmap.Persist()
	}

	id := FromSegFilenName(sortedFileNames[0])
	activeSegment, err := Load(id)
	if err != nil {
		return err
	}

	err = activeSegment.Put(key, value)
	if err == nil {
		return activeSegment.hashmap.Persist()
	}

	if err.Error() == "exceed max segment items" {
		fmt.Println("Excceeed max segment items, create a new one")
		newSegment, _ := NewSegment(id + 1)
		e := newSegment.Put(key, value)
		if e == nil {
			return activeSegment.hashmap.Persist()
		}
		return e
	}

	return err
}

// Get will need to serach all segment until find the value for the key
// this is done by traverse all keys in hashmaps until find one
func (b *BufferManager) Get(key string) (string, error) {
	sortedFileNames := ListSegmentFileNames(".", true)
	for i, n := range sortedFileNames {
		fmt.Printf("Try %d: name %s\n", i, n)
		id := FromSegFilenName(n)
		seg, err := Load(id)
		if err != nil {
			return "", err
		}

		if _, ok := seg.hashmap.mem[key]; ok {
			return seg.Get(key)
		}
	}
	fmt.Printf("No such key")
	return "", errors.New("no such key")
}
