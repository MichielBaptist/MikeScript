package interp

import (
	"fmt"
	"mikescript/src/mstype"
	"strings"
)

////////////////////////////////////////////////////////////////////////
// Evaluation result
////////////////////////////////////////////////////////////////////////

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

// func (er *EvalResult) ValidReturnType() bool {
// 	switch er.rt {
// 	case mstype.RT_INT:
// 		_, ok := er.val.(int)
// 		return ok
// 	case mstype.RT_FLOAT:
// 		_, ok := er.val.(float64)
// 		return ok
// 	case mstype.RT_STRING:
// 		_, ok := er.val.(string)
// 		return ok
// 	case mstype.RT_BOOL:
// 		_, ok := er.val.(bool)
// 		return ok
// 	case mstype.RT_TUPLE:
// 		_, ok := er.val.([]EvalResult)
// 		return ok
// 	case mstype.RT_FUNCTION:
// 		return true
// 	case mstype.RT_INVALID:
// 		return false
// 	case mstype.RT_NOTHING:
// 		return true
// 	default:
// 		fmt.Println("Error trying to check valid return type for unknown type:", er.rt)
// 		return false
// 	}
// }

func (er EvalResult) String() string {

	// Check if the EvalResult has any errors,
	// if so, we print the errors
	if !er.Valid() {
		return fmt.Sprintf("Error: %v", er.err)
	}

	// We expect to see just the value, so we delegate
	// all string to .val

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

////////////////////////////////////////////////////////////////////////
// Result type
////////////////////////////////////////////////////////////////////////

// type ResultType uint8
// const (
// 	mstype.RT_INVALID ResultType = iota
// 	mstype.RT_NOTHING
// 	mstype.RT_INT
// 	mstype.RT_FLOAT
// 	mstype.RT_STRING
// 	mstype.RT_BOOL
// 	mstype.RT_TUPLE
// 	mstype.RT_FUNCTION
// 	mstype.RT_BREAK
// 	mstype.RT_CONTINUE
// )

// func (res ResultType) String() string {
// 	switch res {

// 	// Real types
// 	case mstype.RT_INT:		return "int"
// 	case mstype.RT_FLOAT:		return "float"
// 	case mstype.RT_STRING:		return "string"
// 	case mstype.RT_BOOL:		return "bool"

// 	// composite types
// 	case mstype.RT_TUPLE:		return "tuple"
// 	case mstype.RT_FUNCTION:	return "function"

// 	// statement types (not a real type, but used by interp)
// 	case mstype.RT_INVALID:	return "invalid"
// 	case mstype.RT_NOTHING:	return "nothing"
// 	case mstype.RT_BREAK:		return "break"
// 	case mstype.RT_CONTINUE:	return "continue"
// 	default:			return "UNKNOWN -- check expression.go"
// 	}
// }

func evalErr(msg string) EvalResult {
	return EvalResult{err: []error{&EvalError{msg}}}
}
