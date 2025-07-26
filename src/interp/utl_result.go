package interp

import (
	"fmt"
	"mikescript/src/mstype"
	"strings"
)

type EvalResult struct {
	rt    mstype.MSType		  	// type of the result
	val   any 					// Container for the result
	err   []error    			// Error message on evaluation fail
}

func (er *EvalResult) Valid() bool {
	return len(er.err) == 0
}

func (er *EvalResult) Expect(rt mstype.MSType) bool {
	return er.rt.Eq(&rt)
}

func (er *EvalResult) IsType(t *mstype.MSType) bool {
	return er.rt.Eq(t)
}

func (er EvalResult) String() string {

	// Check if the EvalResult has any errors,
	// if so, we print the errors
	if !er.Valid() {
		return fmt.Sprintf("Error: %v", er.err)
	}
	
	// If we have a nil value for whatever reason we show nothing.
	if er.val == nil {
		return "nothing"
	}

	// based on return type
	switch tp := er.rt.(type) {
	case *mstype.MSSimpleTypeS:
		switch tp.Rt{
		case mstype.RT_STRING:	return fmt.Sprintf("\"%v\"", er.val)
		case mstype.RT_TUPLE:	return tupleToString(er.val.([]EvalResult))
		}
	}
	return fmt.Sprintf("%v", er.val)
}

func tupleToString(t []EvalResult) string {
	s := []string{}
	for _, v := range t {
		s = append(s, v.String())
	}
	return fmt.Sprintf("(%s)", strings.Join(s, ", "))
}

////////////////////////////////////////////////////////////////////////
// Evaluation error
////////////////////////////////////////////////////////////////////////

type EvalError struct {
	message string
}

func (ee *EvalError) Error() string {
	return "Evaluation error: " + ee.message
}

func evalErr(msg string) EvalResult {
	return EvalResult{err: []error{&EvalError{msg}}}
}
