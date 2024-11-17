package interp

import (
	"fmt"
	"mikescript/src/ast"
)

////////////////////////////////////////////////////////////////////////
// Evaluator interface and implementations
////////////////////////////////////////////////////////////////////////

type MSEvaluator struct {
	ast ast.Program			// The AST to evaluate
	err []error				// Evaluation errors
	env *Environment		// Global environment
}

func NewMSEvaluator() *MSEvaluator {
	return &MSEvaluator{env: NewEnvironment(nil)}
}

func (evaluator *MSEvaluator) Eval(ast ast.Program) EvalResult {

	// set the ast
	evaluator.ast = ast

	// evaluate the ast
	return evaluator.executeStatements(&evaluator.ast)
	
}

func (evaluator *MSEvaluator) Errors() []error {
	return evaluator.err
}

func (evaluator *MSEvaluator) PrintEnv() {
	evaluator.env.printEnv()
}

func (evaluator *MSEvaluator) statementError(err []error) {
	evaluator.err = append(evaluator.err, err...)
}

////////////////////////////////////////////////////////////////////////
// Evaluation result
////////////////////////////////////////////////////////////////////////

type EvalResult struct {
	rt    ResultType  	// type of the result
	val   interface{} 	// Container for the result
	err   []error    	// Error message on evaluation fail
}

func (er *EvalResult) Valid() bool {
	// 1. Check if there are any errors
	// 2. Check if the return type is valid wrt the value
	return len(er.err) == 0
}

func (er *EvalResult) ExpectedType(rt ResultType) bool {
	return er.rt == rt && er.ValidReturnType()
}

func (er *EvalResult) ValidReturnType() bool {
	switch er.rt {
	case RT_INT:
		_, ok := er.val.(int)
		return ok
	case RT_FLOAT:
		_, ok := er.val.(float64)
		return ok
	case RT_STRING:
		_, ok := er.val.(string)
		return ok
	case RT_BOOL:
		_, ok := er.val.(bool)
		return ok
	case RT_TUPLE:
		_, ok := er.val.([]EvalResult)
		return ok
	case RT_FUNCTION:
		return true
	case RT_INVALID:
		return true
	case RT_NONE:
		return true
	default:
		fmt.Println("Error trying to check valid return type for unknown type:", er.rt)
		return false
	}
}

func (er *EvalResult) ReturnType() ResultType {
	return er.rt
}

func (er EvalResult) String() string {
	if !er.Valid() {
		return fmt.Sprintf("Error: %v", er.err)
	}
	switch er.rt {
	case RT_INT:		return fmt.Sprintf("%v", er.val)
	case RT_FLOAT:		return fmt.Sprintf("%v", er.val)
	case RT_STRING:		return fmt.Sprintf("\"%v\"", er.val)
	case RT_BOOL:		return fmt.Sprintf("%v", er.val)
	case RT_TUPLE:		return tupleToString(er.val.([]EvalResult))
	case RT_FUNCTION:	return fmt.Sprintf("%v", er.val)
	case RT_INVALID:	return "Invalid"
	case RT_NONE:		return "Nothing"
	case RT_BREAK:		return "Break"
	case RT_CONTINUE:	return "Continue"
	default:			return "Unknown"
	}
}

func tupleToString(t []EvalResult) string {
	s := "("
	for i, v := range t {
		s += v.String()
		if i < len(t)-1 {
			s += ", "
		}
	}
	s += ")"
	return s
}

////////////////////////////////////////////////////////////////////////
// Evaluation error
////////////////////////////////////////////////////////////////////////

type EvalError struct {
	message string
}

func (ee *EvalError) Error() string {
	return "Evaluation error: " + ee.message
}

////////////////////////////////////////////////////////////////////////
// Result type
////////////////////////////////////////////////////////////////////////

type ResultType uint8
const (
	RT_INVALID ResultType = iota
	RT_NONE
	RT_INT
	RT_FLOAT
	RT_STRING
	RT_BOOL
	RT_TUPLE
	RT_FUNCTION
	RT_BREAK
	RT_CONTINUE
)

func (res ResultType) String() string {
	switch res {
	case RT_INVALID:	return "invalid"
	case RT_NONE:		return "none"
	case RT_INT:		return "int"
	case RT_FLOAT:		return "float"
	case RT_STRING:		return "string"
	case RT_BOOL:		return "bool"
	case RT_TUPLE:		return "tuple"
	case RT_FUNCTION:	return "function"
	case RT_BREAK:		return "break"
	case RT_CONTINUE:	return "continue"
	default:			return "unknown"
	}
}

func evalErr(msg string) EvalResult {
	return EvalResult{err: []error{&EvalError{msg}}}
}

////////////////////////////////////////////////////////////////////////
// Helper functions
////////////////////////////////////////////////////////////////////////

func SameType(rt1, rt2 ResultType) bool {
	return rt1 == rt2
}

func StringType(rt ResultType) bool {
	return rt == RT_STRING
}
