package interp

import "mikescript/src/ast"

func (e *MSEvaluator) executeTypeDeclaration(n *ast.TypeDeclarationNode) (MSVal, error) {

	// add the type to env
	err := e.types.NewType(n.Tname.VarName(), n.Type)

	return &MSNothing{}, err
}