package engine

import (
	"testing"
)

func TestNewSegmentFile(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{name: "test create with empty filename", input: "", wantError: true},
		{name: "test with normal string", input: "TestSegFileName", wantError: false},
		{name: "test with exsiting file", input: "seg_file.go", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewSegFile(tt.input)
			if s != nil {
				defer s.Delete()
			}

			gotError := err != nil
			if gotError != tt.wantError {
				t.Errorf("error NewSegFile, want return nil is %v, got %v", tt.wantError, err)
			}
		})
	}
}
