package string

import (
	"cmp"
	"github.com/Confidenceman02/scion-tools/pkg/basics"
	"github.com/Confidenceman02/scion-tools/pkg/bitwise"
	"github.com/Confidenceman02/scion-tools/pkg/char"
	"github.com/Confidenceman02/scion-tools/pkg/list"
	"github.com/Confidenceman02/scion-tools/pkg/maybe"
	"github.com/Confidenceman02/scion-tools/pkg/tuple"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type String string

func (s1 String) App(s2 basics.Appendable[String]) basics.Appendable[String] {
	return s1.T() + s2.T()
}
func (i String) Cmp(y basics.Comparable[String]) int {
	return cmp.Compare(i, y.T())
}
func (i String) T() String {
	return i
}

// Strings

// Determine if a string is empty.
func IsEmpty(x String) bool {
	return x == ""
}

// Get the length of a string.
func Length(x String) basics.Int {
	return basics.Int(utf8.RuneCountInString(string(x)))
}

// Reverse a string.
func Reverse(x String) String {
	r := []rune(x)

	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return String(r)
}

func Repeat(n basics.Int, chunk String) String {
	return repeatHelp(n, chunk, "")
}

func repeatHelp(n basics.Int, chunk String, result String) String {
	if n <= 0 {
		return result
	} else {
		var r String
		if bitwise.And(n, 1) == 0 {
			r = result
		} else {
			r = Append(result, chunk).T()
		}
		return repeatHelp(bitwise.ShiftRightBy(1, n), Append(chunk, chunk).T(), r)
	}
}

// Replace all occurrences of some substring.
func Replace(before String, after String, str String) String {
	return Join(after, Split(before, str))
}

// Building and Splitting

// Append two strings. You can also use basics.Append to do this.
func Append(x String, y String) String {
	return x + y
}

// Concatenate many strings into one.
func Concat(chunks list.List[String]) String {
	return Join("", chunks)
}

// Split a string using a given separator.
func Split(sep String, s String) list.List[String] {
	return list.FromSliceMap(
		func(s string) String { return String(s) },
		strings.Split(string(s), string(sep)),
	)
}

// Put many strings together with a given separator.
func Join(sep String, chunks list.List[String]) String {
	return String(
		strings.Join(
			list.ToSliceMap(func(a String) string { return string(a) },
				chunks,
			),
			string(sep),
		),
	)
}

// Break a string into words, splitting on chunks of whitespace.
func Words(str String) list.List[String] {
	return list.FromSliceMap(
		func(s string) String { return String(s) },
		regexp.MustCompile("\\s+").Split(strings.Trim(string(str), " "), -1),
	)
}

// Break a string into lines, splitting on newlines.
func Lines(str String) list.List[String] {
	return list.FromSliceMap(
		func(s string) String { return String(s) },
		regexp.MustCompile("\\r\\n|\\r|\\n").Split(string(str), -1),
	)

}

// Get Substrings

// Take a substring given a start and end index. Negative indexes are taken starting from the end of the list.
func Slice(start basics.Int, end basics.Int, str String) String {
	var start1 int
	var end1 int
	count := utf8.RuneCountInString(string(str))

	// resolve start
	if math.Signbit(float64(start)) {
		start1 = count + int(start)
		if start1 < 0 {
			start1 = 0
		}
	} else {
		start1 = int(start)
	}

	// resolve end
	if math.Signbit(float64(end)) {
		// Add negative to count
		end1 = count + int(end)
	} else {
		end1 = int(end)
		if end1 > count {
			end1 = count
		}
	}

	if start1 >= end1 || start1 > count {
		return ""
	} else {
		return str[start1:end1]
	}
}

// Take *n* characters from the left side of a string.
func Left(n basics.Int, str String) String {
	if n < 1 {
		return ""
	} else {
		return Slice(0, n, str)
	}
}

// Take *n* characters from the right side of a string.
func Right(n basics.Int, str String) String {
	if n < 1 {
		return ""
	} else {
		return Slice(-n, Length(str), str)
	}
}

// Drop *n* characters from the left side of a string.
func DropLeft(n basics.Int, str String) String {
	if n < 1 {
		return ""
	} else {
		return Slice(n, Length(str), str)
	}
}

// Drop *n* characters from the right side of a string.
func DropRight(n basics.Int, str String) String {
	if n < 1 {
		return ""
	} else {
		return Slice(0, -n, str)
	}
}

// Check For Substrings

// See if the second string contains the first one.
func Contains(sub String, str String) bool {
	return strings.Index(string(str), string(sub)) > -1
}

// See if the second string starts with the first one.
func StartsWith(sub String, str String) bool {
	return strings.Index(string(str), string(sub)) == 0
}

// See if the second string ends with the first one.
func EndsWith(sub String, str String) bool {
	return strings.LastIndex(string(str), string(sub)) == int(Length(str)-Length(sub))
}

// Get all of the indexes for a substring in another string.
func Indexes(sub String, str String) list.List[basics.Int] {
	subLen := Length(sub)
	if subLen < 1 {
		return list.Empty[basics.Int]()
	}
	var s = str
	var is []basics.Int = []basics.Int{}
	var i = strings.Index(string(s), string(sub))

	for i > -1 {
		isLen := len(is)
		ii := basics.Int(i + 1)
		if 0 < isLen {
			is = append(is, basics.Add(is[isLen-1], ii))
			s = DropLeft(ii, s)
			i = strings.Index(string(s), string(sub))
		} else {
			is = append(is, basics.Int(i))
			s = DropLeft(basics.Int(ii), s)
			i = strings.Index(string(s), string(sub))
		}
	}

	return list.FromSlice(is)
}

// Alias for `indexes`.
func Indices(sub String, str String) list.List[basics.Int] {
	return Indexes(sub, str)
}

// Int Conversions

// Try to convert a string into an int, failing on improperly formatted strings.
func ToInt(x String) maybe.Maybe[basics.Int] {
	var total int
	code0 := x[0]
	var start int = 0

	if code0 == 0x2B /* + */ || code0 == 0x2D /* - */ {
		start = 1
	}

	var i = start

	for i < len(x) {
		var code = x[i]
		if code < 0x30 || 0x39 < code /* 0 - 9 */ {
			return maybe.Nothing{}
		}

		total = 10*total + int(code-0x30)
		i++
	}

	if i == start {
		return maybe.Nothing{}
	} else {
		if code0 == 0x2D /* - */ {
			return maybe.Just[basics.Int]{Value: basics.Int(-total)}
		} else {
			return maybe.Just[basics.Int]{Value: basics.Int(total)}
		}
	}
}

// Convert an Int to a String.
func FromInt(x basics.Int) String {
	return String(strconv.FormatInt(int64(x), 10))
}

// Float conversions

// Try to convert a string into a float, failing on improperly formatted strings.
func ToFloat(x String) maybe.Maybe[basics.Float] {
	f, err := strconv.ParseFloat(string(x), 32)

	if err != nil {
		return maybe.Nothing{}
	} else {
		return maybe.Just[basics.Float]{Value: basics.Float(f)}
	}
}

// Convert a Float to a String.
func FromFloat(x basics.Float) String {
	return String(strconv.FormatFloat(float64(x), 'g', -1, 32))
}

// Char Conversions

// Create a string from a given character.
func FromChar(char char.Char) String {
	return String(string(char))
}

// Add a character to the beginning of a string.
func Cons(char char.Char, str String) String {
	return String(string(char) + string(str))
}

// Split a non-empty string into its head and tail. This lets you pattern match on strings exactly as you would with lists.
func Uncons(str String) maybe.Maybe[tuple.Tuple2[char.Char, String]] {
	if str == "" {
		return maybe.Nothing{}
	} else {
		word, ln := utf8.DecodeRuneInString(string(str))
		return maybe.Just[tuple.Tuple2[char.Char, String]]{
			Value: tuple.Pair(char.Char(word), String(string(str)[ln:])),
		}
	}
}

// List Conversions

// Convert a string to a list of characters.
func ToList(str String) list.List[char.Char] {
	return Foldr(list.Cons[char.Char], list.Empty[char.Char](), str)
}

// Convert a list of characters into a String. Can be useful if you want to create a string primarily by consing, perhaps for decoding something.
func FromList(chars list.List[char.Char]) String {
	return String(
		strings.Join(
			list.ToSliceMap(func(c char.Char) string { return string(c) }, chars),
			"",
		),
	)
}

// Formatting

// Convert a string to all upper case. Useful for case-insensitive comparisons and VIRTUAL YELLING.
func ToUpper(str String) String {
	return String(strings.ToUpper(string(str)))
}

// Convert a string to all lower case. Useful for case-insensitive comparisons.
func ToLower(str String) String {
	return String(strings.ToLower(string(str)))
}

// Pad a string on both sides until it has a given length.
func Pad(n basics.Int, char char.Char, str String) String {
	half := basics.Fdiv(basics.ToFloat(basics.Sub(n, Length(str))), 2)

	return basics.Append(
		basics.Append(
			Repeat(basics.Ceiling(half), FromChar(char)),
			str,
		),
		Repeat(basics.Floor(half), FromChar(char)),
	).T()
}

// Pad a string on the left until it has a given length.
func PadLeft(n basics.Int, char char.Char, str String) String {
	return basics.Append(
		Repeat(basics.Sub(n, Length(str)), FromChar(char)),
		str,
	).T()
}

// Pad a string on the right until it has a given length.
func PadRight(n basics.Int, char char.Char, str String) String {
	return basics.Append(str, Repeat(basics.Sub(n, Length(str)), FromChar(char))).T()
}

// Get rid of whitespace on both sides of a string.
func Trim(str String) String {
	return String(strings.TrimSpace(string(str)))
}

// Get rid of whitespace on the left of a string.
func TrimLeft(str String) String {
	return String(strings.TrimLeftFunc(string(str), unicode.IsSpace))
}

// Get rid of whitespace on the right of a string.
func TrimRight(str String) String {
	return String(strings.TrimRightFunc(string(str), unicode.IsSpace))
}

// Higher-Order Functions

// Transform every character in a string
func Map(f func(char.Char) char.Char, str String) String {
	var builder strings.Builder

	for _, r := range str {
		builder.WriteRune(rune(f(char.Char(r))))
	}

	return String(builder.String())
}

// Keep only the characters that pass the test.
func Filter(isGood func(char.Char) bool, str String) String {
	var builder strings.Builder

	for _, r := range str {
		if isGood(char.Char(r)) {
			builder.WriteRune(r)
		}
	}

	return String(builder.String())
}

// Reduce a string from the left.
func Foldl[B any](f func(char.Char, B) B, state B, str String) B {
	for _, r := range str {
		state = f(char.Char(r), state)
	}
	return state
}

// Reduce a string from the right.
func Foldr[B any](f func(char.Char, B) B, state B, str String) B {
	runeSlice := []rune(str)
	for i := len(runeSlice); i >= 1; i-- {
		state = f(char.Char(runeSlice[i-1]), state)
	}
	return state
}

// Determine whether any characters pass the test.
func Any(isGood func(char.Char) bool, str String) bool {
	for _, r := range str {
		if isGood(char.Char(r)) {
			return true
		}
	}
	return false
}

// Determine whether all characters pass the test.
func All(isGood func(char.Char) bool, str String) bool {
	for _, r := range str {
		if !isGood(char.Char(r)) {
			return false
		}
	}
	return true
}
