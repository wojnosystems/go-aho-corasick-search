package aho_corasick

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func TestSearch(t *testing.T) {
	cases := map[string]struct {
		keywords        []string
		input           string
		expectedOutputs []Output
	}{
		"aho-corasick": {
			keywords: []string{
				"he",
				"she",
				"his",
				"hers",
			},
			input: "ushers",
			expectedOutputs: []Output{
				{
					KeywordIndex: 1,
				},
				{
					KeywordIndex: 0,
				},
				{
					KeywordIndex: 3,
				},
			},
		},
	}

	for caseName, c := range cases {
		machine, err := NewLowerLetters(c.keywords)
		require.NoError(t, err, caseName)
		actuals := NewAsyncResults(10)
		err = machine.Search(bytes.NewBufferString(c.input), actuals)
		require.NoError(t, err, caseName)
		for _, output := range c.expectedOutputs {
			actual, ok := actuals.Next()
			assert.True(t, ok, caseName)
			assert.Equal(t, output.KeywordIndex, actual.KeywordIndex, caseName)
		}
	}
}

// Test that the search algorithm ends properly when a stream is closed
func TestLowerLetters_SearchClosed(t *testing.T) {
	machine, err := NewLowerLetters([]string{
		"he",
		"she",
		"his",
		"hers",
	})
	require.NoError(t, err)
	tmp, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	_ = tmp.Close()
	defer func() { _ = os.Remove(tmp.Name()) }()
	err = machine.Search(tmp, NewAsyncResults(10))
	require.NoError(t, err)
}
