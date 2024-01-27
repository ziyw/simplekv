package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type SegList struct {
	list []*Segment
}

func NewSegList() *SegList {
	sl := &SegList{
		list: make([]*Segment, 0),
	}
	sl.LoadSegments()

	go sl.PersistSnapshots()
	return sl
}

func (sl SegList) Size() int {
	return len(sl.list)
}

func (sl *SegList) AddSegment(s *Segment) {
	sl.list = append(sl.list, s)
}

func (sl *SegList) RemoveSegment(target int) {
	for i, s := range sl.list {
		if s.id != target {
			continue
		}

		if s.id == target {
			sl.list = append(sl.list[:i], sl.list[i+1:]...)
			return
		}
	}
}

func (sl SegList) GetNextId() int {
	if len(sl.list) == 0 {
		return 1
	}

	last := sl.list[len(sl.list)-1]
	return last.id + 1
}

func (sl *SegList) GetCurrentSegment() *Segment {
	if len(sl.list) == 0 {
		return nil
	}
	return sl.list[len(sl.list)-1]
}

// Load persisted segment into active segment list.
func (sl *SegList) LoadSegments() {
	fmt.Println("Load Segments")

	files, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.Contains(f.Name(), "snapshot_") {
			fmt.Println("Load Segments from snapshot" + f.Name())
			sl.AddSegment(NewSegmentFromSnapshot(f.Name()))
		}
	}
	fmt.Println("Finish Load Segments")
}

// Periodic function to persist all snapshots of all active segments
func (sl *SegList) PersistSnapshots() {
	for {
		time.Sleep(time.Second)
		fmt.Println("Persist all snapshots")
		for _, s := range sl.list {
			s.CreateSnapshot()
		}
	}
}
