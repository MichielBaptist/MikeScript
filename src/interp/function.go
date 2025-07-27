package interp

type FunctionResult interface {
	Call(e *MSEvaluator) EvalResult				// native or foreign, executes code
	Bind(args []EvalResult)	EvalResult			// binds arguments, returns bound function
	Arity() int									// # args left to bind
}