package interp

import (
	"mikescript/src/ast"
	"mikescript/src/mstype"
)


func (evaluator *MSEvaluator) executeFuncDeclStatement(node *ast.FuncDeclNodeS) (MSVal, error) {

	// resolve return type and function arg types
	resolvedParams, err := evaluator.resolveParamList(node.Params)

	if err != nil {
		return nil, err
	}

	resolvedReturn, err := evaluator.resolveType(node.Rt)

	if err != nil {
		return nil, err
	}

	// Make new ast node with resolved params
	resolvedNode := &ast.FuncDeclNodeS{Fname: node.Fname, Params: resolvedParams, Rt: resolvedReturn, Body: node.Body}
	
	// Wrap the decl with a callable
	callable := NewMSFunction(resolvedNode, evaluator.env)

	// Add EvalResult to current scope
	fname := node.Fname.Name.Lexeme
	err = evaluator.env.NewVar(fname, callable)

	// Throw env error
	if err != nil {
		return MSNothing{}, err
	}

	// The result of a function declartion is Nothing
	return MSNothing{}, nil

}

func (e *MSEvaluator) resolveParamList(ps []ast.FuncParamS) ([]ast.FuncParamS, error) {

	var resolvedParam mstype.MSType
	var err error

	resolvedParams := make([]ast.FuncParamS, len(ps))
	for i, p := range ps {

		// resolve 
		resolvedParam, err = e.resolveType(p.Type)

		if err != nil {
			return nil, err
		}

		resolvedParams[i] = ast.FuncParamS{Type: resolvedParam, Iden: p.Iden}
	}

	return resolvedParams, nil
}