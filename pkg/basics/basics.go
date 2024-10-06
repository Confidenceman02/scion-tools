// Package basics has a heap of useful functions inspired by the Elm basics module.
package basics

import (
	"math"
)

type Number interface {
	Int | Float
}

type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Float interface {
	~float64
}

// Add two numbers. The number type variable means this operation can be specialized to any Number type.
func Add[T Number](a, b T) T {
	return a + b
}

// ToFloat - Convert an integer into a float. Useful when mixing Int and Float values.
func ToFloat[T Int](x T) float64 {
	return float64(x)
}

// Round a number to the nearest integer.
func Round(x float64) int {
	return int(math.Round(x))
}

// Floor function, rounding down.
func Floor(x float64) int {
	return int(math.Floor(x))
}

// Ceiling function, rounding up.
func Ceiling(x float64) int {
	return int(math.Ceil(x))
}

// Truncate a number, rounding towards zero
func Truncate(x float64) int {
	return int(math.Trunc(x))
}
