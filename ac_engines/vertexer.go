package ac_engines

type vertexer interface {
	nextState(edge int64) stateIndex
	setNextState(edge int64, si stateIndex)
	setInvalidEdgesTo(si stateIndex)
}
