package interp

import (
	"fmt"
	"math"
	"mikescript/src/utils"
)

// Defines call
type FunctionResult interface {
	call(evaluator *MSEvaluator, args []EvalResult) EvalResult
	max_arity() uint8	// Returns maximum amount of arguments
	min_arity() uint8 	// Returns minimum amount of arguments
}

///////////////////////////////////////////////////////////////
// mikescript builtins
///////////////////////////////////////////////////////////////
func MSBuiltinPrint() EvalResult {
	return EvalResult{
		rt: RT_FUNCTION,
		val: &PrintFunction{},
		err: []error{},
	}
}

///////////////////////////////////////////////////////////////
// Print function
///////////////////////////////////////////////////////////////

// Is a FunctionResult placeholder for print
type PrintFunction struct {}

func (pf *PrintFunction) call(_evaluator *MSEvaluator, args []EvalResult) EvalResult {

	// Get all result values in a slice
	// and convert them to strings
	strs := make([]string, len(args))
	for i, arg := range args {
		strs[i] = fmt.Sprint(arg.val)
	}

	// Print the joined strings, this is where
	// the connection between mikesript and go is made.
	fmt.Println(utils.StrJoin(strs, ", "))

	// Done without return value
	return EvalResult{rt: RT_NONE}
}

func (pf *PrintFunction) min_arity() uint8 {
	return 0
}

func (pf *PrintFunction) max_arity() uint8 {
	return math.MaxUint8
}

func (pf *PrintFunction) String() string {
	return "{...}"
}