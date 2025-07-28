package interp

type MSCallable interface {
	Call(e *MSEvaluator) (MSVal, error)	// native or foreign, executes code
	Bind(args []MSVal)	(MSVal, error)	// binds arguments, returns bound function
	Arity() int							// # args left to bind
}