package interp

import (
	"fmt"
	"mikescript/src/ast"
)

func (e *MSEvaluator) evaluateTuple(n *ast.TupleNodeS) (MSVal, error) {

	values, err := e.evaluateExpressionsStarSensitive(n.Expressions)

	if err != nil {
		return nil, err
	}

	return MSTuple{Values: values}, nil
}

func valueUnpack(val MSVal) ([]MSVal, error) {
	unpackable, ok := val.(MSIterable)
	
	if !ok {
		return nil, &EvalError{fmt.Sprintf("Tried unpacking %s, which is not iterable", val.Type())}
	}

	return unpackable.Elems()
}

func (e *MSEvaluator) evaluateExpressionsStarSensitive(exprs []ast.ExpNodeI) ([]MSVal, error) {

	// Init slice
	var values []MSVal

	// Loop all expressions and check if starred
	for _, expr := range exprs {

		// unpack starred expressions
		switch t := expr.(type) {
		case *ast.StarredExpNodeS:

			// Need to evaluate the inner expression
			starredVal, err := e.evaluateExpression(t.Node)
			
			if err != nil {
				return nil, err
			}

			// Check if it's iterable (array, tuple, ...)
			elems, err := valueUnpack(starredVal)

			if err != nil {
				return nil, err
			}

			// Append all elements to values
			values = append(values, elems...)
		
		default:

			// normal expression, evaluate as usual

			val, err := e.evaluateExpression(expr)

			if err != nil {
				return nil, err
			}

			values = append(values, val)
		}
	}

	return values, nil
}
