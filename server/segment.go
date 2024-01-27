package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const SEGMENT_NAME_FMT = "segment_%d"
const SNAPSHOT_NAME_FMT = "snapshot_%d"

// Segment is how key-value are persisted on disk.
// content of the key-values are on segmentFile
// Snapshot is the in memory hash map of (key, offset), where key-value saved on disk starts from offset.
type Segment struct {
	id          int
	segmentFile string

	hashmap      map[string]uint32 // in-memory map (key, offset)
	snapshotFile string            // persisted version of the hashmap
}

func NewSegment(nxtId int) *Segment {
	return &Segment{
		id:           nxtId,
		segmentFile:  fmt.Sprintf(SEGMENT_NAME_FMT, nxtId),
		hashmap:      make(map[string]uint32),
		snapshotFile: fmt.Sprintf(SNAPSHOT_NAME_FMT, nxtId),
	}
}

// Remove duplicate records in segment file
// read from back of the segment file to the front, get all keys
// then write the new key-value to new file, update hashmap
func (s Segment) Compress(nxtId int) *Segment {
	newSegment := NewSegment(nxtId)
	for k := range s.hashmap {
		v, _ := s.GetValue(k)
		newSegment.Append(k, v)
	}
	return newSegment
}

// Merge two segment file into a third Segment
func Merge(s1, s2 Segment, nxtId int) *Segment {
	newSegment := NewSegment(nxtId)

	// latest segment is fst, older is snd
	fst, snd := s1, s2
	if s1.id < s2.id {
		fst, snd = s2, s1
	}

	for k := range fst.hashmap {
		v, _ := fst.GetValue(k)
		newSegment.Append(k, v)
	}

	for k := range snd.hashmap {
		// skip values that are already in first segment
		_, ok := fst.GetValue(k)
		if ok {
			continue
		}
		v, _ := snd.GetValue(k)
		newSegment.Append(k, v)
	}

	return newSegment
}

// Append will append the key value byte to segment file and update hashmap
func (s *Segment) Append(key, value string) {
	// TODO: improve by keep filesize in memory, avoid fetching stats every time
	// Get the key-offset pair for hashmap
	var offset uint32
	info, _ := os.Stat(s.segmentFile)
	if info != nil {
		offset = uint32(info.Size())
	} else {
		offset = 0
	}

	// append new key-value to the end of the segment file
	f, err := os.OpenFile(s.segmentFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	newRecord := []byte(key + "," + value + "\n")
	if _, err := f.Write(newRecord); err != nil {
		log.Fatal(err)
	}

	// add key-offset to hashmap after writing to disk
	s.hashmap[key] = offset
}

// Check if value exist in current segment
func (s *Segment) GetValue(key string) (value string, ok bool) {
	offset, ok := s.hashmap[key]
	if !ok {
		return "", false // hash key doesn't exist in current map
	}

	b, err := os.ReadFile(s.segmentFile)
	if err != nil {
		log.Fatal(err)
	}

	// from offset to ,
	keyRange := offset
	for b[keyRange] != ',' {
		keyRange++
	}

	// from , to \n
	valueOffset, valueRange := keyRange+1, keyRange+1
	for b[valueRange] != '\n' {
		valueRange++
	}

	return string(b[valueOffset:valueRange]), true
}

// Encode Hashmap and persist it to a file
func (s Segment) CreateSnapshot() {
	// check if snapshot file already exist, delete the previous one
	f, _ := os.Stat(s.snapshotFile)
	if f != nil {
		if e := os.Remove(s.snapshotFile); e != nil {
			log.Fatalf("error deleting previous snapshot")
		}
	}
	// encode hashmap to byte array, then write to snapshot file
	snapshot := EncodeHashMapToSnapshot(s.hashmap)
	err := os.WriteFile(s.snapshotFile, snapshot, 0644)
	if err != nil {
		log.Fatalf("error creating hashmap snapshot")
	}
}

// Load a snapshot file back to hashmap
func (s *Segment) LoadSnapshot() {
	_, err := os.Stat(s.snapshotFile)
	if err != nil {
		log.Fatalf("load snapshot failed due to snapshot file doesn't exist: %v", err)
	}

	snapshot, err := os.ReadFile(s.snapshotFile)
	if err != nil {
		log.Fatalf("load snapshot failed: %v", err)
	}

	s.hashmap = DecodeHashMapSnapshot(snapshot)
}

func NewSegmentFromSnapshot(file string) *Segment {
	idStr := strings.Split(file, "snapshot_")[1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatalf("error pasring snapshot id :%s", file)
	}

	s := NewSegment(id)
	s.LoadSnapshot()
	return s
}

func EncodeHashMapToSnapshot(hashmap map[string]uint32) []byte {
	var snapshot bytes.Buffer
	encoder := gob.NewEncoder(&snapshot)
	if err := encoder.Encode(hashmap); err != nil {
		log.Fatalf("error encoding hashmap: %v", err)
	}
	return snapshot.Bytes()
}

func DecodeHashMapSnapshot(snapshot []byte) map[string]uint32 {
	reader := bytes.NewReader(snapshot)
	decoder := gob.NewDecoder(reader)

	hashmap := make(map[string]uint32)
	if err := decoder.Decode(&hashmap); err != nil {
		log.Fatalf("error decoding hashmap snapshot: %v", err)
	}
	return hashmap
}

func CreateFile(filename string) {
	// Check if file exist, if yes, panic
	info, _ := os.Stat(filename)
	if info != nil {
		log.Fatal("segment file already exist")
	}

	_, err := os.Create(filename)
	if err != nil {
		log.Fatalf("error creating segment file: %v", err)
	}
}

func DeleteFile(filename string) {
	err := os.Remove(filename)
	if err != nil {
		log.Fatalf("error deleting file: %v", err)
	}
}
