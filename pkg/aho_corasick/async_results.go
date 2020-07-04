package aho_corasick

type asyncResults struct {
	outputs chan Output
}

func NewAsyncResults(bufferSize int) ResultReadWriter {
	return &asyncResults{
		outputs: make(chan Output, bufferSize),
	}
}

func (r *asyncResults) Next() (out Output, ok bool) {
	out, ok = <-r.outputs
	return
}

func (r *asyncResults) Emit(out Output) {
	r.outputs <- out
}

func (r *asyncResults) Close() error {
	close(r.outputs)
	return nil
}
