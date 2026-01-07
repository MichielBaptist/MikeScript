package interp

import (
	"fmt"
	"mikescript/src/mstype"
)

///////////////////////////////////////////////////////////////
// mikescript builtins
///////////////////////////////////////////////////////////////
func MSBuiltinLen() MSVal {
	return LenFunction{}
}

func NewLenFunction(lst MSIterable) MSVal {
	return LenFunction{lst: lst}
}


///////////////////////////////////////////////////////////////
// Len function
///////////////////////////////////////////////////////////////

// Is a FunctionResult placeholder for print
type LenFunction struct {
	lst MSIterable
}

// --------------------------------------------------------
// Implements MSValue
// --------------------------------------------------------

func (pf LenFunction) Type() mstype.MSType {
	return &mstype.MSOperationTypeS{Left: []mstype.MSType{}, Right: mstype.MS_INT}
}

func (pf LenFunction) String() string {
	return ">> len -> int"
}

func (pf LenFunction) Nullable() bool {
	return false
}

func (i LenFunction) NullVal() MSVal {
	return nil
}

// --------------------------------------------------------
// Implements FunctionResult
// --------------------------------------------------------

func (pf LenFunction) Call(_evaluator *MSEvaluator) (MSVal, error) {
	return pf.lst.Len()
}

func (pf LenFunction) Bind(args []MSVal) (MSVal, error) {

	// Check arity
	if len(args) != 1 {
		return nil, &BindingError{msg: "LenFunction expects exactly 1 argument"}
	}

	// Get the argument
	arg := args[0]

	// Check if arg is iterable
	iterable, ok := arg.(MSIterable)
	
	if !ok {
		msg := fmt.Sprintf("LenFunction expected argument of iterable type, got '%s'", arg.Type())
		return nil, &EvalError{message: msg}
	}
	
	// Bind the argument
	return NewLenFunction(iterable), nil
}

func (pf LenFunction) Arity() int {
	return 1
}