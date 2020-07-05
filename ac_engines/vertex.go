package ac_engines

type vertex struct {
	nextState []stateIndex
	failState stateIndex
	output    []int
}

func newVertex(numberOfStates int) vertex {
	v := vertex{
		nextState: make([]stateIndex, numberOfStates),
		failState: -1,
		output:    nil,
	}
	// By default, all states are invalid
	for i := range v.nextState {
		v.nextState[i] = invalidState
	}
	return v
}

func (v *vertex) appendOutputIndex(i []int) {
	if v.output == nil {
		v.output = make([]int, 0, 5)
	}
	v.output = append(v.output, i...)
}

func (v *vertex) outputs() []int {
	if v.output == nil {
		return make([]int, 0, 0)
	}
	return v.output
}

func (v vertex) hasOutput() bool {
	return v.output != nil && len(v.output) != 0
}
