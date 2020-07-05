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
