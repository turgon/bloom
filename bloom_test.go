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
