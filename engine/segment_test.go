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

func TestLoadExistingSegment(t *testing.T) {
	s1, err := NewSegment("seg_1")
	if err != nil {
		t.Error("segment creation failed")
	}
	defer s1.Delete()
	s1.Put("hello", "world")
	s1.Put("Whatever", "second line content")
	s1.hashmap.Persist()

	s2, err := Load("seg_1")
	if err != nil {
		t.Error("segment load failed", err)
	}

	fmt.Println("Loaded hashmap is")
	fmt.Println(s2.hashmap.mem)

	v, err := s2.Get("hello")
	if err != nil {
		t.Error("load value from new segment failed", err)
		return
	}

	k1, _, _ := s1.GetAll()
	k2, _, _ := s2.GetAll()
	if v != "world" {
		t.Errorf("load value error, want %v, got %v", k1, k2)
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

func TestSegmentPut_OverMaxItems(t *testing.T) {
	s, _ := NewSegment("test_max_items")
	defer s.Delete()

	for i := 0; i < MAX_ITEM_NUM; i++ {
		key, value := fmt.Sprintf("key is %d", i), fmt.Sprintf("value is %d", i)
		s.Put(key, value)
	}
	fmt.Println(s.GetAll())

	fmt.Println("Print one more")
	err := s.Put("MoreThanOne", "LastOne")
	if err == nil {
		t.Error("should return exceed error", err)
	} else {
		fmt.Println("error is " + err.Error())
	}
}

func TestCompressSegment(t *testing.T) {
	s, err := NewSegment("seg_1")
	if err != nil {
		t.Error("error create new segment")
	}
	defer s.Delete()

	s.Put("hello", "1")
	s.Put("hello", "2")
	s.Put("hello", "3")

	s2, err := s.Compress("seg_2")
	if err != nil {
		t.Error("compress seg_1 failed", err)
	}

	v, err := s2.Get("hello")
	if err != nil {
		t.Error("compressed result is not correct", err)
	}
	defer s2.Delete()
	if v != "3" {
		t.Errorf("compressed result value is not correct, want %v, got %v", "3", v)
	}
}

func TestMergeSegments(t *testing.T) {
	s1, _ := NewSegment("s1")
	s2, _ := NewSegment("s2")

	s1.Put("hello", "s1")
	s1.Put("World", "s1")
	s2.Put("hello", "s2")
	s2.Put("Whatever", "s2")

	defer s1.Delete()
	defer s2.Delete()

	s3, err := Merge(s1, s2, "s3")
	if err != nil {
		t.Errorf("error mergin s1 and s2, %v", err)
	}
	defer s3.Delete()

	v, err := s3.Get("hello")
	if err != nil {
		t.Errorf("error getting key from s3, %v", err)
	}

	if v != "s2" {
		t.Errorf("value is wrong for key hello, want s2, get %v", v)
	}

	v, err = s3.Get("World")
	if err != nil {
		t.Errorf("get value error, %v", err)
	}
	if v != "s1" {
		t.Errorf("value is wrong for key World, want s1, got %v", v)
	}
}
