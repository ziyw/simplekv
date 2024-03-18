package engine

import (
	"errors"
)

const (
	SEG_FILE_PREFIX  = "SEG_FILE_"
	HASH_FILE_PREFIX = "MAP_FILE_"
)

type Segment struct {
	id           string
	segFileName  string
	hashFileName string
	hashmap      *HashMap
	segFile      *SegFile
}

func NewSegment(id string) (*Segment, error) {
	segFileName := SEG_FILE_PREFIX + id
	hashFileName := HASH_FILE_PREFIX + id

	if CheckExist(segFileName) || CheckExist(hashFileName) {
		return nil, errors.New("already exist error")
	}

	sf, err := NewSegFile(segFileName)
	if err != nil {
		return nil, err
	}

	hf, err := NewHashMap(hashFileName)
	if err != nil {
		return nil, err
	}

	return &Segment{
		id:           id,
		segFileName:  segFileName,
		hashFileName: hashFileName,
		hashmap:      hf,
		segFile:      sf,
	}, nil
}

// Load from existing segment
// func Load(id string) *Segment {
// 	segFileName := SEG_FILE_PREFIX + id
// 	hashFileName := HASH_FILE_PREFIX + id
// sf, err := LoadSegFile(segFileName)
// }

// func (s Segment) PersistHashmap() {}
// func (s Segment) LoadHashmap()    {}

func (s Segment) Delete() {
	s.segFile.Delete()
	s.hashmap.Delete()
}

// Put key-value pair to both segment file and hashmap
func (s Segment) Put(key, value string) error {
	offset, err := s.segFile.Append(key, value)
	if err != nil {
		return err
	}
	s.hashmap.mem[key] = offset
	return nil
}

// Get value using key from hashmap directly
func (s Segment) Get(key string) (string, error) {
	offset, ok := s.hashmap.mem[key]
	if !ok {
		return "", errors.New("key not exist")
	}
	_, value, err := s.segFile.Fetch(offset)
	return value, err
}

func (s Segment) GetAll() ([]string, []string, error) {
	var keys, values []string
	for k, offset := range s.hashmap.mem {
		_, value, err := s.segFile.Fetch(offset)
		if err != nil {
			return nil, nil, err
		}
		keys = append(keys, k)
		values = append(values, value)
	}
	return keys, values, nil
}
