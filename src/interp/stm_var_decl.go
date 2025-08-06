package interp

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/mstype"
)


func (evaluator *MSEvaluator) executeDeclarationStatement(node *ast.VarDeclNodeS) (MSVal, error) {

	// First resolve the type, in case the type is
	// a named type.
	typ, err := evaluator.resolveType(node.Vartype)

	if err != nil {
		return nil, err
	}

	// Get the default value for the type
	val := evaluator.typeToVal(typ, node.Identifier)

	if val == nil {
		_ = []int{}[0]
	}

	// Declare variable in env
	err = evaluator.env.NewVar(node.VarName(), val)

	return MSNothing{}, err
}

////////////////////////////////////////////////////////////////////////
// type to default value
////////////////////////////////////////////////////////////////////////

func (e *MSEvaluator) typeToVal(tk mstype.MSType, name *ast.VariableExpNodeS) MSVal {
	// cases:
	// 1. simple type
	// 2. composite
	// 3. function
	// 4. array

	switch t := tk.(type){
	case *mstype.MSSimpleTypeS: 	return e.simpleTypeToVal(t)
	case *mstype.MSCompositeTypeS:	return e.compositeTypeToVal(t, name)
	case *mstype.MSOperationTypeS:	return MSFunctionFromType(t, name, e.env)
	case *mstype.MSArrayType:		return e.arrayTypeToVal(t)
	default:						fmt.Printf("Found unknown type: '%s'\n", t)
	}
	return nil
}

func (e *MSEvaluator) simpleTypeToVal(rt *mstype.MSSimpleTypeS) MSVal {
	switch rt.Rt {
	case mstype.RT_INT:		return MSInt{0}
	case mstype.RT_FLOAT:	return MSFloat{0.0}
	case mstype.RT_STRING:	return MSString{""}
	case mstype.RT_BOOL:	return MSBool{false}
	default:				return nil
	}
}

func (e *MSEvaluator) compositeTypeToVal(ct *mstype.MSCompositeTypeS, name *ast.VariableExpNodeS) MSVal {
	vals := make([]MSVal, len(ct.Types))
	for i, t := range ct.Types {
		vals[i] = e.typeToVal(t, name)
	}
	return MSTuple{Values: vals}
}

func (e *MSEvaluator) arrayTypeToVal(t *mstype.MSArrayType) MSVal {
	return MSArray{Values: make([]MSVal, 0), VType: t.Type}
}


////////////////////////////////////////////////////////////////////////
// named type resolution
////////////////////////////////////////////////////////////////////////

func (e *MSEvaluator) resolveType(typ mstype.MSType) (mstype.MSType, error) {
	switch t := typ.(type) {
	case *mstype.MSNamedTypeS:		return e.resolveNamedType(t)
	case *mstype.MSArrayType:		return e.resolveArrayType(t)
	case *mstype.MSCompositeTypeS:	return e.resolveCompositeType(t)
	case *mstype.MSOperationTypeS:	return e.resolveOperationType(t)
	case *mstype.MSSimpleTypeS:		return t, nil
	}

	return nil, &TypeResolutionError{msg: fmt.Sprintf("Unknown type '%v'", typ)}
}

func (e *MSEvaluator) resolveNamedType(t *mstype.MSNamedTypeS) (mstype.MSType, error) {
	rt, err := e.types.GetType(t.Name)
	if err != nil {
		return nil, err
	}
	return e.resolveType(rt)
}

func (e *MSEvaluator) resolveArrayType(t *mstype.MSArrayType) (mstype.MSType, error) {
	resolved, err := e.resolveType(t.Type)
	return &mstype.MSArrayType{Type: resolved}, err
}

func (e *MSEvaluator) resolveCompositeType(t *mstype.MSCompositeTypeS) (mstype.MSType, error) {
	resolved, err := e.resolveTypeList(t.Types)
	return &mstype.MSCompositeTypeS{Types: resolved}, err
}

func (e *MSEvaluator) resolveOperationType(t *mstype.MSOperationTypeS) (mstype.MSType, error) {
	resolvedLeft, err := e.resolveTypeList(t.Left)
	if err != nil {
		return nil, err
	}
	resolvedRight, err := e.resolveType(t.Right)
	if err != nil {
		return nil, err
	}

	return &mstype.MSOperationTypeS{Left: resolvedLeft, Right: resolvedRight}, nil
}

func (e *MSEvaluator) resolveTypeList(t []mstype.MSType) ([]mstype.MSType, error) {
	resolved := make([]mstype.MSType, len(t))
	for i, ct := range t {
		r, err := e.resolveType(ct)
		if err != nil {
			return nil, err
		}
		resolved[i] = r
	}
	return resolved, nil
}

type TypeResolutionError struct {
	msg string
}

func (e *TypeResolutionError) Error() string {
	return "Type resolution error:" + e.msg
}