# onii-chan

Package itertools re-implements Python's excellent itertools module in Go.

### Introduction

Functions provided in this package allows more succinct without sacrificing
performance (at least not too much).

For example, a basic producer-consumer pattern may be implemented as such:

    // The two Consume goroutines below will both receive every single
    // element from sourceCh, sharing an underlying buffer that will grow
    // and shrink as necessary. Channels tee1 and tee2 will also be properly
    // closed when depleted.

    var sourceCh <-chan MyDataStruct = Produce()
    tee1, tee2 := sourceCh.Tee()
    go Consume(tee1)
    go Consume(tee2)

### Psuedo-Variadic Support

Where relevant, this package provide variadic support by having multiple
functions defined for the same operation.

For example, cartesian products of two and three input channels may be
obtained as such:

    pdt2 := Product(ChanFrom(1, 2, 3),
                    ChanFrom("A", "B", "C")) // Product is an alias of Product2
    fmt.Println(Slice(pt2))

    // output:
    // [[1 A] [1 B] [1 C] [2 A] [2 B] [2 C] [3 A] [3 B] [3 C]]

    pdt3 := Product3(ChanFrom(1, 2, 3),
                     ChanFrom("A", "B", "C"),
                     ChanFrom(1.1, 2.2, 3.3))
    fmt.Println(Slice(pt3))

    // output (linebreaks added for clarity):
    // [
    //     [1 A 1.1] [1 A 2.2] [1 A 3.3] [1 B 1.1] [1 B 2.2] [1 B 3.3] [1 C 1.1] [1 C 2.2] [1 C 3.3]
    //     [2 A 1.1] [2 A 2.2] [2 A 3.3] [2 B 1.1] [2 B 2.2] [2 B 3.3] [2 C 1.1] [2 C 2.2] [2 C 3.3]
    //     [3 A 1.1] [3 A 2.2] [3 A 3.3] [3 B 1.1] [3 B 2.2] [3 B 3.3] [3 C 1.1] [3 C 2.2] [3 C 3.3]
    // ]

### Usage Notes

The functions in the package generally consumes the input channel. If the
input channel is needed elsewhere, consider using the Tee function provided
to make a copy.

### Sister Library

Refer to the tuples sister library.
