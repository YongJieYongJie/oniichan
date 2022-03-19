package oniichan

import (
	"fmt"
	"sync"
	"testing"

	"golang.org/x/exp/constraints"

	"github.com/YongJieYongJie/tttuples/atuple"
)

func TestAll(t *testing.T) {
	testCases := []*singleReturnTestCase{
		{
			sut:  All[int64],
			name: "basic (int64, returns true)",
			arguments: []interface{}{
				ChanFrom[int64](2, 4, 6),
				func(n int64) bool { return n%2 == 0 },
			},
			expects: true,
		},
		{
			sut:  All[[]string],
			name: "basic ([]string, returns false)",
			arguments: []interface{}{
				ChanFrom([]string{"a", "a"}, []string{"b", "b"}, []string{"c", "c", "c"}),
				func(n []string) bool { return len(n)%2 == 0 }},
			expects: false,
		},
	}
	for _, c := range testCases {
		RunTestCase(t, c)
	}
}

func TestAny(*testing.T) {} // Skipped because Any delegates to All.

func TestFilter(*testing.T) {} // Skipped because Filter delegates to FilterFalse.

func TestMax(t *testing.T) {
	testCases := []*singleReturnTestCase{
		{
			sut:  Max[int64],
			name: "basic (int64)",
			arguments: []interface{}{
				ChanFrom[int64](2, 4, 6, 3),
			},
			expects: int64(6),
		},
		{
			sut:  Max[float64],
			name: "basic (float64)",
			arguments: []interface{}{
				ChanFrom[float64](2, 4, 6, 3),
			},
			expects: float64(6),
		},
	}

	for _, c := range testCases {
		RunTestCase(t, c)
	}
}

func TestMaxFunc(t *testing.T) {
	testCases := []*singleReturnTestCase{
		{
			sut:  MaxFunc[string, int],
			name: "",
			arguments: []interface{}{
				ChanFrom("a", "aaa", "aa"),
				func(s string) int { return len(s) },
			},
			expects: "aaa",
		},
	}

	for _, c := range testCases {
		RunTestCase(t, c)
	}
}

func TestMin(t *testing.T) {
	testCases := []*singleReturnTestCase{
		{
			sut:  Min[int64],
			name: "basic (int64)",
			arguments: []interface{}{
				ChanFrom[int64](2, 4, 6, 3),
			},
			expects: int64(2),
		},
		{
			sut:  Min[float64],
			name: "basic (float64)",
			arguments: []interface{}{
				ChanFrom[float64](2, 4, 6, 3),
			},
			expects: float64(2),
		},
	}

	for _, c := range testCases {
		RunTestCase(t, c)
	}
}

func TestMinFunc(t *testing.T) {
	testCases := []*singleReturnTestCase{
		{
			sut:  MinFunc[string, int],
			name: "basic (string, int)",
			arguments: []interface{}{
				ChanFrom("a", "aaa", "aa"),
				func(s string) int { return len(s) },
			},
			expects: "a",
		},
	}

	for _, c := range testCases {
		RunTestCase(t, c)
	}
}

func TestSum(t *testing.T) {
	testCases := []*singleReturnTestCase{
		{
			sut:  Sum[int64],
			name: "basic (int64)",
			arguments: []interface{}{
				ChanFrom[int64](1, 2, 3),
			},
			expects: int64(6),
		},
	}

	for _, c := range testCases {
		RunTestCase(t, c)
	}
}

func TestReduce(t *testing.T) {
	testCases := []*singleReturnTestCase{
		{
			sut:  Reduce[int64],
			name: "basic (int64)",
			arguments: []interface{}{
				ChanFrom[int64](1, 2, 3, 4),
				func(v1, v2 int64) int64 { return v1 * v2 },
			},
			expects: int64(24),
		},
	}

	for _, c := range testCases {
		RunTestCase(t, c)
	}
}

func TestSorted(t *testing.T) {
	testCases := []*singleReturnTestCase{
		{
			sut:  Sorted[string, int],
			name: "basic (string, int)",
			arguments: []interface{}{
				ChanFrom("a", "aaa", "aa"),
				func(s string) int { return len(s) },
			},
			expects: ChanFrom("a", "aa", "aaa"),
		},
	}

	for _, c := range testCases {
		RunTestCase(t, c)
	}
}

func TestReversed(t *testing.T) {
	testCases := []*singleReturnTestCase{
		{
			sut:  Reversed[int64],
			name: "basic (int64)",
			arguments: []interface{}{
				ChanFrom[int64](1, 2, 3),
			},
			expects: ChanFrom[int64](3, 2, 1),
		},
	}
	for _, c := range testCases {
		RunTestCase(t, c)
	}
}

func TestEnumerate(t *testing.T) {
	testCases := []*singleReturnTestCase{
		{
			sut:  Enumerate[int64],
			name: "basic (int64)",
			arguments: []interface{}{
				ChanFrom[int64](6, 6, 6),
			},
			expects: ChanFrom(
				atuple.Pack2[int, int64](0, 6),
				atuple.Pack2[int, int64](1, 6),
				atuple.Pack2[int, int64](2, 6),
			),
		},
	}

	for _, c := range testCases {
		RunTestCase(t, c)
	}
}

func TestRangeFromTo(t *testing.T) {
	type testCase struct {
		name       string
		inputStart int64
		inputStop  int64
		expect     []int64
	}

	case1 := testCase{
		name:       "basic",
		inputStart: 5,
		inputStop:  10,
		expect:     []int64{5, 6, 7, 8, 9},
	}

	actual := RangeFromTo(case1.inputStart, case1.inputStop)
	for _, e := range case1.expect {
		a := <-actual
		if a != e {
			t.Errorf("Expected %v but got %v", e, a)
		}
	}
}

func TestRangeStep(t *testing.T) {
	type testCase struct {
		name       string
		inputStart int64
		inputStop  int64
		inputStep  int64
		expect     []int64
	}

	case1 := testCase{
		name:       "basic",
		inputStart: 5,
		inputStop:  10,
		inputStep:  2,
		expect:     []int64{5, 7, 9},
	}

	actual := RangeStep(case1.inputStart, case1.inputStop, case1.inputStep)
	for _, e := range case1.expect {
		a := <-actual
		if a != e {
			t.Errorf("Expected %v but got %v", e, a)
		}
	}
}

func TestMap(t *testing.T) {
	type testCase struct {
		name   string
		inputF func(int64) string
		inputC <-chan int64
		expect []string
	}

	case1 := testCase{
		name:   "basic",
		inputF: func(v int64) string { return fmt.Sprint(v) },
		inputC: RangeFromTo[int64](5, 10),
		expect: []string{"5", "6", "7", "8", "9"},
	}

	actual := Map(case1.inputF, case1.inputC)
	for _, e := range case1.expect {
		a := <-actual
		if a != e {
			t.Errorf("Expected %v but got %v", e, a)
		}
	}
}

func TestCount(t *testing.T) {
	type testCaseCount[T constraints.Ordered] struct {
		name       string
		inputStart T
		inputStep  T
		expect     []T
	}

	case1 := testCaseCount[int64]{
		name:       "basic",
		inputStart: 3,
		inputStep:  -2,
		expect:     []int64{3, 1, -1, -3},
	}

	ch := CountStep(case1.inputStart, case1.inputStep)
	for _, v := range case1.expect {
		actual := <-ch
		if v == actual {
			continue
		}
		t.Errorf("Case %s: Expected element %v (in %v) but got %v.",
			case1.name, v, case1.expect, actual)
	}
}

func TestCycle(t *testing.T) {
	type testCaseCycle struct {
		name   string
		input  <-chan int64
		expect []int64
	}

	case1 := testCaseCycle{
		name:   "basic",
		input:  Chan([]int64{2, 4, 6}),
		expect: []int64{2, 4, 6, 2, 4, 6, 2},
	}

	ch := Cycle(case1.input)
	for _, v := range case1.expect {
		actual := <-ch
		if v == actual {
			continue
		}
		t.Errorf("Case %s: Expected element %v (in %v) but got %v.",
			case1.name, v, case1.expect, actual)
	}
}

func TestRepeat(t *testing.T) {
	type testCaseRepeat[T any] struct {
		name   string
		input  T
		expect []T
	}

	case1 := testCaseRepeat[int64]{
		name:   "basic",
		input:  4,
		expect: []int64{4, 4, 4, 4},
	}

	ch := Repeat(case1.input)
	for _, v := range case1.expect {
		actual := <-ch
		if v == actual {
			continue
		}
		t.Errorf("Case %s: Expected element %v (in %v) but got %v.",
			case1.name, v, case1.expect, actual)
	}
}

func TestAccumulate(t *testing.T) {
	type testCase struct {
		name   string
		input  <-chan int64
		expect []int64
	}

	case1 := testCase{
		name:   "basic",
		input:  Chan([]int64{1, 1, 1, 1}),
		expect: []int64{1, 2, 3, 4},
	}

	actual := Accumulate(case1.input)
	for _, e := range case1.expect {
		a := <-actual
		if a != e {
			t.Errorf("Expected %v but got %v", e, a)
		}
	}
}

func TestChain(t *testing.T) {
	type testCase struct {
		name   string
		input  []<-chan int64
		expect []int64
	}
	case1 := testCase{
		name: "basic",
		input: []<-chan int64{
			Chan([]int64{1, 1, 1}),
			Chan([]int64{2, 2}),
			Chan([]int64{3, 3, 3}),
		},
		expect: []int64{1, 1, 1, 2, 2, 3, 3, 3},
	}
	actual := Chain(case1.input...)
	for _, e := range case1.expect {
		a := <-actual
		if a != e {
			t.Errorf("Expected %v but got %v", e, a)
		}
	}
}

func TestChainFromChan(t *testing.T) {
	type testCase struct {
		name   string
		input  <-chan <-chan int64
		expect []int64
	}
	case1 := testCase{
		name: "basic",
		input: Chan([]<-chan int64{
			Chan([]int64{1, 1, 1}),
			Chan([]int64{2, 2}),
			Chan([]int64{3, 3, 3}),
		}),
		expect: []int64{1, 1, 1, 2, 2, 3, 3, 3},
	}
	actual := ChainFromChan(case1.input)
	for _, e := range case1.expect {
		a := <-actual
		if a != e {
			t.Errorf("Expected %v but got %v", e, a)
		}
	}
}

func TestCompress(t *testing.T) {
	type testCase struct {
		name   string
		input1 <-chan int64
		input2 <-chan bool
		expect []int64
	}

	case1 := testCase{
		name:   "basic",
		input1: Chan([]int64{1, 2, 3, 4, 5, 6, 7}),
		input2: Chan([]bool{true, false, true, false, true}),
		expect: []int64{1, 3, 5},
	}

	actual := Compress(case1.input1, case1.input2)
	for _, e := range case1.expect {
		a := <-actual
		if a != e {
			t.Errorf("Expected %v but got %v", e, a)
		}
	}
}

func TestDropWhile(t *testing.T) {
	type testCase struct {
		name      string
		inputPred func(int64) bool
		inputC    <-chan int64
		expect    []int64
	}

	case1 := testCase{
		name:      "basic",
		inputPred: func(v int64) bool { return v%2 == 0 },
		inputC:    Chan([]int64{1, 2, 3, 4, 5, 6}),
		expect:    []int64{2, 4, 6},
	}

	actual := DropWhile(case1.inputPred, case1.inputC)
	for _, e := range case1.expect {
		a := <-actual
		if a != e {
			t.Errorf("Expected %v but got %v", e, a)
		}
	}
}

func TestFilterFalse(t *testing.T) {
	type testCase struct {
		name      string
		inputPred func(int64) bool
		inputC    <-chan int64
		expect    []int64
	}

	case1 := testCase{
		name:      "basic",
		inputPred: func(v int64) bool { return v%2 == 0 },
		inputC:    RangeFromTo[int64](5, 10),
		expect:    []int64{5, 7, 9},
	}

	actual := FilterFalse(case1.inputPred, case1.inputC)
	for _, e := range case1.expect {
		a := <-actual
		if a != e {
			t.Errorf("Expected %v but got %v", e, a)
		}
	}
}

func TestGroupBy(t *testing.T) {
	type testCase struct {
		name   string
		input  <-chan int64
		expect []atuple.Packed2[int64, <-chan int64]
	}

	case1 := testCase{
		name:  "basic",
		input: Chan([]int64{1, 1, 1, 2, 2, 3, 3, 1, 1}),
		expect: []atuple.Packed2[int64, <-chan int64]{
			atuple.Pack2[int64](1, Chan([]int64{1, 1, 1})),
			atuple.Pack2[int64](2, Chan([]int64{2, 2})),
			atuple.Pack2[int64](3, Chan([]int64{3, 3})),
			atuple.Pack2[int64](1, Chan([]int64{1, 1})),
		},
	}

	actual := GroupBy(case1.input)
	for i, e := range case1.expect {
		a := <-actual
		actualKey := a.Item1()
		expectedKey := e.Item1()
		if actualKey != expectedKey {
			t.Errorf("%d Expected %v but got %v", i, expectedKey, actualKey)
		}
		for expectedElem := range e.Item2() {
			actualElem := <-a.Item2()
			if actualElem != expectedElem {
				t.Errorf("%d Expected %v but got %v", i, expectedElem, actualElem)
			}
		}
	}
}

func TestGroupByKey(t *testing.T) {
	type testCase struct {
		name         string
		inputKeyFunc func(int64) string
		input        <-chan int64
		expect       []atuple.Packed2[string, <-chan int64]
	}

	case1 := testCase{
		name:         "basic",
		inputKeyFunc: func(v int64) string { return fmt.Sprint(v) },
		input:        Chan([]int64{1, 1, 1, 2, 2, 3, 3, 1, 1}),
		expect: []atuple.Packed2[string, <-chan int64]{
			atuple.Pack2("1", Chan([]int64{1, 1, 1})),
			atuple.Pack2("2", Chan([]int64{2, 2})),
			atuple.Pack2("3", Chan([]int64{3, 3})),
			atuple.Pack2("1", Chan([]int64{1, 1})),
		},
	}

	actual := GroupByKey(case1.input, case1.inputKeyFunc)
	for _, e := range case1.expect {
		a := <-actual
		actualKey := a.Item1()
		expectedKey := e.Item1()
		if actualKey != expectedKey {
			t.Errorf("Expected %v but got %v", expectedKey, actualKey)
		}
		for expectedElem := range e.Item2() {
			actualElem := <-a.Item2()
			if actualElem != expectedElem {
				t.Errorf("Expected %v but got %v", expectedElem, actualElem)
			}
		}
	}
}

func TestSliceTo(t *testing.T) {
	type testCase struct {
		name      string
		inputC    <-chan int64
		inputStop int
		expect    []int64
	}

	case1 := testCase{
		name:      "basic",
		inputC:    RangeFromTo[int64](2, 8),
		inputStop: 4,
		expect:    []int64{2, 3, 4, 5},
	}

	actual := SliceTo(case1.inputC, case1.inputStop)
	for _, e := range case1.expect {
		a := <-actual
		if a != e {
			t.Errorf("Expected %v but got %v", e, a)
		}
	}
	for e := range actual {
		t.Errorf("Expect channel to be closed, but got %v", e)
	}
}

func TestPairwise(t *testing.T) {
	type testCase struct {
		name   string
		input  <-chan int64
		expect []atuple.Packed2[int64, int64]
	}

	case1 := testCase{
		name:  "basic",
		input: RangeFromTo[int64](10, 15),
		expect: []atuple.Packed2[int64, int64]{
			atuple.Pack2[int64, int64](10, 11),
			atuple.Pack2[int64, int64](11, 12),
			atuple.Pack2[int64, int64](12, 13),
			atuple.Pack2[int64, int64](13, 14),
		},
	}

	actual := Pairwise(case1.input)
	for _, e := range case1.expect {
		a := <-actual
		elem1 := a.Item1()
		if elem1 != e.Item1() {
			t.Errorf("Mismatched first elem: Expected %v but got %v", e, a)
		}
		if a.Item2() != e.Item2() {
			t.Errorf("Mismatched second elem: Expected %v but got %v", e, a)
		}
	}
}

func TestStarMap(t *testing.T) {
	type testCase struct {
		name   string
		inputF func(float64, int64) string
		inputC <-chan atuple.Packed2[float64, int64]
		expect []string
	}

	case1 := testCase{
		name:   "basic",
		inputF: func(v1 float64, v2 int64) string { return fmt.Sprint(v1 * float64(v2)) },
		inputC: Chan([]atuple.Packed2[float64, int64]{
			atuple.Pack2[float64, int64](1.1, 2),
			atuple.Pack2[float64, int64](2.2, 4),
			atuple.Pack2[float64, int64](3.3, 4),
		}),
		expect: []string{"2.2", "8.8", "13.2"},
	}

	actual := StarMap2(case1.inputF, case1.inputC)
	for _, e := range case1.expect {
		a := <-actual
		if a != e {
			t.Errorf("Expected %v but got %v", e, a)
		}
	}
}

func TestStarMapN(t *testing.T) {
	type testCase struct {
		name   string
		inputF func(float64, int64) string
		inputC <-chan []any
		expect []string
	}

	case1 := testCase{
		name:   "basic",
		inputF: func(v1 float64, v2 int64) string { return fmt.Sprint(v1 * float64(v2)) },
		inputC: Chan([][]any{
			{1.1, int64(2)},
			{2.2, int64(4)},
			{3.3, int64(4)},
		}),
		expect: []string{"2.2", "8.8", "13.2"},
	}

	actual := StarMapN[string](case1.inputF, case1.inputC)
	for _, e := range case1.expect {
		a := <-actual
		if a != e {
			t.Errorf("Expected %v but got %v", e, a)
		}
	}
}

func TestTakeWhile(t *testing.T) {
	type testCase struct {
		name      string
		inputPred func(int64) bool
		inputC    <-chan int64
		expect    []int64
	}

	case1 := testCase{
		name:      "basic",
		inputPred: func(v int64) bool { return v < 8 },
		inputC:    RangeFromTo[int64](5, 10),
		expect:    []int64{5, 6, 7},
	}

	actual := TakeWhile(case1.inputPred, case1.inputC)
	for _, e := range case1.expect {
		a := <-actual
		if a != e {
			t.Errorf("Expected %v but got %v", e, a)
		}
	}
}

func TestTee(t *testing.T) {
	type testCase struct {
		name   string
		input  <-chan int64
		expect []int64
	}

	case1 := testCase{
		name:   "basic",
		input:  RangeFromTo[int64](5, 10),
		expect: []int64{5, 6, 7, 8, 9},
	}

	actual1, actual2 := Tee(case1.input)
	for _, e := range case1.expect {
		var a1, a2 int64
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			a1 = <-actual1
		}()
		go func() {
			defer wg.Done()
			a2 = <-actual2
		}()
		wg.Wait()
		if a1 != e {
			t.Errorf("Expected %v but got %v", e, a1)
		}
		if a2 != e {
			t.Errorf("Expected %v but got %v", e, a2)
		}
	}
}

func TestZip(t *testing.T) {
	type testCase struct {
		name   string
		input1 <-chan int
		input2 <-chan string
		expect []atuple.Packed2[int, string]
	}

	case1 := testCase{
		name:   "basic",
		input1: Chan([]int{1, 2, 3, 4}),
		input2: Chan([]string{"a", "b", "c"}),
		expect: []atuple.Packed2[int, string]{
			atuple.Pack2(1, "a"), // Report issues with int / int64 type inferene
			atuple.Pack2(2, "b"),
			atuple.Pack2(3, "c")},
	}

	actual := Zip(case1.input1, case1.input2)
	for _, e := range case1.expect {
		a := <-actual
		if e[0] != a.Item1() || e[1] != a.Item2() {
			t.Errorf("Expected %v but got %v", e, a)
		}
	}
}

func TestZipLongest(t *testing.T) {
	type testCase struct {
		name            string
		input1          <-chan int
		input2          <-chan float64
		inputFillValues atuple.Packed2[int, float64]
		expect          []atuple.Packed2[int, float64]
	}

	case1 := testCase{
		name:            "basic",
		input1:          Chan([]int{1, 2, 3}),
		input2:          Chan([]float64{1.1, 2.2, 3.3, 4.4, 5.5}),
		inputFillValues: atuple.Pack2(9, 9.9),
		expect: []atuple.Packed2[int, float64]{
			atuple.Pack2(1, 1.1),
			atuple.Pack2(2, 2.2),
			atuple.Pack2(3, 3.3),
			atuple.Pack2(9, 4.4),
			atuple.Pack2(9, 5.5),
		},
	}

	actual := ZipLongest2(case1.input1, case1.input2, case1.inputFillValues)
	for _, e := range case1.expect {
		a := <-actual
		if a.Item1() != e.Item1() {
			t.Errorf("Mismatched first elem: Expected %v but got %v", e, a)
		}
		if a.Item2() != e.Item2() {
			t.Errorf("Mismatched second elem: Expected %v but got %v", e, a)
		}
	}
}

func TestProduct(t *testing.T) {
	type testCase struct {
		name   string
		input1 <-chan int
		input2 <-chan float64
		expect []atuple.Packed2[int, float64]
	}

	case1 := testCase{
		name:   "basic",
		input1: Chan([]int{1, 2, 3}),
		input2: Chan([]float64{1.1, 2.2}),
		expect: []atuple.Packed2[int, float64]{
			atuple.Pack2(1, 1.1),
			atuple.Pack2(1, 2.2),
			atuple.Pack2(2, 1.1),
			atuple.Pack2(2, 2.2),
			atuple.Pack2(3, 1.1),
			atuple.Pack2(3, 2.2),
		},
	}

	actual := Product(case1.input1, case1.input2)
	for _, e := range case1.expect {
		a := <-actual
		if a.Item1() != e.Item1() {
			t.Errorf("Mismatched first elem: Expected %v but got %v", e, a)
		}
		if a.Item2() != e.Item2() {
			t.Errorf("Mismatched second elem: Expected %v but got %v", e, a)
		}
	}
}

func TestPermutations(t *testing.T) {
	type testCase struct {
		name   string
		input  <-chan int
		expect []atuple.Packed2[int, int]
	}
	case1 := testCase{
		name:  "basic",
		input: Chan([]int{1, 2, 3}),
		expect: []atuple.Packed2[int, int]{
			atuple.Pack2(1, 2),
			atuple.Pack2(1, 3),
			atuple.Pack2(2, 1),
			atuple.Pack2(2, 3),
			atuple.Pack2(3, 1),
			atuple.Pack2(3, 2),
		},
	}
	actual := Permutations(case1.input)
	for _, e := range case1.expect {
		a := <-actual
		if a.Item1() != e.Item1() {
			t.Errorf("Mismatched first elem: Expected %v but got %v", e, a)
		}
		if a.Item2() != e.Item2() {
			t.Errorf("Mismatched second elem: Expected %v but got %v", e, a)
		}
	}
}

func TestCombinations(t *testing.T) {
	type testCase struct {
		name   string
		input  <-chan int
		expect []atuple.Packed2[int, int]
	}
	case1 := testCase{
		name:  "basic",
		input: Chan([]int{1, 2, 3}),
		expect: []atuple.Packed2[int, int]{
			atuple.Pack2(1, 2),
			atuple.Pack2(1, 3),
			atuple.Pack2(2, 3),
		},
	}
	actual := Combinations(case1.input)
	for _, e := range case1.expect {
		a := <-actual
		if a.Item1() != e.Item1() {
			t.Errorf("Mismatched first elem: Expected %v but got %v", e, a)
		}
		if a.Item2() != e.Item2() {
			t.Errorf("Mismatched second elem: Expected %v but got %v", e, a)
		}
	}
}

func TestCombinationsWithReplacement(t *testing.T) {
	type testCase struct {
		name   string
		input  <-chan int
		expect []atuple.Packed2[int, int]
	}
	case1 := testCase{
		name:  "basic",
		input: Chan([]int{1, 2, 3}),
		expect: []atuple.Packed2[int, int]{
			atuple.Pack2(1, 1),
			atuple.Pack2(1, 2),
			atuple.Pack2(1, 3),
			atuple.Pack2(2, 2),
			atuple.Pack2(2, 3),
			atuple.Pack2(3, 3),
		},
	}
	actual := CombinationsWithReplacement(case1.input)
	for _, e := range case1.expect {
		a := <-actual
		if a.Item1() != e.Item1() {
			t.Errorf("Mismatched first elem: Expected %v but got %v", e, a)
		}
		if a.Item2() != e.Item2() {
			t.Errorf("Mismatched second elem: Expected %v but got %v", e, a)
		}
	}
}

func TestToChan(t *testing.T) {
	actual := Chan([]int64{1, 2, 3})
	expected := []int64{1, 2, 3}

	for _, e := range expected {
		a := <-actual
		if a != e {
			t.Errorf("Expected %v but got %v", e, a)
		}
	}
}
