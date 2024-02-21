package engine

import "testing"

func TestCheckExist(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"test exist file", "file_test.go", true},
		{"test non-exist file", "file_not_exist", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckExist(tt.input)
			if got != tt.want {
				t.Errorf("%s: want: %v, got: %v", tt.name, tt.want, got)
			}
		})
	}
}

func TestCreateFile(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"test create existing file", "file_test.go", true},
		{"test create non-exist file", "new_test_file", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateFile(tt.input)
			if (got != nil) != tt.wantError {
				t.Errorf("%s: want: %v, got: %v", tt.name, tt.wantError, got)
			}

			// clean up
			if got == nil {
				DeleteFile(tt.input)
			}
		})
	}
}

func TestDeleteFile(t *testing.T) {
	err := DeleteFile("non-exist-file")
	if err == nil {
		t.Error("test delete non-exist file: want error, got: nil")
	}
}

func TestOpenFile(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"test open existing file", "file.go", false},
		{"test open non-exist file", "new_test_file", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, gotErr := OpenFile(tt.input)
			if (gotErr != nil) != tt.wantError {
				t.Errorf("%s: want: %v, got: %v", tt.name, tt.wantError, gotErr)
			}
			defer f.Close()
		})
	}

}
