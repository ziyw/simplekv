package main

type SegList struct {
	list []*Segment
}

func NewSegList() *SegList {
	return &SegList{
		list: make([]*Segment, 0),
	}
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
