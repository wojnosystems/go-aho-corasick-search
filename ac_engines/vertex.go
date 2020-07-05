package ac_engines

type vertex struct {
	failState stateIndex
	output    []int
}

func newVertex() vertex {
	return vertex{
		failState: invalidState,
		output:    nil,
	}
}

func (v *vertex) setFailState(state stateIndex) {
	v.failState = state
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
