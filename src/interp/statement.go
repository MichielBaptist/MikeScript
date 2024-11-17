package interp

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/token"
)


func (evaluator *MSEvaluator) executeStatement(node *ast.StmtNodeI) EvalResult {
	switch node := (*node).(type) {
	case ast.Program:			return evaluator.executeStatements(&node)
	case ast.ExStmtNodeS:		return evaluator.executeExpressionStatement(&node)
	case ast.DeclarationNodeS:	return evaluator.executeDeclarationStatement(&node)
	case ast.BlockNodeS:		return evaluator.executeBlock(&node)
	case ast.IfNodeS:			return evaluator.executeIfstatement(&node)
	case ast.WhileNodeS:		return evaluator.executeWhileStatement(&node)
	case ast.ContinueNodeS:		return EvalResult{rt: RT_CONTINUE}
	case ast.BreakNodeS:		return EvalResult{rt: RT_BREAK}
	default:					return evalErr(fmt.Sprintf("Unknown statement type: %v", node))
	}
}

func (evaluator *MSEvaluator) executeBlock(node *ast.BlockNodeS) EvalResult {

	// Create a new environment with the current 
	// environment as the enclosing environment
	evaluator.env = NewEnvironment(evaluator.env)

	// Pop the environment when we are done with block
	defer func() {
		evaluator.env = evaluator.env.enclosing
	}()

	return evaluator.executeStatements(&ast.Program{Statements: node.Statements})
}

func (evaluator *MSEvaluator) executeWhileStatement(node *ast.WhileNodeS) EvalResult {

	var res EvalResult

	for {

		// Evaluate expression
		cond := evaluator.evaluateExpression(&node.Condition)
		if !cond.Valid() {
			return cond
		}
		if cond.rt != RT_BOOL {
			return evalErr(fmt.Sprintf("Condition must be of type bool, got %v", cond.rt))
		}

		// Get value of the bool
		condb, ok := cond.val.(bool)
		if !ok {
			return evalErr(fmt.Sprintf("Condition value is not a bool: %v", cond.val))
		}

		// Here the condition should be a boolean
		// If the condition is false, we break out of the loop
		if !condb {
			break
		}

		// Execute the body of the while loop
		res = evaluator.executeBlock(&node.Body)

		// Check if result has an error
		if !res.Valid() {
			return res
		}

		// Check if result is break, on break we exit
		if res.rt == RT_BREAK {
			return EvalResult{rt: RT_NONE}
		}
	}

	return res
}

func (evaluator *MSEvaluator) executeIfstatement(node *ast.IfNodeS) EvalResult {

	// Evaluate the condition
	cond := evaluator.evaluateExpression(&node.Condition)

	// Sanity checks
	if !cond.Valid() {
		return cond
	}
	if cond.rt != RT_BOOL {
		return evalErr(fmt.Sprintf("Condition must be of type bool, got %v", cond.rt))
	}
	if cond.val == nil {
		return evalErr("Condition value is nil.")
	}

	// Execute the then or else statement based on the condition
	switch cond.val.(type) {
	case bool:
		if cond.val.(bool) {
			return evaluator.executeStatement(&node.ThenStmt)
		} else if node.ElseStmt != nil {
			return evaluator.executeStatement(&node.ElseStmt)
		}
	default:
		return evalErr(fmt.Sprintf("Incompatible result type and value: %v: %v", cond.rt, cond.val))
	}

	return EvalResult{rt: RT_NONE}

}

func (evaluator *MSEvaluator) executeStatements(node *ast.Program) EvalResult {

	var res EvalResult
	
	// A ast.Program is a list of statements
	// So we just loop through the statements and evaluate them
	for _, stmt := range node.Statements {
		res = evaluator.executeStatement(&stmt)

		// If we notice an error during execution of this 
		// Satetement, we print the error and return
		if !res.Valid() {
			evaluator.statementError(res.err)
			return res
		}

		// Check if the result is break or continue
		// then exit the block statement. You don't need
		// to be in a loop to use break or continue.
		if res.rt == RT_BREAK || res.rt == RT_CONTINUE {
			break
		}
	}

	// The result of a ast.Program is the result of the last statement
	return res
}

func (evaluator *MSEvaluator) executeExpressionStatement(node *ast.ExStmtNodeS) EvalResult {
	return evaluator.evaluateExpression(&node.Ex)
}

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

func declaredTypeDefaultValue(tk token.Token) interface{} {
	switch tk.Type {
	case token.INT_TYPE:		return 0
	case token.FLOAT_TYPE:		return 0.0
	case token.STRING_TYPE:		return ""
	case token.BOOLEAN_TYPE:	return false
	default:					return nil
	}
}