package bloom

import (
	//"fmt"
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
