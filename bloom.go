package bloom

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
