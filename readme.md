based on https://golang.org/doc/codewalk/markov/

The main change I made was to make multiple markov chains with different prefix lengths at the same time.  This uses more memory, but I found it gave better results.  A long prefex length tends to repeat lines exactly, and a short prefix length is even more incoherent than markov chains usually are.

If you have a large amount of data, then just using a long prefix will probably be fine and you don't need this code.
