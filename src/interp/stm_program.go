package interp

import (
	"mikescript/src/ast"
)

func (evaluator *MSEvaluator) executeStatements(node *ast.Program) (MSVal, error) {

	var res MSVal
	var err error

	for _, stmt := range node.Statements {
		
		res, err = evaluator.executeStatement(&stmt)

		// On error, break
		if err != nil {
			return MSNothing{}, err
		}

		switch res.(type){
		case MSReturn:		return res, nil
		case MSBreak:		return res, nil
		case MSContinue:	return res, nil
		}
	}

	// Ignore result of the last statement
	return res, err
}