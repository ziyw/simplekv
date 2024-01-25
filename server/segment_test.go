package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestEncodeAndDecodeHashmap(t *testing.T) {
	want := map[string]uint32{"Hello": 12, "World": 45}

	input := map[string]uint32{"Hello": 12, "World": 45}
	got := DecodeHashMapSnapshot(EncodeHashMapToSnapshot(input))
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("test expected: %v, got: %v\n", want, got)
	}
}

func TestSegmentAppend(t *testing.T) {
	s := Segment{
		id:      "1",
		hashmap: make(map[string]uint32),
	}
	s.Append("Whatever", "54321")
	fmt.Print(s.hashmap)
}

func TestSegmentGetValue(t *testing.T) {
	s := Segment{
		id:      "2",
		hashmap: make(map[string]uint32),
	}

	s.Append("Hello", "12345")
	s.Append("Whatever", "54321")

	val, _ := s.GetValue("Hello")
	fmt.Println(val)
}

func TestSegmentCompress(t *testing.T) {
	s := Segment{
		id:      "3",
		hashmap: make(map[string]uint32),
	}

	s.Append("Hello", "12345")
	s.Append("Whatever", "54321")
	s.Append("Hello", "ThisIsNew")
	s.Append("Whatever", "This is newer")

	newSegment := s.Compress()
	fmt.Println("SegmentCompress result")
	fmt.Print(newSegment.hashmap)
	for key := range newSegment.hashmap {
		v, _ := newSegment.GetValue(key)
		fmt.Printf("(%s: %s)", key, v)
	}

	DeleteFile(newSegment.id)

}
