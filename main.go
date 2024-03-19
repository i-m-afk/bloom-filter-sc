/*
* Author: Rishav Kumar
 */
package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	fnvhash "github.com/i-m-afk/bloom-filter/internal/fnv-hash"
)

const (
	BitArraySize = 1559967
	HashNum      = 15
	BitsPerByte  = 8
	HeaderSize   = 12
)

func main() {
	dict := flag.String("dict", "dict.txt", "dictionary file")
	bloomfilter := flag.String("bf", "words.bf", "bloom filter file-name")
	flag.Parse()

	bf := make([]byte, (BitArraySize+BitsPerByte-1)/BitsPerByte) // bloom filter

	_, err := os.Stat(*bloomfilter)

	if err == nil {
		err := loadBloomFilter(*bloomfilter, bf)
		if err != nil {
			log.Println("failed to load bloom filter")
			log.Fatal(err)
		}
		log.Println("Bloom Filter loaded successfuly")
		spellcheckUserInput(bf)
		return
	}

	err = loadDict(*dict, bf)
	if err != nil {
		log.Fatal(err)
	}
	test(*dict, bf)
	saveBloomFilter(*bloomfilter, bf)
	spellcheckUserInput(bf)
}

func loadDict(path string, arr []byte) error {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// fill the bit array
	for scanner.Scan() {
		insertItems([]byte(scanner.Text()), arr)
	}
	err = scanner.Err()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// insert items to bool array (bloom filter)
// fnvhash produces a uint64 number.
// this number can be divided by the total size of the bit array.
// to prevent overflow
// see readme to know about index calculation
func insertItems(data []byte, arr []byte) {
	for i := 0; i < HashNum; i++ {
		hash := fnvhash.Fnv1(data, i)
		index := hash % BitArraySize
		byteIndex := index / BitsPerByte
		bitIndex := index % BitsPerByte
		arr[byteIndex] |= 1 << bitIndex
	}
}

// testing the recall of bloom-filter
func test(path string, arr []byte) error {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	countTrue := 0
	countFalse := 0
	// check
	for scanner.Scan() {
		exist := isWordInDictionary(scanner.Text(), arr)
		if exist {
			countTrue++
		} else {
			countFalse++
		}
	}
	fmt.Printf("True: %v False Negatives: %v Total: %d\n", countTrue, countFalse, countFalse+countTrue)
	err = scanner.Err()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil

}

// cli for user to check there spellings
func spellcheckUserInput(arr []byte) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter a word to check its spelling: ")

	for scanner.Scan() {
		word := scanner.Text()

		if !isWordInDictionary(word, arr) {
			fmt.Printf("'%s' is misspelled.\n", word)
		} else {
			fmt.Printf("'%s' is spelled correctly.\n", word)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func isWordInDictionary(word string, arr []byte) bool {
	// the word is not in the dictionary if any index produced by hash is false in bloom-filter
	for i := 0; i < HashNum; i++ {
		hash := fnvhash.Fnv1([]byte(word), i)
		index := hash % BitArraySize
		byteIndex := index / BitsPerByte
		bitIndex := index % BitsPerByte
		if (arr[byteIndex] & (1 << bitIndex)) == 0 {
			return false
		}
	}
	return true
}

// save the bloom filter data in bytes
func saveBloomFilter(filename string, arr []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	header := createHeader()
	_, err = writer.Write(header)
	if err != nil {
		return err
	}

	_, err = writer.Write(arr)
	if err != nil {
		return err
	}
	return writer.Flush()
}

func loadBloomFilter(filename string, arr []byte) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	err = verifyHeader(file)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	_, err = reader.Read(arr)
	if err != nil {
		return err
	}
	return nil
}

// Headers as per challenge description
func createHeader() []byte {
	head := make([]byte, 0)
	head = append(head, []byte("CCBF")...)
	// version
	version := 1
	head = append(head, byte(version>>8))
	head = append(head, byte(version))

	// number of hash
	numHash := HashNum
	head = append(head, byte(numHash>>8))
	head = append(head, byte(numHash))

	var size int64 = BitArraySize
	// big endian notation (4 bytes)
	head = append(head, byte(size>>24))
	head = append(head, byte(size>>16))
	head = append(head, byte(size>>8))
	head = append(head, byte(size))
	return head
}

func verifyHeader(file *os.File) error {
	header := make([]byte, HeaderSize)
	_, err := file.Read(header)
	if err != nil {
		return err
	}
	// Identifier
	expect := []byte("CCBF")
	if !bytes.Equal(header[:4], expect) {
		return errors.New("Invalid Identifier")
	}
	// version
	version := 1
	expectedVersion := make([]byte, 2)
	expectedVersion = append(expectedVersion, byte(version>>8))
	expectedVersion = append(expectedVersion, byte(version))

	if bytes.Equal(expectedVersion, header[4:6]) {
		return errors.New("Invalid Version")
	}
	// num hash
	hashNum := HashNum
	expectedHashNum := make([]byte, 2)
	expectedHashNum = append(expectedHashNum, byte(hashNum>>8))
	expectedHashNum = append(expectedHashNum, byte(hashNum))

	if bytes.Equal(expectedHashNum, header[6:8]) {
		return errors.New("Invalid Number of hash Function")
	}
	expectedBitArraySize := make([]byte, 0)
	var size int64 = BitArraySize
	expectedBitArraySize = append(expectedBitArraySize, byte(size>>24))
	expectedBitArraySize = append(expectedBitArraySize, byte(size>>16))
	expectedBitArraySize = append(expectedBitArraySize, byte(size>>8))
	expectedBitArraySize = append(expectedBitArraySize, byte(size))
	if !bytes.Equal(header[8:], expectedBitArraySize) {
		return errors.New("Invalid bit array size")
	}
	return nil
}
