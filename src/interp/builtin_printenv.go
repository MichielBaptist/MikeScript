package interp

import (
	"mikescript/src/mstype"
)

///////////////////////////////////////////////////////////////
// mikescript builtins
///////////////////////////////////////////////////////////////
func MSBuiltinPrintEnv() MSVal {
	return NewPrintEnvFunction([]MSVal{})
}

func NewPrintEnvFunction(args []MSVal) MSVal {
	return PrintEnvFunction{args: args}
}

///////////////////////////////////////////////////////////////
// Print function
///////////////////////////////////////////////////////////////

// Is a FunctionResult placeholder for print
type PrintEnvFunction struct {
	args []MSVal
}

// --------------------------------------------------------
// Implements MSValue
// --------------------------------------------------------

func (pf PrintEnvFunction) Type() mstype.MSType {
	return &mstype.MSOperationTypeS{Left: []mstype.MSType{}, Right: mstype.MS_NOTHING}
}

func (pf PrintEnvFunction) String() string {
	return ">> print_env -> nothing"
}

func (pf PrintEnvFunction) Nullable() bool {
	return false
}

// --------------------------------------------------------
// Implements FunctionResult
// --------------------------------------------------------

func (pf PrintEnvFunction) Call(_evaluator *MSEvaluator) (MSVal, error) {
	_evaluator.env.printEnv()
	return MSNothing{}, nil
}

func (pf PrintEnvFunction) Bind(args []MSVal) (MSVal, error) {
	return MSNothing{}, BindingError{msg: "Cannot bind print_env function."}
}

func (pf PrintEnvFunction) Arity() int {
	return 0
}