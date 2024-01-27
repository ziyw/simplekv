package main

import (
	"testing"
)

func TestNewSegList(t *testing.T) {
	sl := NewSegList()
	if sl.Size() != 0 {
		t.Fatalf("test failed: new seglist size is not zero")
	}
}

func TestAddSegment(t *testing.T) {
	sl := NewSegList()

	sl.AddSegment(NewSegment(1))
	if sl.Size() != 1 {
		t.Fatalf("test failed: add segment didn't increase size")
	}

	sl.AddSegment(NewSegment(2))
	if sl.Size() != 2 {
		t.Fatalf("test failed: add segment didn't increase size")
	}
}

func TestCurrentSegment(t *testing.T) {
	sl := NewSegList()

	if sl.GetCurrentSegment() != nil {
		t.Fatal("test failed: current segment should be nil when no segment active")
	}

	sl.AddSegment(NewSegment(1))
	sl.AddSegment(NewSegment(2))
	sl.AddSegment(NewSegment(3))

	if sl.GetCurrentSegment().id != 3 {
		t.Fatal("test failed, current segment is not the latest")
	}

	sl.RemoveSegment(2)
	if sl.GetCurrentSegment().id != 3 {
		t.Fatal("test failed, after removal, current segment is not the latest")
	}

	sl.RemoveSegment(3)
	if sl.GetCurrentSegment().id != 1 {
		t.Fatal("test failed, after remove all, current segment is not 1")
	}

	sl.RemoveSegment(1)
	if sl.GetCurrentSegment() != nil {
		t.Fatal("test failed: current segment should be nil when no segment active")
	}
}

func TestRemoveSegment(t *testing.T) {
	sl := NewSegList()
	sl.AddSegment(NewSegment(1))
	sl.AddSegment(NewSegment(2))
	sl.AddSegment(NewSegment(3))

	sl.RemoveSegment(2)
	if sl.Size() != 2 {
		t.Fatalf("test failed: add segment didn't increase size")
	}

	sl.RemoveSegment(2)
	if sl.Size() != 2 {
		t.Fatalf("test failed: removal is not idempotent")
	}

	sl.RemoveSegment(1)
	sl.RemoveSegment(3)
	if sl.Size() != 0 {
		t.Fatal("test failed: didn't remove all segs")
	}
}

func TestNextId(t *testing.T) {
	sl := NewSegList()
	if sl.GetNextId() != 1 {
		t.Fatal("test failed: first next id is not zero")
	}

	sl.AddSegment(NewSegment(12))
	if sl.GetNextId() != 13 {
		t.Fatal("test failed: next id should be 13")
	}
}
