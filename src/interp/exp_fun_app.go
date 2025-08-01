package interp

import (
	"fmt"
	"mikescript/src/ast"
)

func (evaluator *MSEvaluator) evaluateFunctionApplication(node *ast.FuncAppNodeS) (MSVal, error) {

	//////////////////////////////////////////////////
	// evaluate right side or "x, y, z >> f";
	//////////////////////////////////////////////////
	fn, err := evaluator.evaluateExpression(node.Fun)

	// Check for errors
	if err != nil {
		return MSNothing{}, err
	}

	// Cast the value to FunctionResult (interface)
	callable, ok := fn.(MSCallable)

	// Cast to function, if not ok error
	if !ok {
		err = &EvalError{fmt.Sprintf("Function application is not implemented for type '%s'", fn)}
		return MSNothing{}, err
	}

	//////////////////////////////////////////////////
	// evaluate left side or "x, y, z >> f";
	//////////////////////////////////////////////////

	// Check if the arity of the function supports binding
	if callable.Arity() < len(node.Args) {
		err := fmt.Sprintf("Exceeded arity of '%s' expected maximum %v arguments but received %v", callable, callable.Arity(), len(node.Args))
		return MSNothing{}, &BindingError{msg: err}
	}

	// First evaluate all arguments, keep track of any errors.
	args := make([]MSVal, len(node.Args))
	for i, arg := range node.Args {
		arg, err := evaluator.evaluateExpression(arg)

		if err != nil {
			return MSNothing{}, err
		}

		args[i] = arg
	}

	//////////////////////////////////////////////////
	// bind function
	//////////////////////////////////////////////////
	return callable.Bind(args)
}

func (evaluator *MSEvaluator) evaluateFunctionCall(node *ast.FuncCallNodeS) (MSVal, error) {
	fn, err := evaluator.evaluateExpression(node.Fun)

	// Check for errors
	if err != nil {
		return MSNothing{}, err
	}

	// Cast the value to FunctionResult (interface)
	callable, ok := fn.(MSCallable)

	if !ok {
		err = &EvalError{fmt.Sprintf("Function call is not implemented for type '%s'", fn)}
		return MSNothing{}, err
	}

	
	return callable.Call(evaluator)
}
