package engine

import (
	"fmt"
	"testing"
)

func TestPut_NoExistingSegment_shouldCreateNewFile(t *testing.T) {
	b := BufferManager{}
	b.Put("hello", "world")

	segFileName := ToSegFileName(1)
	if !CheckExist(segFileName) {
		t.Errorf("BufferManager should created segfile %s", segFileName)
	}

	seg, err := Load(1)
	defer seg.Delete()

	if err != nil {
		t.Errorf("New segfile should be able to load %s, err: %v", segFileName, err)
	}

	keys, values, err := seg.GetAll()
	if err != nil {
		t.Errorf("Error get all from new segment: %v", err)
	}

	if len(keys) != 1 || len(values) != 1 {
		t.Errorf("Error getting all keys and values from segment, keys: %v, values: %v", keys, values)
	}

	if keys[0] != "hello" || values[0] != "world" {
		t.Errorf("error getting values keys: %v, values: %v", keys, values)
	}

}

func TestPut_HaveActiveSegment_NotFull(t *testing.T) {
	// if there are more than one segment, BufferManager should select the latest
	// even when the previous two are all empty
	s1, _ := NewSegment(1)
	s2, _ := NewSegment(2)
	want, _ := NewSegment(3)
	defer s1.Delete()
	defer s2.Delete()
	defer want.Delete()

	b := BufferManager{}
	b.Put("hello", "world")
	b.Put("foo", "bar")

	// There is a difference between persisted hashmap and hashmap in memeory
	got, _ := Load(3)
	keys, values, err := got.GetAll()
	if err != nil {
		t.Errorf("error getting keys and values: %v", err)
	}
	if len(keys) != 2 || len(values) != 2 {
		t.Errorf("error keys: %v, values: %v", keys, values)
	}
}

// if exceed max item, put should create a new segment
func TestPut_ExceedMaxItem(t *testing.T) {
	b := BufferManager{}

	for i := 0; i <= MAX_ITEM_NUM+2; i++ {
		b.Put(fmt.Sprintf("key_%d", i), fmt.Sprintf("value_%d", i))
	}

	if !CheckExist(ToSegFileName(1)) || !CheckExist(ToSegFileName(2)) {
		t.Errorf("BufferManage did not create a new segment file")
	}

	s1, _ := Load(1)
	defer s1.Delete()
	keys, values, _ := s1.GetAll()
	if len(keys) != MAX_ITEM_NUM || len(values) != MAX_ITEM_NUM {
		t.Errorf("error keys:%v, values: %v", keys, values)
	}

	s2, _ := Load(2)
	defer s2.Delete()
	keys, values, _ = s2.GetAll()
	if len(keys) != 2 || len(values) != 2 {
		t.Errorf("error keys:%v, values: %v", keys, values)
	}

}

func TestGet_SimpleCase(t *testing.T) {
	s1, _ := NewSegment(1)
	s2, _ := NewSegment(2)
	s3, _ := NewSegment(3)
	defer s1.Delete()
	defer s2.Delete()
	defer s3.Delete()

	s1.Put("dup", "s1")
	s1.Put("unique", "s1")
	s2.Put("dup", "s2")
	s3.Put("dup", "s3")

	b := BufferManager{}
	value, err := b.Get("dup")
	if err != nil {
		t.Errorf("error getting value: %v", err)
	}
	if value != "s3" {
		t.Errorf("want: k3, got: %v", value)
	}

	value, err = b.Get("unique")
	if err != nil {
		t.Errorf("error getting value: %v", err)
	}
	if value != "s1" {
		t.Errorf("want :s1, got: %v", value)
	}
}
