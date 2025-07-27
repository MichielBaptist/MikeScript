package interp

import (
	"fmt"
	"mikescript/src/mstype"
	"strings"
)

type EvalResult struct {
	Rt    mstype.MSType		  	// type of the result
	Val   any 					// Container for the result
	Err   []error    			// Error message on evaluation fail
}

func (er *EvalResult) Valid() bool {
	return len(er.Err) == 0
}

func (er *EvalResult) Expect(rt mstype.MSType) bool {
	return er.Rt.Eq(&rt)
}

func (er *EvalResult) IsType(t *mstype.MSType) bool {
	return er.Rt.Eq(t)
}

func (er EvalResult) String() string {

	// Check if the EvalResult has any errors,
	// if so, we print the errors
	if !er.Valid() {
		return fmt.Sprintf("Error: %v", er.Err)
	}
	
	// If we have a nil value for whatever reason we show nothing.
	if er.Val == nil {
		return "nothing"
	}

	// based on return type
	switch tp := er.Rt.(type) {
	case *mstype.MSSimpleTypeS:
		switch tp.Rt{
		case mstype.RT_STRING:	return fmt.Sprintf("\"%v\"", er.Val)
		case mstype.RT_TUPLE:	return tupleToString(er.Val.([]EvalResult))
		}
	}
	return fmt.Sprintf("%v", er.Val)
}

func tupleToString(t []EvalResult) string {
	s := []string{}
	for _, v := range t {
		s = append(s, v.String())
	}
	return fmt.Sprintf("(%s)", strings.Join(s, ", "))
}