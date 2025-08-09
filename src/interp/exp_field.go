package interp

import (
	"mikescript/src/ast"
)

func (e *MSEvaluator) evaluateFieldAccess(n *ast.FieldAccessNodeS) (MSVal, error) {
	// target.field

	var fieldName string = n.Field.VarName()

	target, err := e.evaluateExpression(n.Target)

	if err != nil {
		return nil, err
	}

	fieldable := target.(MSFieldable)

	// if !ok {
	// 	msg := fmt.Sprintf("Object '%v' of type '%v' does not have field property '%v'", target, target.Type(), n.Field)
	// 	return nil, &EvalError{message: msg}
	// }

	// Check if field is valid
	if err := fieldable.ValidField(fieldName) ; err != nil {
		return nil, err
	}

	return fieldable.Get(fieldName)
}

func (e *MSEvaluator) evaluateFieldAssign(n *ast.FieldAssignmentNode) (MSVal, error) {

	target, err := e.evaluateExpression(n.Target)

	if err != nil {
		return nil, err
	}

	// cast to Fieldable
	fieldable := target.(MSFieldable)

	// eval val
	value, err := e.evaluateExpression(n.Value)

	if err != nil {
		return nil, err
	}

	return fieldable.Set(n.Field.VarName(), value)
}