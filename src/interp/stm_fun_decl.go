package interp

import (
	"mikescript/src/ast"
)


func (evaluator *MSEvaluator) executeFuncDeclStatement(node *ast.FuncDeclNodeS) (MSVal, error) {
	var err error

	// Resolve funcdecl node
	resolvedNode, err := evaluator.resolveFunctionDeclaration(node)

	if err != nil {
		return nil, err
	}

	// Wrap the decl with a callable
	callable := NewMSFunction(resolvedNode, evaluator.env)

	// Add EvalResult to current scope
	fname := node.Fname.Name.Lexeme
	err = evaluator.env.NewVar(fname, callable)

	// Throw env error
	if err != nil {
		return nil, err
	}

	// The result of a function declartion is Nothing
	return MSNothing{}, nil
}
