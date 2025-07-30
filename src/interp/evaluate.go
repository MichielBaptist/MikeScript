package interp

import (
	"mikescript/src/ast"
)

////////////////////////////////////////////////////////////////////////
// Evaluator interface and implementations
////////////////////////////////////////////////////////////////////////

type MSEvaluator struct {
	// Contains program state.
	ast ast.Program			// The AST to evaluate
	err []error				// Evaluation errors
	env *Environment		// Environment of curresnt scope
	glb *Environment 		// Fixed reference to global scope (outermost env)

	locals map[ast.ExpNodeI]int
}

func NewMSEvaluator() *MSEvaluator {

	// Init gobal and local scope
	env := NewEnvironment(nil)
	glb := env

	// Add builtins to glb
	glb.NewVar("print", MSBuiltinPrint())
	glb.NewVar("print_env", MSBuiltinPrintEnv())

	return &MSEvaluator{env: env, glb: glb}
}

func (evaluator *MSEvaluator) Eval(ast ast.Program, locals map[ast.ExpNodeI]int) (MSVal, error) {
	
	// set the ast
	evaluator.ast = ast

	// set locals map
	evaluator.locals = locals

	// evaluate the ast
	return evaluator.executeStatements(&evaluator.ast)
	
}

func (evaluator *MSEvaluator) Errors() []error {
	return evaluator.err
}

func (evaluator *MSEvaluator) PrintEnv() {
	evaluator.env.printEnv()
}
