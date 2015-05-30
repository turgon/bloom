package bloom

import (
	//"fmt"
	"testing"
)

func TestNewBloom(t *testing.T) {
	s := NewBloom(1, 1)
	if s.k != 1 {
		t.Errorf("NewBloom returned wrong hashes for filter")
	}
	if s.m != 1 {
		t.Errorf("NewBloom returned wrong filter size")
	}
	if cap(s.collection) != 1 {
		t.Errorf("NewBloom built wrong sized collection of uint64s")
	}
}

func TestSetBit(t *testing.T) {
	s := NewBloom(1, 1)
	s.SetBit(0)
	if !s.HasBit(0) {
		t.Errorf("Can't set expected bit!")
	}

	s.SetBit(1)
	if s.HasBit(1) {
		t.Errorf("Set unexpected bit!")
	}

	b := NewBloom(65, 1)
	b.SetBit(64)
	if !s.HasBit(64) {
		t.Errorf("Can't set expected bit!")
	}

	b.SetBit(65)
	if s.HasBit(65) {
		t.Errorf("Set unexpected bit!")
	}

}

func TestHash(t *testing.T) {
	x,y := hash([]byte("test"))
	// This is a bad test; I just printed the output and copied it here.
	if x != 12429135405209477533 || y != 11102079182576635266 {
		t.Errorf("murmur3 is broken!")
	}
}

/*
func TestPositions(t *testing.T) {
	b := NewBloom(16, 4)
	for i, v := range b.positions([]byte("test")) {
	}
}
*/

func TestAdd(t *testing.T) {
	s := NewBloom(1, 1)
	s.Add([]byte("test"))
}

func TestHas(t *testing.T) {
	s := NewBloom(1, 100)
	if s.Has([]byte("test")) {
		t.Errorf("empty filter hash a member!")
	}
	s.Add([]byte("test"))
	if !s.Has([]byte("test")) {
		t.Errorf("filter missing a member!")
	}
}

func BenchmarkSetBit(b *testing.B) {
	s := NewBloom(uint64(b.N), 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.SetBit(uint64(i))
	}
}
