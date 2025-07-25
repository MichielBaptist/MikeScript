package interp

import (
	"fmt"
	"math"
	"mikescript/src/mstype"
	"strings"
)

///////////////////////////////////////////////////////////////
// mikescript builtins
///////////////////////////////////////////////////////////////
func MSBuiltinPrint() EvalResult {
	return EvalResult{
		rt: &mstype.MSOperationTypeS{},
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
		strs[i] = arg.String()
	}

	// Print the joined strings, this is where
	// the connection between mikesript and go is made.
	fmt.Println(strings.Join(strs, ", "))

	// Done without return value
	return EvalResult{rt: mstype.MS_NOTHING, val: nil}
}

func (pf *PrintFunction) arity() int {
	return math.MaxUint8
}

func (pf *PrintFunction) String() string {
	return "* >> print"
}