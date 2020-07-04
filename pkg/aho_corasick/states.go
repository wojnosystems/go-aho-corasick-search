package aho_corasick

type states []vertex

func newStates() states {
	return make([]vertex, 0, 10)
}
