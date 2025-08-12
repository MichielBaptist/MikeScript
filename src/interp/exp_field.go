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