package basics

import "fmt"

func ExampleAdd() {
	fmt.Println(Add(Int(2), Int(3)))
	// Output: 5
}

func ExampleToFloat() {
	fmt.Printf("%.1f", ToFloat(1))
	// Output: 1.0
}

func ExampleRound() {
	fmt.Println(Round(1.0))
	// Output: 1
}

func ExampleFloor() {
	fmt.Println(Floor(1.0))
	// Output: 1
}

func ExampleCeiling() {
	fmt.Println(Ceiling(1.0))
	// Output: 1
}

func ExampleTruncate() {
	fmt.Println(Truncate(1.0))
	// Output: 1
}
