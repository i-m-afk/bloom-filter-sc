package fnvhash

import (
	"crypto/rand"
	"fmt"
)

func Hashbench() {
	numKeys := 1050000 // total keys
	numFunctions := 15

	hashCounts := make(map[uint64]int) // map to store hash counts

	for i := 0; i < numKeys; i++ {
		key := make([]byte, 10)
		rand.Read(key)
		for j := 0; j < numFunctions; j++ {
			hash := Fnv1(key, j)
			hashCounts[hash]++
		}
	}

	// collisions count
	collisions := 0
	for _, count := range hashCounts {
		fmt.Println(count)
		if count > 1 {
			collisions += count
		}
	}
	fmt.Printf("Total collisions : %d\n", collisions)
}
