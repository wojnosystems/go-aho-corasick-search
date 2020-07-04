package aho_corasick

import "io"

type ResultWriter interface {
	Emit(out Output)
	io.Closer
}

type ResultReader interface {
	Next() (out Output, ok bool)
}

type ResultReadWriter interface {
	ResultReader
	ResultWriter
}
