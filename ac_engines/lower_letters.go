package ac_engines

import (
	"errors"
	"github.com/wojnosystems/go-aho-corasick-search"
	"github.com/wojnosystems/go-aho-corasick-search/fifo"
	"github.com/wojnosystems/go-aho-corasick-search/result"
	"io"
	"os"
)

type LowerLetters struct {
	states states
}

const (
	letterStates = 26
)

func NewLowerLetters(keyWords []string) (m LowerLetters, err error) {
	m = LowerLetters{}
	m.states, err = buildTrie(keyWords)
	if err != nil {
		return m, err
	}
	m.states = buildFails(m.states)
	return m, err
}

// buildTrie AKA the goto function
func buildTrie(keyWords []string) (states states, err error) {
	states = newStates()
	start := newVertex(letterStates)
	states = append(states, start)
	for keywordIndex, word := range keyWords {
		states, err = addKeywordToTrie(states, keywordIndex, word)
		if err != nil {
			return
		}
	}

	// Set all of the start Tries that don't match back to the start
	for stateIndex, state := range start.nextState {
		if state == invalidState {
			start.nextState[stateIndex] = startState
		}
	}
	return
}

func addKeywordToTrie(statesIn states, keywordIndex int, keyword string) (states states, err error) {
	states = statesIn
	currentState := 0
	letterIndex := 0
	for ; letterIndex < len(keyword); letterIndex++ {
		letter := keyword[letterIndex]
		if !isLowerAscii(letter) {
			return states, newInvalidCharsetError(rune(letter))
		}
		letterStateIndex := letter - 'a'
		nextState := states[currentState].nextState[letterStateIndex]
		if nextState == invalidState {
			break
		}
		currentState = nextState
	}
	for ; letterIndex < len(keyword); letterIndex++ {
		letter := keyword[letterIndex]
		letterStateIndex := letter - 'a'
		v := newVertex(letterStates)
		states[currentState].nextState[letterStateIndex] = len(states)
		currentState = len(states)
		states = append(states, v)
	}
	states[currentState].appendOutputIndex([]int{keywordIndex})
	return
}

func isLowerAscii(r uint8) bool {
	return 'a' <= r && r <= 'z'
}

func buildFails(statesIn states) (states states) {
	states = statesIn
	start := states[startState]
	q := fifo.Int{}
	// initialize the states at depth 1
	for _, state := range start.nextState {
		if state != startState {
			// All failure states at depth 1 go back to the start state
			states[state].failState = startState
			q.Push(state)
		}
	}

	for !q.IsEmpty() {
		statePreviousDepth, _ := q.Peek()
		q.Pop()
		for letterIndex, s := range states[statePreviousDepth].nextState {
			if s != invalidState {
				q.Push(s)
				state := states[statePreviousDepth].failState
				for states[state].nextState[letterIndex] == invalidState {
					state = states[state].failState
				}
				// Found a failState that is not invalid
				// set the nextState fail state to the entry for this letter
				states[s].failState = states[state].nextState[letterIndex]
				states[s].appendOutputIndex(states[states[s].failState].output)
			}
		}
	}
	return
}

func (m LowerLetters) Search(input io.Reader, results result.Writer) (err error) {
	currentState := startState
	letter := []byte{0}
	defer func() {
		_ = results.Close()
	}()
	for {
		_, err = input.Read(letter)
		if err != nil {
			if isEOFOrClosed(err) {
				err = nil
			}
			return
		}
		letterIndex := letter[0] - 'a'
		for m.states[currentState].nextState[letterIndex] == invalidState {
			currentState = m.states[currentState].failState
		}
		currentState = m.states[currentState].nextState[letterIndex]
		if m.states[currentState].hasOutput() {
			for _, output := range m.states[currentState].outputs() {
				results.Emit(aho_corasick_search.Output{
					KeywordIndex: output,
				})
			}
		}
	}
}

func isEOFOrClosed(err error) bool {
	if errors.Is(err, io.EOF) ||
		errors.Is(err, io.ErrClosedPipe) {
		return true
	}
	if pathErr, ok := err.(*os.PathError); ok {
		if pathErr.Err.Error() == "file already closed" {
			return true
		}
	}
	return false
}
