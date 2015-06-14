# bloom

A Bloom filter implementation in Golang.

Documentation can be found [on godocs.org][5]

Pull requests and feedback are welcome!

# What is a Bloom filter?

A bloom filter is a data structure that can be used to test set membership. It's probabilistic, so it's not always right, BUT it's always right when it comes to false negatives. That is, when it tells you something isn't a set member, it's telling you the truth.

The [Wikipedia article][1] has a nice overall description of Bloom Filters. I referred to Bloom's original paper, ["Space/Time Trade-offs in Hash Coding with Allowable Errors"][2] as well as ["Building a Better Bloom Filter"][3] by Kirsch and Mitzenmacher. As well, this [fantastic blog post by Jonathan Ellis][4] inspired me to make sure I selected a hashing function with good avalanche characteristics, which is why after a few hours of research I selected Murmur3.

# How does it work?

Let's suppose we want to create a tiny bloom filter of just 64 bits all set to 0. Each bit has a positional address from zero to sixty-three. When we want to add an item to the filter, we can use a hashing function that produces output in the range 0-63, hash the input, and treat the output as the address of the bit we will set to 1. In this simple example, we will only set one bit, but generally we would want to hash mutliple times. Once the bit is set, the filter consists of sixty-three 0 bits and one 1 bit in a seemingly random position.

Now, to test whether an item is in the set, we follow the same procedure of hashing the item. This time, instead of setting the bit in that position to 1, we simply test that it is set to 1 already. If not, then this item cannot be in the set -- if it was, that bit would be 1. If the bit we test is 1, then this item is probably in the set.

It's worth noting that any other input value has a 1/64 chance of being hashed into the same bit as our original item. This is precisely how false positives come to be!


[1]: http://en.wikipedia.org/wiki/Bloom_filter
[2]: http://astrometry.net/svn/trunk/documents/papers/dstn-review/papers/bloom1970.pdf
[3]: http://www.eecs.harvard.edu/~kirsch/pubs/bbbf/esa06.pdf
[4]: http://spyced.blogspot.com/2009/01/all-you-ever-wanted-to-know-about.html
[5]: https://godoc.org/github.com/turgon/bloom
