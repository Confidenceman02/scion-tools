package maybe

func ExampleWithDefault() {
	WithDefault(1, Nothing{}) // -> Just 1
}

func ExampleMap() {
	some_just := Just[int]{Value: 2}
	some_nothing := Nothing{}

	Map(func(i int) int { return i * 2 }, some_just)    // -> Just 4
	Map(func(i int) int { return i * 2 }, some_nothing) // -> Nothing
}

func ExampleMap2() {
	some_just_1 := Just[int]{Value: 2}
	some_just_2 := Just[int]{Value: 2}

	Map2(func(a int, b int) int { return a * b }, some_just_1, some_just_2) // -> Just 4
}

func ExampleMap3() {
	some_just_1 := Just[int]{Value: 2}
	some_just_2 := Just[int]{Value: 2}
	some_just_3 := Just[int]{Value: 2}

	Map3(func(a, b, c int) int { return a * b * c }, some_just_1, some_just_2, some_just_3) // -> Just 8
}

func ExampleMap4() {
	some_just_1 := Just[int]{Value: 2}
	some_just_2 := Just[int]{Value: 2}
	some_just_3 := Just[int]{Value: 2}
	some_just_4 := Just[int]{Value: 2}

	Map4(
		func(a, b, c, d int) int { return a * b * c * d },
		some_just_1,
		some_just_2,
		some_just_3,
		some_just_4) // -> Just 16
}

func ExampleMap5() {
	some_just_1 := Just[int]{Value: 2}
	some_just_2 := Just[int]{Value: 2}
	some_just_3 := Just[int]{Value: 2}
	some_just_4 := Just[int]{Value: 2}
	some_just_5 := Just[int]{Value: 2}

	Map5(
		func(a, b, c, d, e int) int { return a * b * c * d },
		some_just_1,
		some_just_2,
		some_just_3,
		some_just_4,
		some_just_5) // -> Just 32
}

func ExampleAndThen() {
	some_just := Just[int]{Value: 2}

	AndThen(func(i int) Maybe[int] { return Just[int]{Value: i * 2} }, some_just) // -> Just 4
}

func ExampleMaybeWith() {
	some_just := Just[int]{Value: 2}

	MaybeWith(some_just,
		func(j Just[int]) int { return j.Value * 2 },
		func(n Nothing) int { return 0 },
	) // -> 4
}
