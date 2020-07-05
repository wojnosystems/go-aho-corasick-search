package ac_engines

type stateIndex int

const (
	startState   stateIndex = 0
	invalidState stateIndex = -1
)

type states []vertex

func newStates() states {
	return make([]vertex, 0, 10)
}

func lastStateIndex(states states) stateIndex {
	return stateIndex(len(states))
}
