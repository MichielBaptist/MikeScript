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
	return NewPrintFunction([]EvalResult{})
}

func NewPrintFunction(args []EvalResult) EvalResult {
	return EvalResult{
		Rt: &mstype.MSOperationTypeS{Right: mstype.MS_NOTHING},
		Val: &PrintFunction{args: args},
		Err: []error{},
	}
}

///////////////////////////////////////////////////////////////
// Print function
///////////////////////////////////////////////////////////////

// Is a FunctionResult placeholder for print
type PrintFunction struct {
	args []EvalResult
}

// --------------------------------------------------------
// Implements FunctionResult
// --------------------------------------------------------

func (pf *PrintFunction) Call(_evaluator *MSEvaluator) EvalResult {

	// Get all result values in a slice
	// and convert them to strings
	strs := make([]string, len(pf.args))
	for i, arg := range pf.args {
		strs[i] = arg.String()
	}

	// Print the joined strings, this is where
	// the connection between mikesript and go is made.
	fmt.Println(strings.Join(strs, ", "))

	// Done without return value
	return EvalResult{Rt: mstype.MS_NOTHING, Val: nil}
}

func (pf *PrintFunction) Bind(args []EvalResult) EvalResult {
	return NewPrintFunction(args)
}

func (pf *PrintFunction) Arity() int {
	return math.MaxUint8
}

// --------------------------------------------------------
// Implements stringer
// --------------------------------------------------------
func (pf *PrintFunction) String() string {

	fs := ">> print -> nothing"

	if len(pf.args) == 0 {
		return fs
	}

	// Convert all print args to string
	strs := make([]string, len(pf.args))
	for i, arg := range pf.args {
		strs[i] = arg.String()
	}

	return fmt.Sprintf("%s %s", strings.Join(strs, ", "), fs)
}