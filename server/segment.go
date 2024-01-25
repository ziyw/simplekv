package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

const SEGMENT_NAME_FMT = "segment_%d"
const SNAPSHOT_NAME_FMT = "snapshot_%d"

type SegList struct {
	active []Segment
	segMap map[int]int
}

func NewSegList() SegList {
	return SegList{
		active: make([]Segment, 0),
		segMap: make(map[int]int),
	}
}

func (sl *SegList) AddSegment(s Segment) {
	sl.segMap[s.id] = len(sl.active)
	sl.active = append(sl.active, s)
}

func (sl *SegList) RemoveSegment(id int) {
	rm := sl.segMap[id]
	sl.active = append(sl.active[:rm], sl.active[rm+1:]...)
	delete(sl.segMap, id)
}

func (sl SegList) GetNextId() int {
	if len(sl.active) == 0 {
		return 1
	}

	last := sl.active[len(sl.active)-1]
	return last.id + 1
}

func (sl SegList) GetCurrentSegment() *Segment {
	if len(sl.active) == 0 {
		s := NewSegment(sl.GetNextId())
		sl.AddSegment(*s)
	}
	return &sl.active[len(sl.active)-1]
}

// Segment:
// Each segment is a append only log file and has a in-memory hashmap for fetching items
type Segment struct {
	id            int
	segment_name  string
	snapshot_name string
	hashmap       map[string]uint32 // in-memory map (key, offset)
}

func NewSegment(nxtId int) *Segment {
	return &Segment{
		id:            nxtId,
		segment_name:  fmt.Sprintf(SEGMENT_NAME_FMT, nxtId),
		snapshot_name: fmt.Sprintf(SNAPSHOT_NAME_FMT, nxtId),
		hashmap:       make(map[string]uint32),
	}
}

// Remove duplicate records in segment file
// read from back of the segment file to the front, get all keys
// then write the new key-value to new file, update hashmap
func (s Segment) Compress(nxtId int) *Segment {
	// nxtId := getNextId()
	// newSegment := Segment{
	// 	id:            nxtId,
	// 	segment_name:  fmt.Sprintf(SEGMENT_NAME_FMT, nxtId),
	// 	snapshot_name: fmt.Sprintf(SNAPSHOT_NAME_FMT, nxtId),
	// 	hashmap:       make(map[string]uint32),
	// }

	newSegment := NewSegment(nxtId)
	for k := range s.hashmap {
		v, _ := s.GetValue(k)
		newSegment.Append(k, v)
	}
	return newSegment
}

// Merge two segment file into a third Segment
func Merge(s1, s2 Segment, nxtId int) *Segment {
	// nxtId := getNextId(nxtId)
	// newSegment := Segment{
	// 	id:            nxtId,
	// 	segment_name:  fmt.Sprintf(SEGMENT_NAME_FMT, nxtId),
	// 	snapshot_name: fmt.Sprintf(SNAPSHOT_NAME_FMT, nxtId),
	// 	hashmap:       make(map[string]uint32),
	// }

	newSegment := NewSegment(nxtId)

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

// Append will append the key value byte to segment file
// hashmap will also be updated
func (s *Segment) Append(key, value string) {
	// Get the key-offset pair for hashmap
	var offset uint32
	info, _ := os.Stat(s.segment_name)
	if info != nil {
		offset = uint32(info.Size())
	} else {
		offset = 0
	}
	s.hashmap[key] = offset

	// append new key-value to the end of the segment file
	f, err := os.OpenFile(s.segment_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	newRecord := []byte(key + "," + value + "\n")
	if _, err := f.Write(newRecord); err != nil {
		f.Close()
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

// Check if value exist in current segment
func (s *Segment) GetValue(key string) (value string, ok bool) {
	offset, ok := s.hashmap[key]
	if !ok {
		return "", false // hash key doesn't exist in current map
	}

	b, err := os.ReadFile(s.segment_name)
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
	f, _ := os.Stat(s.snapshot_name)
	if f != nil {
		if e := os.Remove(s.snapshot_name); e != nil {
			log.Fatalf("error deleting previous snapshot")
		}
	}
	// encode hashmap to byte array, then write to snapshot file
	snapshot := EncodeHashMapToSnapshot(s.hashmap)
	err := os.WriteFile(s.snapshot_name, snapshot, 0644)
	if err != nil {
		log.Fatalf("error creating hashmap snapshot")
	}
}

// Load a snapshot file back to hashmap
func (s *Segment) LoadSnapshot() {
	_, err := os.Stat(s.snapshot_name)
	if err != nil {
		log.Fatalf("load snapshot failed due to snapshot file doesn't exist: %v", err)
	}

	snapshot, err := os.ReadFile(s.snapshot_name)
	if err != nil {
		log.Fatalf("load snapshot failed: %v", err)
	}

	s.hashmap = DecodeHashMapSnapshot(snapshot)
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
