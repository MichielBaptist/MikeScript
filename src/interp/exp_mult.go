package interp

import (
	"mikescript/src/mstype"
	"mikescript/src/utils"
)

func evalMult(left, right EvalResult) EvalResult {

	// Based on type we need to perform add diff
	leftT := left.Rt
	rightT := right.Rt

	switch l := leftT.(type){
	case *mstype.MSSimpleTypeS:
		switch r := rightT.(type) {
		case *mstype.MSSimpleTypeS:
			return evalMultSimple(left, right, l, r)
		}
	}

	return evalErr(invalidBinop(left, right, "*"))
}

func evalMultSimple(left, right EvalResult, lt, rt *mstype.MSSimpleTypeS) EvalResult{

	switch lt.Rt {
	case mstype.RT_STRING:
		switch rt.Rt {
		case mstype.RT_INT:		return EvalResult{Rt: mstype.MS_STRING, Val: utils.RepeatString(left.Val.(string), right.Val.(int))}
		}
	case mstype.RT_INT:
		switch rt.Rt {
		case mstype.RT_INT:		return EvalResult{Rt: mstype.MS_INT, Val: left.Val.(int) * right.Val.(int)}
		case mstype.RT_FLOAT:	return EvalResult{Rt: mstype.MS_FLOAT, Val: float64(left.Val.(int)) * right.Val.(float64)}
		case mstype.RT_BOOL:	return EvalResult{Rt: mstype.MS_INT, Val: left.Val.(int) * utils.BoolToInt(right.Val.(bool))}
		case mstype.RT_STRING: 	return EvalResult{Rt: mstype.MS_STRING, Val: utils.RepeatString(right.Val.(string), left.Val.(int))}
		}
	case mstype.RT_FLOAT:
		switch right.Rt {
		case mstype.MS_INT:		return EvalResult{Rt: mstype.MS_FLOAT, Val: left.Val.(float64) * float64(right.Val.(int))}
		case mstype.MS_FLOAT:	return EvalResult{Rt: mstype.MS_FLOAT, Val: left.Val.(float64) * right.Val.(float64)}
		case mstype.MS_BOOL:	return EvalResult{Rt: mstype.MS_FLOAT, Val: left.Val.(float64) * float64(utils.BoolToInt(right.Val.(bool)))}
		}
	case mstype.RT_BOOL:
		switch rt.Rt {
		case mstype.RT_INT:		return EvalResult{Rt: mstype.MS_INT, Val: utils.BoolToInt(left.Val.(bool)) * right.Val.(int)}
		case mstype.RT_FLOAT:	return EvalResult{Rt: mstype.MS_FLOAT, Val: float64(utils.BoolToInt(left.Val.(bool))) * right.Val.(float64)}
		case mstype.RT_BOOL:	return EvalResult{Rt: mstype.MS_INT, Val: utils.BoolToInt(left.Val.(bool)) * utils.BoolToInt(right.Val.(bool))}
		}
	}
	return evalErr(invalidBinop(left, right, "*"))
}