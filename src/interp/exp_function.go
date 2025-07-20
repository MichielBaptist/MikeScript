package interp

import (
	"fmt"
	"mikescript/src/ast"
)

func (evaluator *MSEvaluator) evaluateFunctionApplication(node *ast.FuncAppNodeS) EvalResult {

	// Evaluate the function
	fn := evaluator.evaluateExpression(&node.Fun)

	// Check for errors
	if (!fn.Valid()) {
		return fn // Found errors
	}

	// First evaluate all arguments, keep track of any errors.
	args := make([]EvalResult, len(node.Args))
	for i, arg := range node.Args {
		args[i] = evaluator.evaluateExpression(&arg)
	}

	// Accumulate all errors into one
	errs := []error{}
	for _, arg := range args {
		errs = append(errs, arg.err...)
	}
	if len(errs) > 0 {
		return EvalResult{err: errs}
	}

	// Check the type of the evaluation of Fun
	// For now, only func types will be callable.
	if (fn.rt != RT_FUNCTION) {
		return evalErr(fmt.Sprintf("Function application is not implemented for type '%s'", fn.rt))
	}

	// We can now be sure we can cast to FunctionResult
	callable, ok := fn.val.(FunctionResult)

	if (!ok) {
		return evalErr(fmt.Sprintf("Could not cast %s to a FunctionResult", fn.val))
	}

	res := callable.call(evaluator, args)

	return res
}