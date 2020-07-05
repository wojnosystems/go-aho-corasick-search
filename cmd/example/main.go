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
		"he",
		"she",
		"his",
		"hers",
	}
	stateMachineBuilder := ac_engines.NewLowerLetters()
	for _, s := range stringsToFind {
		_ = stateMachineBuilder.AddKeyword(s)
	}
	stateMachine := stateMachineBuilder.Build()
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
		word := stringsToFind[match.KeywordIndex]
		fmt.Printf("Match! %s @ b:%d c:%d\n",
			word,
			match.ByteOffset-uint64(len(word)),
			match.CharacterOffset-uint64(utf8.RuneCountInString(word)))
	}
}
