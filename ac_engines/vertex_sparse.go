package ac_engines

type vertexSparse struct {
	vertex
	edges               map[int64]stateIndex
	invalidEdgesSet     bool
	invalidEdgesToState stateIndex
}

func newVertexSparse(numberOfStates int) vertexSparse {
	v := vertexSparse{
		vertex: newVertex(),
		edges:  make(map[int64]stateIndex, numberOfStates),
	}
	return v
}

func (v *vertexSparse) nextState(edge int64) (next stateIndex, ok bool) {
	next, ok = v.edges[edge]
	if !ok && v.invalidEdgesSet {
		next = v.invalidEdgesToState
		ok = true
	}
	return
}

func (v *vertexSparse) setNextState(edge int64, si stateIndex) {
	v.edges[edge] = si
}

func (v *vertexSparse) setInvalidEdgesTo(si stateIndex) {
	v.invalidEdgesSet = true
	v.invalidEdgesToState = si
}
