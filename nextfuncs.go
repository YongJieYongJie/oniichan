package oniichan

// NextOrDefault returns the next element from the input channel, or default if
// the input channel is closed.
func NextOrDefault[T any](in <-chan T, default_ T) (elem T, has bool) {
	for e := range in {
		return e, true
	}
	return default_, false
}

// NextNonBlocking returns the next element from the input channel, or the zero
// value if the operation will block.
func NextNonBlocking[T any](in <-chan T) (elem T, has bool) {
	select {
	case e := <-in:
		return e, true
	default:
		return
	}
}

// NextNonBlocking returns the next element from the input channel, or the
// provided default value if the operation will block.
func NextNonBlockingOrDefault[T any](in <-chan T, default_ T) (elem T, has bool) {
	select {
	case e := <-in:
		return e, true
	default:
		return default_, false
	}
}

// Next returns the next element from input channel, or the zero value if input
// channel is closed.
func Next[T any](in <-chan T) (elem T, has bool) {
	for e := range in {
		return e, true
	}
	return elem, false
}
