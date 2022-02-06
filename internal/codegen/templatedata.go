package main

import (
	"fmt"
	"strings"
)

// TemplateData provides convenience data and functions for writiing
// text/templates that will each expand into multiple functions, with
// increasing number of certain language element: type parameter, input
// argument, return value, element in return value.
//
// Refer to the documentaton on each of the fields and methods for details.
type TemplateData struct {
	N int // maximum number of elements (will generate functions from 2, 3, ..., to N)

	T1_T2___TN string // T1, T2, ..., TN, for use in template parameter lists

	OneToN []int // []int{ 1, 2, ..., N }, for use with the range keyword in text/template
	TwoToN []int // []int{ 2, ..., N }, for use with the range keyword in text/template
}

// newTemplateData creates a new TemplateData with the various fields
// initialized using N.
func newTemplateData(N int) TemplateData {
	td := TemplateData{N: N}

	typeParams := make([]string, N)
	for i := 0; i < N; i++ {
		typeParams[i] = fmt.Sprintf("T%d", i+1)
	}

	td.T1_T2___TN = strings.Join(typeParams, ", ")

	oneToN := make([]int, N)
	for i := 0; i < N; i++ {
		oneToN[i] = i + 1
	}
	td.OneToN = oneToN
	td.TwoToN = oneToN[1:]

	return td
}

// RepeatOneToN repeats the string s passed in from 1 to N (inclusive),
// replacing occurences of "__N0" with the 0-based index, and "__N1" with the
// 1-based index.
//
// For example, the following line within a template
//     {{ .RepeatOneToN "mySlice[__N0] = myElem__N1" }}
// produces
//     "mySlice[0] = myElem1, mySlice[1] = myElem2, ..., mySlice[N-1] = myElemN "
func (td TemplateData) RepeatOneToN(s string) string {
	out := make([]string, td.N)
	for i := 0; i < td.N; i++ {
		replacedN0 := strings.ReplaceAll(s, "__N0", fmt.Sprint(i))
		replacedN1 := strings.ReplaceAll(replacedN0, "__N1", fmt.Sprint(i+1))
		out[i] = replacedN1
	}
	return strings.Join(out, ", ")
}

// RepeatOneToNWithSep repeats the provided string s in the same way as
// RepeatOneToN, except it uses the provided sep as the separator instead of
// comma.
func (td TemplateData) RepeatOneToNWithSep(s string, sep string) string {
	out := make([]string, td.N)
	for i := 0; i < td.N; i++ {
		replacedN0 := strings.ReplaceAll(s, "__N0", fmt.Sprint(i))
		replacedN1 := strings.ReplaceAll(replacedN0, "__N1", fmt.Sprint(i+1))
		out[i] = replacedN1
	}
	return strings.Join(out, sep)
}

// RepeatTwoToN repeats the provided string s in the same way as
// RepeatOneToN, except it starts from i = 2.
func (td TemplateData) RepeatTwoToN(s string) string {
	out := make([]string, td.N)
	for i := 0; i < td.N; i++ {
		replacedN0 := strings.ReplaceAll(s, "__N0", fmt.Sprint(i))
		replacedN1 := strings.ReplaceAll(replacedN0, "__N1", fmt.Sprint(i+1))
		out[i] = replacedN1
	}
	return strings.Join(out[1:], ", ")
}
