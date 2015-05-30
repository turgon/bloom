package bloom

import (
	"github.com/spaolacci/murmur3"
)

type Bloom struct {
	k          uint
	m          uint64
	collection []uint64
}

func NewBloom(m uint64, k uint) Bloom {

	var length uint64 = m / 64
	if m%64 > 0 {
		length += 1
	}

	b := Bloom{
		k:          k,
		m:          m,
		collection: make([]uint64, length, length),
	}

	return b
}

func singlemask(pos uint64) uint64 {
	var offset uint = uint(pos % 64)
	var t uint64

	t = 1 << offset

	return t
}

func (b *Bloom) SetBit(pos uint64) {

	if pos >= b.m {
		return
	}

	var block uint64 = pos / 64

	t := singlemask(pos)

	b.collection[block] |= t
}

func (b *Bloom) HasBit(pos uint64) bool {
	if pos >= b.m {
		return false
	}

	var block uint64 = pos / 64
	t := singlemask(pos)
	return b.collection[block] & t > 0	
}

func hash(x []byte) (uint64, uint64) {
	var h murmur3.Hash128 = murmur3.New128()
	h.Write([]byte(x))
	v1, v2 := h.Sum128()
	return v1, v2
}

func (b *Bloom) positions (x []byte) []uint64 {
	v1, v2 := hash(x)

	poss := make([]uint64, b.k)

	var i uint
	for i = 0; i < b.k; i++ {
		poss[i] = ((v1 + uint64(i) * v2) % b.m)
	}
	return poss
}

func (b *Bloom) Add(x []byte) {
	positions := b.positions(x)
	for _, v := range positions {
		b.SetBit(v)
	}
}

func (b *Bloom) Has(x []byte) bool {
	positions := b.positions(x)
	for _, v := range positions {
		if !b.HasBit(v) {
			return false
		}
	}
	return true
}

