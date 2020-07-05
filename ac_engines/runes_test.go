package ac_engines

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wojnosystems/go-aho-corasick-search"
	"github.com/wojnosystems/go-aho-corasick-search/result"
	"io/ioutil"
	"os"
	"testing"
)

func TestRunes_Search(t *testing.T) {
	cases := map[string]struct {
		keywords        []string
		input           string
		expectedOutputs []aho_corasick_search.Output
	}{
		"aho-corasick": {
			keywords: []string{
				"he",
				"she",
				"his",
				"hers",
			},
			input: "ushers",
			expectedOutputs: []aho_corasick_search.Output{
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
		"aho-corasick with unicode": {
			keywords: []string{
				"hの",
				"shの",
				"his",
				"hのrs",
			},
			input: "ushのrs",
			expectedOutputs: []aho_corasick_search.Output{
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
		machineBuilder := NewRunes()
		for _, keyword := range c.keywords {
			err := machineBuilder.AddKeyword(keyword)
			require.NoError(t, err, caseName)
		}
		actuals := result.NewSync(10)
		err := machineBuilder.Build().Search(bytes.NewBufferString(c.input), actuals)
		require.NoError(t, err, caseName)
		for _, output := range c.expectedOutputs {
			actual, ok := actuals.Next()
			assert.True(t, ok, caseName)
			assert.Equal(t, output.KeywordIndex, actual.KeywordIndex, caseName)
		}
	}
}

// Test that the search algorithm ends properly when a stream is closed
func TestRunes_SearchClosed(t *testing.T) {
	machineBuilder := NewRunes()
	for _, s := range []string{
		"he",
		"she",
		"his",
		"hers",
	} {
		err := machineBuilder.AddKeyword(s)
		require.NoError(t, err)
	}
	tmp, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	_ = tmp.Close()
	defer func() { _ = os.Remove(tmp.Name()) }()
	err = machineBuilder.Build().Search(tmp, result.NewSync(10))
	require.NoError(t, err)
}

func BenchmarkRunes_Search(b *testing.B) {
	machineBuilder := NewRunes()
	for _, s := range []string{
		"he",
		"she",
		"his",
		"hers",
	} {
		err := machineBuilder.AddKeyword(s)
		if err != nil {
			b.Error(err)
			return
		}
	}

	machine := machineBuilder.Build()

	inputBuf := bytes.Buffer{}
	for i := 0; i < 10000; i++ {
		inputBuf.WriteString("ushers")
	}
	input := bufio.NewReader(&inputBuf)

	for i := 0; i < b.N; i++ {
		_ = machine.Search(input, &dummyResult{})
	}
}
