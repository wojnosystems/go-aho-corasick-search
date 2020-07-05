package ac_engines

type sparseStates []vertexSparse

func newSparseStates() sparseStates {
	return make([]vertexSparse, 0, 10)
}

func (s sparseStates) lastStateIndex() stateIndex {
	return stateIndex(len(s))
}
