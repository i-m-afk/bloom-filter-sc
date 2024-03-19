package main

import (
	"bufio"
	"log"
	"os"
	"testing"

	"github.com/i-m-afk/bloom-filter/internal/fnv-hash"
)

// Test for collisions for all the 15 different hashes
func TestFnv1Collision(t *testing.T) {
	testData := loadTestData("dict.txt")
	hashValues := make(map[uint64]bool)
	count := 0
	for _, data := range testData {
		for i := 0; i < 15; i++ {
			hash := fnvhash.Fnv1([]byte(data), i)
			if hashValues[hash] {
				count++
				t.Errorf("Collision detected for input %q", data)
			}
			hashValues[hash] = true
		}
		t.Log("Total collisions: \n", count)
	}
}

func loadTestData(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var testData []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		testData = append(testData, scanner.Text())
	}
	return testData
}
