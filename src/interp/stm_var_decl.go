package interp

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/token"
)


func (evaluator *MSEvaluator) executeDeclarationStatement(node *ast.DeclarationNodeS) EvalResult {

	// Map the Vartype token to returntype
	var rt ResultType = declaredTypeToReturnType(node.Vartype)

	// Check if we have a valid declaration type
	if rt == RT_INVALID {
		return evalErr(fmt.Sprintf("Unknown type '%v'", node.Vartype.Lexeme))
	}

	// Get the default value for the type
	val := declaredTypeDefaultValue(node.Vartype)

	if val == nil {
		return evalErr(fmt.Sprintf("Unknown type '%v'", node.Vartype.Lexeme))
	}

	// Declare variable in env
	err := evaluator.env.NewVar(node.Identifier.Name.Lexeme, EvalResult{rt: rt, val: val}, rt)

	if err != nil {
		return EvalResult{err: []error{err}}
	}

	return EvalResult{rt: rt}
}

func declaredTypeToReturnType(tk token.Token) ResultType {
	switch tk.Type {
	case token.INT_TYPE:		return RT_INT
	case token.FLOAT_TYPE:		return RT_FLOAT
	case token.STRING_TYPE:		return RT_STRING
	case token.BOOLEAN_TYPE:	return RT_BOOL
	default:					return RT_INVALID
	}
}

func declaredTypeDefaultValue(tk token.Token) any {
	switch tk.Type {
	case token.INT_TYPE:		return 0
	case token.FLOAT_TYPE:		return 0.0
	case token.STRING_TYPE:		return ""
	case token.BOOLEAN_TYPE:	return false
	default:					return nil
	}
}