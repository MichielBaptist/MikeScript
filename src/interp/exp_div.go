package interp

import (
	"mikescript/src/mstype"
	"mikescript/src/utils"
)

func evalDiv(left, right EvalResult) EvalResult {

	switch lt := left.Rt.(type){
	case *mstype.MSSimpleTypeS:
		switch rt := right.Rt.(type){
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
		if right.Val.(int) == 0 {return errDiv}
		return EvalResult{Rt: mstype.MS_FLOAT, Val: float64(left.Val.(int)) / float64(right.Val.(int))}
	case left.IsType(&mstype.MS_INT) && right.IsType(&mstype.MS_FLOAT):
		if right.Val.(float64) == 0 {return errDiv}
		return EvalResult{Rt: mstype.MS_FLOAT, Val: float64(left.Val.(int)) / right.Val.(float64)}

	// Left is float
	case left.IsType(&mstype.MS_FLOAT) && right.IsType(&mstype.MS_INT):
		if right.Val.(int64) == 0 {return errDiv}
		return EvalResult{Rt: mstype.MS_FLOAT, Val: left.Val.(float64) / float64(right.Val.(int))}
	case left.IsType(&mstype.MS_FLOAT) && right.IsType(&mstype.MS_FLOAT):
		if right.Val.(float64) == 0 {return errDiv}
		return EvalResult{Rt: mstype.MS_FLOAT, Val: left.Val.(float64) / right.Val.(float64)}

	// Left is bool
	case left.IsType(&mstype.MS_BOOL) && right.IsType(&mstype.MS_INT):
		if right.Val.(int64) == 0 {return errDiv}
		return EvalResult{Rt: mstype.MS_FLOAT, Val: utils.BoolToFloat(left.Val.(bool)) / float64(right.Val.(int))}
	case left.IsType(&mstype.MS_BOOL) && right.IsType(&mstype.MS_FLOAT):
		if right.Val.(int64) == 0 {return errDiv}
		return EvalResult{Rt: mstype.MS_FLOAT, Val: utils.BoolToFloat(left.Val.(bool)) / right.Val.(float64)}
	}

	return errRes

}