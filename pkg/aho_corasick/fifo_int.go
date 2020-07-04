package aho_corasick

type fifoInt struct {
	items    []int
	writePos int
	readPos  int
}

func (f *fifoInt) Push(i int) {
	if f.items == nil {
		f.items = make([]int, 0, 10)
	}
	f.items = append(f.items, i)
	f.writePos++
}

func (f fifoInt) Peek() (i int, ok bool) {
	if f.items == nil {
		return -1, false
	}
	if f.IsEmpty() {
		return -1, false
	}
	return f.items[f.readPos], true
}

func (f *fifoInt) Pop() {
	f.readPos++
	if f.IsEmpty() {
		f.Reset()
	}
}

func (f fifoInt) IsEmpty() bool {
	return f.readPos >= f.writePos
}

func (f *fifoInt) Reset() {
	f.readPos = 0
	f.writePos = 0
}
