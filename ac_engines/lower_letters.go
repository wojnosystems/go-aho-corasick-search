package ac_engines

import (
	"errors"
	"github.com/wojnosystems/go-aho-corasick-search"
	"github.com/wojnosystems/go-aho-corasick-search/result"
	"io"
	"os"
)

type LowerLetters struct {
	states denseStates
}

const (
	letterStates = 26
)

func NewLowerLetters(keyWords []string) (m LowerLetters, err error) {
	m = LowerLetters{}
	m.states, err = buildDenseTrie(keyWords)
	if err != nil {
		return m, err
	}
	m.states = buildDenseFails(m.states)
	return m, err
}

// buildTrie AKA the goto function
func buildDenseTrie(keyWords []string) (states denseStates, err error) {
	states = newDenseStates()
	states = append(states, newVertexDense(letterStates))
	for keywordIndex, word := range keyWords {
		states, err = addKeywordToDenseTrie(states, keywordIndex, word)
		if err != nil {
			return
		}
	}

	// Set all of the start Tries that don't match back to the start
	states[startState].setInvalidEdgesTo(startState)
	return
}

func addKeywordToDenseTrie(statesIn denseStates, keywordIndex int, keyword string) (states denseStates, err error) {
	states = statesIn
	currentState := startState
	letterIndex := 0
	for ; letterIndex < len(keyword); letterIndex++ {
		letter := keyword[letterIndex]
		if !isLowerAscii(letter) {
			return states, newInvalidCharsetError(rune(letter))
		}
		letterStateIndex := letter - 'a'
		nextState, ok := states[currentState].nextState(int64(letterStateIndex))
		if !ok {
			break
		}
		currentState = nextState
	}
	for ; letterIndex < len(keyword); letterIndex++ {
		letter := keyword[letterIndex]
		letterStateIndex := letter - 'a'
		v := newVertexDense(letterStates)
		states[currentState].setNextState(int64(letterStateIndex), states.lastStateIndex())
		currentState = states.lastStateIndex()
		states = append(states, v)
	}
	states[currentState].appendOutputIndex([]int{keywordIndex})
	return
}

func isLowerAscii(r uint8) bool {
	return 'a' <= r && r <= 'z'
}

func buildDenseFails(statesIn denseStates) (states denseStates) {
	states = statesIn
	start := states[startState]
	q := stateFifo{}
	// initialize the denseStates at depth 1
	for _, state := range start.edges {
		if state != startState {
			// All failure denseStates at depth 1 go back to the start state
			states[state].failState = startState
			q.Push(state)
		}
	}

	for !q.IsEmpty() {
		statePreviousDepth, _ := q.Peek()
		q.Pop()
		for letterIndex, s := range states[statePreviousDepth].edges {
			if s != invalidState {
				q.Push(s)
				state := states[statePreviousDepth].failState
				_, ok := states[state].nextState(int64(letterIndex))
				for !ok {
					state = states[state].failState
					_, ok = states[state].nextState(int64(letterIndex))
				}
				// Found a failState that is not invalid
				// set the nextState fail state to the entry for this letter
				states[s].failState, _ = states[state].nextState(int64(letterIndex))
				states[s].appendOutputIndex(states[states[s].failState].output)
			}
		}
	}
	return
}

func (m LowerLetters) Search(input io.Reader, results result.Writer) (err error) {
	var bytesReadSoFar uint64
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
		bytesReadSoFar++
		letterIndex := letter[0] - 'a'
		_, ok := m.states[currentState].nextState(int64(letterIndex))
		for !ok {
			currentState = m.states[currentState].failState
			_, ok = m.states[currentState].nextState(int64(letterIndex))
		}
		currentState, _ = m.states[currentState].nextState(int64(letterIndex))
		if m.states[currentState].hasOutput() {
			for _, output := range m.states[currentState].outputs() {
				results.Emit(aho_corasick_search.Output{
					KeywordIndex:    output,
					ByteOffset:      bytesReadSoFar,
					CharacterOffset: bytesReadSoFar,
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
