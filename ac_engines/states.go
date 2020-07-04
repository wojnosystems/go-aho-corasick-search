package ac_engines

type states []vertex

func newStates() states {
	return make([]vertex, 0, 10)
}
