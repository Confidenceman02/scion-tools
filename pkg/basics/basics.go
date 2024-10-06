// Package basics has a heap of useful functions inspired by the Elm basics module.
package basics

import (
	"math"
)

type Int int
type Float float64

type Number interface {
	Int | Float
}

// Add two numbers. The number type variable means this operation can be specialized to any Number type.
func Add[T Number](a, b T) T {
	return a + b
}

// ToFloat - Convert an integer into a float. Useful when mixing Int and Float values.
func ToFloat[T Int](x T) Float {
	return Float(x)
}

// Round a number to the nearest integer.
func Round(x Float) Int {
	return Int(math.Round(float64(x)))
}

// Floor function, rounding down.
func Floor(x Float) Int {
	return Int(math.Floor(float64(x)))
}

// Ceiling function, rounding up.
func Ceiling(x Float) Int {
	return Int(math.Ceil(float64(x)))
}

// Truncate a number, rounding towards zero
func Truncate(x Float) Int {
	return Int(math.Trunc(float64(x)))
}
