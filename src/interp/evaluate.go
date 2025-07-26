package interp

import (
	"mikescript/src/ast"
	"mikescript/src/mstype"
)

////////////////////////////////////////////////////////////////////////
// Evaluator interface and implementations
////////////////////////////////////////////////////////////////////////

type MSEvaluator struct {
	// Contains program state.
	ast ast.Program			// The AST to evaluate
	err []error				// Evaluation errors
	env *Environment		// Environment of current scope
	glb *Environment 		// Fixed reference to global scope (outermost env)
}

func NewMSEvaluator() *MSEvaluator {

	// Init gobal and local scope
	env := NewEnvironment(nil)
	glb := env

	// Add builtins to glb
	glb.NewVar("print", MSBuiltinPrint(), &mstype.MSOperationTypeS{Right: mstype.MS_NOTHING})

	return &MSEvaluator{env: env, glb: glb}
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

func (evaluator *MSEvaluator) statementError(e EvalResult) EvalResult{
	evaluator.err = append(evaluator.err, e.Err...)
	return e
}

type EvalError struct {
	message string
}

func (ee *EvalError) Error() string {
	return "Evaluation error: " + ee.message
}

func evalErr(msg string) EvalResult {
	return EvalResult{Err: []error{&EvalError{msg}}}
}

