package engine

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
)

type hashmap interface {
	PersistSnapshot(io.write) error
	LoadSnapshot() error
}

type offset uint32

// hashmap is key-offset pair. key is a string for the saved key,
// offset is the offset of the value in segment file.
type Hashmap map[string]offset

func (h Hashmap) Put(key string, value offset) {
	h[key] = value
}

func (h Hashmap) Get(key string) (offset, error) {
	v, ok := h[key]
	if !ok {
		return 0, errors.New("key does not exist")
	}
	return v, nil
}

// NewSnapshot will encode hashmap to byte array.
func NewSnapshot(h Hashmap) []byte {
	var snapshot bytes.Buffer
	encode := gob.NewEncoder(&snapshot)
	if err := encode.Encode(h); err != nil {
		log.Fatal("error encoding hashmap: %v", err)
	}
	return snapshot.Bytes()
}

// LoadSnapshot
func LoadSnapshot(snapshot []byte) Hashmap {
	reader := bytes.NewReader(snapshot)
	decoder := gob.NewDecoder(reader)
	h := make(Hashmap)
	if err := decoder.Decode(&h); err != nil {
		log.Fatal("error decoding hashmap snapshot: %v", err)
	}
	return h
}
