package bloom

import (
	"errors"
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

func (b *Bloom) SetBit(pos uint64) error {
	if pos >= b.m {
		return errors.New("pos exceeds b.m")
	}
	var block uint64 = pos / 64
	var offset uint = uint(pos % 64)

	var t uint64
	t = 1 << offset

	b.collection[block] |= t

	return nil
}
