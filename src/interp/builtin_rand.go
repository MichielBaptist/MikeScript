package interp

import (
	"math/rand"
	"mikescript/src/mstype"
)

///////////////////////////////////////////////////////////////
// mikescript builtins
///////////////////////////////////////////////////////////////
func MSBuiltinRand() MSVal {
	return RandFunction{}
}

///////////////////////////////////////////////////////////////
// Print function
///////////////////////////////////////////////////////////////

// Is a FunctionResult placeholder for print
type RandFunction struct {}

// --------------------------------------------------------
// Implements MSValue
// --------------------------------------------------------

func (pf RandFunction) Type() mstype.MSType {
	return &mstype.MSOperationTypeS{Left: []mstype.MSType{}, Right: mstype.MS_FLOAT}
}

func (pf RandFunction) String() string {
	return ">> rand -> float"
}

func (pf RandFunction) Nullable() bool {
	return false
}

func (i RandFunction) NullVal() MSVal {
	return nil
}

// --------------------------------------------------------
// Implements FunctionResult
// --------------------------------------------------------

func (pf RandFunction) Call(_evaluator *MSEvaluator) (MSVal, error) {
	return MSFloat{Val: rand.Float64()}, nil
}

func (pf RandFunction) Bind(args []MSVal) (MSVal, error) {
	return nil, nil
}

func (pf RandFunction) Arity() int {
	return 0
}