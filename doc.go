// Package itertools re-implements Python's excellent itertools module in Go.
//
// Introduction
//
// Tired of writing repetitive Go code? Fret not, for onii-chan is here to help!
//
// Functions provided in this package allows us to write more succinct Go code
// without sacrificing performance (at least not too much).
//
// For example, a basic producer-consumer pattern may be implemented as such:
//
//     // The two Consume goroutines below will both receive every single
//     // element from sourceCh, sharing an underlying buffer that will grow
//     // and shrink as necessary. Channels tee1 and tee2 will also be properly
//     // closed when depleted.
//     //
//     // Without Tee(), elements consumed by one goroutine will not be seen by
//     // the other.
//
//     var sourceCh <-chan MyDataStruct // assume we have a receive-only channel with custom data type
//     tee1, tee2 := sourceCh.Tee()
//     go Consume(tee1)
//     go Consume(tee2)
//
// Psuedo-Variadic Support
//
// Where relevant, this package provide variadic support by having multiple
// functions defined for the same operation.
//
// For example, we use Product2() to calculate cartesian products of two input
// channels:
//
//     pdt2 := Product2(ChanFrom(1, 2, 3),
//                      ChanFrom("A", "B", "C"))
//     fmt.Println(Slice(pdt2))
//
//     // output:
//     // [[1 A] [1 B] [1 C] [2 A] [2 B] [2 C] [3 A] [3 B] [3 C]]
//
// , and we use Product3() to calculate cartesian products of three input
// channels:
//
//     pdt3 := Product3(ChanFrom(1, 2, 3),
//                      ChanFrom("A", "B", "C"),
//                      ChanFrom(1.1, 2.2, 3.3))
//     fmt.Println(Slice(pdt3))
//
//     // output (linebreaks added for clarity):
//     // [
//     //     [1 A 1.1] [1 A 2.2] [1 A 3.3] [1 B 1.1] [1 B 2.2] [1 B 3.3] [1 C 1.1] [1 C 2.2] [1 C 3.3]
//     //     [2 A 1.1] [2 A 2.2] [2 A 3.3] [2 B 1.1] [2 B 2.2] [2 B 3.3] [2 C 1.1] [2 C 2.2] [2 C 3.3]
//     //     [3 A 1.1] [3 A 2.2] [3 A 3.3] [3 B 1.1] [3 B 2.2] [3 B 3.3] [3 C 1.1] [3 C 2.2] [3 C 3.3]
//     // ]
//
// Usage Notes
//
// The functions in the package generally consumes the input channel. If the
// input channel is needed elsewhere, consider using the Tee function provided
// to make a copy.
//
package oniichan
