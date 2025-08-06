package interp

// Eval
type EvalError struct {
	message string
}

func (ee *EvalError) Error() string {
	return "Evaluation error: " + ee.message
}


// Env
type EnvironmentError struct {
	message string
}

func (e *EnvironmentError) Error() string {
	return "Environment Error: " + e.message
}
