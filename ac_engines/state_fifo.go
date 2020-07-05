package ac_engines

type stateFifo struct {
	items    []stateIndex
	writePos int
	readPos  int
}

func (f *stateFifo) Push(i stateIndex) {
	if f.items == nil {
		f.items = make([]stateIndex, 0, 10)
	}
	f.items = append(f.items, i)
	f.writePos++
}

func (f stateFifo) Peek() (i stateIndex, ok bool) {
	if f.items == nil {
		return -1, false
	}
	if f.IsEmpty() {
		return -1, false
	}
	return f.items[f.readPos], true
}

func (f *stateFifo) Pop() {
	f.readPos++
	if f.IsEmpty() {
		f.Reset()
	}
}

func (f stateFifo) IsEmpty() bool {
	return f.readPos >= f.writePos
}

func (f *stateFifo) Reset() {
	f.readPos = 0
	f.writePos = 0
}
