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

	value, err := e.evaluateExpression(n.Value)

	if err != nil {
		return nil, err
	}

	currentVal, err := fieldable.Get(n.Field.VarName())

	if err != nil {
		return nil, err
	}

	if _, ok := value.(MSNothing) ; ok && currentVal.Nullable() {
		value = currentVal.NullVal()
	}

	return fieldable.Set(n.Field.VarName(), value)
}