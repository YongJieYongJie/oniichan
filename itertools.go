package oniichan

// TODO: Make the functions non-blocking. E.g., check usages of Slice(), and
// make sure it's non-blocking.

import (
	"fmt"
	"reflect"
	"sort"
	"sync"

	"github.com/YongJieYongJie/tttuples/atuple"
)

const buffer_size = 257 // Prime number closest to 256

// Ordered is any Ordered type: any type that supports the operators < <= >= >.
//
// Note: This is a temporary shim to constraints.Ordered as proposed in
// https://github.com/golang/go/discussions/47319.
type Ordered interface {
	Integer | Float
}

// Integer is any integer type.
//
// Note: This is a temporary shim to constraints.Integer as proposed in
// https://github.com/golang/go/discussions/47319.
type Integer interface {
	int | int8 | int32 | int64 |
		uint | uint8 | uint32 | uint64
}

// Float is any floating-point type.
//
// Note: This is a temporary shim to constraints.Float as proposed in
// https://github.com/golang/go/discussions/47319.
type Float interface {
	float32 | float64
}

// Complex is any complex number type.
//
// Note: This is a temporary shim to constraints.Complex as proposed in
// https://github.com/golang/go/discussions/47319.
type Complex interface {
	complex64 | complex128
}

// Number is a constraint that permits any number type: any type that
// supports the arithmetic operators + - / *.
type Number interface {
	Integer | Float | Complex
}

// All returns true if all elements in the input channel returns true when pred
// is applied.
func All[T any](in <-chan T, pred func(T) bool) bool {
	for e := range in {
		if !pred(e) {
			return false
		}
	}
	return true
}

// Any returns true if any element in the input channel returns true when pred
// is applied.
func Any[T any](in <-chan T, pred func(T) bool) bool {
	return !All(in, pred)
}

// Filter returns a channel that filters elements from the input channel,
// returning only those for which the predicate is true.
func Filter[T any](pred func(T) bool, in <-chan T) <-chan T {
	notPred := func(e T) bool { return !pred(e) }
	return FilterFalse(notPred, in)
}

// Max returns the largest element from input channel.
func Max[T Ordered](in <-chan T) T {
	out := <-in
	for e := range in {
		if e > out {
			out = e
		}
	}
	return out
}

// MaxFunc returns the element from input channel that returns that largest
// value when f is applied.
func MaxFunc[T1 any, T2 Ordered](in <-chan T1, f func(T1) T2) T1 {
	out := <-in
	fOut := f(out)
	for e := range in {
		fE := f(e)
		if fE > fOut {
			out = e
			fOut = fE
		}
	}
	return out
}

// Min returns the smallest element from input channel.
func Min[T Ordered](in <-chan T) T {
	out := <-in
	for e := range in {
		if e < out {
			out = e
		}
	}
	return out
}

// MinFunc returns the element from input channel that returns that smallest
// value when f is applied.
//
// Design decision: whether to use keyfunc like python, or lessfunc like Go
func MinFunc[T1 any, T2 Ordered](in <-chan T1, f func(T1) T2) T1 {
	out := <-in
	fOut := f(out)
	for e := range in {
		fE := f(e)
		if fE < fOut {
			out = e
			fOut = fE
		}
	}
	return out
}

// Sum returns the sum over all elements received from the input channel.
func Sum[T Number](in <-chan T) T {
	var sum T
	for e := range in {
		sum += e
	}
	return sum
}

// Reduce returns the result of applying the 2-argument function f cumulatively
// to the elements received from the input channel, reducing 2 elements to 1,
// until the channel is exhausted and there is only 1 left.
func Reduce[T any](in <-chan T, f func(T, T) T) T {
	res, has := Next(in)
	if !has {
		return res
	}
	for e := range in {
		res = f(res, e)
	}
	return res
}

// Sorted returns a channel producing elements from input channel in sorted
// order based on the key func.
//
// Note: Input channel should be a finite channel as the elements are first
// read into a slice.
//
// Note 2: API is different sort.Slice.
func Sorted[T1 any, T2 Ordered](in <-chan T1, key func(T1) T2) <-chan T1 {
	out := make(chan T1)
	go func() {
		defer close(out)
		var buf []T1
		var keyBuf []T2
		for e := range in {
			buf = append(buf, e)
			keyBuf = append(keyBuf, key(e))
		}
		sort.Slice(buf, func(i, j int) bool { return keyBuf[i] < keyBuf[j] })
		for _, e := range buf {
			out <- e
		}
	}()
	return out
}

// Reversed returns a channel producing elements from input channel in reversed
// order.
//
// Note: Input channel should be a finite channel as the elements are first
// read into a slice.
func Reversed[T any](in <-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		var buf []T
		for e := range in {
			buf = append(buf, e)
		}
		for i := len(buf) - 1; i >= 0; i-- {
			out <- buf[i]
		}
	}()
	return out
}

// Enumerate returns a channel producing 2-tuple containing a count starting
// from 0, and the element received from the input channel.
func Enumerate[T any](in <-chan T) <-chan atuple.Packed2[int, T] {
	out := make(chan atuple.Packed2[int, T])
	go func() {
		defer close(out)
		var i int
		for e := range in {
			out <- atuple.Pack2(i, e)
			i++
		}
	}()
	return out
}

// Enumerate returns a channel producing 2-tuple containing a count starting
// from the provided argument, and the element received from the input channel.
func EnumerateFrom[T any](in <-chan T, from int) <-chan atuple.Packed2[int, T] {
	out := make(chan atuple.Packed2[int, T])
	go func() {
		defer close(out)
		i := from
		for e := range in {
			out <- atuple.Pack2(i, e)
			i++
		}
	}()
	return out
}

// Range returns a channel that produces numbers from 0 to (but not including)
// stop.
func Range[T Ordered](stop T) <-chan T {
	if stop <= 0 {
		panic("invalid argument to Range: stop must be greater than zero")
	}
	return RangeFromTo(0, stop)
}

// RangeFromTo returns a channel that produces numbers from start to (but not
// including) stop.
func RangeFromTo[T Ordered](start, stop T) <-chan T {
	if start >= stop {
		panic("invalid argument to RangeFromTo: start must be smaller than stop")
	}
	out := make(chan T)
	go func() {
		defer close(out)
		for e := start; e < stop; e += 1 {
			out <- e
		}
	}()
	return out
}

// TODO: consider renaming to RangeFromToStep?
// RangeStep returns a channel that produces numbers from start to (but not
// including) stop, spaced evenly apart by step amount.
func RangeStep[T Ordered](start, stop, step T) <-chan T {
	if step == 0 {
		panic("invalid argument to RangeStep: step cannot be 0")
	}
	if start > stop && step > 0 {
		panic("invalid argument to RangeStep: step must be negative if start > stop")
	}
	if start < stop && step < 0 {
		panic("invalid argument to RangeStep: step must be positive if start < stop")
	}
	out := make(chan T)
	go func() {
		defer close(out)
		if step > 0 {
			for e := start; e < stop; e += step {
				out <- e
			}
		} else {
			for e := start; e > stop; e += step {
				out <- e
			}
		}
	}()
	return out
}

// Map returns a channel that produces each element by applying function f to
// the corresponding element received from the input channel.
func Map[T1, T2 any](f func(T1) T2, in <-chan T1) <-chan T2 {
	out := make(chan T2)
	go func() {
		defer close(out)
		for e := range in {
			out <- f(e)
		}
	}()
	return out
}

//go:generate go run ./internal/codegen/ -N 7 --outfile codegen_map.go
var mapTemplate string = `
// Map{{ .N }} is the {{ .N }}-element version of Map.
func Map{{ .N }}[{{ .T1_T2___TN }}, R any](
	f func({{ .T1_T2___TN }}) R,
	{{ .RepeatOneToN  "in__N1 <-chan T__N1" }},
) <-chan R {
	out := make(chan R)
	go func() {
		defer close(out)
		for {
			{{ range .OneToN }}
			e{{ . }}, has{{ . }} := Next(in{{ . }})
			if !has{{ . }} {
				break
			}
			{{ end }}
			out <- f({{ .RepeatOneToN "e__N1" }})
		}
	}()
	return out
}
`

// Count returns an infinite channel that produces increasing values
// starting with number start.
func Count[v Ordered](start v) <-chan v {
	return CountStep(start, 1)
}

// Count returns an infinite channel that produces evenly spaced values
// starting with number start.
func CountStep[v Ordered](start, step v) <-chan v {
	out := make(chan v)
	go func() {
		defer close(out)
		for e := start; true; e += step {
			out <- e
		}
	}()
	return out
}

// Cycle returns an infinite channel that produces elements from the input
// channel, saving a copy of each and looping back to the start when the input
// channel is exhausted.
func Cycle[T any](in <-chan T) <-chan T {
	buf := make([]T, 0, buffer_size)
	out := make(chan T)
	go func() {
		defer close(out)
		for e := range in {
			out <- e
			buf = append(buf, e)
		}
		for {
			for _, e := range buf {
				out <- e
			}
		}
	}()
	return out
}

// Repeat returns a channel that produces elem over and over again.
func Repeat[T any](elem T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for {
			out <- elem
		}
	}()
	return out
}

// RepeatTimes returns a channel that produces elem up to the number of times as
// provided.
func RepeatTimes[T any](elem T, times int) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for i := 0; i < times; i++ {
			out <- elem
		}
	}()
	return out
}

// Accumulate returns a channel that produces the accumulated sums of elements
// from in.
func Accumulate[T Ordered](in <-chan T) <-chan T {
	out := make(chan T)
	var prev T
	go func() {
		defer close(out)
		for e := range in {
			prev += e
			out <- prev
		}
	}()
	return out
}

// AccumulateFunc returns a channel that produces the accumulated sums of
// elements from in, calculated using the provided accumulator function.
func AccumulateFunc[T any](in <-chan T, accumulator func(T, T) T) <-chan T {
	out := make(chan T)
	var prev T
	go func() {
		defer close(out)
		for e := range in {
			prev = accumulator(prev, e)
			out <- prev
		}
	}()
	return out
}

// Chain returns a channel that sends each element from the first input channel
// until the input channel is exhausted, the proceed to the next input channel,
// until all input channels are exhausted.
func Chain[T any](ins ...<-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for _, in := range ins {
			for e := range in {
				out <- e
			}
		}
	}()
	return out
}

// ChainFromChan is an alternative to Chain that accepts a channel of channels
// as input.
func ChainFromChan[T any](ins <-chan <-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for in := range ins {
			for e := range in {
				out <- e
			}
		}
	}()
	return out
}

// Compress returns a channel that produces elements from input channel where
// there is a corresponding element in selectors that is true.
func Compress[T any](in <-chan T, selectors <-chan bool) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for tuple := range Zip(in, selectors) {
			if !tuple.Item2() {
				continue
			}
			out <- tuple.Item1()
		}
	}()
	return out
}

// DropWhile returns a channel that drops element from the input channel as
// long as the predicate is true; afterwards, returns every element.
func DropWhile[T any](pred func(T) bool, in <-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for e := range in {
			if !pred(e) {
				continue
			}
			out <- e
		}
	}()
	return out
}

// FilterFalse returns a channel that filters elements from the input channel,
// returning only those for which the predicate is false.
func FilterFalse[T any](pred func(T) bool, in <-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for e := range in {
			if pred(e) {
				continue
			}
			out <- e
		}
	}()
	return out
}

// GroupBy calls GroupByKey passing in the identity function as the key
// function.
func GroupBy[T comparable](in <-chan T) <-chan atuple.Packed2[T, <-chan T] {
	return GroupByKey(in, func(v T) T { return v })
	// return GroupByKey(in, Identity[T])
}

// Identity returns the argument provided.
func Identity[T any](x T) T {
	return x
}

// GroupByKey returns a channel that returns consecutive keys and groups from
// the input channel. The input key function is used to compute the key value
// for each element. Generally, the input channel needs to be already sorted on
// the same key function.
//
// The operation of GroupByKey() is similar to the uniq filter in Unix.
func GroupByKey[T1 any, T2 comparable](
	in <-chan T1, key func(T1) T2) <-chan atuple.Packed2[T2, <-chan T1] {
	out := make(chan atuple.Packed2[T2, <-chan T1])
	go func() {
		defer close(out)

		elem, has := Next(in)
		if !has {
			return
		}

		groupElems := make(chan T1, 1)
		groupElems <- elem
		groupKey := key(elem)
		currGroup := atuple.Pack2[T2, <-chan T1](groupKey, groupElems)
		out <- currGroup

		for e := range in {
			if groupKey == key(e) {
				groupElems <- e
			} else {
				close(groupElems)
				groupKey = key(e)
				groupElems = make(chan T1, 1)
				groupElems <- e
				currGroup = atuple.Pack2[T2, <-chan T1](groupKey, groupElems)
				out <- currGroup
			}
		}
	}()
	return out
}

// SliceTo returns a channel that produces elements from the input channel from
// the first up to and before stop, analagous to to the indexing operation
// [0:stop] on a slice.
func SliceTo[T any](in <-chan T, stop int) <-chan T {
	if stop < 0 {
		panic("invalid argument to SliceTo: stop must be greater than or equal to zero")
	}
	return SliceFromTo(in, 0, stop)
}

// SliceFromTo returns a channel that produces elements from the input channel from
// the start up to and before stop, analagous to to the indexing operation
// [start:stop] on a slice.
func SliceFromTo[T any](in <-chan T, start, stop int) <-chan T {
	if start > stop {
		panic("invalid argument to SliceFromTo: start must be smaller or equal to stop")
	}
	if start < 0 {
		panic("invalid argument to SliceFromTo: start must be greater than or equal to zero")
	}
	return SliceFromToStep(in, start, stop, 1)
}

// SliceFromToStep returns a channel that produces elements from the input channel from
// the start up to and before stop, analagous to to the indexing operation
// [start:stop:step] on a slice.
func SliceFromToStep[T any](in <-chan T, start, stop, step int) <-chan T {
	if start > stop {
		panic("invalid argument to SliceFromToStep: start must be smaller or equal to stop")
	}
	if start < 0 {
		panic("invalid argument to SliceFromToStep: start must be greater than or equal to zero")
	}
	if step <= 0 {
		panic("invalid argument to SliceFromToStep: step must be greater than zero")
	}
	out := make(chan T)
	go func() {
		defer close(out)
		enumerated := Enumerate(in)
		for e := range enumerated {
			if e.Item1() < start {
				continue
			}
			if e.Item1()%step == 0 {
				out <- e.Item2()
			}
			break
		}
		for e := range enumerated {
			if e.Item1() >= stop {
				return
			}
			if e.Item1()%step != 0 {
				continue
			}
			out <- e.Item2()
		}
	}()
	return out
}

// Pairwise returns a channel that produces successive overlapping pairs from
// the input channel.
func Pairwise[T any](in <-chan T) <-chan atuple.Packed2[T, T] {
	out := make(chan atuple.Packed2[T, T])
	go func() {
		defer close(out)
		prev, has := Next(in)
		if !has {
			return
		}
		for e := range in {
			out <- atuple.Pack2(prev, e)
			prev = e
		}
	}()
	return out
}

// StarMap returns a channel that computes f using arguments obtained from the
// input channel.
func StarMap[T1, T2, R any](
	f func(T1, T2) R, in <-chan atuple.Packed2[T1, T2]) <-chan R {

	out := make(chan R)
	go func() {
		defer close(out)
		for e := range in {
			out <- f(e.Item1(), e.Item2())
		}
	}()
	return out
}

//go:generate go run ./internal/codegen/ -N 7 --outfile codegen_star_map.go --imports "github.com/YongJieYongJie/tttuples/atuple"
var starMapTemplate string = `
// StarMap{{ .N }} is the {{ .N }}-element version of StarMap.
func StarMap{{ .N }}[{{ .T1_T2___TN }}, R any](
	f func({{.T1_T2___TN}}) R,
	in <-chan atuple.Packed{{ .N }}[{{ .T1_T2___TN }}],
) <-chan R {

	out := make(chan R)
	go func() {
		defer close(out)
		for e := range in {
			out <- f({{ .RepeatOneToN "e.Item__N1()" }})
		}
	}()
	return out
}
`

// StarMap returns a channel that computes f using arguments obtained from the
// input channel.
//
// TODO: Consider whether this function is still required, since I'll be
// generating StarMap2~7 (or more)? This might be necessary for dynamic usage
// (i.e., number of arguments to f not known at compile time).
func StarMapN[R any](f any, in <-chan []any) <-chan R {
	// Design decision:
	//  - Use any as the type of f so existing functions can be passed in
	//    directly.
	// Trade-offs:
	//  - Return type T cannot be inferred, and must be specified.
	// Alternative:
	//  - Use func(...any) T as the type of f. Not chosen because existing
	//    functions cannot be used directly, and any adapter to convert existing
	//    functions into func(...any) T will necessarily need to accept type any.
	typeF := reflect.TypeOf(f)
	if typeF.Kind() != reflect.Func {
		panic("invalid argument to StarMapN: f must be a function")
	}
	if typeF.NumOut() != 1 {
		panic("invalid argument to StarMapN: f must be a function that return 1 " +
			"value")
	}
	var typeParam R
	if typeF.Out(0) != reflect.TypeOf(typeParam) {
		panic("invalid argument to StarMapN: f must be a function that return 1 " +
			"value of type T")
	}

	out := make(chan R)
	go func() {
		defer close(out)

		// Manually fetch the first element so we can check that the len and type
		// of the tuple matches each argument to f.
		e, has := Next(in)
		if !has {
			return
		}

		if typeF.NumIn() != len(e) {
			panic("invalid argument to StarMapN: number of arguments to f must be " +
				"same as length of each element in in")
		}

		for i := 0; i < typeF.NumIn(); i++ {
			if typeF.In(i) != reflect.TypeOf(e[i]) {
				panic(fmt.Sprintf("invalid argument to StarMapN: type of in's element "+
					"at index %d (%v) does not match type of f's argument at index %d "+
					"(%v)", i, typeF.In(i), i, reflect.TypeOf(e[i])))
			}
		}

		args := make([]reflect.Value, typeF.NumIn())
		for i := 0; i < typeF.NumIn(); i++ {
			args[i] = reflect.ValueOf(e[i])
		}
		retVal := reflect.ValueOf(f).Call(args)[0]
		out <- retVal.Interface().(R)

		for e = range in {
			for i := 0; i < typeF.NumIn(); i++ {
				args[i] = reflect.ValueOf(e[i])
			}
			retVal := reflect.ValueOf(f).Call(args)[0]
			out <- retVal.Interface().(R)
		}
	}()
	return out
}

func TakeWhile[T any](pred func(T) bool, in <-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for e := range in {
			if !pred(e) {
				return
			}
			out <- e
		}
	}()
	return out
}

// Tee returns 2 independent channels from the input channel.
//
// Note: Tee is an alias of Tee2.
func Tee[T any](in <-chan T) (<-chan T, <-chan T) {
	tdm1 := newTeeDataManager(in, buffer_size)
	tdm2 := tdm1.Tee()
	return tdm1.downstream, tdm2.downstream
}

//go:generate go run ./internal/codegen/ -N 7 --outfile codegen_tee.go
var teeTemplate string = `
// Tee{{ .N }} is the {{ .N }}-element version of Tee.
func Tee{{ .N }}[T any](in <-chan T) ({{ .RepeatOneToN "<-chan T" }}) {
	tdm1 := newTeeDataManager(in, buffer_size)
	return tdm1.downstream, {{ .RepeatTwoToN "tdm1.Tee().downstream" }}
}
`

// TeeN is the variable-element version of Tee, allowing the number of tee'd
// channels to be determined dynamically at runtime instead of at compile time.
func TeeN[T any](in <-chan T, n int) []<-chan T {
	if n < 1 {
		panic("invalid argument to TeeN: n must be greater than or equal to one")
	}
	tdm1 := newTeeDataManager(in, buffer_size)
	out := make([]<-chan T, n)
	out[0] = tdm1.downstream
	for i := 1; i < n; i++ {
		out[i] = tdm1.Tee().downstream
	}
	return out
}

type teeDataManager[T any] struct {
	// sourceMu is used for synchronizing access to the upstream channel and node
	sourceMu *sync.Mutex
	node     *teeDataNode[T]

	// initialized channel is closed when the first element has been pulled from
	// upstream.
	initialized chan struct{}

	nextIndex       int
	nodeSize        int
	isNodeRotatedMu sync.Mutex
	isNodeRotated   bool

	upstream   <-chan T
	downstream chan T
}

func newTeeDataManager[T any](
	upstream <-chan T, nodeSize int) *teeDataManager[T] {

	tdm := &teeDataManager[T]{
		sourceMu: &sync.Mutex{},
		node: &teeDataNode[T]{
			data: make([]T, 0, nodeSize),
			next: nil,
		},

		initialized: make(chan struct{}),

		nextIndex:     0,
		nodeSize:      nodeSize,
		isNodeRotated: false,

		upstream:   upstream,
		downstream: make(chan T),
	}

	go tdm.Init()
	go tdm.Start()
	return tdm
}

func (tdm *teeDataManager[T]) Init() {
	defer close(tdm.initialized)
	e, has := Next(tdm.upstream)
	if !has {
		return
	}
	tdm.node.data = append(tdm.node.data, e)
}

func (tdm *teeDataManager[T]) Start() {
	defer close(tdm.downstream)
	<-tdm.initialized
	if len(tdm.node.data) == 0 {
		// Special case where upstream if already closed before tee-ing.
		return
	}
	for {
		e, has := tdm.Next()
		if !has {
			break
		}
		tdm.downstream <- e
	}
}

func (tdm *teeDataManager[T]) Next() (elem T, has bool) {
	defer func() { tdm.nextIndex += 1 }()

	// Move to next block of data if necessary
	if tdm.nextIndex == tdm.nodeSize {
		tdm.isNodeRotatedMu.Lock()
		defer tdm.isNodeRotatedMu.Unlock()
		tdm.isNodeRotated = true
		tdm.nextIndex = 0
		tdm.node = tdm.node.next
		if tdm.node == nil {
			tdm.node = &teeDataNode[T]{
				data: make([]T, 0, tdm.nodeSize),
				next: nil,
			}
		}
	}

	// Case 1: Next element already in buffer, we simply return.
	if tdm.nextIndex < len(tdm.node.data) {
		return tdm.node.data[tdm.nextIndex], true
	}

	// If not, we acquire lock to the upstream channel to fetch the next
	// element.
	tdm.sourceMu.Lock()
	defer tdm.sourceMu.Unlock()

	// Case 2: Next element is added to buffer during lock acquisition, we
	// simply return.
	if tdm.nextIndex < len(tdm.node.data) {
		return tdm.node.data[tdm.nextIndex], true
	}

	// Case 3: Fetch another element from upstream channel.
	e, has := Next(tdm.upstream)
	if !has {
		return elem /* zero value */, false
	}
	tdm.node.data = append(tdm.node.data, e)
	return e, true
}

func (tdm *teeDataManager[T]) Tee() *teeDataManager[T] {
	tdm.isNodeRotatedMu.Lock()
	defer tdm.isNodeRotatedMu.Unlock()
	if tdm.isNodeRotated {
		panic("cannot Tee because consumed past the first block of data from source")
	}

	teed := &teeDataManager[T]{
		sourceMu: tdm.sourceMu,
		node:     tdm.node,

		initialized: tdm.initialized,

		nextIndex:     0,
		nodeSize:      tdm.nodeSize,
		isNodeRotated: false,

		upstream:   tdm.upstream,
		downstream: make(chan T),
	}
	go teed.Start()
	return teed
}

type teeDataNode[T any] struct {
	data []T
	next *teeDataNode[T]
}

func Zip[T1, T2 any](in1 <-chan T1, in2 <-chan T2) <-chan atuple.Packed2[T1, T2] {
	out := make(chan atuple.Packed2[T1, T2])
	go func() {
		defer close(out)
		for {
			e1, has := Next(in1)
			if !has {
				return
			}
			e2, has := Next(in2)
			if !has {
				return
			}
			out <- atuple.Pack2(e1, e2)
		}
	}()
	return out
}

//go:generate go run ./internal/codegen/ -N 7 --outfile codegen_zip.go --imports "github.com/YongJieYongJie/tttuples/atuple"
var zipTemplate string = `
func Zip{{ .N }}[{{ .T1_T2___TN }} any](
	{{- .RepeatOneToN "in__N1 <-chan T__N1" -}}
) <-chan atuple.Packed{{ .N }}[{{ .T1_T2___TN }}] {
	out := make(chan atuple.Packed{{ .N }}[{{ .T1_T2___TN }}])
	go func() {
		defer close(out)
		for {
			{{ range .OneToN }}
			e{{ . }}, has := Next(in{{ . }})
			if !has {
				return
			}
			{{ end }}
			out <- atuple.Pack{{ .N }}({{ .RepeatOneToN "e__N1" }})
		}
	}()
	return out
}
`

func ZipN(ins ...<-chan any) <-chan []any {
	out := make(chan []any)
	go func() {
		defer close(out)
		for {
			zipped := make([]any, 0, len(ins))
			for i, in := range ins {
				e, has := Next(in)
				if !has {
					return
				}
				zipped[i] = e
			}
			out <- zipped
		}
	}()
	return out
}

func ZipLongest[T1, T2 any](
	in1 <-chan T1, in2 <-chan T2, fillValues atuple.Packed2[T1, T2],
) <-chan atuple.Packed2[T1, T2] {

	out := make(chan atuple.Packed2[T1, T2])
	go func() {
		defer close(out)
		for {
			e1, has1 := NextOrDefault(in1, fillValues.Item1())
			e2, has2 := NextOrDefault(in2, fillValues.Item2())
			if !has1 && !has2 {
				return
			}
			out <- atuple.Pack2(e1, e2)
		}
	}()
	return out
}

//go:generate go run ./internal/codegen/ -N 7 --outfile codegen_zip_longest.go --imports "github.com/YongJieYongJie/tttuples/atuple"
var zipLongestTemplate string = `
func ZipLongest{{ .N }}[{{ .T1_T2___TN }} any](
	{{ .RepeatOneToN "in__N1 <-chan T__N1" }},
	fillValues atuple.Packed{{ .N }}[{{ .T1_T2___TN }}],
) <-chan atuple.Packed{{ .N }}[{{ .T1_T2___TN }}] {

	out := make(chan atuple.Packed{{ .N }}[{{ .T1_T2___TN }}])
	go func() {
		defer close(out)
		for {
			{{ range .OneToN -}}
			e{{ . }}, has{{ . }} := NextOrDefault(in{{ . }}, fillValues.Item{{ . }}())
			{{ end }}
			if {{ .RepeatOneToNWithSep "!has__N1" " && "}} {
				return
			}
			out <- atuple.Pack{{ .N }}({{ .RepeatOneToN "e__N1" }})
		}
	}()
	return out
}
`

// Procduct returns a channel producing the cartesian products of elements from
// the input channels.
//
// Note: Input channel should be a finite channel as the elements are first
// read into a slice.
//
// Note 2: Product is an alias of Product2.
func Product[T1, T2 any](ch1 <-chan T1, ch2 <-chan T2) <-chan atuple.Packed2[T1, T2] {
	out := make(chan atuple.Packed2[T1, T2])
	s1 := Slice(ch1)
	s2 := Slice(ch2)
	go func() {
		defer close(out)
		for _, e1 := range s1 {
			for _, e2 := range s2 {
				out <- atuple.Pack2(e1, e2)
			}
		}
	}()
	return out
}

//go:generate go run ./internal/codegen/ -N 7 --outfile codegen_product.go --imports "github.com/YongJieYongJie/tttuples/atuple"
var productTemplate string = `
// Product{{ .N }} is the {{ .N }}-element version of Product.
func Product{{ .N }}[{{ .T1_T2___TN }} any](
	{{- .RepeatOneToN "ch__N1 <-chan T__N1" -}}
) <-chan atuple.Packed{{ .N }}[{{ .T1_T2___TN }}] {
	out := make(chan atuple.Packed{{ .N }}[{{ .T1_T2___TN }}])
  {{ range .OneToN -}}
	s{{ . }} := Slice(ch{{ . }})
	{{ end }}
	go func() {
		defer close(out)
		{{ range .OneToN -}}
		for _, e{{ . }} := range s{{ . }} {
		{{ end -}}
				out <- atuple.Pack{{ .N }}({{ .RepeatOneToN "e__N1" }})
		{{ .RepeatOneToNWithSep "}" " " }}
	}()
	return out
}
`

// Permutations returns a channel that produces successive 2-element
// permutations of elements from the input channel.
//
// Note: Input channel should be a finite channel as the elements are first
// read into a slice.
//
// Note 2: Permutations is an alias of Permutations2.
func Permutations[T any](in <-chan T) <-chan atuple.Packed2[T, T] {
	out := make(chan atuple.Packed2[T, T])

	go func() {
		defer close(out)
		permLen := 2 // number of elements in each permutation
		s := Slice(in)
		n := len(s)
		if n < permLen {
			panic("invalid argument to Permutations: size of in must be two or more")
		}
		indices := Slice(Range(n))
		out <- atuple.Pack2(s[indices[0]], s[indices[1]])
		cycles := Slice(RangeStep(n, n-permLen, -1))

		for {
			// based on Python's itertool:
			// https://github.com/python/cpython/blob/70c16468deee9390e34322d32fda57df6e0f46bb/Modules/itertoolsmodule.c#L3408
			// TODO: Find the name of the algorithm and update the variable names etc.
			var i int
			// Decrement rightmost cycle, moving leftward upon zero rollover
			for i = permLen - 1; i >= 0; i-- {
				cycles[i]--
				if cycles[i] == 0 {
					if i == 0 {
						return
					}

					// Move index at i to the end, shifting indices from i+1:end
					// forward by 1.
					indexI := indices[i]
					for j := i; j < n-1; j++ {
						indices[j] = indices[j+1]
					}
					indices[n-1] = indexI

					cycles[i] = n - i
				} else {
					j := cycles[i]
					indices[i], indices[n-j] = indices[n-j], indices[i]
					out <- atuple.Pack2(s[indices[0]], s[indices[1]])
					break
				}
			}
			// If i is negative, then the cycles have all rolled-over and we're done.
			if i < 0 {
				return
			}
		}
	}()

	return out
}

//go:generate go run ./internal/codegen/ -N 7 --outfile codegen_permutations.go --imports "github.com/YongJieYongJie/tttuples/atuple"
var permutationsTemplate string = `
// Permutations{{ .N }} is the {{ .N }}-element version of Permutations.
func Permutations{{ .N }}[T any](in <-chan T) <-chan atuple.Packed{{ .N }}[{{ .RepeatOneToN "T" }}] {
	out := make(chan atuple.Packed{{ .N }}[{{ .RepeatOneToN "T" }}])

	go func() {
		defer close(out)
		permLen := {{ .N }} // number of elements in each permutation
		s := Slice(in)
		n := len(s)
		if n < permLen {
			panic("invalid argument to Permutations{{ .N }}: size of in must greater than or equal to {{ .N }}")
		}
		indices := Slice(Range(n))
		out <- atuple.Pack{{ .N }}({{ .RepeatOneToN "s[indices[__N0]]" }})
		cycles := Slice(RangeStep(n, n-permLen, -1))

		for {
			// based on Python's itertool:
			// https://github.com/python/cpython/blob/70c16468deee9390e34322d32fda57df6e0f46bb/Modules/itertoolsmodule.c#L3408
			// TODO: Find the name of the algorithm and update the variable names etc.
			var i int
			// Decrement rightmost cycle, moving leftward upon zero rollover
			for i = permLen - 1; i >= 0; i-- {
				cycles[i]--
				if cycles[i] == 0 {
					if i == 0 {
						return
					}

					// Move index at i to the end, shifting indices from i+1:end
					// forward by 1.
					indexI := indices[i]
					for j := i; j < n-1; j++ {
						indices[j] = indices[j+1]
					}
					indices[n-1] = indexI

					cycles[i] = n - i
				} else {
					j := cycles[i]
					indices[i], indices[n-j] = indices[n-j], indices[i]
					out <- atuple.Pack{{ .N }}({{ .RepeatOneToN "s[indices[__N0]]" }})
					break
				}
			}
			// If i is negative, then the cycles have all rolled-over and we're done.
			if i < 0 {
				return
			}
		}
	}()

	return out
}
`

// Combinations returns a channel that produces subsequences of 2 elements from
// the input channel.
//
// Note: Input channel should be a finite channel as the elements are first
// read into a pool before generating the 2-element combinations.
//
// Note 2: Combinations is an alias of Combinations2.
func Combinations[T any](in <-chan T) <-chan atuple.Packed2[T, T] {
	out := make(chan atuple.Packed2[T, T])
	go func() {
		defer close(out)
		combiLen := 2
		s := Slice(in)
		n := len(s)
		if n < combiLen {
			panic("invalid argument to Combinations: size of in must be greater than or equal to 2")
		}
		indices := Slice(Range(combiLen))
		out <- atuple.Pack2(s[indices[0]], s[indices[1]])

		for {
			var i int
			// Scan indices right-to-left until finding one that is not at its
			// maximum (i + n - combiLen).
			for i = combiLen - 1; i >= 0 && indices[i] == i+n-combiLen; i-- {
			}

			// If i is negative, then the indices are all at their maximum value and
			// we're done.
			if i < 0 {
				return
			}

			// Increment the current index which we know is not at its maximum. Then
			// move back to the right setting each index to its lowest possible value
			// (one higher than the index to its left -- this maintains the sort
			// order invariant).
			indices[i]++
			for j := i + 1; j < combiLen; j++ {
				indices[j] = indices[j-1] + 1
			}

			out <- atuple.Pack2(s[indices[0]], s[indices[1]])
		}
	}()
	return out
}

//go:generate go run ./internal/codegen/ -N 7 --outfile codegen_combinations.go --imports "github.com/YongJieYongJie/tttuples/atuple"
var combinationsTemplate string = `
// Combinations{{ .N }} is the {{ .N }}-element version of Combinations.
func Combinations{{ .N }}[T any](in <-chan T) <-chan atuple.Packed{{ .N }}[{{ .RepeatOneToN "T" }}] {
	out := make(chan atuple.Packed{{ .N }}[{{ .RepeatOneToN "T" }}])
	go func() {
		defer close(out)
		combiLen := {{ .N }}
		s := Slice(in)
		n := len(s)
		if n < combiLen {
			panic("invalid argument to Combinations: size of in must greater than or equal to {{ .N }}")
		}
		indices := Slice(Range(combiLen))
		out <- atuple.Pack{{ .N }}({{ .RepeatOneToN "s[indices[__N0]]" }})

		for {
			var i int
			// Scan indices right-to-left until finding one that is not at its
			// maximum (i + n - combiLen).
			for i = combiLen - 1; i >= 0 && indices[i] == i+n-combiLen; i-- {
			}

			// If i is negative, then the indices are all at their maximum value and
			// we're done.
			if i < 0 {
				return
			}

			// Increment the current index which we know is not at its maximum. Then
			// move back to the right setting each index to its lowest possible value
			// (one higher than the index to its left -- this maintains the sort
			// order invariant).
			indices[i]++
			for j := i + 1; j < combiLen; j++ {
				indices[j] = indices[j-1] + 1
			}

			out <- atuple.Pack{{ .N }}({{ .RepeatOneToN "s[indices[__N0]]" }})
		}
	}()
	return out
}
`

// CombinationsWithReplacement returns a channel that produces subsequences of
// 2 elements from the input channel, allowing individual elments to be
// repeated more than once.
//
// Note: Input channel should be a finite channel as the elements are first
// read into a pool before generating the 2-element combinations.
//
// Note 2: CombinationsWithReplacement is an alias of CombinationsWithReplacement2.
func CombinationsWithReplacement[T any](in <-chan T) <-chan atuple.Packed2[T, T] {
	out := make(chan atuple.Packed2[T, T])
	go func() {
		defer close(out)
		combiLen := 2
		s := Slice(in)
		n := len(s)
		indices := []int{0, 0}
		out <- atuple.Pack2(s[indices[0]], s[indices[1]])

		for {
			var i int
			// Scan indices right-to-left until finding one that is not at its
			// maximum (n-1).
			for i = combiLen - 1; i >= 0 && indices[i] == n-1; i-- {
			}

			// If i is negative, then the indices are all at their maximum value and
			// we're done.
			if i < 0 {
				return
			}

			// Increment the current index which we know is not at its maximum. Then
			// set all to the right to the same value.
			indices[i]++
			elem := indices[i]
			for j := i + 1; j < combiLen; j++ {
				indices[j] = elem
			}

			out <- atuple.Pack2(s[indices[0]], s[indices[1]])
		}
	}()
	return out
}

//go:generate go run ./internal/codegen/ -N 7 --outfile codegen_combinations_with_replacement.go --imports "github.com/YongJieYongJie/tttuples/atuple"
var combinationsWithReplacementTemplate string = `
// CombinationsWithReplacement{{ .N }} is the {{ .N }}-element version of CombinationsWithReplacement.
func CombinationsWithReplacement{{ .N }}[T any](in <-chan T) <-chan atuple.Packed{{ .N }}[{{ .RepeatOneToN "T" }}] {
	out := make(chan atuple.Packed{{ .N }}[{{ .RepeatOneToN "T" }}])
	go func() {
		defer close(out)
		combiLen := {{ .N }}
		s := Slice(in)
		n := len(s)
		indices := []int{0, 0}
		out <- atuple.Pack{{ .N }}({{ .RepeatOneToN "s[indices[__N0]]" }})

		for {
			var i int
			// Scan indices right-to-left until finding one that is not at its
			// maximum (n-1).
			for i = combiLen - 1; i >= 0 && indices[i] == n-1; i-- {
			}

			// If i is negative, then the indices are all at their maximum value and
			// we're done.
			if i < 0 {
				return
			}

			// Increment the current index which we know is not at its maximum. Then
			// set all to the right to the same value.
			indices[i]++
			elem := indices[i]
			for j := i + 1; j < combiLen; j++ {
				indices[j] = elem
			}

			out <- atuple.Pack{{ .N }}({{ .RepeatOneToN "s[indices[__N0]]" }})
		}
	}()
	return out
}
`
