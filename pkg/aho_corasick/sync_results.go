package aho_corasick

// syncResults
// A result bucket built without Search running in a go routine
// This is not thread-safe and is intended for small, test jobs.
type syncResults struct {
	outputs      []Output
	currentIndex int
}

func NewSyncResults(bufferSize int) ResultReadWriter {
	return &syncResults{
		outputs: make([]Output, 0, bufferSize),
	}
}

func (r *syncResults) Next() (out Output, ok bool) {
	if r.currentIndex >= len(r.outputs) {
		return out, false
	}
	out = r.outputs[r.currentIndex]
	r.currentIndex++
	return out, true
}

func (r *syncResults) Emit(out Output) {
	r.outputs = append(r.outputs, out)
}

func (r *syncResults) Close() error {
	return nil
}
