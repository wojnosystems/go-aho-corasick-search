# Overview

Implements the Aho-Corasick algorithm as described in
 ["Efficient String Matching: An Aid to Bibliographic Search"](https://cr.yp.to/bib/1975/aho.pdf) (1975) 
 by Alfred V. Aho and Margaret J. Corasick at Bell Laboratories.

# How to use it

cd to the root of your go project (which is presumed to be using Go11 vendored modules), and execute the following:

```bash
go get github.com/wojnosystems/go-aho-corasick-search
```

Here's how you use it in your application.

```go
package main

import (
	"bytes"
	"fmt"
	"github.com/wojnosystems/go-aho-corasick-search/ac_engines"
	"github.com/wojnosystems/go-aho-corasick-search/result"
	"unicode/utf8"
)

func main() {
	stringsToFind := []string{
		"hの",
		"shの",
		"his",
		"hのrs",
	}
	// Create a sample input within which to search for the keywords
	input := bytes.NewBufferString("ushのrs")

	// Create a new state machine builder
	stateMachineBuilder := ac_engines.NewRunes()
	for _, s := range stringsToFind {
		// Add in everything you want to look for
		_ = stateMachineBuilder.AddKeyword(s)
	}
	// Create a final state machine, ready for searching
	stateMachine := stateMachineBuilder.Build()

	// Prepare a place for the search engine to put the results
	results := result.NewSync(10)
    _ = stateMachine.Search(input, results)
	for {
		// Get all of the results
		match, ok := results.Next()
		// ok is only true if match was valid. If false, that means there are no more matches
		if !ok {
			break
		}
		word := stringsToFind[match.KeywordIndex]
		fmt.Printf("Match! %s @ b:%d c:%d\n",
			word,
			match.ByteOffset-uint64(len(word)),
			match.CharacterOffset-uint64(utf8.RuneCountInString(word)))
	}
}
```

This will output:

```
Match! she @ b:0 c:0
Match! he @ b:1 c:1
Match! hers @ b:1 c:1
```

See `cmd/rune-example/main.go` for the above example working. The b is the byte offset where the word started, 
the c is the character offset (important as runes don't always correspond to bytes read).

## Multi-threaded approach

If your input source is streaming and may block while waiting for things to be read, such as over a network or 
a slow file system, the implementations support channels for thread-safety. You simply need to use them similarly 
to the below:

```go
package main

import (
	"bytes"
	"fmt"
	"github.com/wojnosystems/go-aho-corasick-search/ac_engines"
	"github.com/wojnosystems/go-aho-corasick-search/result"
	"unicode/utf8"
)

func main() {
	stringsToFind := []string{
		"hの",
		"shの",
		"his",
		"hのrs",
	}
	// Create a sample input within which to search for the keywords
	input := bytes.NewBufferString("ushのrs")

	// Create a new state machine builder
	stateMachineBuilder := ac_engines.NewRunes()
	for _, s := range stringsToFind {
		// Add in everything you want to look for
		_ = stateMachineBuilder.AddKeyword(s)
	}
	// Create a final state machine, ready for searching
	stateMachine := stateMachineBuilder.Build()

	// Prepare a place for the search engine to put the results
	results := result.NewAsync(10)
	go func() {
		// Perform the search, use a Go-Routine in case your input comes from a buffered source, like a file or socket
		// You don't have to use a Go-Routine, you can parse these results without it.
		// In that case, either ensure that your bufferSize is set to the maximum number of matches (which you'd have to 
		// know before hand) or use the result.NewSync buffer, which is just an ever-expanding slice of values.
		// When using result.NewAsync, the buffer is a channel, which means Search will stop processing once the channel
		// is filled with results. Another Go-Routine, like we have below, will need to read the results at that point.
		_ = stateMachine.Search(input, results)
	}()
	for {
		// Get all of the results
		match, ok := results.Next()
		// ok is only true if match was valid. If false, that means there are no more matches
		if !ok {
			break
		}
		word := stringsToFind[match.KeywordIndex]
		fmt.Printf("Match! %s @ b:%d c:%d\n",
			word,
			match.ByteOffset-uint64(len(word)),
			match.CharacterOffset-uint64(utf8.RuneCountInString(word)))
	}
}
```

This will output the same results as the above, but will be able to handle streams of data you may not wish to load 
into memory before using.

The only difference is that the results load asynchronously into results via a channel instead of being crammed into 
a slice.

# A bit about the algorithm

This algorithm is useful for searching for multiple strings known _a priori_, within a large string, in near-linear 
time. The running time for this algorithm is O(n*z) where n is the number of characters to read during the search 
while `z = SUM(len(keyword[i]))` the sum of the lengths of all keywords. As long as your key words are both short 
and there aren't too many of them, this algorithm runs in linear time.

## Where I departed from the algorithm

This implementation does not eliminate the redundant failure transitions as described in Section 6. However, this does 
implement this algorithm using a stream processor and GoLang's pattern of emitting data through a channel, thus 
allowing Search to be run in a Go Routine, while the output processing code is able to process in the main routine, 
pausing only while matches are still being found. This reduces the amount of time waiting for I/O.

## ASCII-only Processing

The implementation "LowerLetters" only supports the ASCII lower-case English letters: a - z for simplicity (and because 
the HackerRank question only dealt with this range). This implementation using a vertexDense, where in each set of next 
states is a simple slice of states to transition to next or -1 (invalidState) if it is unset.

The Runes implementation improves upon this simple design by replacing character bytes with runes (int32). Instead of 
using flat arrays, which would take up a large amount of memory because to the search space size for all runes (2^32), 
this implementation using maps in the vertexSparse struct. However, maps are not free. The memory complexity of this 
implementation would grow much larger and run a bit more slowly.

## Reusable State Machine

Once you've created the state machine, you can run multiple searches on it using many different inputs. Search does not 
alter the state machine in any way and merely runs the search through it. That means, once built (call to 
NewLowerLetters or NewRunes), the state machine is thread-safe and running multiple Searches at the same time is 
perfectly OK.

# Future Work

 * Eliminating the redundant Fail State transitions
 * Augment output to include line count, column position, and character count

# Inspiration

I was tooling around with a HackerRank problem regarding 
["Determining DNA Health"](https://www.hackerrank.com/challenges/determining-dna-health/problem) and after stumbling 
around with my own tries, and inventing a convoluted version of this algorithm, I figured I'd actually try it with the 
real thing.

The part I got hung up on was proving that my failure recovery (which was actually a BFS, just as with Aho-Corasick), 
actually represented the proper prefix of the failure matches. Aho-Corasick converts this to suffix propers, which 
makes things a bit easier to build a state machine with.

I also wasn't combining my outputs, so I would not detect the prefix matches, only the end-states.
