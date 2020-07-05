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
	stateMachine, _ := ac_engines.NewRunes(stringsToFind)
	results := result.NewAsync(10)

	input := bytes.NewBufferString("ushers")
	go func() {
		_ = stateMachine.Search(input, results)
	}()
	for {
		match, ok := results.Next()
		if !ok {
			break
		}
		fmt.Printf("Match! %s\n", stringsToFind[match.KeywordIndex])
	}
}
