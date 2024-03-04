package engine

import (
	"bytes"
	"encoding/gob"
	"errors"
)

// HashMap is the in-memory key-offset pair.
// It includes a in memory map and persist the hashmap content to file.

type Offset uint32

type HashMap struct {
	mem      map[string]Offset
	filename string
}

func NewHashMap(filename string) (*HashMap, error) {
	if CheckExist(filename) {
		return nil, errors.New("hashmap file already exist")
	}

	err := CreateFile(filename)
	if err != nil {
		return nil, err
	}

	return &HashMap{
		mem:      make(map[string]Offset),
		filename: filename,
	}, nil
}

func LoadHashMap(filename string) (*HashMap, error) {
	if !CheckExist(filename) {
		return nil, errors.New("hashmap file does not exist")
	}
	hashmap := &HashMap{
		mem:      make(map[string]Offset),
		filename: filename,
	}

	hashmap.Load()
	return hashmap, nil
}

// Put and Get, normal hashmap behavior
func (h *HashMap) Put(key string, value Offset) {
	h.mem[key] = value
}

func (h *HashMap) Get(key string) (Offset, error) {
	if v, ok := h.mem[key]; !ok {
		return 0, errors.New("key doesn't exist")
	} else {
		return v, nil
	}
}

// Persist in-memory map to file.
func (h *HashMap) Persist() error {
	encoded, err := encode(h.mem)
	if err != nil {
		return err
	}
	return WriteFile(h.filename, encoded)
}

// Load file content back to in-memory map.
func (h *HashMap) Load() error {
	content, err := ReadFile(h.filename)
	if err != nil {
		return err
	}

	decoded, err := decode(content)
	if err != nil {
		return err
	}
	h.mem = decoded
	return nil
}

// Encode in-memory map to byte array
func encode(hashmap map[string]Offset) ([]byte, error) {
	var snapshot bytes.Buffer
	encoder := gob.NewEncoder(&snapshot)
	if err := encoder.Encode(hashmap); err != nil {
		return nil, err
	}
	return snapshot.Bytes(), nil
}

// Decode byte array to in-memory map
func decode(snapshot []byte) (map[string]Offset, error) {
	reader := bytes.NewReader(snapshot)
	decoder := gob.NewDecoder(reader)
	h := make(map[string]Offset)
	if err := decoder.Decode(&h); err != nil {
		return nil, err
	}
	return h, nil
}
