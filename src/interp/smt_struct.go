package interp

import (
	"mikescript/src/ast"
	"mikescript/src/mstype"
)

func (e *MSEvaluator) executeStructDeclaration(n *ast.StructDeclarationNodeS) (MSVal, error) {
	sname := n.Name.VarName()
	fields := make(map[string]mstype.MSType)

	for fname, field := range n.Fields {
		fields[fname.VarName()] = field
	}

	// add new type
	err := e.types.NewType(sname, &mstype.MSStructTypeS{Name: sname, Fields: fields})

	return MSNothing{}, err
}

