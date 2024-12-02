package basics

func ExampleAdd() {
	Add(Int(2), Int(3)) // 5
}

func ExampleSub() {
	Sub(Int(4), Int(3)) // 1
}

func ExampleMul() {
	Mul(Int(2), Int(3)) // 6
}

func ExampleFdiv() {
	Fdiv(Float(10), Float(4)) // 2.5
}

func ExampleToFloat() {
	ToFloat(1) // 1.0
}

func ExampleRound() {
	Round(1.0) // 1
}

func ExampleFloor() {
	Floor(1.0) // 1
}

func ExampleCeiling() {
	Ceiling(1.0) // 1
}

func ExampleTruncate() {
	Truncate(1.0) // 1
}

func ExampleEq() {
	Eq("Hello", "Hello") // true
}

func ExampleLt() {
	Lt(Int(2), Int(3)) // true
}

func ExampleGt() {
	Gt(Int(2), Int(3)) // false
}

func ExampleLe() {
	Le(Int(2), Int(2)) // true
}

func ExampleGe() {
	Ge(Int(2), Int(2)) // true
}

func ExampleMax() {
	Max(Int(42), Int(12345678)) // 12345678
}

func ExampleMin() {
	Min(Int(42), Int(12345678)) // 42
}

func ExampleCompare() {
	Compare(Int(3), Int(4)) // Order => LT{}
}

func ExampleNot() {
	Not(true)  // false
	Not(false) // true
}

func ExampleModBy() {
	ModBy(Int(2), Int(2)) // 0
}

func ExampleNegate() {
	Negate(Int(42)) // -42
}

func ExampleSqrt() {
	Sqrt(Float(4))  // 2
	Sqrt(Float(9))  // 3
	Sqrt(Float(16)) // 4
	Sqrt(Float(25)) // 5
}

func ExampleIdentity() {
	Identity(Float(4)) // 4
}

func ExampleComposeL() {
	composed := ComposeL(func(i Int) Int { return i + 1 }, Identity)
	composed(1) // 2
}
