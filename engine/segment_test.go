package engine

import (
	"fmt"
	"testing"

	"golang.org/x/exp/slices"
)

func TestNewSegment_DuplicateFileNameError(t *testing.T) {
	s1, err := NewSegment("seg_1")
	if err != nil {
		t.Error("segment creation failed", err)
	}
	defer s1.Delete()

	_, err = NewSegment("seg_1")
	if err == nil {
		t.Error("duplicate segment name should result in error")
	}
}

func TestDelete_Normal(t *testing.T) {
	s1, err := NewSegment("seg_1")
	if err != nil {
		t.Error("segment creation failed")
	}
	s1.Delete()

	s2, err := NewSegment("seg_1")
	if err != nil {
		t.Error("deletion failed, segment creation failed")
	}
	s2.Delete()
}

// TODO: good way to return error for illegal creation case
// func TestInvalidSegmentCreation(t *testing.T) {
// 	s := Segment{
// 		id: "1",
// 	}
// 	if s.segFile == nil {
// 		t.Error("segment should not be created from structu")
// 	}
// }

func TestSegmentPutAndGet(t *testing.T) {
	s, _ := NewSegment("1")
	defer s.Delete()

	s.Put("hello", "world")
	s.Put("Second", "This is the second line")
	s.Put("3", "3 lines in total")

	fmt.Println(s.hashmap.mem)
	if _, ok := s.hashmap.mem["hello"]; !ok {
		t.Error("hashmap should contain key")
	}

	got, _ := s.Get("hello")
	if got != "world" {
		t.Error("get function error, should return world")
	}

	got, _ = s.Get("3")
	if got != "3 lines in total" {
		t.Error("get function error, should return 3 lines in total")
	}

	_, err := s.Get("NotExist")
	if err == nil {
		t.Error("should return error for not exist keys")
	}
}

func TestSegmentGetAll(t *testing.T) {
	s, _ := NewSegment("1")
	defer s.Delete()

	s.Put("hello", "world")
	s.Put("Second", "This is the second line")
	s.Put("3", "3 lines in total")

	keys, values, err := s.GetAll()
	if err != nil {
		t.Errorf("error getAll: %v", err)
	}
	if !slices.Contains(keys, "hello") || !slices.Contains(keys, "Second") {
		t.Errorf("getAll doesn't contain all keys, %v", keys)
	}
	if !slices.Contains(values, "world") || !slices.Contains(values, "This is the second line") {
		t.Errorf("getAll doesn't contain all keys, %v", values)
	}
}
