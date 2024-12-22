package list

import (
	"github.com/Confidenceman02/scion-tools/pkg/basics"
	"github.com/Confidenceman02/scion-tools/pkg/tuple"
)

func ExampleEmpty() {
	Empty[int]() // []
}

func ExampleSingleton() {
	Singleton(1234) // [1234]
	Singleton("hi") // ["hi"]
}

func ExampleRepeat() {
	Repeat(2, "hi") // ["hi", "hi"]
}

func ExampleRange() {
	Range(3, 6) // [3,4,5,6]
}

func ExampleCons() {
	Cons(1, Singleton(2)) // [1,2]
}

func ExampleMap() {
	Map(basics.Sqrt, Singleton(basics.Float(4))) // [2]
}

func ExampleIndexedMap() {
	IndexedMap(tuple.Pair, FromSlice([]string{"Tom", "Sue", "Bob"})) // [(0, "Tom"),(1, "Sue"),(2, "Bob")]
}

func ExampleFoldl() {
	Foldl(basics.Add, 0, FromSlice([]basics.Int{1, 2, 3}))                         // 6
	Foldl(Cons[basics.Int], Empty[basics.Int](), FromSlice([]basics.Int{1, 2, 3})) // [3,2,1]
}

func ExampleFoldr() {
	Foldr(basics.Add, 0, FromSlice([]basics.Int{1, 2, 3}))                         // 6
	Foldr(Cons[basics.Int], Empty[basics.Int](), FromSlice([]basics.Int{1, 2, 3})) // [1,2,3]
}

func ExampleFilter() {
	xs := Range(1, 6)
	isEven := func(i basics.Int) bool { return basics.ModBy(2, i) == 0 }

	Filter(isEven, xs) // [2,4,6]
}

func ExampleLength() {
	Length(FromSlice([]basics.Int{1, 2, 3})) // 3
}

func ExampleReverse() {
	Reverse(FromSlice([]basics.Int{1, 2, 3, 4})) // [4,3,2,1]
}

func ExampleMember() {
	Member(9, FromSlice([]basics.Int{1, 2, 3, 4})) // false
	Member(4, FromSlice([]basics.Int{1, 2, 3, 4})) // true
}

func ExampleMaximum() {
	Maximum(FromSlice([]basics.Comparable[basics.Int]{basics.Int(1), basics.Int(2)})) // Just 2
	Maximum_UNSAFE(FromSlice([]basics.Int{1, 2}))                                     // Just 2
}

func ExampleMinimum() {
	Minimum(FromSlice([]basics.Comparable[basics.Int]{basics.Int(1), basics.Int(2)})) // Just 1
	Maximum_UNSAFE(FromSlice([]basics.Int{1, 2}))                                     // Just 2
}

func ExampleSum() {
	Sum(FromSlice([]basics.Int{1, 2, 3})) // 6
}

func ExampleProduct() {
	Product(FromSlice([]basics.Int{2, 2, 2})) // 8
}

func ExampleAppend() {
	x1 := FromSlice([]basics.Int{1, 1, 2})
	x2 := FromSlice([]basics.Int{3, 5, 8})
	Append(x1, x2) // [1,1,2,3,5,8]
}

func ExampleConcat() {
	x1 := FromSlice([]List[basics.Int]{FromSlice([]basics.Int{1, 2}), FromSlice([]basics.Int{3})})
	Concat(x1) // [1,2,3]
}

func ExampleIntersperse() {
	Intersperse("on", FromSlice([]string{"turtles", "turtles", "turtles"})) // ["turtles","on","turtles","on","turtles]
}

func ExampleSort() {
	Sort(FromSlice([]basics.Comparable[basics.Int]{basics.Int(3), basics.Int(2), basics.Int(1)})) // [1,2,3]
	Sort_UNSAFE(FromSlice([]basics.Int{3, 2, 1}))                                                 // [1,2,3]
}

func ExampleHead() {
	Head(FromSlice([]basics.Int{3, 2, 1})) // Just 3
}

func ExampleTail() {
	Tail(FromSlice([]basics.Int{3, 2, 1})) // Just [2,1]
}

func ExampleTake() {
	Take(2, FromSlice([]basics.Int{1, 2, 3, 4})) // Just [1,2]
}

func ExampleDrop() {
	Drop(2, FromSlice([]basics.Int{1, 2, 3, 4})) // Just [3,4]
}

func ExamplePartition() {
	Partition(func(i basics.Int) bool { return i < 3 }, FromSlice([]basics.Int{0, 1, 2, 3, 4, 5})) // ([0,1,2], [3,4,5])}
}

func ExampleUnzip() {
	zipper := FromSlice([]tuple.Tuple2[basics.Int, bool]{tuple.Pair[basics.Int](0, true), tuple.Pair[basics.Int](17, false)})
	Unzip(zipper) // ([0,17], [true,false])}
}
