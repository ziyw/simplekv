package engine

import (
	"bytes"
	"encoding/gob"
)

// hashmap is key-Offset pair. key is a string for the saved key,
// Offset is the Offset of the value in segment file.
type Offset uint32
type HashMap map[string]Offset

func NewHashMap() HashMap {
	return make(map[string]Offset)
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
