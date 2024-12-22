package dict

import (
	"fmt"
	. "github.com/Confidenceman02/scion-tools/pkg/basics"
)

func ExampleEmpty() {
	Empty[int, int]()
}

func ExampleSingleton() {
	var key Int = 1
	Singleton(key, 24)
}

func ExampleInsert() {
	var key Int = 1
	Insert(key, 24, Empty[Int, Int]())
}

func ExampleRemove() {
	some_dict := Singleton(Int(1), 24)

	var key Int = 1
	Remove(key, some_dict)
}

func ExampleIsEmpty() {
	var key Int = 1
	fmt.Println(IsEmpty(Empty[int, int]()))
	fmt.Println(IsEmpty(Singleton(key, 24)))

	// Output:
	// true
	// false
}

func ExampleMember() {
	var key Int = 1
	fmt.Println(Member(key, Singleton(key, 24)))

	// Output:
	// true
}

func ExampleGet() {
	var key Int = 1
	Get(key, Singleton(key, 24)) // -> Just 24
}
