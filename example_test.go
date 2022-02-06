package oniichan_test

import (
	"fmt"

	. "github.com/YongJieYongJie/oniichan"
	"github.com/YongJieYongJie/tttuples/atuple"
)

func ExampleRange() {
	for i := range Range(10) {
		fmt.Println(i)
	}
	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9
}

// This example demonstrates the use of StarMap to run tests.
func ExampleStarMap() {
	unitTest := func(input1 int, input2 string) string {
		return fmt.Sprintf("Called with input1: %d and input2 %s.", input1, input2)
	}

	testData := []atuple.Packed2[int, string]{ // note: type parameters required here
		{1, "A"},
		{2, "B"},
		{3, "A"},
	}

	results := StarMap(unitTest, ChanFrom(testData...))
	for r := range results {
		fmt.Println(r)
	}
	// Output:
	// Called with input1: 1 and input2 A.
	// Called with input1: 2 and input2 B.
	// Called with input1: 3 and input2 A.
}
