package interp

import "mikescript/src/ast"

func (e *MSEvaluator) executeTypeDeclaration(n *ast.TypeDefStatementS) (MSVal, error) {

	// add the type to env
	err := e.env.NewType(n.Tname.VarName(), n.Type)

	return &MSNothing{}, err
}