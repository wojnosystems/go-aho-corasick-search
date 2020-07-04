# Overview

Implements the Aho-Corasick algorithm as described in ["Efficient String Matching: An Aid to Bibliographic Search"](https://cr.yp.to/bib/1975/aho.pdf) (1975) by Alfred V. Aho and Margaret J. Corasick at Bell Laboratories.

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
)

func main() {
  stringsToFind := []string{
    "he",
    "she",
    "his",
    "hers",
  }  
  stateMachine, _ := ac_engines.NewLowerLetters(stringsToFind)
  results := result.NewAsync(10)

  input := bytes.NewBufferString("ushers")
  go func() {
    _ = stateMachine.Search( input, results)
  }()
  for {
    match, ok := results.Next()
    if !ok {
      break
    }
    fmt.Printf("Match! %s\n", stringsToFind[match.KeywordIndex])
  }
}
```

This will output:

```
Match! she
Match! he
Match! hers
```

See `cmd/example/main.go` for the above example working.

# A bit about the algorithm

This algorithm is useful for searching for multiple strings known _a priori_, within a large string, in near-linear time. The running time for this algorithm is O(n*z) where n is the number of characters to read during the search while `z = SUM(len(keyword[i]))` the sum of the lengths of all keywords. As long as your key words are both short and there aren't too many of them, this algorithm runs in linear time.

## Where I departed from the algorithm

This implementation does not eliminate the redundant failure transitions as described in Section 6. However, this does implement this algorithm using a stream processor and GoLang's pattern of emitting data through a channel, thus allowing Search to be run in a Go Routine, while the output processing code is able to process in the main routine, pausing only while matches are still being found. This reduces the amount of time waiting for I/O.

This implementation only supports the ASCII lower-case English letters: a - z for simplicity (and because the HackerRank question only dealt with this range). This could be, somewhat easily, adjusted to support full rune-scanning by replacing the slices of states with maps of states. However, maps are not free. The memory complexity of this implementation would grow much larger and run a bit more slowly.

## Reusable State Machine

Once you've created the state machine, you can run multiple searches on it using many different inputs. Search does not alter the state machine in any way and merely runs the search through it. That means, once built, the state machine is thread-safe.

# Future Work

 * Eliminating the redundant Fail State transitions
 * Supporting runes using maps
 * Augment output to include line count, column position, and character count

# Inspiration

I was tooling around with a HackerRank problem regarding ["Determining DNA Health"](https://www.hackerrank.com/challenges/determining-dna-health/problem) and after stumbling around with my own tries, and inventing a convoluted version of this algorithm, I figured I'd actually try it with the real thing.

The part I got hung up on was proving that my failure recovery (which was actually a BFS, just as with Aho-Corasick), actually represented the proper prefix of the failure matches. Aho-Corasick converts this to suffix propers, which makes things a bit easier to build a state machine with.

I also wasn't combining my outputs, so I would not detect the prefix matches, only the end-states.
