package interp

import (
	"mikescript/src/ast"
	"mikescript/src/mstype"
)

////////////////////////////////////////////////////////////////////////
// Evaluator interface and implementations
////////////////////////////////////////////////////////////////////////

type MSEvaluator struct {
	ast *ast.Program				// The AST to evaluate
	err []error						// Evaluation errors
	env *Environment				// Environment of curresnt scope
	glb *Environment 				// Fixed reference to global scope (outermost env)
	vlocals map[*ast.VariableExpNodeS]int	// How deep do we need to go to resolve variables?
	tlocals map[*mstype.MSNamedTypeS]int 	// How deep do we need to go to resolve types?
}

func NewMSEvaluator() *MSEvaluator {

	env := NewEnvironment(nil)
	glb := env

	// Add builtins to glb
	glb.NewVar("print", MSBuiltinPrint())
	glb.NewVar("print_env", MSBuiltinPrintEnv())
	glb.NewVar("rand", MSBuiltinRand())

	return &MSEvaluator{
		env: env,
		glb: glb,
		vlocals: make(map[*ast.VariableExpNodeS]int),
		tlocals: make(map[*mstype.MSNamedTypeS]int),
	}
}

func (evaluator *MSEvaluator) UpdateVLocals(vlocals map[*ast.VariableExpNodeS]int) {
	for k, v := range vlocals {
		evaluator.vlocals[k] = v
	}
}

func (evaluator *MSEvaluator) UpdateTLocals(vlocals map[*mstype.MSNamedTypeS]int) {
	for k, v := range vlocals {
		evaluator.tlocals[k] = v
	}
}

func (evaluator *MSEvaluator) Eval(ast *ast.Program) (MSVal, error) {

	evaluator.ast = ast

	return evaluator.executeStatements(evaluator.ast)
}

func (evaluator *MSEvaluator) Errors() []error {
	return evaluator.err
}

func (evaluator *MSEvaluator) PrintEnv() {
	evaluator.env.printEnv()
}
