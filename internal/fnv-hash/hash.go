/*
* Author : Rishav Kumar
* references: https://datatracker.ietf.org/doc/html/draft-eastlake-fnv-17.html#section-6.1.1
 */
package fnvhash

const FNV_offset_basis uint64 = 0xcbf29ce484222325
const FNV_prime uint64 = 1099511628211

func Fnv1(data []byte, index int) uint64 {
	hash := FNV_offset_basis
	for _, b := range data {
		hash *= FNV_prime
		hash ^= uint64(b)
	}
	return hash*Fnv1a(data, index) + uint64(index)
}

func Fnv1a(data []byte, index int) uint64 {
	hash := FNV_offset_basis
	for _, b := range data {
		hash ^= uint64(b)
		hash *= FNV_prime
	}
	return hash + uint64(index)
}
