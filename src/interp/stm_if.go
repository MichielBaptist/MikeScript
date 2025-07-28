package interp

import (
	"fmt"
	"mikescript/src/ast"
)


func (evaluator *MSEvaluator) executeIfstatement(node *ast.IfNodeS) (MSVal, error) {

	// Evaluate the condition
	cond, err := evaluator.evaluateExpression(&node.Condition)

	// Sanity checks
	if err != nil {
		return MSNothing{}, err
	}

	bcond, ok := cond.(MSBool)

	if !ok {
		return MSNothing{}, &EvalError{fmt.Sprintf("Condition must be of type bool, got '%v'", cond.Type())}
	}

	// Execute the then or else statement based on the condition
	if bcond.Val {
		return evaluator.executeStatement(&node.ThenStmt)
	} else if node.ElseStmt != nil {
		return evaluator.executeStatement(&node.ElseStmt)
	}

	
	return MSNothing{}, err

}
