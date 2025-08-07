package interp

import "mikescript/src/ast"

func (e *MSEvaluator) evaluateTuple(n *ast.TupleNodeS) (MSVal, error) {

	values := make([]MSVal, len(n.Expressions))

	for i, expr := range n.Expressions {
		val, err := e.evaluateExpression(expr)

		if err != nil {
			return nil, err
		}

		values[i] = val
	}

	return MSTuple{Values: values}, nil
}