package interp

import (
	"fmt"
	"mikescript/src/ast"
)

func (evaluator *MSEvaluator) executeWhileStatement(node *ast.WhileNodeS) (MSVal, error) {

	for {

		cond, err := evaluator.evaluateExpression(node.Condition)
		if err != nil {
			return nil, err
		}

		bcond, ok := cond.(MSBool)

		if !ok {
			return nil, &EvalError{fmt.Sprintf("Condition must be of type bool, got %v", cond.Type())}
		}

		if !bcond.Val {
			break
		}

		res, err := evaluator.executeBlock(node.Body, NewEnvironment(evaluator.env))

		if err != nil {
			return nil, err
		}

		// Check for break or return
		switch res.(type) {
		case MSBreak:	return MSNothing{}, nil
		case MSReturn:	return res, nil
		}
	}

	return MSNothing{}, nil
}