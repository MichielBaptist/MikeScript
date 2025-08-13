package interp

import (
	"mikescript/src/utils"
)


func evalAdd(lval, rval MSVal) (MSVal, error) {
	// we can assume that left and right are valid.

	var err error

	switch l := lval.(type){
	case MSInt:
		switch r := rval.(type){
		case MSInt:		return MSInt{Val: l.Val + r.Val}, err
		case MSFloat:	return MSFloat{Val: float64(l.Val) + r.Val}, err
		case MSBool:	return MSInt{Val: l.Val + utils.BoolToInt(r.Val)}, err
		}
	case MSFloat:
		switch r := rval.(type){
		case MSInt:		return MSFloat{Val: l.Val + float64(r.Val)}, err
		case MSFloat:	return MSFloat{Val: l.Val + r.Val}, err
		case MSBool:	return MSFloat{Val: l.Val + utils.BoolToFloat(r.Val)}, err
		}
	case MSBool:
		// cast to int
		lint := utils.BoolToInt(l.Val)

		switch r := rval.(type){
		case MSInt:		return MSInt{Val: lint + r.Val}, err
		case MSFloat:	return MSFloat{Val: float64(lint) + r.Val}, err
		case MSBool:	return MSInt{Val: lint + utils.BoolToInt(r.Val)}, err
		}
	case MSString:
		switch r := rval.(type){
		case MSString:	return MSString{Val: l.Val + r.Val}, err
		}
	}

	return nil, &EvalError{invalidBinop(lval, rval, "+")}
}

func evalSub(_, _ MSVal) (MSVal, error) {
	// This code should never happen as there is no
	// way in which the parser can generate a subtraction
	// instead it will generate an addition with a negative unary
	msg := "Subtraction should be replaced with addition and unary negation. If you see this error, the parser is broken."
	return nil, &EvalError{message: msg}
}