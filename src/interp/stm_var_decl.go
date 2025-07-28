package interp

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/mstype"
)


func (evaluator *MSEvaluator) executeDeclarationStatement(node *ast.VarDeclNodeS) (MSVal, error) {

	// // Map the Vartype token to returntype
	// var rt ResultType = declaredTypeToReturnType(node.Vartype)

	// // Check if we have a valid declaration type
	// if rt == mstype.RT_INVALID {
	// 	return evalErr(fmt.Sprintf("Tried to initialize an unknown type '%v'", node.Vartype.Lexeme))
	// }

	println("qslkdqklsjdlkqsj")

	// Get the default value for the type
	val := declaredTypeDefaultValue(&node.Vartype, &node.Identifier)

	println("qslkdqklsjdlkqsj")
	
	if val == nil {
		_ = []int{}[0]
	}

	name := node.Identifier.Name.Lexeme

	fmt.Println("Setting variable...")

	// Declare variable in env
	err := evaluator.env.NewVar(name, val)

	// Done
	return MSNothing{}, err
}

func declaredTypeDefaultValue(tk *mstype.MSType, name *ast.VariableExpNodeS) MSVal {
	// 3 cases:
	// 1. simple type --> direct map possible
	// 2. composite (undeclarable)
	// 3. function (undeclarable)

	switch t := (*tk).(type){
	case *mstype.MSSimpleTypeS: 	return resultTypeDefaultValue(t)
	case *mstype.MSCompositeTypeS:	return compositeTypeDefaultValue(t, name)
	case *mstype.MSOperationTypeS:	return MSFunctionFromType(t, name)
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

func compositeTypeDefaultValue(ct *mstype.MSCompositeTypeS, name *ast.VariableExpNodeS) MSVal {
	vals := make([]MSVal, len(ct.Types))
	for i, t := range ct.Types {
		vals[i] = declaredTypeDefaultValue(&t, name)
	}
	return MSTuple{Values: vals}
}