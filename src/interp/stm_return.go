package interp

import (
	"mikescript/src/ast"
	"mikescript/src/mstype"
)

func (evaluator *MSEvaluator) executeReturnStatement(node *ast.ReturnNodeS) EvalResult {
	
	// Evaluate return values (if exists)
	var res EvalResult
	if node.HasReturnValue() {
		res = evaluator.evaluateExpression(&node.Node)
	} else {
		res = EvalResult{Rt: mstype.MS_NOTHING}
	}

	// Check for errors in res
	if !res.Valid() {
		return res
	}

	// wrap the result in a return EvalResult
	return EvalResult{
		Rt: mstype.MS_RETURN,
		Val: res,
		Err: nil,
	}
}
