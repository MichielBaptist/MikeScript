package interp

import (
	"mikescript/src/ast"
)


func (evaluator *MSEvaluator) executeFuncDeclStatement(node *ast.FuncDeclNodeS) (MSVal, error) {
	
	// When declaring a function, all it does is:
	// 1. Create an EvalResult containing the function callable function object
	// 2. Store the EvalResult in the current environment

	// Wrap the decl with a callable
	//callable := MSFunction{decl: node}
	callable := NewMSFunction(node)

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
