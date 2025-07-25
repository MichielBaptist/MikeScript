package interp

import (
	"mikescript/src/mstype"
	"mikescript/src/utils"
)

func evalMult(left, right EvalResult) EvalResult {

	// Based on type we need to perform add diff
	leftT := left.rt
	rightT := right.rt

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
		case mstype.RT_INT:		return EvalResult{rt: mstype.MS_STRING, val: utils.RepeatString(left.val.(string), right.val.(int))}
		}
	case mstype.RT_INT:
		switch rt.Rt {
		case mstype.RT_INT:		return EvalResult{rt: mstype.MS_INT, val: left.val.(int) * right.val.(int)}
		case mstype.RT_FLOAT:	return EvalResult{rt: mstype.MS_FLOAT, val: float64(left.val.(int)) * right.val.(float64)}
		case mstype.RT_BOOL:	return EvalResult{rt: mstype.MS_INT, val: left.val.(int) * utils.BoolToInt(right.val.(bool))}
		case mstype.RT_STRING: 	return EvalResult{rt: mstype.MS_STRING, val: utils.RepeatString(right.val.(string), left.val.(int))}
		}
	case mstype.RT_FLOAT:
		switch right.rt {
		case mstype.MS_INT:		return EvalResult{rt: mstype.MS_FLOAT, val: left.val.(float64) * float64(right.val.(int))}
		case mstype.MS_FLOAT:	return EvalResult{rt: mstype.MS_FLOAT, val: left.val.(float64) * right.val.(float64)}
		case mstype.MS_BOOL:	return EvalResult{rt: mstype.MS_FLOAT, val: left.val.(float64) * float64(utils.BoolToInt(right.val.(bool)))}
		}
	case mstype.RT_BOOL:
		switch rt.Rt {
		case mstype.RT_INT:		return EvalResult{rt: mstype.MS_INT, val: utils.BoolToInt(left.val.(bool)) * right.val.(int)}
		case mstype.RT_FLOAT:	return EvalResult{rt: mstype.MS_FLOAT, val: float64(utils.BoolToInt(left.val.(bool))) * right.val.(float64)}
		case mstype.RT_BOOL:	return EvalResult{rt: mstype.MS_INT, val: utils.BoolToInt(left.val.(bool)) * utils.BoolToInt(right.val.(bool))}
		}
	}
	return evalErr(invalidBinop(left, right, "*"))
}