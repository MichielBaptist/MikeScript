package interp

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/mstype"
)


func (evaluator *MSEvaluator) executeDeclarationStatement(node *ast.VarDeclNodeS) (MSVal, error) {

	// Get the default value for the type
	val := declaredTypeDefaultValue(&node.Vartype, node.Identifier, evaluator.env)

	if val == nil {
		_ = []int{}[0]
	}

	// Declare variable in env
	err := evaluator.env.NewVar(node.VarName(), val)

	return MSNothing{}, err
}

// TODO make evaluator receiver

func declaredTypeDefaultValue(tk *mstype.MSType, name *ast.VariableExpNodeS, env *Environment) MSVal {
	// 3 cases:
	// 1. simple type --> direct map possible
	// 2. composite (undeclarable)
	// 3. function (undeclarable)

	switch t := (*tk).(type){
	case *mstype.MSSimpleTypeS: 	return resultTypeDefaultValue(t)
	case *mstype.MSCompositeTypeS:	return compositeTypeDefaultValue(t, name, env)
	case *mstype.MSOperationTypeS:	return MSFunctionFromType(t, name, env)
	default:						fmt.Printf("Found unknown type: '%s'\n", t)
	}
	return nil
}

func resultTypeDefaultValue(rt *mstype.MSSimpleTypeS) MSVal {
	switch rt.Rt {
	case mstype.RT_INT:		return MSInt{0}
	case mstype.RT_FLOAT:	return MSFloat{0.0}
	case mstype.RT_STRING:	return MSString{""}
	case mstype.RT_BOOL:	return MSBool{false}
	default:				return nil
	}
}

func compositeTypeDefaultValue(ct *mstype.MSCompositeTypeS, name *ast.VariableExpNodeS, env *Environment) MSVal {
	vals := make([]MSVal, len(ct.Types))
	for i, t := range ct.Types {
		vals[i] = declaredTypeDefaultValue(&t, name, env)
	}
	return MSTuple{Values: vals}
}