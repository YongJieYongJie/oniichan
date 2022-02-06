// Code generated by codegen.go DO NOT EDIT.
// Version: 0.0.0
// Checksum: 1288455508127950766

package oniichan

import (
	"github.com/YongJieYongJie/tttuples/atuple"
)

// Product2 is the 2-element version of Product.
func Product2[T1, T2 any](ch1 <-chan T1, ch2 <-chan T2) <-chan atuple.Packed2[T1, T2] {
	out := make(chan atuple.Packed2[T1, T2])
  s1 := Slice(ch1)
	s2 := Slice(ch2)
	
	go func() {
		defer close(out)
		for _, e1 := range s1 {
		for _, e2 := range s2 {
		out <- atuple.Pack2(e1, e2)
		} }
	}()
	return out
}

// Product3 is the 3-element version of Product.
func Product3[T1, T2, T3 any](ch1 <-chan T1, ch2 <-chan T2, ch3 <-chan T3) <-chan atuple.Packed3[T1, T2, T3] {
	out := make(chan atuple.Packed3[T1, T2, T3])
  s1 := Slice(ch1)
	s2 := Slice(ch2)
	s3 := Slice(ch3)
	
	go func() {
		defer close(out)
		for _, e1 := range s1 {
		for _, e2 := range s2 {
		for _, e3 := range s3 {
		out <- atuple.Pack3(e1, e2, e3)
		} } }
	}()
	return out
}

// Product4 is the 4-element version of Product.
func Product4[T1, T2, T3, T4 any](ch1 <-chan T1, ch2 <-chan T2, ch3 <-chan T3, ch4 <-chan T4) <-chan atuple.Packed4[T1, T2, T3, T4] {
	out := make(chan atuple.Packed4[T1, T2, T3, T4])
  s1 := Slice(ch1)
	s2 := Slice(ch2)
	s3 := Slice(ch3)
	s4 := Slice(ch4)
	
	go func() {
		defer close(out)
		for _, e1 := range s1 {
		for _, e2 := range s2 {
		for _, e3 := range s3 {
		for _, e4 := range s4 {
		out <- atuple.Pack4(e1, e2, e3, e4)
		} } } }
	}()
	return out
}

// Product5 is the 5-element version of Product.
func Product5[T1, T2, T3, T4, T5 any](ch1 <-chan T1, ch2 <-chan T2, ch3 <-chan T3, ch4 <-chan T4, ch5 <-chan T5) <-chan atuple.Packed5[T1, T2, T3, T4, T5] {
	out := make(chan atuple.Packed5[T1, T2, T3, T4, T5])
  s1 := Slice(ch1)
	s2 := Slice(ch2)
	s3 := Slice(ch3)
	s4 := Slice(ch4)
	s5 := Slice(ch5)
	
	go func() {
		defer close(out)
		for _, e1 := range s1 {
		for _, e2 := range s2 {
		for _, e3 := range s3 {
		for _, e4 := range s4 {
		for _, e5 := range s5 {
		out <- atuple.Pack5(e1, e2, e3, e4, e5)
		} } } } }
	}()
	return out
}

// Product6 is the 6-element version of Product.
func Product6[T1, T2, T3, T4, T5, T6 any](ch1 <-chan T1, ch2 <-chan T2, ch3 <-chan T3, ch4 <-chan T4, ch5 <-chan T5, ch6 <-chan T6) <-chan atuple.Packed6[T1, T2, T3, T4, T5, T6] {
	out := make(chan atuple.Packed6[T1, T2, T3, T4, T5, T6])
  s1 := Slice(ch1)
	s2 := Slice(ch2)
	s3 := Slice(ch3)
	s4 := Slice(ch4)
	s5 := Slice(ch5)
	s6 := Slice(ch6)
	
	go func() {
		defer close(out)
		for _, e1 := range s1 {
		for _, e2 := range s2 {
		for _, e3 := range s3 {
		for _, e4 := range s4 {
		for _, e5 := range s5 {
		for _, e6 := range s6 {
		out <- atuple.Pack6(e1, e2, e3, e4, e5, e6)
		} } } } } }
	}()
	return out
}

// Product7 is the 7-element version of Product.
func Product7[T1, T2, T3, T4, T5, T6, T7 any](ch1 <-chan T1, ch2 <-chan T2, ch3 <-chan T3, ch4 <-chan T4, ch5 <-chan T5, ch6 <-chan T6, ch7 <-chan T7) <-chan atuple.Packed7[T1, T2, T3, T4, T5, T6, T7] {
	out := make(chan atuple.Packed7[T1, T2, T3, T4, T5, T6, T7])
  s1 := Slice(ch1)
	s2 := Slice(ch2)
	s3 := Slice(ch3)
	s4 := Slice(ch4)
	s5 := Slice(ch5)
	s6 := Slice(ch6)
	s7 := Slice(ch7)
	
	go func() {
		defer close(out)
		for _, e1 := range s1 {
		for _, e2 := range s2 {
		for _, e3 := range s3 {
		for _, e4 := range s4 {
		for _, e5 := range s5 {
		for _, e6 := range s6 {
		for _, e7 := range s7 {
		out <- atuple.Pack7(e1, e2, e3, e4, e5, e6, e7)
		} } } } } } }
	}()
	return out
}