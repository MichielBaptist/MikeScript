package interp

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/mstype"
)

func (e *MSEvaluator) evalArrayIndexExpression(n *ast.ArrayIndexNodeS) (MSVal, error) {

	fmt.Printf("arr: %+s\n", n.Target)
	fmt.Printf("idx: %+s\n", n.Index)

	// Evaluate the target first
	val, err := e.evaluateExpression(n.Target)

	if err != nil {
		return MSNothing{}, err
	}

	// Check if the resulting value is indexable
	listVal := val.(*MSArray)

	// if !ok {
	// 	msg := fmt.Sprintf("Value '%s' of type '%s' is not indexable.", val, val.Type())
	// 	return MSNothing{}, &EvalError{message: msg}
	// }

	// Evaluate index
	idx, err := e.evaluateExpression(n.Index)

	if err != nil {
		return MSNothing{}, err
	}

	// Check if index is MSInt
	idxInt, ok := idx.(MSInt)

	if !ok {
		msg := fmt.Sprintf("Cannot use '%s' of type '%s' as an index, expected type '%s'.", idx, idx.Type(), mstype.MS_INT)
		return MSNothing{}, &EvalError{message: msg}
	}

	// Check if the index is in range or not
	if idxInt.Val < 0 || idxInt.Val >= len(listVal.Values) {
		msg := fmt.Sprintf("Array index out of bounds: '%d', expected value in '[%d, %d]'", idxInt.Val, 0, len(listVal.Values))
		return MSNothing{}, &EvalError{message: msg}
	}

	// Can now safely index the values
	return listVal.Values[idxInt.Val], err

}

func (e *MSEvaluator) evaluateArrayConstructor(n *ast.ArrayConstructorNodeS) (MSVal, error) {

	// Eval all expressions and check type compatibility
	var vals []MSVal

	for _, v := range n.Vals {
		val, err := e.evaluateExpression(v)

		// nil check
		if err != nil {
			return nil, err
		}

		// type check
		if !val.Type().Eq(n.Type) {
			msg := fmt.Sprintf("Array value '%s' has type '%s' but expected '%s'", val, val.Type(), n.Type)
			return nil, &EvalError{message: msg}
		}

		// Valid val, append to array
		vals = append(vals, val)
	}

	return &MSArray{Values: vals, VType: n.Type}, nil
}