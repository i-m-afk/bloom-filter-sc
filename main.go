/*
* Author: Rishav Kumar
 */
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	fnvhash "github.com/i-m-afk/bloom-filter/internal/fnv-hash"
)

// Calculated size and number of hashfunctions
// See readme
const BitArraySize = 1507404
const HashNum = 15

func main() {
	dict := flag.String("dict", "dict.txt", "dictionary file")
	flag.Parse()

	bf := make([]bool, BitArraySize) // bloom filter

	err := loadDict(*dict, bf)
	if err != nil {
		log.Fatal(err)
	}
	test(*dict, bf)
	spellcheckUserInput(bf)
}

func loadDict(path string, arr []bool) error {
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
func insertItems(data []byte, arr []bool) {
	for i := 0; i < HashNum; i++ {
		hash := fnvhash.Fnv1(data, i)
		arr[hash%BitArraySize] = true
	}
}

// testing the recall of bloom-filter
func test(path string, arr []bool) error {
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
func spellcheckUserInput(arr []bool) {
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

func isWordInDictionary(word string, arr []bool) bool {
	// the word is not in the dictionary if any index if false in bloom-filter
	for i := 0; i < HashNum; i++ {
		hash := fnvhash.Fnv1([]byte(word), i)
		index := hash % BitArraySize
		if !arr[index] {
			return false
		}
	}
	return true
}
