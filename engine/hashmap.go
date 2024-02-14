package engine

import (
	"bytes"
	"encoding/gob"
)

// hashmap is key-offset pair. key is a string for the saved key,
// offset is the offset of the value in segment file.
type offset uint32
type HashMap map[string]offset

func NewHashMap() HashMap {
	return make(map[string]offset)
}

func Encode(h HashMap) ([]byte, error) {
	var snapshot bytes.Buffer
	encoder := gob.NewEncoder(&snapshot)
	if err := encoder.Encode(h); err != nil {
		return nil, err
	}
	return snapshot.Bytes(), nil
}

func Decode(snapshot []byte) (HashMap, error) {
	reader := bytes.NewReader(snapshot)
	decoder := gob.NewDecoder(reader)
	h := make(HashMap)
	if err := decoder.Decode(&h); err != nil {
		return nil, err
	}
	return h, nil
}
