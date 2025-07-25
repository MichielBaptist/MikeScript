package interp

import (
	"mikescript/src/mstype"
	"mikescript/src/token"
)

func evalGreater(left, right EvalResult, op token.TokenType) EvalResult {

	switch {
	case left.IsType(&mstype.MS_INT) && right.IsType(&mstype.MS_INT):
		return EvalResult{rt: mstype.MS_BOOL, val: left.val.(int) > right.val.(int)}
	case left.IsType(&mstype.MS_INT) && right.IsType(&mstype.MS_FLOAT):
		return EvalResult{rt: mstype.MS_BOOL, val: float64(left.val.(int)) > right.val.(float64)}
	case left.IsType(&mstype.MS_FLOAT) && right.IsType(&mstype.MS_INT):
		return EvalResult{rt: mstype.MS_BOOL, val: left.val.(float64) > float64(right.val.(int))}
	case left.IsType(&mstype.MS_FLOAT) && right.IsType(&mstype.MS_FLOAT):
		return EvalResult{rt: mstype.MS_BOOL, val: left.val.(float64) > right.val.(float64)}
	}
	return evalErr(invalidBinop(left, right, op.String()))
}

func evalGreaterEq(left, right EvalResult, op token.TokenType) EvalResult {

	switch {
	case left.IsType(&mstype.MS_INT) && right.IsType(&mstype.MS_INT):
		return EvalResult{rt: mstype.MS_BOOL, val: left.val.(int) >= right.val.(int)}
	case left.IsType(&mstype.MS_INT) && right.IsType(&mstype.MS_FLOAT):
		return EvalResult{rt: mstype.MS_BOOL, val: float64(left.val.(int)) >= right.val.(float64)}
	case left.IsType(&mstype.MS_FLOAT) && right.IsType(&mstype.MS_INT):
		return EvalResult{rt: mstype.MS_BOOL, val: left.val.(float64) >= float64(right.val.(int))}
	case left.IsType(&mstype.MS_FLOAT) && right.IsType(&mstype.MS_FLOAT):
		return EvalResult{rt: mstype.MS_BOOL, val: left.val.(float64) >= right.val.(float64)}
	}
	return evalErr(invalidBinop(left, right, op.String()))

}

func evalEq(left EvalResult, right EvalResult, op token.TokenType) EvalResult {
	switch {
	case left.IsType(&mstype.MS_NOTHING) || right.IsType(&mstype.MS_NOTHING):
		return EvalResult{rt: mstype.MS_BOOL, val: false}
	case left.IsType(&mstype.MS_STRING) && right.IsType(&mstype.MS_STRING):
		return EvalResult{rt: mstype.MS_BOOL, val: left.val.(string) == right.val.(string)}
	case left.IsType(&mstype.MS_STRING) || right.IsType(&mstype.MS_STRING):
		return EvalResult{rt: mstype.MS_BOOL, val: false}
	case left.IsType(&mstype.MS_BOOL) && right.IsType(&mstype.MS_BOOL):
		return EvalResult{rt: mstype.MS_BOOL, val: left.val.(bool) == right.val.(bool)}
	case left.IsType(&mstype.MS_BOOL) || right.IsType(&mstype.MS_BOOL):
		return EvalResult{rt: mstype.MS_BOOL, val: false}
	case left.IsType(&mstype.MS_INT) && right.IsType(&mstype.MS_INT):
		return EvalResult{rt: mstype.MS_BOOL, val: left.val.(int) == right.val.(int)}
	case left.IsType(&mstype.MS_FLOAT) && right.IsType(&mstype.MS_FLOAT):
		return EvalResult{rt: mstype.MS_BOOL, val: left.val.(float64) == right.val.(float64)}
	case left.IsType(&mstype.MS_INT) && right.IsType(&mstype.MS_FLOAT):
		return EvalResult{rt: mstype.MS_BOOL, val: float64(left.val.(int)) == right.val.(float64)}
	case left.IsType(&mstype.MS_FLOAT) && right.IsType(&mstype.MS_INT):
		return EvalResult{rt: mstype.MS_BOOL, val: left.val.(float64) == float64(right.val.(int))}
	default:
		return evalErr(invalidBinop(left, right, op.String()))
	}
}

func evalNeq(_, _ EvalResult, _ token.TokenType) EvalResult {
	return evalErr("Not equal should be replaced with equality and unary negation. If you see this error, the parser is broken.")
}

func evalGrGr(_, _ EvalResult) EvalResult {
	return evalErr("Trying to evaluate function call as binary operator. If you see this error, the parser is broken.")
}