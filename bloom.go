package bloom

import (
	"github.com/spaolacci/murmur3"
	"math"
)

type Bloom struct {
	k          uint
	m          uint64
	collection []uint64
	hasher     BloomHasher
}

type BloomHasher func(b *Bloom, value []byte) []uint64

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

func murmur3_128(value []byte) (uint64, uint64) {
	var h murmur3.Hash128 = murmur3.New128()
	h.Write(value)
	return h.Sum128()
}

var hasher = func(b *Bloom, value []byte) []uint64 {
	v1, v2 := murmur3_128(value)

	positions := make([]uint64, b.k)
	for i := range positions {
		positions[i] = ((v1 + uint64(i)*v2) % b.m)
	}
	return positions
}

func (b *Bloom) Insert(value []byte) {
	positions := b.hasher(b, value)

	for _, pos := range positions {
		var slot = pos / 64
		var offset = uint((pos - slot*64) % 64)

		b.collection[slot] |= (1 << offset)
	}
}

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

func EstimateFalsePositives(k uint, m uint64, numItems uint64) float64 {
	exponent := -1.0 * float64(k) * float64(numItems) / float64(m)
	return math.Pow((1.0 - math.Exp(exponent)), float64(k))
}

func (b *Bloom) EstimateFalsePositives(numItems uint64) float64 {
	return EstimateFalsePositives(b.k, b.m, numItems)
}
