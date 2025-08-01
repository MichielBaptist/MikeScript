package interp

import (
	"mikescript/src/ast"
)

func (evaluator *MSEvaluator) executeReturnStatement(node *ast.ReturnNodeS) (MSVal, error) {
	
	// Evaluate return values (if exists)
	var res MSVal
	var err error

	if node.HasReturnValue() {
		res, err = evaluator.evaluateExpression(node.Node)
	} else {
		res = MSNothing{}
	}

	// Check for errors in res
	if err != nil {
		return MSNothing{}, err
	}

	// wrap the result in a return EvalResult
	return MSReturn{Val: res}, nil
}
