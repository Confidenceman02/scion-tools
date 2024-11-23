// Package basics has a heap of useful functions inspired by the Elm basics module.
package basics

import (
	"cmp"
	"math"
	"reflect"
)

type Number interface {
	Int | Float
}

type Int int
type Float float32

func (i Int) Cmp(y Comparable[Int]) int {
	return cmp.Compare(i, y.T())
}
func (i Float) Cmp(y Comparable[Float]) int {
	return cmp.Compare(i, y.T())
}
func (i Int) T() Int {
	return i
}
func (i Float) T() Float {
	return i
}

type Comparable[T any] interface {
	Cmp(Comparable[T]) int
	T() T
}

type Appendable[T any] interface {
	App(Appendable[T]) Appendable[T]
	T() T
}

// Math

// Add two numbers. The number type variable means this operation can be specialized to any Number type.
func Add[T Number](a, b T) T {
	return a + b
}

// Subtract numbers like 4 - 3 == 1.
func Sub[T Number](a, b T) T {
	return a - b
}

// Multiply numbers like `2 * 3 == 6`.
func Mul[T Number](a, b T) T {
	return a * b
}

// Floating-point division:
func Fdiv(a Float, b Float) Float {
	return Float(float32(a) / float32(b))
}

// Int to Float / Float to Int

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

// EQUALITY

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

// COMPARISON

// <
func Lt[T any](x Comparable[T], y Comparable[T]) bool {
	return x.Cmp(y) < 0
}

// >
func Gt[T any](x Comparable[T], y Comparable[T]) bool {
	return x.Cmp(y) > 0
}

// <=
func Le[T any](x Comparable[T], y Comparable[T]) bool {
	return x.Cmp(y) <= 0
}

// >=
func Ge[T any](x Comparable[T], y Comparable[T]) bool {
	return x.Cmp(y) >= 0
}

func Max[T any](x Comparable[T], y Comparable[T]) Comparable[T] {
	if Gt(x, y) {
		return x
	} else {
		return y
	}
}

func Min[T any](x Comparable[T], y Comparable[T]) Comparable[T] {
	if Lt(x, y) {
		return x
	} else {
		return y
	}
}

func Compare[T any](x Comparable[T], y Comparable[T]) Order {
	n := x.Cmp(y)
	if n < 0 {
		return LT{}
	} else if n == 0 {
		return EQ{}
	} else {
		return GT{}
	}
}

type Order interface {
	_order() order
}

type order struct{}

func (ord order) _order() order {
	return ord
}

type LT struct {
	order
}

type EQ struct {
	order
}

type GT struct {
	order
}

// BOOLEANS

// Negate a boolean value.
func Not(pred bool) bool {
	return !pred
}

// Append Strings and Lists

func Append[T any](a Appendable[T], b Appendable[T]) Appendable[T] {
	return a.App(b)
}

// Fancier Math

// Perform arithmetic.
// A common trick is to use (n mod 2) to detect even and odd numbers:
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

// Negate a number.
func Negate[A Number](n A) A {
	return -n
}

// Take the square root of a number.
func Sqrt(n Float) Float {
	return Float(math.Sqrt(float64(n)))
}

// Function helpers

func Identity[A any](x A) A {
	return x
}

// Function composition, passing results along to the left direction.
func ComposeL[A any, B any, C any](g func(B) C, f func(A) B) func(A) C {
	return func(x A) C { return g(f(x)) }
}
