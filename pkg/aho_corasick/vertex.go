package aho_corasick

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
	for i := range v.nextState {
		v.nextState[i] = -1
	}
	return v
}

func (v *vertex) addOutputIndex(i int) {
	if v.output == nil {
		v.output = make([]int, 0, 5)
	}
	v.output = append(v.output, i)
}
