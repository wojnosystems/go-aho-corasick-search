package result

import (
	"github.com/wojnosystems/go-aho-corasick-search"
)

type async struct {
	outputs chan aho_corasick_search.Output
}

func NewAsync(bufferSize int) ReadWriter {
	return &async{
		outputs: make(chan aho_corasick_search.Output, bufferSize),
	}
}

func (r *async) Next() (out aho_corasick_search.Output, ok bool) {
	out, ok = <-r.outputs
	return
}

func (r *async) Emit(out aho_corasick_search.Output) {
	r.outputs <- out
}

func (r *async) Close() error {
	close(r.outputs)
	return nil
}
