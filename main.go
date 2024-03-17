package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	fnvhash "github.com/i-m-afk/bloom-filter/internal/fnv-hash"
)

const BitArraySize = 1000000

func main() {
	bf := make([]bool, BitArraySize) // bloom filter

	err := loadDict("./dict.txt", bf)
	if err != nil {
		log.Fatal(err)
	}
	test("./dict.txt", bf)
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

// insert items to bool array
// fnvhash produces a uint64 number.
// this number can be divided by the total size of the bit array.
// to prevent overflow
// TODO: this is highly redundant, because of hash functions.
// hash functions can also be produced by variation introduced here
func insertItems(data []byte, arr []bool) {
	o := fnvhash.Fnv1(data) % BitArraySize
	a := fnvhash.Fnv1a(data) % BitArraySize
	b := fnvhash.Fnv1a(data) % BitArraySize
	c := fnvhash.Fnv1c(data) % BitArraySize
	d := fnvhash.Fnv1d(data) % BitArraySize
	e := fnvhash.Fnv1e(data) % BitArraySize
	f := fnvhash.Fnv1f(data) % BitArraySize
	g := fnvhash.Fnv1g(data) % BitArraySize
	arr[o] = true
	arr[a] = true
	arr[b] = true
	arr[c] = true
	arr[d] = true
	arr[e] = true
	arr[f] = true
	arr[g] = true
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
	// Check each hash function's index in the bit array
	indices := []uint64{
		fnvhash.Fnv1([]byte(word)) % BitArraySize,
		fnvhash.Fnv1a([]byte(word)) % BitArraySize,
		fnvhash.Fnv1b([]byte(word)) % BitArraySize,
		fnvhash.Fnv1c([]byte(word)) % BitArraySize,
		fnvhash.Fnv1d([]byte(word)) % BitArraySize,
		fnvhash.Fnv1e([]byte(word)) % BitArraySize,
		fnvhash.Fnv1f([]byte(word)) % BitArraySize,
		fnvhash.Fnv1g([]byte(word)) % BitArraySize,
	}

	// the word is not in the dictionary if any index if false in bloom-filter
	for _, index := range indices {
		if !arr[index] {
			return false
		}
	}
	return true
}
