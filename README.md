# scion-tools

[![Actions status](https://Confidenceman02.github.io/djelm/workflows/CI/badge.svg)](https://github.com/Confidenceman02/scion-tools/actions)
[![](https://img.shields.io/badge/license-MIT-blue)](https://github.com/Confidenceman02/scion-tools/blob/main/LICENSE)

> [!NOTE]
> There are some missing packages and functions from the Elm core library, You can see a list of these
> in the [issues](https://github.com/Confidenceman02/scion-tools/issues).

# Elm inspired functional programing in Go

---

## Table of Content

> [!NOTE]
> For brevity we will use `[]` to signify a `List[T]` and `(,)` to signify a `Tuple[A,B]`.
>
> You can use the `FromSlice` utility to create a `List[T]` from a slice.

```
// example
[1,2,3] => FromSlice([]int{1,2,3})
(true,1) => Pair(true,1)
```

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
    <ul>
        <li>
            <a href="#and">And</a>
        </li>
        <li>
            <a href="#shiftrightby">ShiftRightBy</a>
        </li>
    </ul>
  </details>
- <details>
    <summary><a href="#char">Char</a></summary>
    <ul>
        <li>
            <a href="#isdigit">IsDigit</a>
        </li>
    </ul>
  </details>
- <details>
    <summary><a href="#dict">Dict</a></summary>
    <ul>
        <li>
            <a href="#empty">Empty</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#singleton">Singleton</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#remove">Remove</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#Insert">Insert</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#Remove">Remove</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#IsEmpty">IsEmpty</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#memberdict">Member</a>
        </li>
    </ul>
  </details>
- <details>
    <summary><a href="#list">List</a></summary>
    <ul>
        <li>
            <a href="#emptylist">Empty</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#Singletonlist">Singleton</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#repeat">Repeat</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#range">Range</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#cons">Cons</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#map">Map</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#foldl">Foldl</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#foldr">Foldr</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#filter">Filter</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#filtermap">FilterMap</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#length">Length</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#reverse">Reverse</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#memberlist">Member</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#all">All</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#any">Any</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#maximum">Maximum</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#maximumunsafe">Maximum_UNSAFE</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#minimum">Minimum</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#sum">Sum</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#product">Product</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#append">Append</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#concat">Concat</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#concatmap">ConcatMap</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#intersperse">Intersperse</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#map2">Map2</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#map3">Map3</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#map4">Map4</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#map5">Map5</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#sort">Sort</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#sortunsafe">Sort_UNSAFE</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#sortby">SortBy</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#sortwith">SortWith</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#isempty">IsEmpty</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#head">Head</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#tail">Tail</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#take">Take</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#drop">Drop</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#partition">Partition</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#unzip">Unzip</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#listwith">ListWith</a>
        </li>
    </ul>
  </details>
- <details>
    <summary><a href="#maybe">Maybe</a></summary>
    <ul>
        <li>
            <a href="#withdefault">WithDefault</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#map">Map</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#map2">Map2</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#map3">Map3</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#map4">Map4</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#map5">Map5</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#andthen">AndThen</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#maybewith">MaybeWith</a>
        </li>
    </ul>
  </details>
- <details>
    <summary><a href="#string">String</a></summary>
    <ul>
        <li>
            <a href="#isemptystring">IsEmpty</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#lengthstring">Length</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#reverse">Reverse</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#repeat">Repeat</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#replace">Replace</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#appendstring">Append</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#concat">Concat</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#split">Split</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#join">Join</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#words">Words</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#lines">Lines</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#slice">Slice</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#left">Left</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#right">Right</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#dropleft">DropLeft</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#dropright">DropRight</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#contains">Contains</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#startswith">StartsWith</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#endswith">EndsWith</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#indexes">Indexes</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#indices">Indices</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#toint">ToInt</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#fromint">FromInt</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#tofloat">ToFloat</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#fromfloat">FromFloat</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#fromchar">FromChar</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#consstring">Cons</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#uncons">Uncons</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#tolist">ToList</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#fromlist">FromList</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#toupper">ToUpper</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#tolower">ToLower</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#pad">Pad</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#padleft">PadLeft</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#padright">PadRight</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#trim">Trim</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#trimleft">TrimLeft</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#trimright">TrimRight</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#mapstring">Map</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#filter">Filter</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#foldl">Foldl</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#foldr">Foldr</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#any">Any</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#allstring">All</a>
        </li>
    </ul>
  </details>
- <details>
    <summary><a href="#tuple">Tuple</a></summary>
    <ul>
        <li>
            <a href="#pair">Pair</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#first">First</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#second">Second</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#mapfirst">MapFirst</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#mapsecond">MapSecond</a>
        </li>
    </ul>
    <ul>
        <li>
            <a href="#mapboth">MapBoth</a>
        </li>
    </ul>
  </details>

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

- Function composition
- Immutability
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
not(true) // false
not(false) // true
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

```go
isEven := func(i Float) bool { return ModBy(2, Int(i)) == 0 }
composed := ComposeL(isEven, Sqrt)


composed(4) // true
composed(3) // false
```

[Back to top](#table-of-content)

# Bitwise

```go
import "github.com/Confidenceman02/scion-tools/pkg/bitwise"
```

Package for bitwise operations.

## And

`func And(a, b Int) Int`

Bitwise AND

[Back to top](#table-of-content)

## ShiftRightBy

`func ShiftRightBy(offset Int, a Int) Int`

Shift bits to the right by a given offset, filling new bits with whatever is the topmost bit.
This can be used to divide numbers by powers of two.

```go
ShiftRightBy(1,32) // 16
ShiftRightBy(2,32) // 8
ShiftRightBy(1,-32) // -16
```

[Back to top](#table-of-content)

# Char

```go
import "github.com/Confidenceman02/scion-tools/pkg/char"
```

Functions for working with runes. Rune literals are enclosed in 'a' pair of single quotes.

## IsDigit

`func IsDigit(c Char) bool`

Detect digits 0123456789

```go
isDigit('0') // True
isDigit('1') // True
isDigit('9') // True
isDigit('a') // False
isDigit('b') // False
isDigit('A') // False
```

[Back to top](#table-of-content)

# Dict

```go
import "github.com/Confidenceman02/scion-tools/pkg/dict"
```

A dictionary mapping unique keys to values. The keys can be any [comparable](#comparable) type.
This includes Int, Float, Char, String, and tuples or lists of comparable types.
Insert, remove, and query operations all take O(log n) time.

## Empty

`func Empty[K any, V any]() Dict[Comparable[K], V]`

Create an empty dictionary.

```go
Empty()
```

[Back to top](#table-of-content)

## Singleton

`func Singleton[K any, V any](key Comparable[K], value V) Dict[Comparable[K], V]`

Create a dictionary with one key-value pair.

```go
Singleton("hello", "world")
```

[Back to top](#table-of-content)

## Insert

`func Insert[K any, V any](key Comparable[K], v V, d Dict[Comparable[K], V]) Dict[Comparable[K], V]`

Insert a key-value pair into a dictionary. Replaces value when there is a collision.

```go
Insert(2, "two", Singleton(1,"one"))
```

[Back to top](#table-of-content)

## Remove

`func Remove[K any, V any](key Comparable[K], d Dict[Comparable[K], V]) Dict[Comparable[K], V]`

Remove a key-value pair from a dictionary. If the key is not found, no changes are made.

```go
Remove(1, Singleton(1,"one"))
```

[Back to top](#table-of-content)

## IsEmpty

`func IsEmpty[K any, V any](d Dict[Comparable[K], V]) bool`

Determine if a dictionary is empty.

```go
IsEmpty(Empty()) // true
```

[Back to top](#table-of-content)

## Member(Dict)

`func Member[K any, V any](k Comparable[K], d Dict[Comparable[K], V]) bool`

Determine if a key is in a dictionary.

```go
Member(1, Singleton(1, "one")) // true
```

[Back to top](#table-of-content)

# List

```go
import "github.com/Confidenceman02/scion-tools/pkg/list"
```

You can create a List from Go slices with the `FromSlice` function. This module has a bunch of functions to help you work with them!

## Empty

`func Empty[T any]() List[T]`

Create a list with no elements.

```go
Empty() // []
```

## Singleton

`func Singleton[T any](val T) List[T]`

Create a list with only one element.

```go
Singleton(1234) // [1234]
```

[Back to top](#table-of-content)

## Repeat

`func Repeat[T any](n basics.Int, val T) List[T]`

Create a list with _n_ copies of a value.

```go
Repeat(3,1) // [1,1,1]
```

[Back to top](#table-of-content)

## Range

`func Range(low basics.Int, hi basics.Int) List[basics.Int]`

Create a list of numbers, every element increasing by one. You give the lowest and highest number that should be in the list.

```go
Range(3,6) // [3, 4, 5, 6]
Range(3,3) // [3]
Range(6,3) // []
```

[Back to top](#table-of-content)

## Cons

`func Cons[T any](val T, l List[T]) List[T]`

Add an element to the front of a list.

```go
Cons(1,Singleton(2)) // [1,2]
Cons(1,Empty())      // [1]
```

[Back to top](#table-of-content)

## Map

`func Map[A any, B any](f func(A) B, xs List[A]) List[B]`

Apply a function to every element of a list.

```go
Map(Sqrt, [1,4,9])          // [1,2,3]
Map(Not, [true,false,true]) // [false,true,false]
```

[Back to top](#table-of-content)

## IndexedMap

`func IndexedMap[A any, B any](f func(basics.Int, A) B, xs List[A]) List[B]`

Same as map but the function is also applied to the index of each element (starting at zero).

```go
IndexedMap(Tuple.pair, ["Tom","Sue","Bob"]) // [(0,"Tom"),(1,"Sue"),(2,"Bob")]
```

[Back to top](#table-of-content)

## Foldl

`func Foldl[A any, B any](f func(A, B) B, acc B, ls List[A]) B`

Reduce a list from the left.

```go
Foldl(Cons,Empty(), [1,2,3]) // [3,2,1]
```

[Back to top](#table-of-content)

## Foldr

`func Foldr[A any, B any](fn func(A, B) B, acc B, ls List[A]) B`

Reduce a list from the right.

```go
Foldr(Cons,Empty(),[1,2,3]) == [1,2,3]
```

[Back to top](#table-of-content)

## Filter

`func Filter[T any](isGood func(T) bool, list List[T]) List[T]`

Keep elements that satisfy the test.

```go
Filter(isEven, [1,2,3,4,5,6]) // [2,4,6]
```

[Back to top](#table-of-content)

## FilterMap

`func FilterMap[A any, B any](f func(A) maybe.Maybe[B], xs List[A]) List[B]`

Filter out certain values. For example, maybe you have a bunch of strings from an
untrusted source and you want to turn them into numbers:

```go
func numbers() List[Int] {
  return FilterMap(ToInt, ["3","hi","12","4th","May"])
}

// numbers == [3, 12]
```

[Back to top](#table-of-content)

## Length

`func Length[T any](ls List[T]) basics.Int`

Determine the length of a list.

```go
Length([1,2,3]) // 3
```

[Back to top](#table-of-content)

## Reverse

`func Reverse[T any](ls List[T]) List[T]`

Reverse a list.

```go
Reverse([1,2,3,4]) // [4,3,2,1]
```

[Back to top](#table-of-content)

## Member(List)

`func Member[T any](val T, l List[T]) bool`

Figure out whether a list contains a value.

```go
Member(9, [1,2,3,4]) // false
```

[Back to top](#table-of-content)

## All

`func All[T any](isOkay func(T) bool, l List[T]) bool`

Determine if all elements satisfy some test.

```go
All(isEven, [2,4]) // true
```

[Back to top](#table-of-content)

## Any

`func Any[T any](isOkay func(T) bool, l List[T]) bool`

Determine if any elements satisfy some test.

```go
Any(isEven, [2,3]) // true
```

[Back to top](#table-of-content)

## Maximum

`func Maximum[T any](xs List[basics.Comparable[T]]) maybe.Maybe[T]`

Find the maximum element in a non-empty list.

```go
Maximum([1,4,2]) == Just 4
```

[Back to top](#table-of-content)

## Maximum_UNSAFE

`func Maximum[T any](xs List[basics.Comparable[T]]) maybe.Maybe[T]`

Find the maximum element in a non-empty list of generic elements.
If elements do not implement the [Comparable[T]](#comparable) interface the function will panic.

```go
Maximum([1,4,2]) == Just 4
```

[Back to top](#table-of-content)

## Minimum

`func Minimum[T any](xs List[basics.Comparable[T]]) maybe.Maybe[T]`

Find the minimum element in a non-empty list.

```go
Minimum([1,4,2]) == Just 1
```

[Back to top](#table-of-content)

## Sum

`func Sum[T basics.Number](xs List[T]) T`

Get the sum of the list elements.

```go
Sum([1,2,3]) // 6
```

[Back to top](#table-of-content)

## Product

`func Product[T basics.Number](xs List[T]) T`

Get the product of the list elements.

```go
Product([2,2,2]) == 8
```

[Back to top](#table-of-content)

## Append

`func Append[T any](xs List[T], ys List[T]) List[T]`

Put two lists together.

```go
Append(['a','b'], ['c']) // ['a','b','c']
```

[Back to top](#table-of-content)

## Concat

`func Concat[T any](lists List[List[T]]) List[T]`

Concatenate a bunch of lists into a single list:

```go
Concat([[1,2], [3], [4,5]]) // [1,2,3,4,5]
```

[Back to top](#table-of-content)

## Intersperse

`func Intersperse[T any](sep T, xs List[T]) List[T]`

Places the given value between all members of the given list.

```go
Intersperse("on",["turtles","turtles","turtles"]) == ["turtles","on","turtles","on","turtles"]
```

[Back to top](#table-of-content)

## Map2

`func Map2[A any, B any, result any](f func(A, B) result, xs List[A], ys List[B]) List[result]`

// Combine two lists, combining them with the given function.
// If one list is longer, the extra elements are dropped.

[Back to top](#table-of-content)

## Map3

`func Map3[A any, B any, C any, result any](f func(A, B, C) result, xs List[A], ys List[B], zs List[C]) List[result]`

[Back to top](#table-of-content)

## Map4

`func Map4[A any, B any, C any, D any, result any](f func(A, B, C, D) result, xs List[A], ys List[B], zs List[C], ws List[D]) List[result]`

[Back to top](#table-of-content)

## Map5

`func Map5[A any, B any, C any, D any, E any, result any](f func(A, B, C, D, E) result, vs List[A], ws List[B], xs List[C], ys List[D], zs List[E]) List[result]`

[Back to top](#table-of-content)

## Sort

`func Sort[T any](xs List[basics.Comparable[T]]) List[basics.Comparable[T]]`

Sort values from lowest to highest.

```go
Sort([3,1,5]) == [1,3,5]
```

[Back to top](#table-of-content)

## Sort_UNSAFE

`func Sort_UNSAFE[T any](xs List[T]) List[T]`

Sort generic values from lowest to highest.
This function will panic if T is not a Comparable[T]

```go
Sort_UNSAFE([3,1,5]) == [1,3,5]
```

[Back to top](#table-of-content)

## SortBy

`func SortBy[A any](f func(A) basics.Comparable[A], xs List[A]) List[A]`

Sort values by a derived property.

```go
SortBy(String.length,["mouse","cat"]) // ["cat","mouse"]
```

[Back to top](#table-of-content)

## SortBy_UNSAFE

`func SortBy_UNSAFE[A any](f func(A) A, xs List[A]) List[A]`

Sort values by a derived property.
This function will panic if the passed in func doesn't return a Comparable.

```go
SortBy_UNSAFE(String.length,["mouse","cat"]) // ["cat","mouse"]
```

[Back to top](#table-of-content)

## SortWith

`func SortWith[A any](f func(a A, b A) basics.Order, xs List[A]) List[A]`

Sort values with a custom comparison function.

[Back to top](#table-of-content)

## IsEmpty

`func IsEmpty[T any](l List[T]) bool`

Determine if a list is empty.

```go
IsEmpty([]) // True
```

[Back to top](#table-of-content)

## Head

`func Head[T any](l List[T]) maybe.Maybe[T]`

Extract the first element of a list.

```go
Head([1,2,3]) // Just 1
```

[Back to top](#table-of-content)

## Tail

`func Tail[T any](l List[T]) maybe.Maybe[List[T]]`

Extract the rest of the list.

```go
Tail([1,2,3]) // Just [2,3]
```

[Back to top](#table-of-content)

## Take

`func Take[T any](n basics.Int, list List[T]) List[T]`

Take the first n members of a list.

```go
Take(2,[1,2,3,4]) == [1,2]
```

[Back to top](#table-of-content)

## Drop

`func Drop[T any](n basics.Int, list List[T]) List[T]`

Drop the first n members of a list.

```go
Drop(2,[1,2,3,4]) == [3,4]
```

[Back to top](#table-of-content)

## Partition

`func Partition[A any](pred func(A) bool, list List[A]) Tuple2[List[A], List[A]]`

// Partition a list based on some test. The first list contains all values
// that satisfy the test, and the second list contains all the value that do not.

```go
Partition(isEven,[0,1,2,3,4,5]) // ([0,2,4], [1,3,5])
```

[Back to top](#table-of-content)

## Unzip

`func Unzip[A, B any](pairs List[Tuple2[A, B]]) Tuple2[List[A], List[B]]`

Decompose a list of tuples into a tuple of lists.

```go
Unzip([(0, true),(17,false),(1337,true)]) // ([0,17,1337], [true,false,true])
```

[Back to top](#table-of-content)

# Maybe

```go
import "github.com/Confidenceman02/scion-tools/pkg/maybe"
```

Represent values that may or may not exist. It can be useful if you have a record field
that is only filled in sometimes. Or if a function takes a value sometimes, but does
not absolutely need it.

```go
type Just[A any]
type Nothing
```

## WithDefault

`func WithDefault[A any](a A, m Maybe[A]) A`

Provide a default value, turning an optional value into a normal value.
This comes in handy when paired with functions like [Dict.Get](#get) which gives back a Maybe.

```go
WithDefault(100,Just[int]{42}) // 42
```

[Back to top](#table-of-content)

## Map

`func Map[A any, B any](f func(A) B, m Maybe[A]) Maybe[B]`

Transform a Maybe value with a given function:

```go
Map(Sqrt,(Just[Int]{9})) // Just 3
```

[Back to top](#table-of-content)

## Map2

`func Map2[A any, B any, value any](f func(a A, b B) value, m1 Maybe[A], m2 Maybe[B]) Maybe[value]`

Apply a function if all the arguments are Just a value.

```go
Map2(Add,Just[Int]{3}, Just[Int]{4}) // Just 7
```

[Back to top](#table-of-content)

## Map3

`func Map3[A any, B any, C any, value any](f func(a A, b B, c C) value, m1 Maybe[A], m2 Maybe[B], m3 Maybe[C]) Maybe[value]`

[Back to top](#table-of-content)

## Map4

`func Map4[A any, B any, C any, D any, value any](f func(a A, b B, c C, d D) value, m1 Maybe[A], m2 Maybe[B], m3 Maybe[C], m4 Maybe[D]) Maybe[value]`

[Back to top](#table-of-content)

## Map5

`func Map5[A any, B any, C any, D any, E any, value any](f func(a A, b B, c C, d D, e E) value, m1 Maybe[A], m2 Maybe[B], m3 Maybe[C], m4 Maybe[D], m5 Maybe[E]) Maybe[value]`

[Back to top](#table-of-content)

## AndThen

`func AndThen[A any, B any](f func(A) Maybe[B], m Maybe[A]) Maybe[B]`

Chain together many computations that may fail.

[Back to top](#table-of-content)

## MaybeWith

`func MaybeWith[V any, R any](m Maybe[V],j func(Just[V]) R,n func(Nothing) R) R`

Provide functions for a Maybe's Just and Nothing variants

[Back to top](#table-of-content)

# String

```go
import "github.com/Confidenceman02/scion-tools/pkg/string"
```

A built-in representation for efficient string manipulation.
The String type is a wrapper for Go's `string`.

```go
type String string
```

## IsEmpty(String)

`func IsEmpty(x String) bool`

Determine if a string is empty.

```go
IsEmpty("") // true
```

[Back to top](#table-of-content)

## Length(String)

`func Length(x String) basics.Int`

Get the length of a string.

```go
Length("innumerable") // 11
```

[Back to top](#table-of-content)

## Reverse(String)

`func Reverse(x String) String`

Reverse a string.

```go
Reverse("stressed") // "desserts"
```

[Back to top](#table-of-content)

## Repeat(String)

`func Repeat(n basics.Int, chunk String) String`

Repeat a string n times.

```go
Repeat(3,"ha") // "hahaha"
```

[Back to top](#table-of-content)

## Replace

`func Replace(before String, after String, str String) String`

Replace all occurrences of some substring.

```go
Replace(",","/","a,b,c,d,e")           == "a/b/c/d/e"
```

[Back to top](#table-of-content)

## Append

`func Append(x String, y String) String`

Append two strings. You can also use basics.Append to do this.

```go
Append("butter","fly") // "butterfly"
```

[Back to top](#table-of-content)

## Concat

`func Concat(chunks list.List[String]) String`

Concatenate many strings into one.

```go
Concat(["never","the","less"] ) // "nevertheless"
```

[Back to top](#table-of-content)

## Split

`func Split(sep String, s String) list.List[String]`

Split a string using a given separator.

```go
Split(",","cat,dog,cow") // ["cat","dog","cow"]
```

[Back to top](#table-of-content)

## Join

`func Join(sep String, chunks list.List[String]) String`

Put many strings together with a given separator.

```go
Join("a",["H","w","ii","n"]) // "Hawaiian"
```

[Back to top](#table-of-content)

## Words

`func Words(str String) list.List[String]`

Break a string into words, splitting on chunks of whitespace.

```go
Words("How are \t you? \n Good?") // ["How","are","you?","Good?"]
```

[Back to top](#table-of-content)

## Lines

`func Lines(str String) list.List[String]`

Break a string into lines, splitting on newlines.

```go
Lines("How are you?\nGood?") // ["How are you?", "Good?"]
```

[Back to top](#table-of-content)

## Slice

`func Slice(start basics.Int, end basics.Int, str String) String`

Take a substring given a start and end index. Negative indexes are taken starting from the end of the list.

```go
Slice(7,9,"snakes on a plane!") // "on"
```

[Back to top](#table-of-content)

## Left

`func Left(n basics.Int, str String) String`

Take _n_ characters from the left side of a string.

```go
Left(2,"Mulder") // "Mu"
```

[Back to top](#table-of-content)

## Right

`func Right(n basics.Int, str String) String`

Take _n_ characters from the right side of a string.

```go
Right(2,"Scully") // "ly"
```

[Back to top](#table-of-content)

## DropLeft

`func DropLeft(n basics.Int, str String) String`

Drop _n_ characters from the left side of a string.

```go
DropLeft(2,"The Lone Gunmen") // "e Lone Gunmen"
```

[Back to top](#table-of-content)

## DropRight

`func DropRight(n basics.Int, str String) String`

Drop _n_ characters from the right side of a string.

```go
DropRight(2,"Cigarette Smoking Man") // "Cigarette Smoking M"
```

[Back to top](#table-of-content)

## Contains

`func Contains(sub String, str String) bool`

See if the second string contains the first one.

```go
Contains("the","theory") // true
```

[Back to top](#table-of-content)

## StartsWith

`func StartsWith(sub String, str String) bool`

See if the second string starts with the first one.

```go
StartsWith("the","theory") // true
```

[Back to top](#table-of-content)

## EndsWith

`func EndsWith(sub String, str String) bool`

See if the second string ends with the first one.

```go
EndsWith("the","theory") // false
```

[Back to top](#table-of-content)

## Indexes

`func Indexes(sub String, str String) list.List[basics.Int]`

Get all of the indexes for a substring in another string.

```go
Indexes("i","Mississippi") // [1,4,7,10]
```

[Back to top](#table-of-content)

## Indices

`func Indices(sub String, str String) list.List[basics.Int]`

Alias for `indexes`.

```go
Indexes("i","Mississippi") // [1,4,7,10]
```

[Back to top](#table-of-content)

## ToInt

`func ToInt(x String) maybe.Maybe[basics.Int]`

Try to convert a string into an int, failing on improperly formatted strings.

```go
ToInt("123") // Just 123
```

[Back to top](#table-of-content)

## FromInt

`func FromInt(x basics.Int) String`

Convert an Int to a String.

```go
FromInt(123) // "123"
```

[Back to top](#table-of-content)

## ToFloat

`func ToFloat(x String) maybe.Maybe[basics.Float]`

Try to convert a string into a float, failing on improperly formatted strings.

```go
ToFloat("123") // Just 123.0
```

[Back to top](#table-of-content)

## FromFloat

`func FromFloat(x basics.Float) String`

Convert a Float to a String.

```go
FromFloat(123) // "123"
```

[Back to top](#table-of-content)

## FromChar

`func FromChar(char char.Char) String`

Create a string from a given character.

```go
FromFloat(123) // "123"
```

[Back to top](#table-of-content)

## Cons(String)

`func Cons(char char.Char, str String) String`

Add a character to the beginning of a string.

```go
Cons('T',"he truth is out there") // "The truth is out there"
```

[Back to top](#table-of-content)

## Uncons

`func Uncons(str String) maybe.Maybe[tuple.Tuple2[char.Char, String]]`

Split a non-empty string into its head and tail. This lets you pattern match on strings exactly as you would with lists.

```go
Uncons("abc") // Just ('a',"bc")
```

[Back to top](#table-of-content)

## ToList

`func ToList(str String) list.List[char.Char]`

Convert a string to a list of characters.

```go
ToList("abc") // ['a','b','c']
```

[Back to top](#table-of-content)

## FromList

`func FromList(chars list.List[char.Char]) String`

Convert a list of characters into a String. Can be useful if you want to create a string primarily by consing, perhaps for decoding something.

```go
FromList(['a','b','c']) // "abc"
```

[Back to top](#table-of-content)

## ToUpper

`func ToUpper(str String) String`

Convert a string to all upper case. Useful for case-insensitive comparisons and VIRTUAL YELLING.

```go
ToUpper("skinner") // "SKINNER"
```

[Back to top](#table-of-content)

## ToLower

`func ToLower(str String) String {`

Convert a string to all lower case. Useful for case-insensitive comparisons.

```go
ToLower("X-FILES") // "x-files"
```

[Back to top](#table-of-content)

## Pad

`func Pad(n basics.Int, char char.Char, str String) String`

Pad a string on both sides until it has a given length.

```go
Pad(5,' ',"1" ) == "  1  "
```

[Back to top](#table-of-content)

## PadLeft

`func PadLeft(n basics.Int, char char.Char, str String) String {`

Pad a string on the left until it has a given length.

```go
PadLeft(5,'.',"1") == "....1"
```

[Back to top](#table-of-content)

## PadRight

`func PadRight(n basics.Int, char char.Char, str String) String`

Pad a string on the right until it has a given length.

```go
PadRight(5,'.',"1") // "1...."
```

[Back to top](#table-of-content)

## Trim

`func Trim(str String) String`

Get rid of whitespace on both sides of a string.

```go
Trim("  hats  \n") // "hats"
```

[Back to top](#table-of-content)

## TrimLeft

`func TrimLeft(str String) String`

Get rid of whitespace on the left of a string.

```go
TrimLeft("  hats  \n") // "hats  \n"
```

[Back to top](#table-of-content)

## TrimRight

`func TrimRight(str String) String`

Get rid of whitespace on the right of a string.

```go
TrimRight("  hats  \n") // "  hats"
```

[Back to top](#table-of-content)

## Map(String)

`func Map(f func(char.Char) char.Char, str String) String`

Transform every character in a string

[Back to top](#table-of-content)

## Filter

`func Filter(isGood func(char.Char) bool, str String) String`

Keep only the characters that pass the test.

```go
Filter(IsDigit,"R2-D2") // "22"
```

[Back to top](#table-of-content)

## Foldl(String)

`func Foldl[B any](f func(char.Char, B) B, state B, str String) B`

Reduce a string from the left.

```go
Foldl(Cons,"","time") // "emit"
```

[Back to top](#table-of-content)

## Foldr(String)

`func Foldr[B any](f func(char.Char, B) B, state B, str String) B`

Reduce a string from the right.

```go
Foldr(cons,"","time") // "time"
```

[Back to top](#table-of-content)

## Any(String)

`func Any(isGood func(char.Char) bool, str String) bool`

Determine whether any characters pass the test.

```go
Any(IsDigit,"90210") // true
```

[Back to top](#table-of-content)

## All(String)

`func All(isGood func(char.Char) bool, str String) bool`

Determine whether all characters pass the test.

```go
All(IsDigit,"90210") == true
```

[Back to top](#table-of-content)

# Tuple

```go
import "github.com/Confidenceman02/scion-tools/pkg/tuple"
```

This package is a bunch of helpers for working with 2-tuples.

## Pair

`func Pair[A, B any](a A, b B) Tuple2[A, B]`

Create a 2-tuple.

```go
Pair(3,4) // (3,4)
```

[Back to top](#table-of-content)

## First

`func First[A, B any](t Tuple2[A, B]) A`

Extract the first value from a tuple.

```go
First((3,4)) // 3
```

[Back to top](#table-of-content)

## Second

`func Second[A, B any](t Tuple2[A, B]) B`

Extract the second value from a tuple.

```go
Second((3, 4)) // 4
```

[Back to top](#table-of-content)

## MapFirst

`func MapFirst[A, B, C any](f func(A) B, t Tuple2[A, C]) Tuple2[B, C]`

Transform the first value in a tuple.

```go
MapFirst(string.reverse,("stressed", 16)) // ("desserts", 16)
```

[Back to top](#table-of-content)

## MapSecond

`func MapSecond[A, B, C any](f func(B) A, t Tuple2[C, B]) Tuple2[C, A]`

Transform the second value in a tuple.

```go
MapSecond(Sqrt,("stressed", 16)) // ("stressed", 4)
```

[Back to top](#table-of-content)

## MapBoth

`func MapBoth[A, B, C, D any](f func(A) C, g func(B) D, t Tuple2[A, B]) Tuple2[C, D]`

Transform both parts of a tuple.

```go
MapBoth(string.reverse,Sqrt,("stressed", 16)) // ("desserts", 4)
```

[Back to top](#table-of-content)
