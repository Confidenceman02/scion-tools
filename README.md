# scion-tools

[![Actions status](https://Confidenceman02.github.io/djelm/workflows/CI/badge.svg)](https://github.com/Confidenceman02/scion-tools/actions)
[![](https://img.shields.io/badge/license-MIT-blue)](https://github.com/Confidenceman02/scion-tools/blob/main/LICENSE)

> [!NOTE]
> There are some missing packages and functions from the Elm core library, You can see a list of these
> in the [issues](https://github.com/Confidenceman02/scion-tools/issues).

# Elm inspired functional programing in Go

---

## Table of Content

- [The why](#the-why)
- [The when](#the-when)
- <details>
    <summary><a href="#basics">Basics</a></summary>
    <ul>
        <li>
            <a href="#int">Int</a>
        </li>
        <li>
            <a href="#float">Float</a>
        </li>
        <li>
            <a href="#number">Number</a>
        </li>
        <li>
            <a href="#comparable">Comparable</a>
        </li>
        <li>
            <a href="#add">Add</a>
        </li>
        <li>
            <a href="#sub">Sub</a>
        </li>
        <li>
            <a href="#fdiv">Fdiv</a>
        </li>
        <li>
            <a href="#mul">Mul</a>
        </li>
        <li>
            <a href="#tofloat">ToFloat</a>
        </li>
        <li>
            <a href="#round">Round</a>
        </li>
        <li>
            <a href="#floor">Floor</a>
        </li>
        <li>
            <a href="#ceiling">Ceiling</a>
        </li>
        <li>
            <a href="#truncate">Truncate</a>
        </li>
        <li>
            <a href="#lt">Lt</a>
        </li>
        <li>
            <a href="#gt">Gt</a>
        </li>
        <li>
            <a href="#le">Le</a>
        </li>
        <li>
            <a href="#ge">Ge</a>
        </li>
        <li>
            <a href="#max">Max</a>
        </li>
        <li>
            <a href="#min">Min</a>
        </li>
        <li>
            <a href="#compare">Compare</a>
        </li>
        <li>
            <a href="#order">Order</a>
        </li>
        <li>
            <a href="#not">Not</a>
        </li>
        <li>
            <a href="#append">Append</a>
        </li>
        <li>
            <a href="#modby">ModBy</a>
        </li>
        <li>
            <a href="#negate">Negate</a>
        </li>
        <li>
            <a href="#sqrt">Sqrt</a>
        </li>
        <li>
            <a href="#identity">Identity</a>
        </li>
        <li>
            <a href="#composel">ComposeL</a>
        </li>
    </ul>
  </details>
- <details>
    <summary><a href="#bitwise">Bitwise</a></summary>
- <details>
    <summary><a href="#char">Char</a></summary>
- <details>
    <summary><a href="#dict">Dict</a></summary>
- <details>
    <summary><a href="#list">List</a></summary>
- <details>
    <summary><a href="#maybe">Maybe</a></summary>
- <details>
    <summary><a href="#string">String</a></summary>
- <details>
    <summary><a href="#tuple">Tuple</a></summary>

## The Why

Go is an incredibly powerful and performant language that promotes quality and maintainability through its
extensive standard library and emphasis on simplicity. Sometimes however, it can be tricky to express certain
functional patterns in an immutable and reusable fashion that is both concise and elegant.

Elm is a statically and strongly typed language with an approachable syntax that provides exceptional programming ergonomics.
Elm programs, much like Go, are delightful to write and maintain due to the expressiveness one can harness from its functional
programming style.

The scion-tools goal is to give Gophers an expressive functional approach by providing pure Go analogues of Elm's
core library, including its immutable data structures.

It's just Go that looks a little bit like Elm.

[Back to top](#table-of-content)

## The when

With scion-tools, Gophers have the ability to leverage functional programming techniques that harmonize with Go's idiomatic
patterns that we all know and love.

The following scenarios provide some good uses cases where you may benefit from using scion-tools
or other functional programming packages:

- Data processing pipelines
- Error handling
- Concurrency
- Complex algorithms
- Code reusability

[Back to top](#table-of-content)

# Basics

```go
import "github.com/Confidenceman02/scion-tools/pkg/basics"
```

Tons of useful functions.

## Int

A wrapped Go `int`.

```go
type Int int
```

[Back to top](#table-of-content)

## Float

A wrapped Go `float32`.

```go
type Float float32
```

[Back to top](#table-of-content)

## Number

A type alias for [Int](#int) and [Float](#float)

```go
type Number interface {
  Int | Float
}
```

[Back to top](#table-of-content)

## Comparable

A special type that represents all types that can be compared:

- [Number](#number)
- [String](#string)
- [Char](#char)
- [List](#list) of Comparable
- [Tuple](#tuple) of Comparable

[Back to top](#table-of-content)

## Add

`func Add[T Number](a, b T) T`

Add two numbers. The number type variable means this operation can be specialized to any [Number](#number) type.

```go
var n1 Int = 20
var n2 Int = 11

Add(n1, n2) // 31
```

[Back to top](#table-of-content)

## Sub

`func Sub[T Number](a, b T) T`

Subtract numbers.

```go
var n1 Int = 4
var n2 Int = 3

Sub(n1, n2) // 1
```

[Back to top](#table-of-content)

## Fdiv

`func Fdiv(a Float, b Float) Float`

Floating-point division:

```go
var n1 Float = 10
var n2 Float = 4
var n3 Float = -1

Fdiv(n1, n2) // 2.5
Fdiv(n3, n2) // -0.25
```

[Back to top](#table-of-content)

## Mul

`func Mul[T Number](a, b T) T`

Multiply numbers.

```go
var n1 Int = 2
var n2 Int = 6

Mul(n1, n2) // 6
```

[Back to top](#table-of-content)

## ToFloat

`func ToFloat[T Int](x T) Float`

ToFloat - Convert an integer into a float. Useful when mixing Int and Float values.

```go

func halfOf(num Int) Float {
  return ToFloat(number / 2)
}
```

[Back to top](#table-of-content)

## Round

`func Round(x Float) Int`

Round a number to the nearest integer.

```go

Round(1.0) // 1
Round(1.2) // 1
Round(1.5) // 2
Round(1.8) // 2
Round(-1.2) // -1
Round(-1.5) // -1
Round(-1.8) // -2
```

[Back to top](#table-of-content)

## Floor

`func Floor(x Float) Int`

Floor function, rounding down.

```go

Floor(1.0) // 1
Floor(1.2) // 1
Floor(1.5) // 1
Floor(1.8) // 1
Floor(-1.2) // -2
Floor(-1.5) // -2
Floor(-1.8) // -2
```

[Back to top](#table-of-content)

## Ceiling

`func Ceiling(x Float) Int`

Ceiling function, rounding up.

```go

Ceiling(1.0) // 1
Ceiling(1.2) // 2
Ceiling(1.5) // 2
Ceiling(1.8) // 2
Ceiling(-1.2) // -1
Ceiling(-1.5) // -1
Ceiling(-1.8) // -1
```

[Back to top](#table-of-content)

## Truncate

`func Truncate(x Float) Int`

Truncate a number, rounding towards zero

```go

Truncate(1.0) // 1
Truncate(1.2) // 1
Truncate(1.5) // 1
Truncate(1.8) // 1
Truncate(-1.2) // -1
Truncate(-1.5) // -1
Truncate(-1.8) // -1
```

[Back to top](#table-of-content)

## Eq

`func Eq[T any](x, y T) bool`

Check if values are structurally &ldquo;the same&rdquo;.

```go
var arg1 List[string] = FromSlice([]string{"a", "b"}))
var arg2 List[string] = FromSlice([]string{"a", "b"}))

Eq(arg1, arg2) // true
```

[Back to top](#table-of-content)

## Lt

`func Lt[T any](x Comparable[T], y Comparable[T]) bool`

(<)

```go
var arg1 String = String("123")
var arg2 String = String("456")

Lt(arg1, arg2) // true
```

## Gt

`func Gt[T any](x Comparable[T], y Comparable[T]) bool`

(>)

```go
var arg1 String = String("123")
var arg2 String = String("456")

Gt(arg1, arg2) // false
```

[Back to top](#table-of-content)

## Le

`func Le[T any](x Comparable[T], y Comparable[T]) bool`

(<=)

```go
var arg1 String = String("123")
var arg2 String = String("456")

Le(arg1, arg2) // True
```

[Back to top](#table-of-content)

## Ge

`func Ge[T any](x Comparable[T], y Comparable[T]) bool`

(>=)

```go
var arg1 String = String("123")
var arg2 String = String("456")

Ge(arg1, arg2) // False
```

[Back to top](#table-of-content)

## Max

`func Max[T any](x Comparable[T], y Comparable[T]) Comparable[T]`

Find the larger of two comparables.

```go
Max(42,12345678) // 12345678
Max("abc","xyz") // "xyz"
```

[Back to top](#table-of-content)

## Min

`func Min[T any](x Comparable[T], y Comparable[T]) Comparable[T]`

Find the smaller of two comparables.

```go
Min(42,12345678) // 42
Min("abc","xyz") // "abc"
```

[Back to top](#table-of-content)

## Compare

`func Compare[T any](x Comparable[T], y Comparable[T]) Order`

Compare any two comparable values. Comparable values include String, Char,
Int, Float, or a list or tuple containing comparable values. These are also the
only values that work as Dict keys or Set members.

```go

Compare (3,4) // LT
Compare(4,4) // EQ
Compare(5,4) // GT
```

[Back to top](#table-of-content)

## Order

Represents the relative ordering of two things. The relations are less than, equal to,
and greater than.

```go
LT{}
EQ{}
GT{}
```

[Back to top](#table-of-content)

## Not

`func Not(pred bool) bool`

Negate a boolean value.

```go
not(True) // false
not(False) // true
```

[Back to top](#table-of-content)

## Append

`func Append[T any](a Appendable[T], b Appendable[T]) Appendable[T]`

Put two appendable things together. This includes strings and lists.

```go
Append(String("hello"), String("world")) // "helloworld"
```

[Back to top](#table-of-content)

## ModBy

`func ModBy(modulus Int, x Int) Int`

Perform arithmetic.
A common trick is to use (n mod 2) to detect even and odd numbers:

```go
ModBy(2,0) // 0
ModBy(2,1) // 1
ModBy(2,2) // 0
ModBy(2,3) // 1
```

[Back to top](#table-of-content)

## Negate

`func Negate[A Number](n A) A`

Negate a number.

```go
Negate(42) // -42
Negate(-42) // 42
Negate(0) // 0
```

[Back to top](#table-of-content)

## Sqrt

`func Sqrt(n Float) Float`

Take the square root of a number.

```go
Sqrt(4) // 2
Sqrt(9) // 3
Sqrt(16) // 4
Sqrt(25) // 5
```

[Back to top](#table-of-content)

## Identity

`func Identity[A any](x A) A`

Given a value, returns exactly the same value. This is called the identity function.

[Back to top](#table-of-content)

## ComposeL

`func ComposeL[A any, B any, C any](g func(B) C, f func(A) B) func(A) C`

Function composition, passing results along to the left direction.

[Back to top](#table-of-content)

# Bitwise

# Char

# Dict

# List

# Maybe

# String

# Tuple
