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

	if err != nil {
		return nil, err
	}

	// Cast the value to FunctionResult (interface)
	callable, ok := fn.(MSCallable)

	if !ok {
		err = &EvalError{fmt.Sprintf("Function application is not implemented for type '%s'", fn)}
		return nil, err
	}

	//////////////////////////////////////////////////
	// evaluate left side or "x, y, z >> f";
	//////////////////////////////////////////////////

	// Check if the arity of the function supports binding
	if callable.Arity() < len(node.Args) {
		err := fmt.Sprintf("Exceeded arity of '%s' expected maximum %v arguments but received %v", callable, callable.Arity(), len(node.Args))
		return nil, &BindingError{msg: err}
	}

	args := make([]MSVal, len(node.Args))
	for i, arg := range node.Args {
		arg, err := evaluator.evaluateExpression(arg)

		if err != nil {
			return nil, err
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

	if err != nil {
		return nil, err
	}

	// Cast the value to FunctionResult (interface)
	callable, ok := fn.(MSCallable)

	if !ok {
		err = &EvalError{fmt.Sprintf("Function call is not implemented for type '%s'", fn)}
		return nil, err
	}

	
	return callable.Call(evaluator)
}
