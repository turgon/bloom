package bloom

import (
	"fmt"
	"testing"
)

func TestNewBloom(t *testing.T) {
	s := NewBloom(64, 1)
	if s.k != 1 {
		t.Errorf("NewBloom returned wrong hashes for filter")
	}
	if s.m != 64 {
		t.Errorf("NewBloom returned wrong filter size")
	}
	if cap(s.collection) != 1 {
		t.Errorf("NewBloom built wrong sized collection of uint64s")
	}
}

func TestJustTest(t *testing.T) {
	b := NewBloom(128, 1)
	if b.Test([]byte("test")) {
		t.Errorf("Insert or Test broken")
	}
}

func TestInsertAndTest(t *testing.T) {
	b := NewBloom(128, 1)
	b.hasher = func (b *Bloom, value []byte) ([]uint64) {
		return make([]uint64, 1)
	}
	b.Insert([]byte("test"))
	if !b.Test([]byte("test")) {
		t.Errorf("Insert or Test broken")
	}
}

func TestAvalanche(t *testing.T) {
	rawWord := "testing 123 testing test test yes"
	word := []byte(rawWord)

	a1, a2 := murmur3_128(word)

	totalBits := 0
	matchBits := 0

	for k := range word {
		for i := 0; i < 8; i++ {
			wordx := []byte(rawWord)
			wordx[k] ^= (1 << uint(i))

			// now word and wordx have one bit set differently

			b1, b2 := murmur3_128(wordx)

			c1 := a1 ^ b1
			c2 := a2 ^ b2

			for m := 0; m < 64; m++ {
				totalBits += 1
				matchBits += int(c1 & 1)
				c1 >>= 1
			}
			for m := 0; m < 64; m++ {
				totalBits += 1
				matchBits += int(c2 & 1)
				c2 >>= 1
			}
		}
	}
	fmt.Println("total", matchBits, totalBits, 100.0 * float64(matchBits) / float64(totalBits))
}



