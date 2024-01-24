package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

const METADATA_SIZE = 4
const HEADER = "PAGE"

type Page struct {
	id      string
	header  string
	hashmap map[string]string
}

// Load read page from disk to memory
func Load(page_id string) *Page {
	pageFileName := "page_" + page_id

	// if file doesn't exist, return an empty hashmap directly
	_, err := os.Stat(pageFileName)
	if err != nil && os.IsNotExist(err) {
		return &Page{
			id:      page_id,
			header:  HEADER,
			hashmap: make(map[string]string),
		}
	}

	byteArr := read(pageFileName)

	if string(byteArr[:4]) != HEADER {
		log.Fatalf("page header doesn't allign")
	}

	return &Page{
		id:      page_id,
		header:  HEADER,
		hashmap: decodeByteArrToHashMap(byteArr[4:]),
	}
}

// Flush persists a page to disk
func (p Page) Flush() {
	pageFileName := "page_" + p.id

	// check if file exist
	_, err := os.Stat(pageFileName)
	if err == nil {
		fmt.Print("file exist, remove file")
		os.Remove(pageFileName)
	} else if os.IsNotExist(err) {
		fmt.Print("file doesn't exist, skip actions")
	}

	file, err := os.Create(pageFileName)
	if err != nil {
		log.Fatalf("error creating page file: %v", err)
	}
	defer file.Close()

	bw := bufio.NewWriter(file)
	byteArr := []byte(p.header)
	byteArr = append(byteArr, encodeHashMap(p.hashmap)...)
	bw.Write(byteArr)

	if err := bw.Flush(); err != nil {
		log.Fatalf("error flushing byte array: %v", err)
	}
}

func read(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("error opening page file %s, : %v", filename, err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Fatalf("error getting file stats %s: %v", filename, err)
	}

	byteArr := make([]byte, stat.Size())

	if _, err = bufio.NewReader(file).Read(byteArr); err != nil && err != io.EOF {
		log.Fatalf("error read from file %s: %v", filename, err)
	}

	return byteArr
}

// Encode string-to-string hashmap to byte array
func encodeHashMap(hashMap map[string]string) []byte {
	byteArr := []byte{}
	byteArr = append(byteArr, encodeNumToBytes(uint32(len(hashMap)), METADATA_SIZE)...)

	for key, value := range hashMap {
		b := encodeStrToByteArr(key)
		b = append(b, encodeStrToByteArr(value)...)
		byteArr = append(byteArr, b...)
	}
	return byteArr
}

// Decode byte array to string-string hashmap
func decodeByteArrToHashMap(byteArr []byte) map[string]string {
	hashMap := make(map[string]string)

	n := int(decodeBytesToNum(byteArr[:METADATA_SIZE]))
	offset := METADATA_SIZE
	for i := 0; i < n; i++ {
		key, nxt := decodeByteArrToString(byteArr, offset)
		value, nxt := decodeByteArrToString(byteArr, nxt)
		hashMap[key] = value
		offset = nxt
	}
	return hashMap
}

// Encode and decode  string to byte array
func encodeStrToByteArr(in string) []byte {
	// String Metadat: 4 bytes content length, variable content size
	content := []byte(in)
	meta := encodeNumToBytes(uint32(len(in)), METADATA_SIZE)
	return append(meta, content...)
}

func decodeByteArrToString(byteArr []byte, offset int) (string, int) {
	meta := decodeBytesToNum(byteArr[offset : offset+METADATA_SIZE])
	contentOffset := offset + METADATA_SIZE
	contentStop := contentOffset + int(meta)
	str := string(byteArr[offset+METADATA_SIZE : contentStop])
	return str, contentStop
}

// Encode and decode int to byte array
func encodeIntToByteArr(in int) []byte {
	// metadata: 4 bytes content lenght, 1 byte sign flag, 4 byte unsign int value
	// TODO: change this to use one bit
	var flag byte
	var content []byte
	if in > 0 {
		flag = 0
		content = encodeNumToBytes(uint32(in), METADATA_SIZE)
	} else {
		flag = 255
		content = encodeNumToBytes(uint32(-in), METADATA_SIZE)
	}
	meta := encodeNumToBytes(uint32(len(content)), METADATA_SIZE)
	meta = append(meta, flag)
	return append(meta, content...)
}

func decodeByteArrToInt(byteArr []byte, offset int) (int, int) {
	meta := decodeBytesToNum(byteArr[offset : offset+METADATA_SIZE])
	flag := byteArr[offset+METADATA_SIZE : offset+METADATA_SIZE+1][0]

	contentOffset := offset + METADATA_SIZE + 1
	contentStop := contentOffset + int(meta)

	content := int(decodeBytesToNum(byteArr[contentOffset:contentStop]))
	if flag == 255 {
		return -1 * content, contentStop
	}
	return content, contentStop
}

// Fundamental: encode uint32 to bytes and decode bytes to uint32
// These two are used in all meta data handling
func encodeNumToBytes(num uint32, size int) []byte {
	output := make([]byte, size)
	binary.LittleEndian.PutUint32(output, uint32(num))
	return output
}

func decodeBytesToNum(input []byte) uint32 {
	var num uint32
	buf := bytes.NewReader(input)
	err := binary.Read(buf, binary.LittleEndian, &num)
	if err != nil {
		log.Fatal(err)
	}
	return num
}

// func main() {

// 	hmap := make(map[string]string)
// 	hmap["hello"] = "world"

// 	p_1 := Page{
// 		id:      "1",
// 		header:  HEADER,
// 		hashmap: hmap,
// 	}
// 	p_1.Flush()

// 	newPage := Load("1")
// 	fmt.Println("Page content is ")
// 	fmt.Print(newPage.hashmap)

// 	newPage.hashmap["What"] = "ever"
// 	newPage.Flush()

// 	anotherPage := Load("1")
// 	fmt.Println("Updated Page content is ")
// 	fmt.Println(anotherPage.hashmap)
// }
