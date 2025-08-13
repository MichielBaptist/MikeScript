package interp

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/mstype"
)

func (e *MSEvaluator) evalArrayIndexExpression(n *ast.ArrayIndexNodeS) (MSVal, error) {
	// target[index]

	// Evaluate the target first
	val, err := e.evaluateExpression(n.Target)

	if err != nil {
		return nil, err
	}

	// Check if the resulting value is indexable
	indexable, ok := val.(MSIndexable)

	if !ok {
		msg := fmt.Sprintf("Value '%s' of type '%s' is not indexable.", val, val.Type())
		return nil, &EvalError{message: msg}
	}

	// Evaluate index
	idx, err := e.evaluateExpression(n.Index)

	if err != nil {
		return nil, err
	}

	return indexable.Get(idx)

}

func (e *MSEvaluator) evaluateArrayConstructor(n *ast.ArrayConstructorNodeS) (MSVal, error) {
	// [n]type{vals...}


	// Cases:
	// 1. provided a size
	// 2. provided an initializer
	// 3. Error: provided a size and initializer

	nvals := len(n.Vals)

	// Case 3: impossible case, parser needs fixing
	if n.N != nil && nvals > 0 {
		msg := "Evaluator received an initializer and size, the parser is broken."
		return nil, &EvalError{message: msg}
	}

	// Case 1: size provided
	if n.N != nil {
		return e.evaluateArrayConstructorWithSize(n)
	}
	
	// Case 2: no size provided
	return e.evaluateArrayConstructorWithInitializer(n)
}

func (e *MSEvaluator) evaluateArrayConstructorWithInitializer(n *ast.ArrayConstructorNodeS) (MSVal, error) {

	var vals []MSVal	// Final list of values
	var val MSVal		// Current value
	var err error		// err

	// Resolve the base type
	// resolvedType, err := e.resolveType(n.Type)
	resolvedType := n.Type

	// if err != nil {
	// 	return nil, err
	// }

	for _, v := range n.Vals {
		val, err = e.evaluateExpression(v)

		// nil check
		if err != nil {
			return nil, err
		}

		// type check
		if !val.Type().Eq(resolvedType) {
			msg := fmt.Sprintf("Array value '%s' has type '%s' but expected '%s'", val, val.Type(), n.Type)
			return nil, &EvalError{message: msg}
		}

		// Valid val, append to array
		vals = append(vals, val)
	}

	return MSArray{Values: vals, VType: n.Type}, nil
}

func (e *MSEvaluator) evaluateArrayConstructorWithSize(n *ast.ArrayConstructorNodeS) (MSVal, error) {

	var vals []MSVal	// Final list of values
	var size MSVal		// Current value
	var sizeInt MSInt	// Size in Int val
	var err error		// err

	// evaluate N
	size, err = e.evaluateExpression(n.N)

	if err != nil {
		return nil, err
	}

	// Cast size to int
	sizeInt, ok := size.(MSInt)

	if !ok {
		msg := fmt.Sprintf("Value '%s' is of type '%s', expected type '%s'", size, size.Type(), mstype.MS_INT)
		return nil, &EvalError{message: msg}
	}

	// Check if size is valid (positive)
	if sizeInt.Val < 0 {
		msg := fmt.Sprintf("Cannot initialize arrays of negative size, received '%d'", sizeInt.Val)
		return nil, &EvalError{message: msg}
	}

	// resolve the type first
	resolvedType, err := e.resolveType(n.Type)
	//resolvedType := n.Type

	if err != nil {
		return nil, err
	}

	// create an array of proper size and init values
	vals = make([]MSVal, sizeInt.Val)
	for i := 0 ; i < sizeInt.Val ; i++ {
		vals[i] = e.typeToVal(resolvedType, false)
	}

	return MSArray{Values: vals, VType: resolvedType}, nil
}

func (e *MSEvaluator) evaluateArrayAssignment(n *ast.ArrayAssignmentNodeS) (MSVal, error) {

	// Evaluate the target first
	target, err := e.evaluateExpression(n.Target)

	if err != nil {
		return nil, err
	}

	// Check if the resulting value is indexable
	indexable, ok := target.(MSIndexable)

	if !ok {
		msg := fmt.Sprintf("Value '%s' of type '%s' is not indexable.", target, target.Type())
		return nil, &EvalError{message: msg}
	}

	// Evaluate index
	idx, err := e.evaluateExpression(n.Index)

	if err != nil {
		return nil, err
	}

	// Evaluate val
	val, err := e.evaluateExpression(n.Value)

	if err != nil {
		return nil, err
	}

	return indexable.Set(idx, val)

}