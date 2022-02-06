package oniichan_test

import (
	"fmt"

	. "github.com/YongJieYongJie/oniichan"
)

// Return first n items of the iterable as a list
func Take[T any](in <-chan T, n int) <-chan T {
	return SliceFromTo(in, 0, n)
}

// Prepend a single value in front of an iterator
func Prepend[T any](elem T, in <-chan T) <-chan T {
	return Chain(ChanFrom(elem), in)
}

// Return function(0), function(1), ...
func Tabulate[T Ordered, R any](f func(T) R, start T) <-chan R {
	return Map(f, Count(start))
}

// TODO: Move to main itertools
func Tail[T any](in <-chan T, n int) <-chan T {
	out := make(chan T, n)
	go func() {
		defer close(out)
		for {
			e, has := Next(in)
			if !has {
				break
			}
			select {
			case out <- e:
				// out channel isn't full yet
			default:
				<-out
				out <- e
			}
		}
	}()
	return out
}

// Advance the iterator n-steps ahead. If n is None, consume entirely.
func Consume[T any](in <-chan T, n int) {
	SliceTo(in, n)
}

// Returns the nth item or a default value
func Nth[T any](in <-chan T, n int, default_ T) T {
	e, _ := NextOrDefault(SliceTo(in, n), default_)
	return e
}

func AllEqual[T comparable](in <-chan T) bool {
	g := GroupBy(in)
	_, hasFirst := Next(g)
	_, hasSecond := Next(g)
	return !hasFirst || !hasSecond
}

// Quantify counts how many times the predicate is true.
func Quantify[T1 any, T2 Number](in <-chan T1, pred func(T1) T2) T2 {
	return Sum(Map(pred, in))
}

// NCycles returns a channel producing the elements from the input channel n
// times.
func NCycles[T any](in <-chan T, n int) <-chan T {
	// return ChainFromChan[T](Chan(TeeN(in, n)))
	return ChainFromChan(Map(Chan[T], RepeatTimes(Slice(in), n)))
	// TODO Should ChainFromChan have a ChanFromSlice; ChainFromChan should really be called flatten?
	// TODO For ChainFromChan, how to handle various input such as:
	//  - channel of channels --> currently supported (perhaps rename to Flatten)
	//  - slice of channels --> perhaps should be a FlattenFromSlice, and also return slices
	//
	//  - slice of slices --> should clearly be a slices.Flatten function, and return a slice
	//  - channel of slices --> perhaps should be slices.FlattenFromChan
	// TODO Create a companion slices packages?
}

func DotProduct[N Number](vec1 <-chan N, vec2 <-chan N) <-chan N {
	pdtFunc := func(n1, n2 N) N { return n1 * n2 }
	return Map2(pdtFunc, vec1, vec2)
}

func Example_recipes() {
	fmt.Println(
		`Slice(DotProduct(ChanFrom(1, 2, 3), ChanFrom(2, 3, 4)))`)
	fmt.Println(
		Slice(DotProduct(
			ChanFrom(1, 2, 3),
			ChanFrom(2, 3, 4))))
	// output:
	// Slice(DotProduct(ChanFrom(1, 2, 3), ChanFrom(2, 3, 4)))
	// [2 6 12]
}
