package interp

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/mstype"
)


func (evaluator *MSEvaluator) executeIfstatement(node *ast.IfNodeS) EvalResult {

	// Evaluate the condition
	cond := evaluator.evaluateExpression(&node.Condition)

	// Sanity checks
	if !cond.Valid() {
		return cond
	}
	if !cond.IsType(&mstype.MS_BOOL) {
		return evalErr(fmt.Sprintf("Condition must be of type bool, got %v", cond.Rt))
	}

	// Execute the then or else statement based on the condition
	switch cond.Val.(type) {
	case bool:
		if cond.Val.(bool) {
			return evaluator.executeStatement(&node.ThenStmt)
		} else if node.ElseStmt != nil {
			return evaluator.executeStatement(&node.ElseStmt)
		}
	default:
		return evalErr(fmt.Sprintf("Incompatible result type and value: %v: %v", cond.Rt, cond.Val))
	}

	return EvalResult{Rt: mstype.MS_NOTHING}

}
