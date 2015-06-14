package bloom

import (
	"math/rand"
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
	b.hasher = func(b *Bloom, value []byte) []uint64 {
		return make([]uint64, 1)
	}
	b.Insert([]byte("test"))
	if !b.Test([]byte("test")) {
		t.Errorf("Insert or Test broken")
	}
}

func randByteSlice(len int) []byte {
	rbs := make([]byte, len)

	for i := 0; i < len; i++ {
		rbs[i] = byte(rand.Int())
	}
	return rbs
}

func TestAvalanche(t *testing.T) {
	wordLen := 32
	pairs := 1000000

	// let's build pairs comprised of a random hash key and its copy with a single bit twiddled!

	inputs := make([][]byte, pairs)
	tweaked := make([][]byte, pairs)
	inputsc := make([][]uint64, pairs)
	tweakedc := make([][]uint64, pairs)
	avalanche := make([][]uint64, pairs)

	for i := 0; i < pairs; i++ {
		inputs[i] = make([]byte, wordLen)
		tweaked[i] = make([]byte, wordLen)
		inputsc[i] = make([]uint64, 2)
		tweakedc[i] = make([]uint64, 2)
		avalanche[i] = make([]uint64, 2)
		for j := 0; j < wordLen; j++ {
			inputs[i][j] = byte(rand.Int())
			tweaked[i][j] = inputs[i][j]
		}
		tweaked[i][rand.Intn(wordLen)] ^= (1 << uint(rand.Intn(8)))
		inputsc[i][0], inputsc[i][1] = murmur3_128(inputs[i])
		tweakedc[i][0], tweakedc[i][1] = murmur3_128(tweaked[i])

		avalanche[i][0] = inputsc[i][0] ^ tweakedc[i][0]
		avalanche[i][1] = inputsc[i][1] ^ tweakedc[i][1]
	}

	B := make([]uint64, pairs)

	for i := 0; i < pairs; i++ {
		for j := 0; j < 2; j++ {
			c := avalanche[i][j]
			for m := 0; m < 64; m++ {
				B[i] += (c & 1)
				c >>= 1
			}
		}
	}

	// For each experiment/trial, we've counted the number of avalanche bits.
	// Each is equivalent to 128 coin flips, so we can use Fisher's Method
	// and some simple probability to find upper and lower bounds on the
	// number of 1-bit per experiment. Because they're a series of coin
	// flips, we know they ought to be distributed approximately normally,
	// and that the probability of exceeding our chosen threshold is a
	// two-tailed p-value is:
	//   2 * 1/2^128 * sum (128 choose (128-j)), j=0 to 32 ~= 1.28418*10^-8
	// I picked 32 because it has the nice property of being 1/4 of 128,
	// which means any 1-bit or 0-bit count that exceeds half of the total
	// bits implies that we reject the hypothesis.

	var minOnes uint64 = 32
	var maxOnes uint64 = 96

	for i := 0; i < pairs; i++ {
		if B[i] > maxOnes || B[i] < minOnes {
			t.Errorf("Hash failed Avalanche test with %v bits\n", B[i])
		}
	}
}

func TestEstimateFalsePositives(t *testing.T) {
	// If we overload the hash any input should appear to be a set member
	if EstimateFalsePositives(2, 8, 256) != 1.0 {
		t.Errorf("Bad false positive probability estimate")
	}

	// If we've put nothing in the hash, there can't be false positives
	if EstimateFalsePositives(2, 8, 0) != 0.0 {
		t.Errorf("Bad false positive probability estimate")
	}

	// If we put four items into an 8 bit filter, it's possible that with
	// one hash round all four items are assigned distinct bits, however,
	// it isn't expected.
	if EstimateFalsePositives(1, 8, 4) > 0.5 {
		t.Errorf("Bad false positive probability estimate")
	}

	// Likewise, if we put eight items in, it's possible that all of
	// them get assigned to the same bit, but it isn't expected.
	if EstimateFalsePositives(1, 8, 8) < 0.5 {
		t.Errorf("Bad false positive probability estimate")
	}

	// This is a bit contrived, but by way of hand calculations,
	// these estimations should straddle 0.5
	lower := EstimateFalsePositives(1, 64, 44)
	upper := EstimateFalsePositives(1, 64, 45)
	if lower >= 0.5 || upper <= 0.5 {
		t.Errorf("Bad false positive probability estimate")
	}
}

func TestBloomEstimateFalsePositives(t *testing.T) {
	// This is a very simple test of convergence. We ought to be able to
	// build a bloom filter and estimate its false positive probability,
	// then show that many applications of Test using random input
	// forms a hit-rate that approaches the estimate.

	var m uint64 = 65536

	trials := 100000
	hits := 0
	members := m / 3

	b := NewBloom(m, 1)
	for i := 0; uint64(i) < members; i++ {
		b.Insert(randByteSlice(50))
	}

	estimate := b.EstimateFalsePositives(members)

	for i := 0; i < trials; i++ {
		val := randByteSlice(100)
		if b.Test(val) {
			hits++
		}
		// If we've iterated enough and the hit rate is within a percent of the estimate, stop.
		if i > 100 && (((float64(hits)/float64(i))-estimate)/estimate) < 0.01 {
			return
		}
	}
	t.Errorf("Failed to approach the estimate")
}
