package interp

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/mstype"
	"mikescript/src/token"
)


func (evaluator *MSEvaluator) executeDeclarationStatement(node *ast.VarDeclNodeS) EvalResult {

	// // Map the Vartype token to returntype
	// var rt ResultType = declaredTypeToReturnType(node.Vartype)

	// // Check if we have a valid declaration type
	// if rt == mstype.RT_INVALID {
	// 	return evalErr(fmt.Sprintf("Tried to initialize an unknown type '%v'", node.Vartype.Lexeme))
	// }

	// Get the default value for the type
	val, err := declaredTypeDefaultValue(&node.Vartype)

	if err != nil {
		return EvalResult{err: []error{err}}
	}

	// Create var
	res := EvalResult{rt: node.Vartype, val: val}
	name := node.Identifier.Name.Lexeme

	// Declare variable in env
	err = evaluator.env.NewVar(name, res, node.Vartype)

	if err != nil {
		return EvalResult{err: []error{err}}
	}

	// The result of a variable declaration is Nothing
	// this is not the value declared.
	return EvalResult{rt: mstype.MS_NOTHING}
}

func declaredTypeToReturnType(tk token.Token) mstype.ResultType {
	switch tk.Type {
	case token.INT_TYPE:		return mstype.RT_INT
	case token.FLOAT_TYPE:		return mstype.RT_FLOAT
	case token.STRING_TYPE:		return mstype.RT_STRING
	case token.BOOLEAN_TYPE:	return mstype.RT_BOOL
	default:					return mstype.RT_INVALID
	}
}

func declaredTypeDefaultValue(tk *mstype.MSType) (any, error) {
	// 3 cases:
	// 1. simple type --> direct map possible
	// 2. composite (undeclarable)
	// 3. function (undeclarable)

	switch t := (*tk).(type){
	case *mstype.MSSimpleTypeS: return resultTypeDefaultValue(t.Rt), nil
	default:					return nil, &mstype.MSTypeError{Msg: fmt.Sprintf("Type '%s' has no default value.", t.String())}
	}

}

func resultTypeDefaultValue(rt mstype.ResultType) any {
	switch rt {
	case mstype.RT_INT:		return 0
	case mstype.RT_FLOAT:	return 0.0
	case mstype.RT_STRING:	return ""
	case mstype.RT_BOOL:	return false
	default:				return nil
	}
}