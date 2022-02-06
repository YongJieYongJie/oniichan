package oniichan

// Slice returns a slice with elements received from the input channel.
// TODO: Perhaps rename to Collect?
func Slice[T any](in <-chan T) []T {
	var s []T
	for e := range in {
		s = append(s, e)
	}
	return s
}

// SliceFrom returns a slice containing the input arguments.
func SliceFrom[T any](in ...T) []T {
	s := make([]T, len(in))
	for i, e := range in {
		s[i] = e
	}
	return s
}

// Chan creates a channel that will produce elements from the input slice.
func Chan[T any](in []T) <-chan T {
	ch := make(chan T)
	go func() {
		for _, v := range in {
			ch <- v
		}
		close(ch)
	}()
	return ch
}

// Chan creates a channel that will produce elements from the input arguments.
func ChanFrom[T any](in ...T) <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for _, v := range in {
			ch <- v
		}
	}()
	return ch
}
