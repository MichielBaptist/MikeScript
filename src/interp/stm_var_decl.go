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
	val := evaluator.typeToVal(typ, false)
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

func (e *MSEvaluator) typeToVal(tk mstype.MSType, context bool) MSVal {
	switch t := tk.(type){
	case *mstype.MSSimpleTypeS: 	return e.simpleTypeToVal(t)
	case *mstype.MSCompositeTypeS:	return e.compositeTypeToVal(t, context)
	case *mstype.MSOperationTypeS:	return MSFunctionFromType(t, e.env)
	case *mstype.MSArrayType:		return e.arrayTypeToVal(t)
	case *mstype.MSStructTypeS:		return e.structTypeToVal(t, context)
	case *mstype.MSNamedTypeS:		return e.namedTypeToVal(t, context)
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

func (e *MSEvaluator) compositeTypeToVal(ct *mstype.MSCompositeTypeS, context bool) MSVal {
	vals := make([]MSVal, len(ct.Types))
	for i, t := range ct.Types {
		vals[i] = e.typeToVal(t, context)
	}
	return MSTuple{Values: vals}
}

func (e *MSEvaluator) arrayTypeToVal(t *mstype.MSArrayType) MSVal {
	return MSArray{Values: make([]MSVal, 0), VType: t.Type}
}

func (e *MSEvaluator) structTypeToVal(st *mstype.MSStructTypeS, context bool) MSVal {

	if context {
		return MSStruct{Name: st.Name, Fields: nil, SType: st}
	}

	values := make(map[string]MSVal)
	for name, field := range st.Fields {
		values[name] = e.typeToVal(field, true)
	}

	return MSStruct{Name: st.Name, Fields: values, SType: st}
}

func (e *MSEvaluator) namedTypeToVal(nt *mstype.MSNamedTypeS, context bool) MSVal {

	var resolved mstype.MSType
	var err error

	// Look up type in type locals
	depth, ok := e.tlocals[nt]

	if ok {
		resolved, err = e.env.GetType(nt.Name, depth)
	} else {
		resolved, err = e.glb.GetType(nt.Name, depth)
	}

	// Yikes, todo better error handling
	if err != nil {
		return nil
	}

	return e.typeToVal(resolved, context)
}