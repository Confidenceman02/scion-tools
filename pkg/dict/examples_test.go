package dict

import "fmt"

func ExampleEmpty() {
	Empty[int, int]() // -> nil -- root
}

func ExampleSingleton() {
	Singleton(1, 24) // -> 1:24 -- root
	//                      / \
	//                    nil  nil
}

func ExampleInsert() {

	Insert(1, 24, Empty[int, int]()) // -> 1:24 -- root
	//                                      / \
	//                                    nil  nil
}

func ExampleRemove() {
	some_dict := Singleton(1, 24)

	Remove(1, some_dict) // -> nil -- root
}

func ExampleIsEmpty() {
	fmt.Println(IsEmpty(Empty[int, int]()))
	fmt.Println(IsEmpty(Singleton(1, 24)))

	// Output:
	// true
	// false
}

func ExampleMember() {
	fmt.Println(Member(1, Singleton(1, 24)))

	// Output:
	// true
}

func ExampleGet() {
	Get(1, Singleton(1, 24)) // -> Just{Value: 24}
}
