package main

import (
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

// Create Snapshot for an existing hashmap
func (s Segment) CreateSnapshot(hashmap map[string]string) []byte {
	return nil
}

// Load a snapshot file back to hashmap
func LoadSnapshot(input []byte) map[string]string {
	return nil
}

// Append key-value to a file, return offset
func Append(item []byte, file string) uint32 {
	return 0
}

// Compress segment file
func Compress(input []byte) []byte {
	return nil
}

// Merge two segment file into a third Segment
func Merge(seg1, seg2 string) string {
	return "whatever"
}
