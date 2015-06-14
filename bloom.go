package bloom

import (
	"github.com/spaolacci/murmur3"
	"math"
)

// Bloom is a structure that represents a bloom filter.
// It contains:
// 	k: the number of hashes
//	m: the size of the filter in bits
//	collection: the internal m bits of the filter
//	hasher: a function that returns the bit positions for some input.
type Bloom struct {
	k          uint
	m          uint64
	collection []uint64
	hasher     BloomHasher
}

type BloomHasher func(b *Bloom, value []byte) []uint64

// NewBloom returns a properly constructed Bloom structure given a filter
// size of m and k hashes. The default hasher is Murmur3, but a custom
// hasher can be built that implements BloomHasher and overrides the filter's
// hashing function.
func NewBloom(m uint64, k uint) Bloom {

	var length uint64 = m / 64
	if m%64 > 0 {
		length += 1
	}

	b := Bloom{
		k:          k,
		m:          m,
		collection: make([]uint64, length),
		hasher:     hasher,
	}

	return b
}

// A helper function that hashes input using murmur3 and returns the
// uint64 pair.
func murmur3_128(value []byte) (uint64, uint64) {
	var h murmur3.Hash128 = murmur3.New128()
	h.Write(value)
	return h.Sum128()
}

// hasher implements BloomHasher on a Bloom structure, taking input and
// returning a list of bit positions for the filter.
var hasher = func(b *Bloom, value []byte) []uint64 {
	v1, v2 := murmur3_128(value)

	positions := make([]uint64, b.k)
	for i := range positions {
		positions[i] = ((v1 + uint64(i)*v2) % b.m)
	}
	return positions
}

// Insert takes a byte slice and adds it to the bloom filter. After this is
// called, the filter's Test method will return True for the same input.
func (b *Bloom) Insert(value []byte) {
	positions := b.hasher(b, value)

	for _, pos := range positions {
		var slot = pos / 64
		var offset = uint((pos - slot*64) % 64)

		b.collection[slot] |= (1 << offset)
	}
}

// Test takes a byte slice and checks to see if is in the filter. If it returns
// true, then the input is probably a member of the set the filter reresents.
// If it returns false, then the  input is definitely not a member of the set.
func (b *Bloom) Test(value []byte) bool {
	positions := b.hasher(b, value)

	for _, pos := range positions {
		var slot = pos / 64
		var offset = uint((pos - slot*64) % 64)

		if (b.collection[slot] & (1 << offset)) == 0 {
			return false
		}
	}
	return true
}

// EstimateFalsePositives takes filter parameters and the number of items
// you expect to insert (numItems) and returns an estimation of the false
// positive probability.
//
// This function is useful when the filter size and hashes are predetermined
// and you need to know how the filter will behave as you insert items.
func EstimateFalsePositives(k uint, m uint64, numItems uint64) float64 {
	exponent := -1.0 * float64(k) * float64(numItems) / float64(m)
	return math.Pow((1.0 - math.Exp(exponent)), float64(k))
}

// This version of EstimateFalsePositives is bound to an instance of a Bloom
// structure and will use its k and m parameters to build the estimate.
func (b *Bloom) EstimateFalsePositives(numItems uint64) float64 {
	return EstimateFalsePositives(b.k, b.m, numItems)
}

// OptimalHashNumber takes the size of a bloom filter and the number of items
// you expect to insert and computes the number of hashes that minimizes the
// false positive probability.
func OptimalHashNumber(m uint64, numItems uint64) uint {

	// Unfortunately, we need a whole number of hashes.
	// Pick the better of floor or ceiling.
	best := uint((float64(m) / float64(numItems)) * math.Log(2.0))
	bestEstimate := EstimateFalsePositives(best, m, numItems)
	nextBestEstimate := EstimateFalsePositives(best + 1, m, numItems)
	if bestEstimate <= nextBestEstimate {
		return best
	}
	return best+1
}

// OptimalFilterSize takes the number of items you expect to insert into a
// bloom filter and the maximum false positive probability you are willing
// to tolerate, and computes the minimum number of bits the filter requires.
func OptimalFilterSize(numItems uint64, maxFalseProbability float64) uint64 {
	numer := float64(numItems) * math.Log(maxFalseProbability)
	denom := math.Pow(math.Log(2.0), 2.0)
	return uint64(math.Ceil(-1.0 * numer / denom))
}
