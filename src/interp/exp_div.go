package interp

import (
	"mikescript/src/utils"
)

func evalDiv(lval, rval MSVal) (MSVal, error) {

	var err error
	var num float64
	var den float64
	numv := true
	denv := true

	switch l := lval.(type){
	case MSInt:		num = float64(l.Val)
	case MSFloat:	num = l.Val
	case MSBool:	num = utils.BoolToFloat(l.Val)
	default:		numv = false
	}

	switch r := rval.(type){
	case MSInt:		den = float64(r.Val)
	case MSFloat:	den = r.Val
	case MSBool:	den = utils.BoolToFloat(r.Val)
	default:		denv = false
	}

	if !(numv && denv) {
		// Type issue
		return MSNothing{}, &EvalError{invalidBinop(lval, rval, "/")}
	} 

	if den == 0.0 {
		return MSNothing{}, &EvalError{"Division by zero."}
	}

	return MSFloat{Val: num / den}, err

}