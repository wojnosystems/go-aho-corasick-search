package ac_engines

import "fmt"

type InvalidCharsetError struct {
	Char rune
}

func newInvalidCharsetError(char rune) error {
	return &InvalidCharsetError{
		Char: char,
	}
}

func (e *InvalidCharsetError) Error() string {
	return fmt.Sprintf("encountered rune: '%c' but it is not supported by this machine", e.Char)
}
