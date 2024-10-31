package main

import (
	"fmt"
	"strconv"
)

////////////////////////////////////////////////////////////////////////
// Evaluator interface and implementations
////////////////////////////////////////////////////////////////////////

type Evaluator interface {
	eval(node ExpNodeI) EvalResult
}

type MSEvaluator struct {
	ast Program			// The AST to evaluate
	err []error			// Error messages
	env Environment
}

func NewMSEvaluator() *MSEvaluator {
	return &MSEvaluator{env: *NewEnvironment()}
}

func (evaluator *MSEvaluator) Eval(ast Program) EvalResult {

	// set the ast
	evaluator.ast = ast

	// evaluate the ast
	res := evaluator.executeStatements(&evaluator.ast)

	return res

}

func (evaluator *MSEvaluator) statementError(err []error) {
	evaluator.err = append(evaluator.err, err...)
}

////////////////////////////////////////////////////////////////////////
// Execute statements
////////////////////////////////////////////////////////////////////////


func (evaluator *MSEvaluator) executeStatement(node *StmtNodeI) EvalResult {
	switch node := (*node).(type) {
	case Program:			return evaluator.executeStatements(&node)
	case ExStmtNodeS:		return evaluator.executeExpressionStatement(&node)
	case DeclarationNodeS:	return evaluator.executeDeclarationStatement(&node)
	default:				return EvalErr("Unknown statement type.")
	}
}

func (evaluator *MSEvaluator) executeStatements(node *Program) EvalResult {

	var res EvalResult
	
	// A program is a list of statements
	// So we just loop through the statements and evaluate them
	for _, stmt := range node.Statements {
		res = evaluator.executeStatement(&stmt)

		// If we notice an error during execution of this 
		// Satetement, we print the error and return
		if !res.Valid() {
			evaluator.statementError(res.err)
			return res
		}
	}

	// The result of a program is the result of the last statement
	return res
}

func (evaluator *MSEvaluator) executeExpressionStatement(node *ExStmtNodeS) EvalResult {
	return evaluator.evaluateExpression(&node.Ex)
}

func (evaluator *MSEvaluator) executeDeclarationStatement(node *DeclarationNodeS) EvalResult {

	// Map the Vartype token to returntype
	var rt ResultType = declaredTypeToReturnType(node.Vartype)

	// Check if we have a valid declaration type
	if rt == RT_INVALID {
		return EvalErr(fmt.Sprintf("Unknown type '%v'", node.Vartype.Lexeme))
	}

	// Get the default value for the type
	val := declaredTypeDefaultValue(node.Vartype)

	if val == nil {
		return EvalErr(fmt.Sprintf("Unknown type '%v'", node.Vartype.Lexeme))
	}

	// Declare variable in env
	err := evaluator.env.NewVar(node.Identifier.Name.Lexeme, EvalResult{rt: rt, val: val}, rt)

	if err != nil {
		return EvalResult{err: []error{err}}
	}

	return EvalResult{rt: rt}
}

func declaredTypeToReturnType(tk Token) ResultType {
	switch tk.Type {
	case INT_TYPE:			return RT_INT
	case FLOAT_TYPE:		return RT_FLOAT
	case STRING_TYPE:		return RT_STRING
	case BOOLEAN_TYPE:		return RT_BOOL
	default:				return RT_INVALID
	}
}

func declaredTypeDefaultValue(tk Token) interface{} {
	switch tk.Type {
	case INT_TYPE:			return 0
	case FLOAT_TYPE:		return 0.0
	case STRING_TYPE:		return ""
	case BOOLEAN_TYPE:		return false
	default:				return nil
	}
}

////////////////////////////////////////////////////////////////////////
// Evaluate expressions
////////////////////////////////////////////////////////////////////////

func (evaluator *MSEvaluator) evaluateExpression(node *ExpNodeI) EvalResult {
	switch node := (*node).(type) {
	case BinaryExpNodeS:	return evaluator.evaluateBinaryExpression(&node)
	case UnaryExpNodeS:		return evaluator.evaluateUnaryExpression(&node)
	case LiteralExpNodeS:	return evaluator.evaluateLiteralExpression(&node)
	case GroupExpNodeS:		return evaluator.evaluateGroupExpression(&node)
	case AssignmentNodeS:	return evaluator.evaluateAssignmentExpression(&node)
	case FuncAppNodeS:		return evaluator.evaluateFunctionApplication(&node)
	case VariableExpNodeS:	return evaluator.evalVariable(&node)
	default:				return EvalErr(fmt.Sprintf("Unknown expression type: %v", node))
	}
}

func (evaluator *MSEvaluator) evalVariable(node *VariableExpNodeS) EvalResult {

	// Get the value from the environment
	val, err := evaluator.env.Get(node.Name.Lexeme)

	if err != nil {
		return EvalErr(fmt.Sprintf("Variable '%v' is not defined.", node.Name.Lexeme))
	}

	return val
	
}

func (evaluator *MSEvaluator) evaluateBinaryExpression(node *BinaryExpNodeS) EvalResult {

	// Evaluate left and right side of the expression
	left := evaluator.evaluateExpression(&node.Left)
	right := evaluator.evaluateExpression(&node.Right)

	// If either of the evaluations failed, return the error
	if !left.Valid() || !right.Valid() {
		return EvalResult{err: append(left.err, right.err...)}
	}

	switch node.Op.Type {
	case PLUS: 				return evalAdd(left, right)
	case MINUS:				return evalSub(left, right)
	case MULT:				return evalMult(left, right)
	case SLASH:				return evalDiv(left, right)
	case GREATER:			return evalGreater(left, right, node.Op.Type)
	case LESS:				return evalGreater(right, left, node.Op.Type)
	case GREATER_EQ:		return evalGreaterEq(left, right, node.Op.Type)
	case LESS_EQ:			return evalGreaterEq(right, left, node.Op.Type)
	case EQ_EQ:				return evalEq(left, right, node.Op.Type)
	case EXCLAMATION_EQ:	return evalNeq(left, right, node.Op.Type)
	case GREATER_GREATER: 	return evalGrGr(left, right)
	case COMMA:				return evalTuple(left, right)
	default:				return EvalErr("Unknown binary operator.")
	}
}

func (evaluator *MSEvaluator) evaluateUnaryExpression(node *UnaryExpNodeS) EvalResult {
	// evaluate the node
	res := evaluator.evaluateExpression(&node.Node)

	// check if the evaluation was.Valid()
	if !res.Valid() {
		return res
	}

	// handle unary operators
	switch node.Op.Type {
	case MINUS:
		switch res.rt {
		case RT_INT:	return EvalResult{rt: RT_INT, val: -res.val.(int)}
		case RT_FLOAT:	return EvalResult{rt: RT_FLOAT, val: -res.val.(float64)}
		default:		return EvalErr(fmt.Sprintf("Cannot negate type %v", res.rt))
		}
	case EXCLAMATION:
		switch res.rt {
		case RT_BOOL:	return EvalResult{rt: RT_BOOL, val: !res.val.(bool)}
		default:		return EvalErr(fmt.Sprintf("Cannot negate type %v", res.rt))
		}
	default:
		return EvalErr(fmt.Sprintf("Operator %v is not defined for type %v", node.Op.Lexeme, res.rt))
	}
}

func (evaluator *MSEvaluator) evaluateLiteralExpression(node *LiteralExpNodeS) EvalResult {
	switch node.Tk.Type {
	case NUMBER_INT:	return evalIntLiteral(node)
	case NUMBER_FLOAT:	return evalFloatLiteral(node)
	case STRING:		return evalStringLiteral(node)
	case TRUE:			return EvalResult{rt: RT_BOOL, val: true}
	case FALSE:			return EvalResult{rt: RT_BOOL, val: false}
	case IDENTIFIER:	return EvalErr(fmt.Sprintf("Trying to evaluate identifier '%v' as a literal.", node.Tk.Lexeme))
	default:			return EvalErr(fmt.Sprintf("Literal type '%v' is not defined.", node.Tk.Type))
	}
}

func (evaluator *MSEvaluator) evaluateGroupExpression(node *GroupExpNodeS) EvalResult {
	return evaluator.evaluateExpression(&node.Node)
}

func (evaluator *MSEvaluator) evaluateAssignmentExpression(node *AssignmentNodeS) EvalResult {
	// Evaluate the expression
	res := evaluator.evaluateExpression(&node.Exp)

	// set the variable in the environment
	err := evaluator.env.Set(node.Identifier.Name.Lexeme, res)

	if err != nil {
		return EvalResult{err: []error{err}}
	}

	return EvalResult{rt: res.rt, val: res.val}
}

func (evaluator *MSEvaluator) evaluateFunctionApplication(_ *FuncAppNodeS) EvalResult {
	return EvalErr("Function application is not implemented yet.")
}



////////////////////////////////////////////////////////////////////////
// Evaluation result
////////////////////////////////////////////////////////////////////////

type EvalResult struct {
	rt    ResultType  	// type of the result
	val   interface{} 	// Container for the result
	err   []error    	// Error message on evaluation fail
}

func (er EvalResult) Valid() bool {
	return len(er.err) == 0
}

func (er EvalResult) ReturnType() ResultType {
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
	default:			return "unknown"
	}
}

func EvalErr(msg string) EvalResult {
	err := []error{&EvalError{msg}}
	return EvalResult{err: err}
}



////////////////////////////////////////////////////////////////////////
// Evaluation literals
////////////////////////////////////////////////////////////////////////

func evalIntLiteral(node *LiteralExpNodeS) EvalResult {
	// convert the lexeme to an int
	val, err := strconv.Atoi(node.Tk.Lexeme)

	if err != nil {
		return EvalErr(fmt.Sprintf("Could not convert '%v' to int.", node.Tk.Lexeme))
	}

	return EvalResult{rt: RT_INT, val: val}
}

func evalFloatLiteral(node *LiteralExpNodeS) EvalResult {
	// convert the lexeme to a float
	val, err := strconv.ParseFloat(node.Tk.Lexeme, 64)

	if err != nil {
		return EvalErr(fmt.Sprintf("Could not convert '%v' to float64.", node.Tk.Lexeme))
	}

	return EvalResult{rt: RT_FLOAT, val: val}
}

func evalStringLiteral(node *LiteralExpNodeS) EvalResult {
	return EvalResult{rt: RT_STRING, val: node.Tk.Lexeme}
}

////////////////////////////////////////////////////////////////////////
// Evaluation Binary operations                 					  //
////////////////////////////////////////////////////////////////////////

func invalidBinop(left, right EvalResult, op string) string {
	return fmt.Sprintf("Operator '%v' is not defined for types '%v' and '%v'", op, left.rt, right.rt)
}

func evalAdd(left, right EvalResult) EvalResult {

	// Prepare error message
	errRes := EvalErr(invalidBinop(left, right, "+"))

	// Handle all numeric/boolean cases
	switch left.rt {
	case RT_INT:
		switch right.rt {
		case RT_INT:	return EvalResult{rt: RT_INT, val: left.val.(int) + right.val.(int)}
		case RT_FLOAT:	return EvalResult{rt: RT_FLOAT, val: float64(left.val.(int)) + right.val.(float64)}
		case RT_BOOL:	return EvalResult{rt: RT_INT, val: left.val.(int) + boolToInt(right.val.(bool))}
		}
	case RT_FLOAT:
		switch right.rt {
		case RT_INT:	return EvalResult{rt: RT_FLOAT, val: left.val.(float64) + float64(right.val.(int))}
		case RT_FLOAT:	return EvalResult{rt: RT_FLOAT, val: left.val.(float64) + right.val.(float64)}
		case RT_BOOL:	return EvalResult{rt: RT_FLOAT, val: left.val.(float64) + float64(boolToInt(right.val.(bool)))}
		}
	case RT_BOOL:
		// Cast left value to int
		boolint := boolToInt(left.val.(bool))
		switch right.rt {
		case RT_INT:	return EvalResult{rt: RT_INT, val: boolint + right.val.(int)}
		case RT_FLOAT:	return EvalResult{rt: RT_FLOAT, val: float64(boolint) + right.val.(float64)}
		case RT_BOOL:	return EvalResult{rt: RT_INT, val: boolint + boolToInt(right.val.(bool))}
		}
	case RT_STRING:
		switch right.rt {
		case RT_STRING:	return EvalResult{rt: RT_STRING, val: left.val.(string) + right.val.(string)}
		}
	case RT_TUPLE:
		switch right.rt{
		case RT_TUPLE:	return EvalResult{rt: RT_TUPLE, val: append(left.val.([]EvalResult), right.val.([]EvalResult)...)}
		}
	}
	return errRes
}

func evalSub(_, _ EvalResult) EvalResult {
	// This code should never happen as there is no
	// way in which the parser can generate a subtraction
	// instead it will generate an addition with a negative unary
	return EvalErr("Subtraction should be replaced with addition and unary negation. If you see this error, the parser is broken.")
}

func evalMult(left, right EvalResult) EvalResult {

	switch left.rt {
	case RT_STRING:
		switch right.rt {
		case RT_INT:	return EvalResult{rt: RT_STRING, val: repeatString(left.val.(string), right.val.(int))}
		}
	case RT_INT:
		switch right.rt {
		case RT_INT:	return EvalResult{rt: RT_INT, val: left.val.(int) * right.val.(int)}
		case RT_FLOAT:	return EvalResult{rt: RT_FLOAT, val: float64(left.val.(int)) * right.val.(float64)}
		case RT_BOOL:	return EvalResult{rt: RT_INT, val: left.val.(int) * boolToInt(right.val.(bool))}
		case RT_STRING: return EvalResult{rt: RT_STRING, val: repeatString(right.val.(string), left.val.(int))}
		}
	case RT_FLOAT:
		switch right.rt {
		case RT_INT:	return EvalResult{rt: RT_FLOAT, val: left.val.(float64) * float64(right.val.(int))}
		case RT_FLOAT:	return EvalResult{rt: RT_FLOAT, val: left.val.(float64) * right.val.(float64)}
		case RT_BOOL:	return EvalResult{rt: RT_FLOAT, val: left.val.(float64) * float64(boolToInt(right.val.(bool)))}
		}
	case RT_BOOL:
		switch right.rt {
		case RT_INT:	return EvalResult{rt: RT_INT, val: boolToInt(left.val.(bool)) * right.val.(int)}
		case RT_FLOAT:	return EvalResult{rt: RT_FLOAT, val: float64(boolToInt(left.val.(bool))) * right.val.(float64)}
		case RT_BOOL:	return EvalResult{rt: RT_INT, val: boolToInt(left.val.(bool)) * boolToInt(right.val.(bool))}
		}
	}
	return EvalErr(invalidBinop(left, right, "*"))
}

func evalDiv(left, right EvalResult) EvalResult {

	// prepare error message
	errRes := EvalErr(invalidBinop(left, right, "/"))
	errDiv := EvalErr("Division by zero.")

	// Div with dynamic type checking
	switch {
	case left.rt == RT_NONE || right.rt == RT_NONE:
		return errRes
	case left.rt == RT_STRING || right.rt == RT_STRING:
		return errRes
	}

	// Handle all numeric/boolean cases
	switch left.rt {
	case RT_INT:
		switch right.rt {
		case RT_INT:
			if right.val.(int) == 0 {
				return errDiv
			}
			return EvalResult{rt: RT_FLOAT, val: float64(left.val.(int)) / float64(right.val.(int))}
		case RT_FLOAT:
			if right.val.(float64) == 0 {
				return errDiv
			}
			return EvalResult{rt: RT_FLOAT, val: float64(left.val.(int)) / right.val.(float64)}
		case RT_BOOL:
			if boolToInt(right.val.(bool)) == 0 {
				return errDiv
			}
			return EvalResult{rt: RT_FLOAT, val: left.val.(int) / boolToInt(right.val.(bool))}
		default:
			return EvalErr(fmt.Sprintf("Cannot divide '%v' with type '%v'", left.rt, right.rt))
		}
	case RT_FLOAT:
		switch right.rt {
		case RT_INT:
			if right.val.(int) == 0 {
				return errDiv
			}
			return EvalResult{rt: RT_FLOAT, val: left.val.(float64) / float64(right.val.(int))}
		case RT_FLOAT:
			if right.val.(float64) == 0 {
				return errDiv
			}
			return EvalResult{rt: RT_FLOAT, val: left.val.(float64) / right.val.(float64)}
		case RT_BOOL:
			if boolToInt(right.val.(bool)) == 0 {
				return errDiv
			}
			return EvalResult{rt: RT_FLOAT, val: left.val.(float64) / float64(boolToInt(right.val.(bool)))}
		}
	case RT_BOOL:
		switch right.rt {
		case RT_INT:
			if right.val.(int) == 0 {
				return EvalErr("Division by zero.")
			}
			return EvalResult{rt: RT_FLOAT, val: boolToInt(left.val.(bool)) / right.val.(int)}
		case RT_FLOAT:
			if right.val.(float64) == 0 {
				return EvalErr("Division by zero.")
			}
			return EvalResult{rt: RT_FLOAT, val: float64(boolToInt(left.val.(bool))) / right.val.(float64)}
		case RT_BOOL:
			if boolToInt(right.val.(bool)) == 0 {
				return EvalErr("Division by zero.")
			}
			return EvalResult{rt: RT_FLOAT, val: boolToInt(left.val.(bool)) / boolToInt(right.val.(bool))}
		default:
			return EvalErr(invalidBinop(left, right, "/"))
		}
	}
	return EvalErr(invalidBinop(left, right, "/"))
}

func evalGreater(left, right EvalResult, op TokenType) EvalResult {

	switch left.rt {
	case RT_INT:
		switch right.rt {
		case RT_INT:	return EvalResult{rt: RT_BOOL, val: left.val.(int) > right.val.(int)}
		case RT_FLOAT:	return EvalResult{rt: RT_BOOL, val: float64(left.val.(int)) > right.val.(float64)}
		}
	case RT_FLOAT:
		switch right.rt {
		case RT_INT:	return EvalResult{rt: RT_BOOL, val: left.val.(float64) > float64(right.val.(int))}
		case RT_FLOAT:	return EvalResult{rt: RT_BOOL, val: left.val.(float64) > right.val.(float64)}
		}
	}
	return EvalErr(invalidBinop(left, right, op.String()))
}

func evalGreaterEq(left, right EvalResult, op TokenType) EvalResult {

	switch left.rt {
	case RT_INT:
		switch right.rt {
		case RT_INT:	return EvalResult{rt: RT_BOOL, val: left.val.(int) >= right.val.(int)}
		case RT_FLOAT:	return EvalResult{rt: RT_BOOL, val: float64(left.val.(int)) >= right.val.(float64)}
		}
	case RT_FLOAT:
		switch right.rt {
		case RT_INT:	return EvalResult{rt: RT_BOOL, val: left.val.(float64) >= float64(right.val.(int))}
		case RT_FLOAT:	return EvalResult{rt: RT_BOOL, val: left.val.(float64) >= right.val.(float64)}
		}
	}
	return EvalErr(invalidBinop(left, right, op.String()))
}

func evalEq(left, right EvalResult, op TokenType) EvalResult {
	switch {
	case left.rt == RT_NONE || right.rt == RT_NONE:		return EvalResult{rt: RT_BOOL, val: false}
	case left.rt == RT_STRING && right.rt == RT_STRING:	return EvalResult{rt: RT_BOOL, val: left.val.(string) == right.val.(string)}
	case left.rt == RT_STRING || right.rt == RT_STRING:	return EvalResult{rt: RT_BOOL, val: false}
	case left.rt == RT_BOOL && right.rt == RT_BOOL:		return EvalResult{rt: RT_BOOL, val: left.val.(bool) == right.val.(bool)}
	case left.rt == RT_BOOL || right.rt == RT_BOOL:		return EvalResult{rt: RT_BOOL, val: false}
	case left.rt == RT_INT && right.rt == RT_INT:		return EvalResult{rt: RT_BOOL, val: left.val.(int) == right.val.(int)}
	case left.rt == RT_FLOAT && right.rt == RT_FLOAT:	return EvalResult{rt: RT_BOOL, val: left.val.(float64) == right.val.(float64)}
	case left.rt == RT_INT && right.rt == RT_FLOAT:		return EvalResult{rt: RT_BOOL, val: float64(left.val.(int)) == right.val.(float64)}
	case left.rt == RT_FLOAT && right.rt == RT_INT:		return EvalResult{rt: RT_BOOL, val: left.val.(float64) == float64(right.val.(int))}
	default:											return EvalErr(invalidBinop(left, right, op.String()))
	}
}

func evalNeq(_, _ EvalResult, _ TokenType) EvalResult {
	return EvalErr("Not equal should be replaced with equality and unary negation. If you see this error, the parser is broken.")
}

func evalGrGr(left, right EvalResult) EvalResult {
	// Currently only 'print' is implemented
	// so return an error for all other values
	if right.rt != RT_FUNCTION {
		return EvalErr(fmt.Sprintf("Tried to invoke non-function type '%v'", right.val))
	}

	// Interpretor should find the function and call it
	// with the left value as the argument. If the left
	// value is a tuple, the function should be called
	// where the tuple entries are the arguments.
	switch right.val.(type) {
	case string:
		if right.val.(string) == "print" {
			fmt.Println(left.val)
			return EvalResult{rt: RT_NONE}
		}
		return EvalErr(fmt.Sprintf("Function '%v' is not implemented.", right.val))
	default:
		return EvalErr("Function invocation is not implemented yet.")
	}
	
}

func evalTuple(left, right EvalResult) EvalResult {
	return EvalResult{rt: RT_TUPLE, val: []EvalResult{left, right}}
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
