package fifo

type Int struct {
	items    []int
	writePos int
	readPos  int
}

func (f *Int) Push(i int) {
	if f.items == nil {
		f.items = make([]int, 0, 10)
	}
	f.items = append(f.items, i)
	f.writePos++
}

func (f Int) Peek() (i int, ok bool) {
	if f.items == nil {
		return -1, false
	}
	if f.IsEmpty() {
		return -1, false
	}
	return f.items[f.readPos], true
}

func (f *Int) Pop() {
	f.readPos++
	if f.IsEmpty() {
		f.Reset()
	}
}

func (f Int) IsEmpty() bool {
	return f.readPos >= f.writePos
}

func (f *Int) Reset() {
	f.readPos = 0
	f.writePos = 0
}
