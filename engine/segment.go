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
// TODO: cover hashmap file not exist cases
func Load(id string) (*Segment, error) {
	segFileName := SEG_FILE_PREFIX + id
	hashFileName := HASH_FILE_PREFIX + id

	sf, err := LoadSegFile(segFileName)
	if err != nil {
		return nil, err
	}

	hf, err := LoadHashMap(hashFileName)
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

// TODO: Merge(filename1, file2) -> create a new segment from two previous segment files
// TODO: Persist hashmap -> periodic persist hashmap

func (s Segment) Compress(nxtId string) (*Segment, error) {
	compressed, err := NewSegment(nxtId)
	if err != nil {
		return nil, err
	}

	for k := range s.hashmap.mem {
		v, err := s.Get(k)
		if err != nil {
			return nil, err
		}
		compressed.Put(k, v)
	}

	compressed.hashmap.Persist()
	return compressed, nil
}

func Merge(first, second *Segment, nxtId string) (*Segment, error) {
	nxt, err := NewSegment(nxtId)
	if err != nil {
		return nil, err
	}

	// combine hashmap to one
	combined := make(map[string]Offset)
	for k, offset := range first.hashmap.mem {
		combined[k] = offset
	}
	for k, offset := range second.hashmap.mem {
		combined[k] = offset
	}

	for k := range combined {
		if _, ok := second.hashmap.mem[k]; ok {
			v, err := second.Get(k)
			if err != nil {
				return nil, err
			}
			nxt.Put(k, v)
		} else {
			v, err := first.Get(k)
			if err != nil {
				return nil, err
			}
			nxt.Put(k, v)
		}
	}
	return nxt, nil
}
