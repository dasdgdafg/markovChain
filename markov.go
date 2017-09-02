// This code is based on https://golang.org/doc/codewalk/markov/
// // Copyright 2011 The Go Authors. All rights reserved.
// // Use of this source code is governed by a BSD-style
// // license that can be found in the LICENSE file.
// The relevant license can be found at https://golang.org/LICENSE

package markovChain

import (
	"bufio"
	"io"
	"math/rand"
	"strings"
)

// used to indicate that the line should end
const ENDLINE = "____ENDLINE"

// Prefix is a Markov chain prefix of one or more words.
type Prefix []string

// String returns the Prefix as a string (for use as a map key).
func (p Prefix) String() string {
	return strings.Join(p, " ")
}

// Shift removes the first word from the Prefix and appends the given word.
func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

// Chain contains a map ("chain") of prefixes to a list of suffixes.
// A prefix is a string of prefixLen words joined with spaces.
// A suffix is a single word. A prefix can have multiple suffixes.
type Chain struct {
	chain        map[string][]string
	maxPrefixLen int
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func NewChain(maxPrefixLen int) *Chain {
	return &Chain{make(map[string][]string), maxPrefixLen}
}

// Build reads text from the provided Reader and
// parses it into prefixes and suffixes that are stored in Chain.
func (c *Chain) Build(r io.Reader, avoidR io.Reader) {
	avoid := map[string]struct{}{}
	scanner := bufio.NewScanner(avoidR)
	for scanner.Scan() {
		for _, s := range strings.Fields(scanner.Text()) {
			if len(s) >= 3 {
				avoid[strings.ToLower(s)] = struct{}{}
			}
		}
	}

	scanner = bufio.NewScanner(r)
	for scanner.Scan() {
		ps := make([]Prefix, c.maxPrefixLen)
		for i := 0; i < c.maxPrefixLen; i++ {
			ps[i] = make(Prefix, i+1)
		}
		for _, s := range strings.Fields(scanner.Text()) {
			// add a - to words in the avoid list, ie foo -> f-oo
			lowers := strings.ToLower(s)
			lowers = strings.TrimPrefix(lowers, "<")
			lowers = strings.TrimSuffix(lowers, ">")
			_, exists := avoid[lowers]
			exists2 := false
			if len(lowers) > 1 {
				_, exists2 = avoid[lowers[:len(lowers)-1]]
			}
			if exists || exists2 {
				s = s[:1] + "-" + s[1:]
			}
			for i := 0; i < c.maxPrefixLen; i++ {
				key := ps[i].String()
				c.chain[key] = append(c.chain[key], s)
				ps[i].Shift(s)
			}
		}
		for i := 0; i < c.maxPrefixLen; i++ {
			key := ps[i].String()
			c.chain[key] = append(c.chain[key], ENDLINE)
		}
	}
}

// Generate returns a string of at most n words generated from Chain.
func (c *Chain) Generate(n int) string {
	ps := make([]Prefix, c.maxPrefixLen)
	for i := 0; i < c.maxPrefixLen; i++ {
		ps[i] = make(Prefix, i+1)
	}
	var words []string
	for i := 0; i < n; i++ {
		var choices []string
		for j := 0; j < c.maxPrefixLen; j++ {
			choices = append(choices, c.chain[ps[j].String()]...)
		}
		if len(choices) == 0 {
			break
		}
		//fmt.Println("considering", choices)
		next := choices[rand.Intn(len(choices))]
		if next == ENDLINE { // make ending lines less frequent
			next = choices[rand.Intn(len(choices))]
		}
		words = append(words, next)
		for j := 0; j < c.maxPrefixLen; j++ {
			ps[j].Shift(next)
		}
	}
	return strings.Replace(strings.Join(words, " "), ENDLINE, "", -1)
}
