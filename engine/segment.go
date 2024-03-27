package engine

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

const (
	MAX_ITEM_NUM     = 20
	SEG_FILE_PREFIX  = "SEG_FILE_"
	HASH_FILE_PREFIX = "MAP_FILE_"
)

type Segment struct {
	id           int
	segFileName  string
	hashFileName string
	hashmap      *HashMap
	segFile      *SegFile
}

func NewSegment(id int) (*Segment, error) {

	segFileName := ToSegFileName(id)
	hashFileName := ToHashFileName(id)

	if CheckExist(segFileName) || CheckExist(hashFileName) {
		return nil, errors.New("file already exist")
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
func Load(id int) (*Segment, error) {
	segFileName := ToSegFileName(id)
	hashFileName := ToHashFileName(id)

	if !CheckExist(segFileName) {
		return nil, errors.New("segment file does not exist")
	}

	if !CheckExist(hashFileName) {
		return nil, errors.New("hashmap file does not exist")
	}

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
	if len(s.hashmap.mem) >= MAX_ITEM_NUM {
		return errors.New("exceed max segment items")
	}

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

func (s Segment) Compress(nxtId int) (*Segment, error) {
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

// second is the later than first, so second will always override first content
func Merge(first, second *Segment, nxtId int) (*Segment, error) {
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

func ListSegmentFileNames(path string, sorted bool) []string {
	files, err := os.ReadDir(path)
	if err != nil {
		return []string{}
	}

	seg_files := []string{}
	for _, f := range files {
		if strings.Contains(f.Name(), SEG_FILE_PREFIX) {
			seg_files = append(seg_files, f.Name())
		}
	}
	if sorted {
		sort.Strings(seg_files)
		slices.Reverse(seg_files)
		return seg_files
	} else {
		return seg_files
	}
}

// construct segment file name from id or parse id from seg file name.
func ToSegFileName(id int) string {
	return fmt.Sprintf("%s%04d", SEG_FILE_PREFIX, id)
}

func ToHashFileName(id int) string {
	return fmt.Sprintf("%s%04d", HASH_FILE_PREFIX, id)
}

func FromSegFilenName(segFileName string) int {
	parts := strings.Split(segFileName, SEG_FILE_PREFIX)
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return -1
	}
	return id
}
