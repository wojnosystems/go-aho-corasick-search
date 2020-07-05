package ac_engines

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFifoInt_IsEmpty(t *testing.T) {
	cases := map[string]struct {
		build    func() stateFifo
		expected bool
	}{
		"empty": {
			build: func() (f stateFifo) {
				return
			},
			expected: true,
		},
		"push": {
			build: func() (f stateFifo) {
				f.Push(2)
				return
			},
			expected: false,
		},
		"pop": {
			build: func() (f stateFifo) {
				f.Pop()
				return
			},
			expected: true,
		},
		"push-pop": {
			build: func() (f stateFifo) {
				f.Push(2)
				f.Pop()
				return
			},
			expected: true,
		},
	}

	for caseName, c := range cases {
		actual := c.build()
		if c.expected {
			assert.True(t, actual.IsEmpty(), caseName)
		} else {
			assert.False(t, actual.IsEmpty(), caseName)
		}
	}
}

func TestFifoInt_Peek(t *testing.T) {
	cases := map[string]struct {
		build         func() stateFifo
		expectedValue stateIndex
		expected      bool
	}{
		"empty": {
			build: func() (f stateFifo) {
				return
			},
			expectedValue: stateIndex(-1),
			expected:      false,
		},
		"push": {
			build: func() (f stateFifo) {
				f.Push(2)
				return f
			},
			expectedValue: stateIndex(2),
			expected:      true,
		},
		"pop": {
			build: func() (f stateFifo) {
				f.Pop()
				return
			},
			expectedValue: stateIndex(-1),
			expected:      false,
		},
		"push-pop": {
			build: func() (f stateFifo) {
				f.Push(2)
				f.Pop()
				return
			},
			expectedValue: stateIndex(-1),
			expected:      false,
		},
		"push-push-push": {
			build: func() (f stateFifo) {
				f.Push(2)
				f.Push(3)
				f.Push(4)
				return
			},
			expectedValue: stateIndex(2),
			expected:      true,
		},
		"push-push-push-pop-pop": {
			build: func() (f stateFifo) {
				f.Push(2)
				f.Push(3)
				f.Push(4)
				f.Pop()
				f.Pop()
				return
			},
			expectedValue: stateIndex(4),
			expected:      true,
		},
	}

	for caseName, c := range cases {
		actual := c.build()
		actualValue, actualExists := actual.Peek()
		assert.Equal(t, c.expected, actualExists, caseName)
		assert.Equal(t, c.expectedValue, actualValue, caseName)
	}
}
