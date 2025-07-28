package interp

import (
	"fmt"
	"mikescript/src/ast"
)

func (evaluator *MSEvaluator) executeWhileStatement(node *ast.WhileNodeS) (MSVal, error) {

	for {

		// Evaluate expression
		cond, err := evaluator.evaluateExpression(&node.Condition)
		if err != nil {
			return MSNothing{}, err
		}

		// Cast to bool
		bcond, ok := cond.(MSBool)

		if !ok {
			return MSNothing{}, &EvalError{fmt.Sprintf("Condition must be of type bool, got %v", cond.Type())}
		}

		// Here the condition should be a boolean
		// If the condition is false, we break out of the loop
		if !bcond.Val {
			break
		}

		// Execute the body of the while loop
		res, err := evaluator.executeBlock(&node.Body, NewEnvironment(evaluator.env))

		// Check if result has an error
		if err != nil {
			return MSNothing{}, err
		}

		// Check for break or return
		switch res.(type) {
		case MSBreak:	return MSNothing{}, nil
		case MSReturn:	return res, nil
		}
	}

	return MSNothing{}, nil
}