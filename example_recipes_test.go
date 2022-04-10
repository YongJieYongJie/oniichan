package oniichan_test

import (
	"fmt"

	"golang.org/x/exp/constraints"

	. "github.com/YongJieYongJie/oniichan"
)

// printChannel converts a channel to slice before printing it on a line.
func printChannel[T any](c <-chan T) {
	fmt.Println(Slice(c))
}

// Below are examples showing how the various functions provided by oniichan
// can be used as building blocks for other commonly-required operations.
func Example_recipes() {
	// Refer to the output comment below for the output.

	fmt.Print("Example 1: ")
	fmt.Println(`Take(ChanFrom("a", "b", "c"), 2)`)
	printChannel(Take(ChanFrom("a", "b", "c"), 2))

	fmt.Printf("\nExample 2: ")
	fmt.Println(`Prepend(4.4, ChanFrom(1.1, 2.2, 3.3))`)
	printChannel(Prepend(4.4, ChanFrom(1.1, 2.2, 3.3)))

	fmt.Printf("\nExample 3: ")
	fmt.Println(`Take(Tabulate(square, 4), 3)`)
	square := func(x int) int { return x * x }
	printChannel(Take(Tabulate(square, 4), 3))

	fmt.Printf("\nExample 4: ")
	fmt.Println("Consume(ChanFrom(1, 2, 3), 1)")
	chan123 := ChanFrom(1, 2, 3)
	Consume(chan123, 1)
	printChannel(chan123)

	fmt.Printf("\nExample 5: ")
	fmt.Println(`Nth(ChanFrom("A", "B", "C"), 1)`)
	fmt.Println(Nth(ChanFrom("A", "B", "C"), 1))

	fmt.Printf("\nExample 6: ")
	fmt.Println(`AllEqual(ChanFrom("A", "A", "A"))`)
	fmt.Println(AllEqual(ChanFrom("A", "A", "A")))

	isOdd := func(x int) int { return x % 2 }
	fmt.Printf("\nExample 6: ")
	fmt.Println(`Quantify(ChanFrom(1, 2, 3), isOdd)`)
	fmt.Println(Quantify(ChanFrom(1, 2, 3), isOdd))

	fmt.Printf("\nExample 7: ")
	fmt.Println(`DotProduct(ChanFrom(1, 2, 3), ChanFrom(2, 3, 4))`)
	printChannel(DotProduct(ChanFrom(1, 2, 3), ChanFrom(2, 3, 4)))

	// output:
	// Example 1: Take(ChanFrom("a", "b", "c"), 2)
	// [a b]
	//
	// Example 2: Prepend(4.4, ChanFrom(1.1, 2.2, 3.3))
	// [4.4 1.1 2.2 3.3]
	//
	// Example 3: Take(Tabulate(square, 4), 3)
	// [16 25 36]
	//
	// Example 4: Consume(ChanFrom(1, 2, 3), 1)
	// [2 3]
	//
	// Example 5: Nth(ChanFrom("A", "B", "C"), 1)
	// A
	//
	// Example 6: AllEqual(ChanFrom("A", "A", "A"))
	// true
	//
	// Example 6: Quantify(ChanFrom(1, 2, 3), isOdd)
	// 2
	//
	// Example 7: DotProduct(ChanFrom(1, 2, 3), ChanFrom(2, 3, 4))
	// [2 6 12]
}

// Take returns the first n items of the input channel as a list.
func Take[T any](in <-chan T, n int) <-chan T {
	return SliceFromTo(in, 0, n)
}

// Prepend returns a channel with the additional element prepended.
func Prepend[T any](elem T, in <-chan T) <-chan T {
	return Chain(ChanFrom(elem), in)
}

// Tabulate returns a channel producing function f applied to successively
// integer values starting from the start value as specified; i.e., f(start), f(start+1), ...
func Tabulate[T constraints.Integer, R any](f func(T) R, start T) <-chan R {
	return Map(f, Count(start))
}

// Consume advances the input channel n-steps ahead.
func Consume[T any](in <-chan T, n int) {
	for range SliceTo(in, n) {
	}
}

// Nth returns the nth item from the input channel or a default value.
func Nth[T any](in <-chan T, n int, default_ ...T) T {
	if len(default_) > 1 {

	}
	var default__ T
	if len(default_) == 1 {
		default__ = default_[0]
	}
	e, _ := NextOrDefault(SliceTo(in, n), default__)
	return e
}

// AllEqual returns true if all elements in the input channel compares equal
// using the == operator.
func AllEqual[T comparable](in <-chan T) bool {
	g := GroupBy(in)
	_, hasFirst := Next(g)
	_, hasSecond := Next(g)
	return !hasFirst || !hasSecond
}

// Quantify counts how many times the predicate function pred is true for the
// elements in the input channel.
func Quantify[T1 any, T2 Number](in <-chan T1, pred func(T1) T2) T2 {
	return Sum(Map(pred, in))
}

// DotProduct returns a channel producing the dot product of elements from the
// two input channels.
func DotProduct[N Number](vec1 <-chan N, vec2 <-chan N) <-chan N {
	pdtFunc := func(n1, n2 N) N { return n1 * n2 }
	return Map2(pdtFunc, vec1, vec2)
}
