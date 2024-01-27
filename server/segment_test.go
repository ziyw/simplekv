package main

import (
	"testing"
)

func TestNewSegment(t *testing.T) {

}

// func TestNewSegment(t *testing.T) {
// // 	tests :=

// // 	tests := []struct {
// // 		input int
// // 		want  int
// // 	}{
// // 		{input: 12, want: 12},
// // 		{input: 0, want: 0},
// // 		{input: 255, want: 255},
// // 	}

// // 	s1 := NewSegment()

// // 	addToSegList(*s1)
// // 	s2 := NewSegment()

// // 	s3 := NewSegment()
// // 	fmt.Println("s3 segment name" + s3.segment_name)
// // }

// func TestSegmentListFunctions(t *testing.T) {
// 	segList := SegList{
// 		active: make([]Segment, 0),
// 		segMap: make(map[int]int),
// 	}

// 	s1 := NewSegment(segList.GetNextId())
// 	if s1.id != 1 {
// 		t.Fatalf("test failed: first segment id is not 1, got: %v", s1.id)
// 	}
// 	segList.AddSegment(*s1)

// 	s2 := NewSegment(segList.GetNextId())
// 	if s2.id != 2 {
// 		t.Fatalf("test failed: first segment id is not 2, got: %v", s2.id)
// 	}
// 	segList.AddSegment(*s2)

// 	segList.RemoveSegment(s1.id)
// 	s3 := NewSegment(segList.GetNextId())
// 	if s3.id != 3 {
// 		t.Fatalf("test failed: first segment id is not 3, got: %v", s3.id)
// 	}

// }

// func TestEncodeAndDecodeHashmap(t *testing.T) {
// 	want := map[string]uint32{"Hello": 12, "World": 45}

// 	input := map[string]uint32{"Hello": 12, "World": 45}
// 	got := DecodeHashMapSnapshot(EncodeHashMapToSnapshot(input))
// 	if !reflect.DeepEqual(want, got) {
// 		t.Fatalf("test expected: %v, got: %v\n", want, got)
// 	}
// }

// func TestSegmentAppend(t *testing.T) {
// 	s := Segment{
// 		id:      "1",
// 		hashmap: make(map[string]uint32),
// 	}
// 	s.Append("Whatever", "54321")
// 	fmt.Print(s.hashmap)
// }

// func TestSegmentGetValue(t *testing.T) {
// 	s := Segment{
// 		id:      "2",
// 		hashmap: make(map[string]uint32),
// 	}

// 	s.Append("Hello", "12345")
// 	s.Append("Whatever", "54321")

// 	val, _ := s.GetValue("Hello")
// 	fmt.Println(val)
// }

// func TestSegmentCompress(t *testing.T) {
// 	s := Segment{
// 		id:      "3",
// 		hashmap: make(map[string]uint32),
// 	}

// 	s.Append("Hello", "12345")
// 	s.Append("Whatever", "54321")
// 	s.Append("Hello", "ThisIsNew")
// 	s.Append("Whatever", "This is newer")

// 	newSegment := s.Compress()
// 	fmt.Println("SegmentCompress result")
// 	fmt.Print(newSegment.hashmap)
// 	for key := range newSegment.hashmap {
// 		v, _ := newSegment.GetValue(key)
// 		fmt.Printf("(%s: %s)", key, v)
// 	}

// 	DeleteFile(newSegment.id)

// }
