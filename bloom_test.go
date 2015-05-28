package bloom

import (
	"fmt"
	"testing"
)

func TestNewBloom(t *testing.T) {
	s := NewBloom(1, 1)
	fmt.Println(s)
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
	var err error

	s := NewBloom(1, 1)
	err = s.SetBit(0)
	if err != nil {
		t.Errorf("Can't set expected bit!")
	}

	err = s.SetBit(1)
	if err == nil {
		t.Errorf("Set unexpected bit!")
	}

	b := NewBloom(65, 1)
	err = b.SetBit(64)
	if err != nil {
		t.Errorf("Can't set expected bit!")
	}

	err = b.SetBit(65)
	if err == nil {
		t.Errorf("Set unexpected bit!")
	}
	fmt.Println(b)

}

func BenchmarkSetBit(b *testing.B) {
	s := NewBloom(uint64(b.N), 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.SetBit(uint64(i))
	}
}
