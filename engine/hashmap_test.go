package engine

import (
	"testing"
)

func TestNewHashMap(t *testing.T) {
	hmap, err := NewHashMap("ThisIsNew")
	defer DeleteFile(hmap.filename)

	if err != nil {
		t.Error(err)
	}

	_, err = NewHashMap("ThisIsNew")
	if err == nil {
		t.Error("this file should already exist")
	}
}

func TestLoadHashMap(t *testing.T) {
	hashMapFileName := "ThisIsNew"
	CreateFile(hashMapFileName)
	_, err := LoadHashMap(hashMapFileName)
	if err != nil {
		t.Error("test LoadHashMap failed: got error when load file")
	}

	DeleteFile(hashMapFileName)
	_, err = LoadHashMap(hashMapFileName)
	if err == nil {
		t.Error("test LoadHashMap failed: should get error, file already deleted")
	}
}

func TestPutAndGet(t *testing.T) {
	hmap, err := NewHashMap("TestPutAndGet")
	if err != nil {
		t.Errorf("test Put failed: can't open file")
	}
	defer DeleteFile(hmap.filename)

	hmap.Put("hello", 12)
	hmap.Put("world", 13)
	v, err := hmap.Get("hello")
	if v != 12 {
		t.Errorf("test Get failed: value should be 12")
	}
	v, err = hmap.Get("not-exist")
	if err == nil {
		t.Errorf("test Get failed: should return error")
	}
}

func TestPersist(t *testing.T) {
	hmap, err := NewHashMap("TestPersist")
	if err != nil {
		t.Errorf("test Persist failed: can't open file")
	}
	defer DeleteFile(hmap.filename)

	hmap.Put("hello", 12)
	hmap.Put("world", 13)
	err = hmap.Persist()
	if err != nil {
		t.Errorf("test Persist failed: %v", err)
	}

	result, err := LoadHashMap("TestPersist")
	if err != nil {
		t.Errorf("test Persist failed: %v", err)
	}
	err = result.Load()
	if err != nil {
		t.Errorf("test Load failed: %v", err)
	}
	v, err := result.Get("hello")
	if v != 12 {
		t.Errorf("test Load failed: should return 12")
	}
}
