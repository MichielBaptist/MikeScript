package interp

import (
	"mikescript/src/ast"
	"mikescript/src/mstype"
)


func (evaluator *MSEvaluator) executeFuncDeclStatement(node *ast.FuncDeclNodeS) EvalResult {
	
	// When declaring a function, all it does is:
	// 1. Create an EvalResult containing the function callable function object
	// 2. Store the EvalResult in the current environment

	// Wrap the decl with a callable
	//callable := MSFunction{decl: node}
	callable := NewMSFunction(node)

	// Create EvalResult, NOTE: we provide a pointer to MSFunction
	// not the function itself!!
	evalres := EvalResult{
		Rt: node.GetFuncType(),
		Val: callable,
	}

	// Add EvalResult to current scope
	fname := node.Fname.Name.Lexeme
	err := evaluator.env.NewVar(fname, evalres, node.GetFuncType())

	// Throw env error
	if err != nil {
		return EvalResult{Err: []error{err}}
	}

	// The result of a function declartion is Nothing
	return EvalResult{Rt: mstype.MS_NOTHING}

}
