package result

import (
	"github.com/wojnosystems/go-aho-corasick-search"
	"io"
)

type Writer interface {
	Emit(out aho_corasick_search.Output)
	io.Closer
}

type Reader interface {
	Next() (out aho_corasick_search.Output, ok bool)
}

type ReadWriter interface {
	Reader
	Writer
}
