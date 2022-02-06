// Code generated by codegen.go DO NOT EDIT.
// Version: 0.0.0
// Checksum: 15974913744219286196

package oniichan

// Tee2 is the 2-element version of Tee.
func Tee2[T any](in <-chan T) (<-chan T, <-chan T) {
	tdm1 := newTeeDataManager(in, buffer_size)
	return tdm1.downstream, tdm1.Tee().downstream
}

// Tee3 is the 3-element version of Tee.
func Tee3[T any](in <-chan T) (<-chan T, <-chan T, <-chan T) {
	tdm1 := newTeeDataManager(in, buffer_size)
	return tdm1.downstream, tdm1.Tee().downstream, tdm1.Tee().downstream
}

// Tee4 is the 4-element version of Tee.
func Tee4[T any](in <-chan T) (<-chan T, <-chan T, <-chan T, <-chan T) {
	tdm1 := newTeeDataManager(in, buffer_size)
	return tdm1.downstream, tdm1.Tee().downstream, tdm1.Tee().downstream, tdm1.Tee().downstream
}

// Tee5 is the 5-element version of Tee.
func Tee5[T any](in <-chan T) (<-chan T, <-chan T, <-chan T, <-chan T, <-chan T) {
	tdm1 := newTeeDataManager(in, buffer_size)
	return tdm1.downstream, tdm1.Tee().downstream, tdm1.Tee().downstream, tdm1.Tee().downstream, tdm1.Tee().downstream
}

// Tee6 is the 6-element version of Tee.
func Tee6[T any](in <-chan T) (<-chan T, <-chan T, <-chan T, <-chan T, <-chan T, <-chan T) {
	tdm1 := newTeeDataManager(in, buffer_size)
	return tdm1.downstream, tdm1.Tee().downstream, tdm1.Tee().downstream, tdm1.Tee().downstream, tdm1.Tee().downstream, tdm1.Tee().downstream
}

// Tee7 is the 7-element version of Tee.
func Tee7[T any](in <-chan T) (<-chan T, <-chan T, <-chan T, <-chan T, <-chan T, <-chan T, <-chan T) {
	tdm1 := newTeeDataManager(in, buffer_size)
	return tdm1.downstream, tdm1.Tee().downstream, tdm1.Tee().downstream, tdm1.Tee().downstream, tdm1.Tee().downstream, tdm1.Tee().downstream, tdm1.Tee().downstream
}