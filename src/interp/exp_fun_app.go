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

	// Note: in this context, we need to be star-sensitive
	// to properly handle unpacking of tuples/arrays.
	args, err := evaluator.evaluateExpressionsStarSensitive(node.Args)

	if err != nil {
		return nil, err
	}

	//////////////////////////////////////////////////
	// bind function
	//////////////////////////////////////////////////
	return callable.Bind(args)
}

func (evaluator *MSEvaluator) evaluateIterableFunctionApplication(node *ast.IterableFuncAppNodeS) (MSVal, error) {
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

	// evaluate args
	args, err := evaluator.evaluateExpression(node.Args)

	// Cast args to iterable
	iter, ok := args.(MSIterable)

	if !ok {
		err = &EvalError{fmt.Sprintf("Function application arguments are not iterable, got type '%s'", args.Type())}
		return nil, err
	}

	elems, err := iter.Elems()

	if err != nil {
		return nil, err
	}

	// Bind each element
	vals := make([]MSVal, len(elems))
	for i, arg := range elems {

		// Bind the argument
		boundFn, err := callable.Bind([]MSVal{arg})

		if err != nil {
			return nil, err
		}

		vals[i] = boundFn
	}

	// Return new iterable from bound functions
	return iter.From(vals)
}


func (evaluator *MSEvaluator) evaluateIterableFunctionAppAndCall(node *ast.IterableFuncAppAndCallNodeS) (MSVal, error) {
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

	// Evaluate args
	args, err := evaluator.evaluateExpression(node.Args)

	if err != nil {
		return nil, err
	}

	// Cast args to iterable
	iter, ok := args.(MSIterable)

	if !ok {
		err = &EvalError{fmt.Sprintf("Function application arguments are not iterable, got type '%s'", args.Type())}
		return nil, err
	}

	elems, err := iter.Elems()

	if err != nil {
		return nil, err
	}

	// Bind and call each element
	vals := make([]MSVal, len(elems))
	for i, arg := range elems {

		// Bind the argument
		boundFn, err := callable.Bind([]MSVal{arg})

		if err != nil {
			return nil, err
		}

		// Cast to callable
		boundCallable, ok := boundFn.(MSCallable)
		
		if !ok {
			err = &EvalError{fmt.Sprintf("Function call is not implemented for type '%s'", boundFn)}
			return nil, err
		}
		// Call the function
		val, err := boundCallable.Call(evaluator)

		if err != nil {
			return nil, err
		}

		vals[i] = val
	}

	// Return new iterable from bound functions
	return iter.From(vals)
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

func (evaluator *MSEvaluator) evaluateIterableFunctionCall(node *ast.IterableFuncCallNodeS) (MSVal, error) {
	fns, err := evaluator.evaluateExpression(node.Fun)

	println("HELLO")

	if err != nil {
		return nil, err
	}

	// Check if the rhs is iterable, if its not an
	// iterable, we should just call the function as
	// a single element. If it is an array or tuple
	iter, ok := fns.(MSIterable)

	fmt.Printf("%v\n", fns)

	if ok {
		// is an iterable
		elems, err := iter.Elems()

		if err != nil {
			return nil, err
		}

		vals := make([]MSVal, len(elems))
		for i, fn := range elems {

			// Cast to callable, if possible
			callable, ok := fn.(MSCallable)

			if !ok {
				println("fn:", fn)
				err = &EvalError{fmt.Sprintf("Function call is not implemented for type '%s'", fn)}
				return nil, err
			}

			val, err := callable.Call(evaluator)

			if err != nil {
				return nil, err
			}

			vals[i] = val
		}

		// Cast to proper iterable
		return iter.From(vals)
		
	} else {
		// Call as a regular function if possible

		callable, ok := fns.(MSCallable)
		
		if !ok {
			err = &EvalError{fmt.Sprintf("Function call is not implemented for type '%s'", fns)}
			return nil, err
		}

		return callable.Call(evaluator)
	}

}