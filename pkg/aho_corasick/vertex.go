package aho_corasick

const (
	startState   = 0
	invalidState = -1
)

type vertex struct {
	nextState []int
	failState int
	output    []int
}

func newVertex(numberOfStates int) vertex {
	v := vertex{
		nextState: make([]int, numberOfStates),
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
