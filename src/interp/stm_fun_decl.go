package interp

import (
	"mikescript/src/ast"
)


func (evaluator *MSEvaluator) executeFuncDeclStatement(node *ast.FuncDeclNodeS) (MSVal, error) {
	
	// Wrap the decl with a callable
	callable := NewMSFunction(node, evaluator.env)

	// Add EvalResult to current scope
	fname := node.Fname.Name.Lexeme
	err := evaluator.env.NewVar(fname, callable)

	// Throw env error
	if err != nil {
		return MSNothing{}, err
	}

	// The result of a function declartion is Nothing
	return MSNothing{}, nil

}
