package oniichan

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type TestCase interface {
	SUT() interface{}
	Name() string
	Arguments() []interface{}
	Expects() []interface{}
}

type singleReturnTestCase struct {
	sut       interface{}
	name      string
	arguments []interface{}
	expects   interface{}
}

func (tc *singleReturnTestCase) SUT() interface{}         { return tc.sut }
func (tc *singleReturnTestCase) Name() string             { return tc.name }
func (tc *singleReturnTestCase) Arguments() []interface{} { return tc.arguments }
func (tc *singleReturnTestCase) Expects() []interface{}   { return []interface{}{tc.expects} }

func RunTestCase(t *testing.T, tc TestCase) {
	argsReflect := make([]reflect.Value, 0, len(tc.Arguments()))
	for _, arg := range tc.Arguments() {
		argsReflect = append(argsReflect, reflect.ValueOf(arg))
	}

	ret := reflect.ValueOf(tc.SUT()).Call(argsReflect)
	if len(tc.Expects()) != len(ret) {
		t.Errorf("Test failed [%s]: expected number of return values: want %d but got %d",
			tc.Name(), len(tc.Expects()), len(ret))
	}
	for i, retVal := range tc.Expects() {
		if reflect.TypeOf(retVal).Kind() == reflect.Chan {
			err := AssertChanEqual(t, retVal, ret[i].Interface())
			if err != nil {
				t.Errorf("Test failed [%s]: unexpected %d-th return value: %v",
					tc.Name(), i+1, err)
			}
			continue
		}
		// if reflect.TypeOf(retVal) != ret[i].Type() {
		// 	t.Errorf("Test failed [%s]: expected %d-th return value to be have type [%T] but got type [%s]",
		// 		tc.Name(), i+1, retVal, ret[i].Type())
		// }
		if !reflect.DeepEqual(retVal, ret[i].Interface()) {
			t.Errorf("Test failed [%s]: unexpected %d-th return value: want [%v](%T) but got [%v](%s)",
				tc.Name(), i+1, retVal, retVal, ret[i], ret[i].Type())
		}
	}
}

func AssertChanEqual(t *testing.T, expectedCh, actualCh interface{}) error {

	if reflect.TypeOf(expectedCh).Kind() != reflect.Chan {
		panic(fmt.Errorf("invalid argument to AssertChanEqual: expectedCh must "+
			"be a channel, but got %T", expectedCh))
	}

	if reflect.TypeOf(actualCh).Kind() != reflect.Chan {
		return fmt.Errorf("want [%v](%T) but got [%v](%T)",
			expectedCh, expectedCh, actualCh, actualCh)
	}

	var i int
	for {
		eReflect, ok := reflect.ValueOf(expectedCh).Recv()
		if !ok {
			break
		}

		a, err := TryRecvTimeout(actualCh, 5*time.Second)
		if err != nil {
			return fmt.Errorf("expected %d-th element of channel to be [%v](%T) but got error receiving from channel: %w",
				i+1, eReflect, eReflect, err)
		}

		eInterface := eReflect.Interface()
		if !reflect.DeepEqual(eInterface, a) {
			return fmt.Errorf("expected %d-th element of channel to be [%v](%T) but got [%v][%T]",
				i+1, eInterface, eInterface, a, a)
		}
		i++
	}
	return nil
}

func TryRecvTimeout(ch interface{}, timeout time.Duration) (elem interface{}, err error) {
	resultCh := make(chan interface{}, 1)
	alreadyClosed := make(chan struct{})
	go func() {
		actual, ok := reflect.ValueOf(ch).Recv()
		if !ok { // channel is already closed
			close(alreadyClosed)
		}
		resultCh <- actual.Interface()
	}()
	select {
	case elem = <-resultCh:
		return elem, nil
	case <-time.After(timeout):
		return elem, fmt.Errorf("receive timed-out after %v", timeout)
	case <-alreadyClosed:
		return elem, fmt.Errorf("channel already closed")
	}
	panic("unreachable")
}
