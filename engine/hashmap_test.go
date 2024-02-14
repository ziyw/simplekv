package engine

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewHashMap(t *testing.T) {
	hmap := NewHashMap()
	if len(hmap) != 0 {
		t.Error("HashMap init failed")
	}
}

func TestPutAndGet(t *testing.T) {
	hmap := NewHashMap()
	inputKey, inputValue := "Hello", offset(12)
	hmap[inputKey] = offset(inputValue)
	got := hmap[inputKey]
	fmt.Println(hmap)

	if !reflect.DeepEqual(got, inputValue) {
		t.Errorf("HashMap put/get wrong: want: %v, got: %v", inputValue, got)
	}

}

func TestHashMapEncodeAndDecode(t *testing.T) {
	hmap := NewHashMap()
	hmap["hello"] = 12
	hmap["world"] = 13

	encoded, _ := Encode(hmap)
	got, _ := Decode(encoded)

	if !reflect.DeepEqual(hmap, got) {
		t.Errorf("HashMap decode/encode wrong: want %v, got %v", hmap, got)
	}
}

func TestEmptyEncoding(t *testing.T) {
	hmap := NewHashMap()
	encoded, _ := Encode(hmap) // encoded empty map is not empty byte array
	got, _ := Decode(encoded)
	if !reflect.DeepEqual(hmap, got) {
		t.Errorf("HashMap decode/encode wrong: want %v, got %v", hmap, got)
	}
}
