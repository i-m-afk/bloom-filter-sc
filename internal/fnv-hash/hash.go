/*
* Author : Rishav Kumar
* references: https://datatracker.ietf.org/doc/html/draft-eastlake-fnv-17.html#section-6.1.1
 */
package fnvhash

const FNV_offset_basis uint64 = 0xcbf29ce484222325
const FNV_prime uint64 = 1099511628211

func Fnv1(data []byte) uint64 {
	hash := FNV_offset_basis
	for _, b := range data {
		hash *= FNV_prime
		hash ^= uint64(b)
	}
	return hash
}

func Fnv1a(data []byte) uint64 {
	hash := FNV_offset_basis
	for _, b := range data {
		hash ^= uint64(b)
		hash *= FNV_prime
	}
	return hash
}

func Fnv1b(data []byte) uint64 {
	hash := FNV_offset_basis + 0xDEADBEEF
	for _, b := range data {
		hash *= FNV_prime
		hash ^= uint64(b)
	}
	return hash
}

func Fnv1c(data []byte) uint64 {
	hash := FNV_offset_basis + 0xAAAFFF
	for _, b := range data {
		hash *= FNV_prime
		hash ^= uint64(b)
	}
	return hash
}

func Fnv1d(data []byte) uint64 {
	hash := FNV_offset_basis + 0xFEB290
	for _, b := range data {
		hash *= FNV_prime
		hash ^= uint64(b)
	}
	return hash
}

func Fnv1e(data []byte) uint64 {
	hash := FNV_offset_basis
	hash ^= 0xDEADBEEF
	for _, b := range data {
		hash *= FNV_prime
		hash ^= uint64(b)
	}
	return hash
}

func Fnv1f(data []byte) uint64 {
	hash := FNV_offset_basis
	hash ^= 0xABCDEF
	for _, b := range data {
		hash *= FNV_prime
		hash ^= uint64(b)
	}
	return hash
}

func Fnv1g(data []byte) uint64 {
	hash := FNV_offset_basis
	for _, b := range data {
		hash ^= uint64(b)
		hash *= FNV_prime
	}
	hash += 0xABCDEF
	return hash
}
