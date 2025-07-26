package interp

import (
	"mikescript/src/ast"
	"mikescript/src/mstype"
)

func (evaluator *MSEvaluator) executeStatements(node *ast.Program) EvalResult {

	for _, stmt := range node.Statements {
		
		res := evaluator.executeStatement(&stmt)

		// On error, break a
		if !res.Valid() {
			return evaluator.statementError(res)
		}

		// Break out the program on 'break', 'continue' and 'return'
		if 	res.IsType(&mstype.MS_BREAK) || 
			res.IsType(&mstype.MS_CONTINUE)|| 
			res.IsType(&mstype.MS_RETURN){
			return res
		}
	}

	// Ignore result of the last statement
	return EvalResult{Rt: mstype.MS_NOTHING}
}