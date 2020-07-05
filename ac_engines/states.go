package ac_engines

type stateIndex int64

const (
	startState   stateIndex = 0
	invalidState stateIndex = -1
)

type states []vertexDense

func newStates() states {
	return make([]vertexDense, 0, 10)
}

func lastStateIndex(states states) stateIndex {
	return stateIndex(len(states))
}
