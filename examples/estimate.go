package main

import (
	"fmt"
	"github.com/turgon/bloom"
)

func main() {
	// I have 97 items to put into a filter and I can't tolerate
	// more than 3% false positive probability.
	// How big should my filter be?

	bestSize := bloom.OptimalFilterSize(97, 0.03)

	// Now that I know how big the filter should be, how many 
	// hashes should I use to minimize the false positive probability?
	bestHashes := bloom.OptimalHashNumber(bestSize, 97)

	// Great! But, what is the estimated false positive probability
	// given these inputs? 
	estimate := bloom.EstimateFalsePositives(bestHashes, bestSize, 97)

	fmt.Printf("Given 97 items and a 3%% false positive bound, we need a filter with %v bits and %v hashes. It hash a false positive probability of %4.4f%%\n", bestSize, bestHashes, 100.0*estimate)
}
