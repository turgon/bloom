package main

import (
	"fmt"
	"github.com/turgon/bloom"
)

func main() {
	b := bloom.NewBloom(64, 1)

	b.Insert([]byte("test"))

	if b.Test([]byte("test")) {
		fmt.Println("The test member is in the set!")
	}

	if b.Test([]byte("not test")) {
		fmt.Println("A non-set member is  in the set! Oh no!")
	}
}

