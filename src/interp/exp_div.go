package interp

import (
	"mikescript/src/mstype"
	"mikescript/src/utils"
)

func evalDiv(left, right EvalResult) EvalResult {

	switch lt := left.rt.(type){
	case *mstype.MSSimpleTypeS:
		switch rt := right.rt.(type){
		case *mstype.MSSimpleTypeS:
			return evalDivSimple(left, right, lt, rt)
		}
	}

	return evalErr(invalidBinop(left, right, "/"))
}

func evalDivSimple(left, right EvalResult, lt, rt *mstype.MSSimpleTypeS) EvalResult {

	// prepare error message
	errRes := evalErr(invalidBinop(left, right, "/"))
	errDiv := evalErr("Division by zero.")

	// Div not defined for Nothing ot String types
	switch {
	case left.IsType(&mstype.MS_NOTHING) || right.IsType(&mstype.MS_NOTHING):	return errRes
	case left.IsType(&mstype.MS_STRING) || right.IsType(&mstype.MS_STRING):		return errRes
	}

	switch {

	// Left is int
	case left.IsType(&mstype.MS_INT) && right.IsType(&mstype.MS_INT):
		if right.val.(int) == 0 {return errDiv}
		return EvalResult{rt: mstype.MS_FLOAT, val: float64(left.val.(int)) / float64(right.val.(int))}
	case left.IsType(&mstype.MS_INT) && right.IsType(&mstype.MS_FLOAT):
		if right.val.(float64) == 0 {return errDiv}
		return EvalResult{rt: mstype.MS_FLOAT, val: float64(left.val.(int)) / right.val.(float64)}

	// Left is float
	case left.IsType(&mstype.MS_FLOAT) && right.IsType(&mstype.MS_INT):
		if right.val.(int64) == 0 {return errDiv}
		return EvalResult{rt: mstype.MS_FLOAT, val: left.val.(float64) / float64(right.val.(int))}
	case left.IsType(&mstype.MS_FLOAT) && right.IsType(&mstype.MS_FLOAT):
		if right.val.(float64) == 0 {return errDiv}
		return EvalResult{rt: mstype.MS_FLOAT, val: left.val.(float64) / right.val.(float64)}

	// Left is bool
	case left.IsType(&mstype.MS_BOOL) && right.IsType(&mstype.MS_INT):
		if right.val.(int64) == 0 {return errDiv}
		return EvalResult{rt: mstype.MS_FLOAT, val: utils.BoolToFloat(left.val.(bool)) / float64(right.val.(int))}
	case left.IsType(&mstype.MS_BOOL) && right.IsType(&mstype.MS_FLOAT):
		if right.val.(int64) == 0 {return errDiv}
		return EvalResult{rt: mstype.MS_FLOAT, val: utils.BoolToFloat(left.val.(bool)) / right.val.(float64)}
	}

	return errRes

}