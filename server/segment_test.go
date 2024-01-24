package main

import (
	"reflect"
	"testing"
)

func TestEncodeAndDecodeHashmap(t *testing.T) {
	want := map[string]string{"Hello": "12345", "World": "54321"}

	input := map[string]string{"Hello": "12345", "World": "54321"}
	got := DecodeHashMapSnapshot(EncodeHashMapToSnapshot(input))
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("test expected: %v, got: %v\n", want, got)
	}
}
