package interp

import (
	"mikescript/src/utils"
	"strings"
)

func evalMult(lval, rval MSVal) (MSVal, error) {

	var err error

	switch l := lval.(type){
	case MSInt:
		switch r := rval.(type){
		case MSInt:		return MSInt{Val: l.Val * r.Val}, err
		case MSFloat:	return MSFloat{Val: float64(l.Val) * r.Val}, err
		case MSBool:	return MSInt{Val: l.Val * utils.BoolToInt(r.Val)}, err
		}
	case MSFloat:
		switch r := rval.(type){
		case MSInt:		return MSFloat{Val: l.Val * float64(r.Val)}, err
		case MSFloat:	return MSFloat{Val: l.Val * r.Val}, err
		case MSBool:	return MSFloat{Val: l.Val * utils.BoolToFloat(r.Val)}, err
		}
	case MSBool:
		// cast to int
		lint := utils.BoolToInt(l.Val)

		switch r := rval.(type){
		case MSInt:		return MSInt{Val: lint * r.Val}, err
		case MSFloat:	return MSFloat{Val: float64(lint) * r.Val}, err
		case MSBool:	return MSInt{Val: lint * utils.BoolToInt(r.Val)}, err
		}
	case MSString:
		switch r := rval.(type){
		case MSInt:		return MSString{Val: strings.Repeat(l.Val, r.Val)}, err
		}
	}

	return nil, &EvalError{invalidBinop(lval, rval, "*")}
}
