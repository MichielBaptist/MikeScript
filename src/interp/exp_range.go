package interp

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/mstype"
)

func (e *MSEvaluator) evaluateRangeConstructor(node *ast.RangeConstructorNodeS) (MSVal, error) {

	fromVal, err := e.evaluateExpression(node.From)

	if err != nil {
		return nil, err
	}

	toVal, err := e.evaluateExpression(node.To)

	if err != nil {
		return nil, err
	}
	
	// Expect both to be integers
	fromInt, ok := fromVal.(MSInt)
	if !ok {
		return nil, &EvalError{fmt.Sprintf("Range constructor 'from' value must be of type 'int', got '%s'", fromVal.Type())}
	}
	toInt, ok := toVal.(MSInt)
	if !ok {
		return nil, &EvalError{fmt.Sprintf("Range constructor 'to' value must be of type 'int', got '%s'", toVal.Type())}
	}

	if toInt.Val < fromInt.Val {
		return nil, &EvalError{fmt.Sprintf("Range constructor 'to' value must be greater than or equal to 'from' value, got from='%d' to='%d'", fromInt.Val, toInt.Val)}
	}

	// construct range
	nval := toInt.Val - fromInt.Val
	nval = max(0, nval)
	vals := make([]MSVal, nval)
	for i := 0; i < nval; i++ {
		vals[i] = MSInt{Val: fromInt.Val + i}
	}

	return MSArray{Values: vals, VType: mstype.MS_INT}, nil

}