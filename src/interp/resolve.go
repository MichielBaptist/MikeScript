package interp

import (
	"mikescript/src/ast"
	"mikescript/src/mstype"
)

func (e *MSEvaluator) resolveType(t mstype.MSType) (mstype.MSType, error) {
	switch tt := t.(type) {
	case *mstype.MSSimpleTypeS:		return tt, nil
	case *mstype.MSStructTypeS:		return tt, nil
	case *mstype.MSCompositeTypeS:	return e.resolveCompositeType(tt)
	case *mstype.MSArrayType:		return e.resolveArrayType(tt)
	case *mstype.MSNamedTypeS:		return e.resolveNamedType(tt)
	case *mstype.MSOperationTypeS:	return e.resolveOperationType(tt)
	default:						_ = []int{}[0] ; return nil, nil
	}
}


func (e *MSEvaluator) resolveTypes(ts []mstype.MSType) ([]mstype.MSType, error) {
	resolved := make([]mstype.MSType, len(ts))
	for i, t := range ts {
		r, err := e.resolveType(t)
		if err != nil {
			return nil, err
		}
		resolved[i] = r
	}

	return resolved, nil
}

func (e *MSEvaluator) resolveCompositeType(ct *mstype.MSCompositeTypeS) (*mstype.MSCompositeTypeS, error) {
	resolved, err := e.resolveTypes(ct.Types)

	if err != nil {
		return nil, err
	}

	return &mstype.MSCompositeTypeS{Types: resolved}, nil
}

func (e *MSEvaluator) resolveNamedType(nt *mstype.MSNamedTypeS) (mstype.MSType, error) {

	var resolved mstype.MSType
	var err error

	// Look up scope depth
	depth, ok := e.tlocals[nt]

	if ok {
		resolved, err = e.env.GetType(nt.Name, depth)
	} else {
		resolved, err = e.glb.GetType(nt.Name, depth)
	}

	// Yikes, todo better error handling
	if err != nil {
		return nil, err
	}

	return e.resolveType(resolved)
}

func (e *MSEvaluator) resolveArrayType(at *mstype.MSArrayType) (*mstype.MSArrayType, error) {
	resolvedBase, err := e.resolveType(at.Type)
	return &mstype.MSArrayType{Type: resolvedBase}, err
}

func (e *MSEvaluator) resolveOperationType(ot *mstype.MSOperationTypeS) (*mstype.MSOperationTypeS, error) {
	resolvedLeft, err := e.resolveTypes(ot.Left)
	if err != nil {
		return nil, err
	}
	resolvedRight, err := e.resolveType(ot.Right)
	return &mstype.MSOperationTypeS{Left: resolvedLeft, Right: resolvedRight}, err
}

func (e *MSEvaluator) resolveFunctionDeclaration(f *ast.FuncDeclNodeS) (*ast.FuncDeclNodeS, error) {

	// resolve return type and function arg types
	resolvedParams, err := e.resolveParamList(f.Params)

	if err != nil {
		return nil, err
	}

	resolvedReturn, err := e.resolveType(f.Rt)

	if err != nil {
		return nil, err
	}

	resolvedFuncDecl := ast.FuncDeclNodeS{
		Fname: f.Fname,
		Params: resolvedParams,
		Rt: resolvedReturn,
		Body: f.Body,
	}

	// Make new ast node with resolved params
	return &resolvedFuncDecl, nil
	
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