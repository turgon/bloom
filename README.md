# bloom

A Bloom filter implementation in Golang.

Pull requests and feedback are welcome!

# What is a Bloom filter?

A bloom filter is a data structure that can be used to test set membership. It's probabilistic, so it's not always right, BUT it's always right when it comes to false negatives. That is, when it tells you something isn't a set member, it's telling you the truth.

# How does it work?

Let's suppose we want to create a tiny bloom filter of just 64 bits all set to 0. Each bit has a positional address from zero to sixty-three. When we want to add an item to the filter, we can use a hashing function that produces output in the range 0-63, hash the input, and treat the output as the address of the bit we will set to 1. In this simple example, we will only set one bit, but generally we would want to hash mutliple times. Once the bit is set, the filter consists of sixty-three 0 bits and one 1 bit in a seemingly random position.

Now, to test whether an item is in the set, we follow the same procedure of hashing the item. This time, instead of setting the bit in that position to 1, we simply test that it is set to 1 already. If not, then this item cannot be in the set -- if it was, that bit would be 1. If the bit we test is 1, then this item is probably in the set.

# Ok, but how does it REALLY work?

Typically the bitfield will be somewhat larger, and the number of rounds of hashing will be more than one. There's a formula you can use that will tell you how many bits you need based on the upper bound of your set size and the probability you want to target for false positives. There's also a formula you can use to compute th probability given the other variables. Lots of formulas. Lots of good times.


