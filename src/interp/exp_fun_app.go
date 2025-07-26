package interp

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/mstype"
)

func (evaluator *MSEvaluator) evaluateFunctionApplication(node *ast.FuncAppNodeS) EvalResult {

	//////////////////////////////////////////////////
	// evaluate right side or "x, y, z >> f";
	//////////////////////////////////////////////////
	fn := evaluator.evaluateExpression(&node.Fun)

	// Check for errors
	if (!fn.Valid()) {
		return fn // Found errors
	}

	// Check the type of the evaluation of Fun
	// For now, only func types will be callable.
	_, ok := fn.Rt.(*mstype.MSOperationTypeS)
	if !ok {
		return evalErr(fmt.Sprintf("Function application is not implemented for type '%s'", fn.Rt))
	}

	if fn.Val == nil {
		return evalErr("Trying to apply an undefined function")
	}

	// We can now be sure we can cast to FunctionResult
	// This will throw an error if not possible
	callable := fn.Val.(FunctionResult)
	
	// if (!ok) {
	// 	return evalErr(fmt.Sprintf("Could not cast %s to a FunctionResult", fn.val))
	// }

	//////////////////////////////////////////////////
	// evaluate left side or "x, y, z >> f";
	//////////////////////////////////////////////////

	// First evaluate all arguments, keep track of any errors.
	args := make([]EvalResult, len(node.Args))
	for i, arg := range node.Args {
		args[i] = evaluator.evaluateExpression(&arg)
	}

	// Accumulate all errors into one
	errs := []error{}
	for _, arg := range args {
		errs = append(errs, arg.Err...)
	}
	if len(errs) > 0 {
		return EvalResult{Err: errs}
	}

	//////////////////////////////////////////////////
	// bind function
	//////////////////////////////////////////////////
	return callable.Bind(args)
}

func (evaluator *MSEvaluator) evaluateFunctionCall(node *ast.FuncCallNodeS) EvalResult {
	fn := evaluator.evaluateExpression(&node.Fun)

	// Check for errors
	if (!fn.Valid()) {
		return fn // Found errors
	}

	// Check the type of the evaluation of Fun
	// For now, only func types will be callable.
	_, ok := fn.Rt.(*mstype.MSOperationTypeS)
	if !ok {
		return evalErr(fmt.Sprintf("Function call is not implemented for type '%s'", fn.Rt))
	}
	if fn.Val == nil {
		return evalErr("Trying to call an undefined function")
	}

	// We can now be sure we can cast to FunctionResult
	// This will throw an error if not possible
	callable := fn.Val.(FunctionResult)

	return callable.Call(evaluator)
}