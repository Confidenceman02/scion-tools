// Package basics has a heap of useful functions inspired by the Elm basics module.
package basics

import (
	"math"
	"reflect"
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

// Check if values are structurally &ldquo;the same&rdquo;.
func Eq[T any](x, y T) bool {
	switch reflect.TypeOf(x).Kind() {
	case reflect.Func:
		// mimic Elm's behavior
		panic("Can't compare functions")
	default:
		return reflect.DeepEqual(x, y)
	}
}

// Negate a boolean value.
func Not(pred bool) bool {
	return !pred
}

// Take the square root of a number.
func Sqrt(n Float) Float {
	return Float(math.Sqrt(float64(n)))
}

func ModBy(modulus Int, x Int) Int {
	answer := math.Mod(float64(x), float64(modulus))
	if modulus == 0 {
		panic("ModBy: modulus cannot be zero")
	}

	if (answer > 0 && modulus < 0) || (answer < 0 && modulus > 0) {
		return Int(answer) + modulus
	} else {
		return Int(answer)
	}
}

// FUNCTION HELPERS

// Function composition, passing results along to the left direction.
func ComposeL[A any, B any, C any](g func(B) C, f func(A) B) func(A) C {
	return func(x A) C { return g(f(x)) }
}
