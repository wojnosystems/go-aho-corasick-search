package ac_engines

type vertexDense struct {
	vertex
	edges []stateIndex
}

func newVertexDense(numberOfStates int) vertexDense {
	v := vertexDense{
		vertex: newVertex(),
		edges:  make([]stateIndex, numberOfStates),
	}
	// By default, all states are invalid
	for i := range v.edges {
		v.edges[i] = invalidState
	}
	return v
}

func (v *vertexDense) nextState(edge int64) stateIndex {
	return v.edges[edge]
}

func (v *vertexDense) setNextState(edge int64, si stateIndex) {
	v.edges[edge] = si
}

func (v *vertexDense) setInvalidEdgesTo(si stateIndex) {
	for stateIndex, state := range v.edges {
		if state == invalidState {
			v.edges[stateIndex] = startState
		}
	}
}
