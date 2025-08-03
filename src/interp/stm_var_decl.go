package interp

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/mstype"
)


func (evaluator *MSEvaluator) executeDeclarationStatement(node *ast.VarDeclNodeS) (MSVal, error) {

	// Get the default value for the type
	val := evaluator.typeToVal(node.Vartype, node.Identifier)

	if val == nil {
		_ = []int{}[0]
	}

	// Declare variable in env
	err := evaluator.env.NewVar(node.VarName(), val)

	return MSNothing{}, err
}

// TODO make evaluator receiver

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