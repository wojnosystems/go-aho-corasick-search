package ac_engines

type denseStates []vertexDense

func newDenseStates() denseStates {
	return make([]vertexDense, 0, 10)
}

func (s denseStates) lastStateIndex() stateIndex {
	return stateIndex(len(s))
}

const (
	startState   stateIndex = 0
	invalidState stateIndex = -1
)
