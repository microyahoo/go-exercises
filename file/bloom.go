package file

import (
	"hash/fnv"
	"math"
)

// BloomFilter is a space-efficient probabilistic data structure for
// approximate set membership testing. False positives are possible but
// false negatives are not: Contains may wrongly return true, but never
// wrongly returns false.
//
// For 600M items at 1% FPR the filter uses ~720 MB of memory.
type BloomFilter struct {
	bits []uint64
	m    uint64 // total number of bits
	k    uint   // number of hash functions
}

// NewBloomFilter creates a bloom filter sized for expectedItems with
// the given false-positive rate (e.g. 0.01 for 1%).
func NewBloomFilter(expectedItems uint64, fpRate float64) *BloomFilter {
	if expectedItems == 0 {
		expectedItems = 1
	}
	if fpRate <= 0 || fpRate >= 1 {
		fpRate = 0.01
	}

	// optimal bit count: m = -n * ln(p) / (ln2)^2
	m := uint64(math.Ceil(-float64(expectedItems) * math.Log(fpRate) / (math.Ln2 * math.Ln2)))
	// optimal hash count: k = (m/n) * ln2
	k := uint(math.Round(float64(m) / float64(expectedItems) * math.Ln2))
	if k < 1 {
		k = 1
	}

	words := (m + 63) / 64
	m = words * 64

	return &BloomFilter{
		bits: make([]uint64, words),
		m:    m,
		k:    k,
	}
}

// Add inserts an element into the filter. NOT safe for concurrent use.
func (bf *BloomFilter) Add(data []byte) {
	h1, h2 := bf.hash(data)
	for i := uint(0); i < bf.k; i++ {
		idx := (h1 + uint64(i)*h2) % bf.m
		bf.bits[idx/64] |= 1 << (idx % 64)
	}
}

// Contains reports whether data is probably in the set.
// Safe for concurrent use after all Add calls are finished.
func (bf *BloomFilter) Contains(data []byte) bool {
	h1, h2 := bf.hash(data)
	for i := uint(0); i < bf.k; i++ {
		idx := (h1 + uint64(i)*h2) % bf.m
		if bf.bits[idx/64]&(1<<(idx%64)) == 0 {
			return false
		}
	}
	return true
}

// hash returns two independent 64-bit hashes derived from FNV-128a.
// These are combined via double-hashing to produce k hash values.
func (bf *BloomFilter) hash(data []byte) (uint64, uint64) {
	h := fnv.New128a()
	h.Write(data)
	s := h.Sum(nil) // 16 bytes

	h1 := uint64(s[0])<<56 | uint64(s[1])<<48 | uint64(s[2])<<40 | uint64(s[3])<<32 |
		uint64(s[4])<<24 | uint64(s[5])<<16 | uint64(s[6])<<8 | uint64(s[7])
	h2 := uint64(s[8])<<56 | uint64(s[9])<<48 | uint64(s[10])<<40 | uint64(s[11])<<32 |
		uint64(s[12])<<24 | uint64(s[13])<<16 | uint64(s[14])<<8 | uint64(s[15])
	if h2 == 0 {
		h2 = 1
	}
	return h1, h2
}

// MemoryUsageMB returns the memory used by the bit array in megabytes.
func (bf *BloomFilter) MemoryUsageMB() float64 {
	return float64(len(bf.bits)*8) / (1024 * 1024)
}

// NumHashFunctions returns k, the number of hash functions used.
func (bf *BloomFilter) NumHashFunctions() uint {
	return bf.k
}
