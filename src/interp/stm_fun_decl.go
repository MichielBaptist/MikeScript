package interp

import (
	"mikescript/src/ast"
)


func (evaluator *MSEvaluator) executeFuncDeclStatement(node *ast.FuncDeclNodeS) (MSVal, error) {
	var err error

	resolvedNode, err := evaluator.resolveFunctionDeclaration(node)

	if err != nil {
		return nil, err
	}

	// Wrap the decl with a callable
	callable := NewMSFunction(resolvedNode, evaluator.env)

	// Add to current scope
	fname := node.Fname.Name.Lexeme
	err = evaluator.env.NewVar(fname, callable)

	if err != nil {
		return nil, err
	}

	// The result of a function declartion is Nothing
	return MSNothing{}, nil
}
