package interp

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/mstype"
)

func (evaluator *MSEvaluator) executeWhileStatement(node *ast.WhileNodeS) EvalResult {

	var res EvalResult

	for {

		// Evaluate expression
		cond := evaluator.evaluateExpression(&node.Condition)
		if !cond.Valid() {
			return cond
		}
		if !cond.IsType(&mstype.MS_BOOL) {
			return evalErr(fmt.Sprintf("Condition must be of type bool, got %v", cond.rt))
		}

		// Get value of the bool
		condb, ok := cond.val.(bool)
		if !ok {
			return evalErr(fmt.Sprintf("Condition value is not a bool: %v", cond.val))
		}

		// Here the condition should be a boolean
		// If the condition is false, we break out of the loop
		if !condb {
			break
		}

		// Execute the body of the while loop
		res = evaluator.executeBlock(&node.Body, NewEnvironment(evaluator.env))

		// Check if result has an error
		if !res.Valid() {
			return res
		}

		// Check if result is break, on break we exit
		if res.IsType(&mstype.MS_BREAK) {
			return EvalResult{rt: mstype.MS_NOTHING}
		}
	}

	return res
}