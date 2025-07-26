package interp

import (
	"mikescript/src/mstype"
	"mikescript/src/token"
)

func evalGreater(left, right EvalResult, op token.TokenType) EvalResult {

	switch {
	case left.IsType(&mstype.MS_INT) && right.IsType(&mstype.MS_INT):
		return EvalResult{Rt: mstype.MS_BOOL, Val: left.Val.(int) > right.Val.(int)}
	case left.IsType(&mstype.MS_INT) && right.IsType(&mstype.MS_FLOAT):
		return EvalResult{Rt: mstype.MS_BOOL, Val: float64(left.Val.(int)) > right.Val.(float64)}
	case left.IsType(&mstype.MS_FLOAT) && right.IsType(&mstype.MS_INT):
		return EvalResult{Rt: mstype.MS_BOOL, Val: left.Val.(float64) > float64(right.Val.(int))}
	case left.IsType(&mstype.MS_FLOAT) && right.IsType(&mstype.MS_FLOAT):
		return EvalResult{Rt: mstype.MS_BOOL, Val: left.Val.(float64) > right.Val.(float64)}
	}
	return evalErr(invalidBinop(left, right, op.String()))
}

func evalGreaterEq(left, right EvalResult, op token.TokenType) EvalResult {

	switch {
	case left.IsType(&mstype.MS_INT) && right.IsType(&mstype.MS_INT):
		return EvalResult{Rt: mstype.MS_BOOL, Val: left.Val.(int) >= right.Val.(int)}
	case left.IsType(&mstype.MS_INT) && right.IsType(&mstype.MS_FLOAT):
		return EvalResult{Rt: mstype.MS_BOOL, Val: float64(left.Val.(int)) >= right.Val.(float64)}
	case left.IsType(&mstype.MS_FLOAT) && right.IsType(&mstype.MS_INT):
		return EvalResult{Rt: mstype.MS_BOOL, Val: left.Val.(float64) >= float64(right.Val.(int))}
	case left.IsType(&mstype.MS_FLOAT) && right.IsType(&mstype.MS_FLOAT):
		return EvalResult{Rt: mstype.MS_BOOL, Val: left.Val.(float64) >= right.Val.(float64)}
	}
	return evalErr(invalidBinop(left, right, op.String()))

}

func evalEq(left EvalResult, right EvalResult, op token.TokenType) EvalResult {
	switch {
	case left.IsType(&mstype.MS_NOTHING) || right.IsType(&mstype.MS_NOTHING):
		return EvalResult{Rt: mstype.MS_BOOL, Val: false}
	case left.IsType(&mstype.MS_STRING) && right.IsType(&mstype.MS_STRING):
		return EvalResult{Rt: mstype.MS_BOOL, Val: left.Val.(string) == right.Val.(string)}
	case left.IsType(&mstype.MS_STRING) || right.IsType(&mstype.MS_STRING):
		return EvalResult{Rt: mstype.MS_BOOL, Val: false}
	case left.IsType(&mstype.MS_BOOL) && right.IsType(&mstype.MS_BOOL):
		return EvalResult{Rt: mstype.MS_BOOL, Val: left.Val.(bool) == right.Val.(bool)}
	case left.IsType(&mstype.MS_BOOL) || right.IsType(&mstype.MS_BOOL):
		return EvalResult{Rt: mstype.MS_BOOL, Val: false}
	case left.IsType(&mstype.MS_INT) && right.IsType(&mstype.MS_INT):
		return EvalResult{Rt: mstype.MS_BOOL, Val: left.Val.(int) == right.Val.(int)}
	case left.IsType(&mstype.MS_FLOAT) && right.IsType(&mstype.MS_FLOAT):
		return EvalResult{Rt: mstype.MS_BOOL, Val: left.Val.(float64) == right.Val.(float64)}
	case left.IsType(&mstype.MS_INT) && right.IsType(&mstype.MS_FLOAT):
		return EvalResult{Rt: mstype.MS_BOOL, Val: float64(left.Val.(int)) == right.Val.(float64)}
	case left.IsType(&mstype.MS_FLOAT) && right.IsType(&mstype.MS_INT):
		return EvalResult{Rt: mstype.MS_BOOL, Val: left.Val.(float64) == float64(right.Val.(int))}
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