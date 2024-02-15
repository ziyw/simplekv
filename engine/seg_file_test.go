package engine

import (
	"reflect"
	"testing"
)

func TestNewEntry(t *testing.T) {
	tests := []struct {
		name       string
		inputKey   string
		inputValue string
		want       string
	}{
		{name: "test emtpy input", inputKey: "", inputValue: "", want: ",\n"},
		{name: "test normal input", inputKey: "Key", inputValue: "Value", want: "Key,Value\n"},
		{name: "test illegal key", inputKey: "key\n,", inputValue: "value", want: ""},
		{name: "test illegal value", inputKey: "key", inputValue: "value,\n", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewEntry(tt.inputKey, tt.inputValue)
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("test %s error, want %v, got %v", tt.name, tt.want, got)
			}
		})
	}
}

func TestParseNewEntry(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantKey   string
		wantValue string
		wantError bool
	}{
		{name: "test empty input", input: "", wantKey: "", wantValue: "", wantError: true},
		{name: "test illegal input", input: "whatever", wantKey: "", wantValue: "", wantError: true},
		{name: "test normal input", input: "key,value\n", wantKey: "key", wantValue: "value", wantError: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotValue, err := ParseEntry(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("test %s should return error, input: %v", tt.name, tt.input)
			}
			if gotKey != tt.wantKey || gotValue != tt.wantValue {
				t.Errorf("test %s error, want %v %v, got %v %v", tt.name, gotKey, gotValue, tt.wantKey, tt.wantValue)
			}
		})
	}
}

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

func TestAppend(t *testing.T) {
	// if file doesn't exist before Append, get error
	badSegFile := SegFile{
		name: "not_exist",
	}
	_, err := badSegFile.Append("helloworld")
	if err == nil {
		t.Errorf("when segFile is not created, should return error")
	}

	// when file already exist, directly append to file
	goodSegFile, _ := NewSegFile("goodFile")
	//defer goodSegFile.Delete()

	offset, _ := goodSegFile.Append("HelloWorld")
	if offset != 0 {
		t.Errorf("offset is wrong, want %v, got %v", 0, offset)
	}
	newOffset, _ := goodSegFile.Append("HelloAgain")
	if newOffset == 0 {
		t.Errorf("offset is wrong, new offset should not be zero")
	}
}

func TestGet(t *testing.T) {}

func TestGetAll(t *testing.T) {}

func TestCompress(t *testing.T) {}

func TestMerge(t *testing.T) {}
