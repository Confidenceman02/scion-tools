package char

import "unicode"

type Char rune

// Digits

// Detect digits 0123456789
func IsDigit(c Char) bool {
	return unicode.IsDigit(rune(c))
}
