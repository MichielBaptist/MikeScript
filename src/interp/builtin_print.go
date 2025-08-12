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
func MSBuiltinPrint() MSVal {
	return NewPrintFunction([]MSVal{})
}

func NewPrintFunction(args []MSVal) MSVal {
	return PrintFunction{args: args}
}

///////////////////////////////////////////////////////////////
// Print function
///////////////////////////////////////////////////////////////

// Is a FunctionResult placeholder for print
type PrintFunction struct {
	args []MSVal
}

// --------------------------------------------------------
// Implements MSValue
// --------------------------------------------------------

func (pf PrintFunction) Type() mstype.MSType {
	return &mstype.MSOperationTypeS{Left: []mstype.MSType{}, Right: mstype.MS_NOTHING}
}

func (pf PrintFunction) String() string {
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

func (pf PrintFunction) Nullable() bool {
	return false
}

// --------------------------------------------------------
// Implements FunctionResult
// --------------------------------------------------------

func (pf PrintFunction) Call(_evaluator *MSEvaluator) (MSVal, error) {

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
	return MSNothing{}, nil
}

func (pf PrintFunction) Bind(args []MSVal) (MSVal, error) {
	return NewPrintFunction(args), nil
}

func (pf PrintFunction) Arity() int {
	return math.MaxUint8
}