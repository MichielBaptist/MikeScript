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

	strs := make([]string, len(pf.args))
	for i, arg := range pf.args {
		strs[i] = arg.String()
	}

	return fmt.Sprintf("%s %s", strings.Join(strs, ", "), fs)
}

func (pf PrintFunction) Nullable() bool {
	return false
}

func (i PrintFunction) NullVal() MSVal {
	return nil
}

// --------------------------------------------------------
// Implements FunctionResult
// --------------------------------------------------------

func (pf PrintFunction) Call(_evaluator *MSEvaluator) (MSVal, error) {

	strs := make([]string, len(pf.args))
	for i, arg := range pf.args {
		strs[i] = arg.String()
	}

	fmt.Println(strings.Join(strs, ", "))

	return MSNothing{}, nil
}

func (pf PrintFunction) Bind(args []MSVal) (MSVal, error) {
	return NewPrintFunction(args), nil
}

func (pf PrintFunction) Arity() int {
	return math.MaxUint8
}