package interp

import (
	"fmt"
	"mikescript/src/ast"
)

func (evaluator *MSEvaluator) executeForStatement(node *ast.ForNodeS) (MSVal, error) {

	// Evaluate the iterable expression
	iterableVal, err := evaluator.evaluateExpression(node.Iterable)
	if err != nil {
		return MSNothing{}, err
	}

	iterable, ok := iterableVal.(MSIterable)
	if !ok {
		msg := fmt.Sprintf("Value of type '%v' is not iterable", iterableVal.Type())
		return MSNothing{}, &EvalError{message: msg}
	}

	// get elems
	elems, err := iterable.Elems()
	if err != nil {
		return MSNothing{}, err
	}

	for _, val := range elems {

		// Create a new scope for the loop variable
		env := NewEnvironment(evaluator.env)
		env.NewVar(node.LoopVar.VarName(), val)

		// execute the loop body in environment
		res, err := evaluator.executeBlock(node.Body, env)

		if err != nil {
			return MSNothing{}, err
		}

		// Check for break or continue
		switch res.(type) {
		case MSBreak:	return MSNothing{}, nil
		case MSReturn:	return res, nil
		}
	}

	return MSNothing{}, nil
}
