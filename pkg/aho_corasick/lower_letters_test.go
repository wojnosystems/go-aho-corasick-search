package aho_corasick

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSearch(t *testing.T) {
	cases := map[string]struct {
		keywords        []string
		input           string
		expectedOutputs []Output
	}{
		"aho-corasick": {
			keywords: []string{"he", "she", "his", "hers"},
			input:    "ushers",
			expectedOutputs: []Output{
				{
					keywordIndex: 1,
				},
				{
					keywordIndex: 0,
				},
				{
					keywordIndex: 3,
				},
			},
		},
	}

	for caseName, c := range cases {
		machine := New(c.keywords, 26)
		actual := machine.Search(c.input)
		for i, output := range c.expectedOutputs {
			assert.Equal(t, output.keywordIndex, actual[i].keywordIndex, caseName)
		}
	}
}
