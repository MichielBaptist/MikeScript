package interp

import (
	"mikescript/src/ast"
	"mikescript/src/mstype"
)

func (evaluator *MSEvaluator) executeStatements(node *ast.Program) EvalResult {

	// By default the return value of a set of statements
	// is 'nothing' (cannot have value)
	res := EvalResult{Rt: mstype.MS_NOTHING}
	
	// A ast.Program is a list of statements
	// So we just loop through the statements and evaluate them
	for _, stmt := range node.Statements {
		res = evaluator.executeStatement(&stmt)

		// If we notice an error during execution of this 
		// Satetement, we print the error and return
		if !res.Valid() {
			evaluator.statementError(res.Err)
			return res
		}

		// Check if the result is break or continue
		// then exit the block statement. You don't need
		// to be in a loop to use break or continue.
		if res.IsType(&mstype.MS_BREAK) || res.IsType(&mstype.MS_CONTINUE){
			break
		}
	}

	// The result of a ast.Program is the result of the last statement
	return res
}