package aho_corasick

type LowerLetters struct {
	states states
}

const (
	LetterStates = 26
)

func NewLowerLetters(keyWords []string) (LowerLetters, error) {
	m := LowerLetters{
		states: newStates(),
	}
	m.states = append(m.states, newVertex(LetterStates))
	err := m.buildTrie(keyWords)
	if err != nil {
		return m, err
	}
	return m, err
}

func (m *LowerLetters) buildTrie(keyWords []string) (err error) {
	for keywordIndex, word := range keyWords {
		currentNode := m.states[0]
		for _, letter := range word {
			if !isAscii(letter) {
				return newInvalidCharsetError(letter)
			}
			letterState := letter - 'a'
			if currentNode.nextState[letterState] == -1 {
				v := newVertex(LetterStates)
				currentNode.nextState[letterState] = len(m.states)
				m.states = append(m.states, v)
			}
			currentNode = m.states[currentNode.nextState[letterState]]
		}
		currentNode.addOutputIndex(keywordIndex)
	}
	return nil
}

func isAscii(r rune) bool {
	return 'a' <= r && r <= 'z'
}
