package result

import (
	"github.com/wojnosystems/go-aho-corasick-search"
)

// sync
// A result bucket built without Search running in a go routine
// This is not thread-safe and is intended for small, test jobs.
type sync struct {
	outputs      []aho_corasick_search.Output
	currentIndex int
}

func NewSync(bufferSize int) ReadWriter {
	return &sync{
		outputs: make([]aho_corasick_search.Output, 0, bufferSize),
	}
}

func (r *sync) Next() (out aho_corasick_search.Output, ok bool) {
	if r.currentIndex >= len(r.outputs) {
		return out, false
	}
	out = r.outputs[r.currentIndex]
	r.currentIndex++
	return out, true
}

func (r *sync) Emit(out aho_corasick_search.Output) {
	r.outputs = append(r.outputs, out)
}

func (r *sync) Close() error {
	return nil
}
