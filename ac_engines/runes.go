package ac_engines

import (
	"bufio"
	"github.com/wojnosystems/go-aho-corasick-search"
	"github.com/wojnosystems/go-aho-corasick-search/result"
	"io"
)

type Runes struct {
	states sparseStates
}

func NewRunes(keyWords []string) (m Runes, err error) {
	m = Runes{}
	m.states, err = buildSparseTrie(keyWords)
	if err != nil {
		return m, err
	}
	m.states = buildSparseFails(m.states)
	return m, err
}

// buildTrie AKA the goto function
func buildSparseTrie(keyWords []string) (states sparseStates, err error) {
	states = newSparseStates()
	states = append(states, newVertexSparse(letterStates))
	for keywordIndex, word := range keyWords {
		states, err = addKeywordToSparseTrie(states, keywordIndex, word)
		if err != nil {
			return
		}
	}

	// Set all of the start Tries that don't match back to the start
	states[startState].setInvalidEdgesTo(startState)
	return
}

func addKeywordToSparseTrie(statesIn sparseStates, keywordIndex int, keyword string) (states sparseStates, err error) {
	states = statesIn
	currentState := startState
	letterIndex := 0
	for ; letterIndex < len(keyword); letterIndex++ {
		letter := keyword[letterIndex]
		letterStateIndex := letter
		nextState, ok := states[currentState].nextState(int64(letterStateIndex))
		if !ok {
			break
		}
		currentState = nextState
	}
	for ; letterIndex < len(keyword); letterIndex++ {
		letter := keyword[letterIndex]
		letterStateIndex := letter
		v := newVertexSparse(letterStates)
		states[currentState].setNextState(int64(letterStateIndex), states.lastStateIndex())
		currentState = states.lastStateIndex()
		states = append(states, v)
	}
	states[currentState].appendOutputIndex([]int{keywordIndex})
	return
}

func buildSparseFails(statesIn sparseStates) (states sparseStates) {
	states = statesIn
	start := states[startState]
	q := stateFifo{}
	// initialize the states at depth 1
	for _, state := range start.edges {
		if state != startState {
			q.Push(state)
			states[state].failState = startState
		}
	}

	for !q.IsEmpty() {
		statePreviousDepth, _ := q.Peek()
		q.Pop()
		for letterIndex, s := range states[statePreviousDepth].edges {
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

	// All failure states at depth 1 go back to the start state
	for _, state := range start.edges {
		states[state].setInvalidEdgesTo(startState)
	}

	return
}

func (m Runes) Search(input io.Reader, results result.Writer) (err error) {
	bufReader := bufio.NewReader(input)
	currentState := startState
	var letter rune
	defer func() {
		_ = results.Close()
	}()
	for {
		letter, _, err = bufReader.ReadRune()
		if err != nil {
			if isEOFOrClosed(err) {
				err = nil
			}
			return
		}
		letterIndex := letter
		_, ok := m.states[currentState].nextState(int64(letterIndex))
		for !ok {
			currentState = m.states[currentState].failState
			_, ok = m.states[currentState].nextState(int64(letterIndex))
		}
		currentState, _ = m.states[currentState].nextState(int64(letterIndex))
		if m.states[currentState].hasOutput() {
			for _, output := range m.states[currentState].outputs() {
				results.Emit(aho_corasick_search.Output{
					KeywordIndex: output,
				})
			}
		}
	}
}
