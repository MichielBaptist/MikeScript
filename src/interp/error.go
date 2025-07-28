package interp

type EvalError struct {
	message string
}

func (ee *EvalError) Error() string {
	return "Evaluation error: " + ee.message
}

