package engine

import "testing"

func TestHashmapPut(t *testing.T) {
	h := NewHashmap()
	err := h.Put(1, 2)
	if err != nil {
		t.Errof("should not got error: %v", err)
	}

	// expect error
	err = h.Put(1, -12)
	if err == nil {
		t.Errorf("should got error: %v", err)
	}
}
