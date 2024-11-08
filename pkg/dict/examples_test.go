package dict

import (
	"fmt"
	. "github.com/Confidenceman02/scion-tools/pkg/basics"
)

func ExampleEmpty() {
	Empty[int, int]() // -> nil -- root
}

func ExampleSingleton() {
	Singleton(Int(1), 24) // -> 1:24 -- root
	//                      / \
	//                    nil  nil
}

func ExampleInsert() {

	Insert(Int(1), 24, Empty[Int, Int]()) // -> 1:24 -- root
	//                                      / \
	//                                    nil  nil
}

func ExampleRemove() {
	some_dict := Singleton(Int(1), 24)

	Remove(Int(1), some_dict) // -> nil -- root
}

func ExampleIsEmpty() {
	fmt.Println(IsEmpty(Empty[int, int]()))
	fmt.Println(IsEmpty(Singleton(Int(1), 24)))

	// Output:
	// true
	// false
}

func ExampleMember() {
	fmt.Println(Member(Int(1), Singleton(Int(1), 24)))

	// Output:
	// true
}

func ExampleGet() {
	Get(Int(1), Singleton(Int(1), 24)) // -> Just{Value: 24}
}
