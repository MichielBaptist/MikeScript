package interp

import (
	"mikescript/src/ast"
)

func (evaluator *MSEvaluator) executeReturnStatement(node *ast.ReturnNodeS) (MSVal, error) {
	
	var res MSVal
	var err error

	if node.HasReturnValue() {
		res, err = evaluator.evaluateExpression(node.Node)
	} else {
		res = MSNothing{}
	}

	if err != nil {
		return nil, err
	}

	// wrap the result in a return EvalResult
	return MSReturn{Val: res}, nil
}
