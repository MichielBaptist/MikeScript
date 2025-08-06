package interp

type MSCallable interface {
	Call(e *MSEvaluator) (MSVal, error)	// native or foreign, executes code '=f'
	Bind(args []MSVal)	(MSVal, error)	// binds arguments, returns bound function 'a1, a2 >> f'
	Arity() int							// # args left to bind
}