package interp

import (
	"fmt"
	"mikescript/src/ast"
)


func (evaluator *MSEvaluator) executeIfstatement(node *ast.IfNodeS) EvalResult {

	// Evaluate the condition
	cond := evaluator.evaluateExpression(&node.Condition)

	// Sanity checks
	if !cond.Valid() {
		return cond
	}
	if !cond.Expect(RT_BOOL) {
		return evalErr(fmt.Sprintf("Condition must be of type bool, got %v", cond.rt))
	}

	// Execute the then or else statement based on the condition
	switch cond.val.(type) {
	case bool:
		if cond.val.(bool) {
			return evaluator.executeStatement(&node.ThenStmt)
		} else if node.ElseStmt != nil {
			return evaluator.executeStatement(&node.ElseStmt)
		}
	default:
		return evalErr(fmt.Sprintf("Incompatible result type and value: %v: %v", cond.rt, cond.val))
	}

	return EvalResult{rt: RT_NONE}

}
