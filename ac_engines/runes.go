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

type RunesBuilder struct {
	states        sparseStates
	keywordsSoFar int
}

func NewRunes() (m RunesBuilder) {
	m = RunesBuilder{}
	m.states = newSparseStates()
	m.states = append(m.states, newVertexSparse(letterStates))
	return
}

func (m *RunesBuilder) AddKeyword(keyword string) (err error) {
	m.states, err = addKeywordToSparseTrie(m.states, m.keywordsSoFar, keyword)
	m.keywordsSoFar++
	return
}

func (m *RunesBuilder) Build() Runes {
	r := Runes{
		states: m.states,
	}
	m.states = nil
	r.states[startState].setInvalidEdgesTo(startState)
	r.states = buildSparseFails(r.states)
	return r
}

func addKeywordToSparseTrie(statesIn sparseStates, keywordIndex int, keyword string) (states sparseStates, err error) {
	states = statesIn
	currentState := startState
	letterIndex := 0
	keywordRunes := []rune(keyword)
	for ; letterIndex < len(keywordRunes); letterIndex++ {
		letter := keywordRunes[letterIndex]
		letterStateIndex := letter
		nextState, ok := states[currentState].nextState(int64(letterStateIndex))
		if !ok {
			break
		}
		currentState = nextState
	}
	for ; letterIndex < len(keywordRunes); letterIndex++ {
		letter := keywordRunes[letterIndex]
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
	return
}

func (m Runes) Search(input io.Reader, results result.Writer) (err error) {
	var bytesReadSoFar uint64
	var bytesReadThisRune int
	var runesReadSoFar uint64
	bufReader := bufio.NewReader(input)
	currentState := startState
	var letter rune
	defer func() {
		_ = results.Close()
	}()
	for {
		letter, bytesReadThisRune, err = bufReader.ReadRune()
		if err != nil {
			if isEOFOrClosed(err) {
				err = nil
			}
			return
		}
		bytesReadSoFar += uint64(bytesReadThisRune)
		runesReadSoFar++
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
					KeywordIndex:    output,
					CharacterOffset: runesReadSoFar,
					ByteOffset:      bytesReadSoFar,
				})
			}
		}
	}
}
