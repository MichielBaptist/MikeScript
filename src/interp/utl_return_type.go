package interp

import (
	"fmt"
	"strings"
)

////////////////////////////////////////////////////////////////////////
// Evaluation result
////////////////////////////////////////////////////////////////////////

type EvalResult struct {
	rt    ResultType  	// type of the result
	val   any 			// Container for the result
	err   []error    	// Error message on evaluation fail
}

func (er *EvalResult) Valid() bool {
	return len(er.err) == 0
}

func (er *EvalResult) Expect(rt ResultType) bool {
	return er.rt == rt
}

func (er *EvalResult) ExpectedType(rt ResultType) bool {
	return er.rt == rt && er.ValidReturnType()
}

func (er *EvalResult) ValidReturnType() bool {
	switch er.rt {
	case RT_INT:
		_, ok := er.val.(int)
		return ok
	case RT_FLOAT:
		_, ok := er.val.(float64)
		return ok
	case RT_STRING:
		_, ok := er.val.(string)
		return ok
	case RT_BOOL:
		_, ok := er.val.(bool)
		return ok
	case RT_TUPLE:
		_, ok := er.val.([]EvalResult)
		return ok
	case RT_FUNCTION:
		return true
	case RT_INVALID:
		return false
	case RT_NOTHING:
		return true
	default:
		fmt.Println("Error trying to check valid return type for unknown type:", er.rt)
		return false
	}
}

func (er *EvalResult) ReturnType() ResultType {
	return er.rt
}

func (er EvalResult) String() string {

	// Check if the EvalResult has any errors,
	// if so, we print the errors
	if !er.Valid() {
		return fmt.Sprintf("Error: %v", er.err)
	}

	// based on return type
	switch er.rt {
	case RT_INT, RT_FLOAT, RT_BOOL, RT_FUNCTION:
		return fmt.Sprintf("%v", er.val)
	case RT_STRING:
		return fmt.Sprintf("\"%v\"", er.val)
	case RT_TUPLE:
		return tupleToString(er.val.([]EvalResult))
	default:
		return er.rt.String()	// see ReturnType .String()
	}
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

type ResultType uint8
const (
	RT_INVALID ResultType = iota
	RT_NOTHING
	RT_INT
	RT_FLOAT
	RT_STRING
	RT_BOOL
	RT_TUPLE
	RT_FUNCTION
	RT_BREAK
	RT_CONTINUE
)

func (res ResultType) String() string {
	switch res {

	// Real types
	case RT_INT:		return "int"
	case RT_FLOAT:		return "float"
	case RT_STRING:		return "string"
	case RT_BOOL:		return "bool"

	// composite types
	case RT_TUPLE:		return "tuple"
	case RT_FUNCTION:	return "function"

	// statement types (not a real type, but used by interp)
	case RT_INVALID:	return "invalid"
	case RT_NOTHING:	return "nothing"
	case RT_BREAK:		return "break"
	case RT_CONTINUE:	return "continue"
	default:			return "UNKNOWN -- check expression.go"
	}
}

func evalErr(msg string) EvalResult {
	return EvalResult{err: []error{&EvalError{msg}}}
}


/////
// types
////
type MSType interface {
	eq(other *MSType) bool
}

type MSSimpleTypeS struct {
	rt ResultType
}

func (t *MSSimpleTypeS) eq(o *MSType) bool {
	switch other := (*o).(type){
	case *MSSimpleTypeS: 	return t.rt == other.rt
	default:				return false
	}
}

type MSCompositeTypeS struct {
	left MSType
	right MSType
}

func (t *MSCompositeTypeS) eq(o *MSType) bool {
	switch other := (*o).(type) {
	case *MSCompositeTypeS:	return t.left.eq(&other.left) && t.right.eq(&other.right)
	default:				return false
	}
}

type MSOperationTypeS struct {
	left []MSType
	right MSType // result of operation
}

func (t *MSOperationTypeS) eq(o *MSType) bool {
	switch other := (*o).(type) {
	case *MSOperationTypeS:

		// Compare output of op
		e := t.right.eq(&other.right)

		// Compare parameters of op, note this
		// circuit breaks on a single false
		e = e && len(t.left) == len(other.left)
		for i := 0; e && i < len(t.left) ; i++ {
			e = e && t.left[i].eq(&other.left[i])
		}
		return e
	default:
		return false
	}
}