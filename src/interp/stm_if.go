package interp

import (
	"fmt"
	"mikescript/src/ast"
)


func (evaluator *MSEvaluator) executeIfstatement(node *ast.IfNodeS) (MSVal, error) {

	cond, err := evaluator.evaluateExpression(node.Condition)

	if err != nil {
		return nil, err
	}

	bcond, ok := cond.(MSBool)

	if !ok {
		return nil, &EvalError{fmt.Sprintf("Condition must be of type bool, got '%v'", cond.Type())}
	}

	if bcond.Val {
		return evaluator.executeStatement(node.ThenStmt)
	} else if node.ElseStmt != nil {
		return evaluator.executeStatement(node.ElseStmt)
	}

	
	return MSNothing{}, err

}
