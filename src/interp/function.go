package interp

type FunctionResult interface {
	Call(e *MSEvaluator) EvalResult				// native or foreign, executes code
	Arity() int											// # args left to bind
	Bind(args []EvalResult)	EvalResult	// binds arguments, returns bound function
}