package interp

import (
	"mikescript/src/token"
	"mikescript/src/utils"
)

const errNeq string = "Not 'exp1 != exp2' be replaced with '!(exp1 == exp2)'. If you see this error, the parser is broken."
const errGrGr string = "Trying to evaluate function call as binary operator. If you see this error, the parser is broken."

func compareGreater(l, r float64) bool {
	return l > r
}

func compareGreaterEq(l, r float64) bool {
	return l >= r
}

func cvtFloat(val MSVal) (float64, bool) {
	// Converts a ms val to float

	switch t := val.(type){
	case MSInt:		return float64(t.Val), true
	case MSFloat:	return t.Val, true
	case MSBool:	return utils.BoolToFloat(t.Val), true
	}

	return 0.0, false
}

func evalCompare(lval, rval MSVal, op token.TokenType, fn func(float64, float64) bool) (MSVal, error) {
	// on comparison of 2 basic values, we compare their floats.

	lf, lok := cvtFloat(lval)
	rf, rok := cvtFloat(rval)

	if !(lok && rok) {
		// Can't convert to float
		msg := invalidBinop(lval, rval, op.String())
		return MSNothing{}, &EvalError{message: msg}
	}

	return MSBool{Val: fn(lf, rf)}, nil
}


func evalGreater(lval, rval MSVal, op token.TokenType) (MSVal, error) {
	return evalCompare(lval, rval, op, compareGreater)
}

func evalGreaterEq(lval, rval MSVal, op token.TokenType) (MSVal, error) {
	return evalCompare(lval, rval, op, compareGreaterEq)
}

func evalEq(lval, rval MSVal, op token.TokenType) (MSVal, error) {

	var err error

	switch l := lval.(type){
	case MSNothing:
		switch rval.(type) {
		case MSNothing:	return MSBool{Val: true}, err
		default: 		return MSBool{Val: false}, err
		}
	case MSInt:
		switch r := rval.(type){
		case MSInt:		return MSBool{Val: l.Val == r.Val}, err
		case MSFloat:	return MSBool{Val: float64(l.Val) == r.Val}, err
		case MSBool:	return MSBool{Val: l.Val == utils.BoolToInt(r.Val)}, err
		default:		return MSBool{Val: false}, err
		}
	case MSFloat:
		switch r := rval.(type){
		case MSInt:		return MSBool{Val: l.Val == float64(r.Val)}, err
		case MSFloat:	return MSBool{Val: l.Val == r.Val}, err
		case MSBool:	return MSBool{Val: l.Val == utils.BoolToFloat(r.Val)}, err
		default:		return MSBool{Val: false}, err
		}
	case MSString:
		switch r := rval.(type){
		case MSString:	return MSBool{Val: l.Val == r.Val}, err
		default:		return MSBool{Val: false}, err
		}
	case MSBool:
		switch r := rval.(type){
		case MSInt:		return MSBool{Val: utils.BoolToInt(l.Val) == r.Val}, err
		case MSFloat:	return MSBool{Val: utils.BoolToFloat(l.Val) == r.Val}, err
		case MSBool:	return MSBool{Val: l.Val == r.Val}, err
		default:		return MSBool{Val: false}, err
		}
	}

	return MSNothing{}, &EvalError{invalidBinop(lval, rval, "==")}
}

func evalNeq(_, _ MSVal, _ token.TokenType) (MSVal, error) {
	return MSNothing{}, &EvalError{message: errNeq}
}

func evalGrGr(_, _ MSVal) (MSVal, error) {
	return MSNothing{}, &EvalError{errGrGr}
}