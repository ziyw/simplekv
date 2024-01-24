package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"
)

// Segment:
// Each segment is a append only log file and has a in-memory hashmap for fetching items

type Segment struct {
	id      string            // segment file name
	hashmap map[string]string // in-memory map (key, offset)
}

func (s Segment) CreateSegmentFile() {
	// Check if file exist, if yes, panic
	info, _ := os.Stat(s.id)
	if info != nil {
		log.Fatal("segment file already exist")
	}

	_, err := os.Create(s.id)
	if err != nil {
		log.Fatalf("error creating segment file: %v", err)
	}
}

// Compress segment file
func Compress(input []byte) []byte {
	return nil
}

// Merge two segment file into a third Segment
func Merge(seg1, seg2 string) string {
	return "whatever"
}

/**
	Segment file related functions
**/

// Append will append the key value byte to segment file
// hashmap will also be updated
func (s *Segment) Append(key, value string) {
	filename := "segment_" + s.id

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	newRecord := []byte(key + value)
	if _, err := f.Write(newRecord); err != nil {
		f.Close() // ignore error; Write error takes precedence
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	s.hashmap[key] = value
}

/**
	HashMap Related Functions
**/

// Encode Hashmap and persist it to a file
func (s Segment) CreateSnapshot() string {
	filename := "snapshot_" + s.id

	// check if snapshot file already exist, delete the previous one
	f, _ := os.Stat(filename)
	if f != nil {
		if e := os.Remove(filename); e != nil {
			log.Fatalf("error deleting previous snapshot")
		}
	}
	// encode hashmap to byte array, then write to snapshot file
	snapshot := EncodeHashMapToSnapshot(s.hashmap)
	err := os.WriteFile(filename, snapshot, 0644)
	if err != nil {
		log.Fatalf("error creating hashmap snapshot")
	}
	return filename
}

// Load a snapshot file back to hashmap
func (s *Segment) LoadSnapshot() {
	filename := "snapshot_" + s.id

	_, err := os.Stat(filename)
	if err != nil {
		log.Fatalf("load snapshot failed due to snapshot file doesn't exist: %v", err)
	}

	snapshot, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("load snapshot failed: %v", err)
	}

	s.hashmap = DecodeHashMapSnapshot(snapshot)
}

func EncodeHashMapToSnapshot(hashmap map[string]string) []byte {
	var snapshot bytes.Buffer
	encoder := gob.NewEncoder(&snapshot)
	if err := encoder.Encode(hashmap); err != nil {
		log.Fatalf("error encoding hashmap: %v", err)
	}
	return snapshot.Bytes()
}

func DecodeHashMapSnapshot(snapshot []byte) map[string]string {
	reader := bytes.NewReader(snapshot)
	decoder := gob.NewDecoder(reader)

	hashmap := make(map[string]string)
	if err := decoder.Decode(&hashmap); err != nil {
		log.Fatalf("error decoding hashmap snapshot: %v", err)
	}
	return hashmap
}
