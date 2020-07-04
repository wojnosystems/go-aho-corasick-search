package fifo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFifoInt_IsEmpty(t *testing.T) {
	cases := map[string]struct {
		build    func() Int
		expected bool
	}{
		"empty": {
			build: func() Int {
				return Int{}
			},
			expected: true,
		},
		"push": {
			build: func() Int {
				f := Int{}
				f.Push(2)
				return f
			},
			expected: false,
		},
		"pop": {
			build: func() Int {
				f := Int{}
				f.Pop()
				return f
			},
			expected: true,
		},
		"push-pop": {
			build: func() Int {
				f := Int{}
				f.Push(2)
				f.Pop()
				return f
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
		build         func() Int
		expectedValue int
		expected      bool
	}{
		"empty": {
			build: func() Int {
				return Int{}
			},
			expectedValue: -1,
			expected:      false,
		},
		"push": {
			build: func() Int {
				f := Int{}
				f.Push(2)
				return f
			},
			expectedValue: 2,
			expected:      true,
		},
		"pop": {
			build: func() Int {
				f := Int{}
				f.Pop()
				return f
			},
			expectedValue: -1,
			expected:      false,
		},
		"push-pop": {
			build: func() Int {
				f := Int{}
				f.Push(2)
				f.Pop()
				return f
			},
			expectedValue: -1,
			expected:      false,
		},
		"push-push-push": {
			build: func() Int {
				f := Int{}
				f.Push(2)
				f.Push(3)
				f.Push(4)
				return f
			},
			expectedValue: 2,
			expected:      true,
		},
		"push-push-push-pop-pop": {
			build: func() Int {
				f := Int{}
				f.Push(2)
				f.Push(3)
				f.Push(4)
				f.Pop()
				f.Pop()
				return f
			},
			expectedValue: 4,
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
