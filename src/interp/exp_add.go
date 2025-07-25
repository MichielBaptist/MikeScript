package interp

import (
	"mikescript/src/mstype"
	"mikescript/src/utils"
)


func evalAdd(left, right EvalResult) EvalResult {

	// Based on type we need to perform add diff
	leftT := left.rt
	rightT := right.rt

	switch l := leftT.(type){
	case *mstype.MSSimpleTypeS:
		switch r := rightT.(type) {
		case *mstype.MSSimpleTypeS:

			// Simple addition
			return evalAddSimple(left, l, right, r)
		}
	}

	return evalErr(invalidBinop(left, right, "+"))
}

func evalAddSimple(left EvalResult, lT *mstype.MSSimpleTypeS, right EvalResult, rT *mstype.MSSimpleTypeS, ) EvalResult {

	// Handle all numeric/boolean cases
	switch lT.Rt {
	case mstype.RT_INT:
		switch rT.Rt {
		case mstype.RT_INT:		return EvalResult{rt: mstype.MS_INT, val: left.val.(int) + right.val.(int)}
		case mstype.RT_FLOAT:	return EvalResult{rt: mstype.MS_FLOAT, val: float64(left.val.(int)) + right.val.(float64)}
		case mstype.RT_BOOL:	return EvalResult{rt: mstype.MS_INT, val: left.val.(int) + utils.BoolToInt(right.val.(bool))}
		}
	case mstype.RT_FLOAT:
		switch right.rt {
		case mstype.MS_INT:		return EvalResult{rt: mstype.MS_FLOAT, val: left.val.(float64) + float64(right.val.(int))}
		case mstype.MS_FLOAT:	return EvalResult{rt: mstype.MS_FLOAT, val: left.val.(float64) + right.val.(float64)}
		case mstype.MS_BOOL:	return EvalResult{rt: mstype.MS_FLOAT, val: left.val.(float64) + float64(utils.BoolToInt(right.val.(bool)))}
		}
	case mstype.RT_BOOL:

		// Cast left value to int
		boolint := utils.BoolToInt(left.val.(bool))

		// Check right side of add
		switch right.rt {
		case mstype.MS_INT:		return EvalResult{rt: mstype.MS_INT, val: boolint + right.val.(int)}
		case mstype.MS_FLOAT:	return EvalResult{rt: mstype.MS_FLOAT, val: float64(boolint) + right.val.(float64)}
		case mstype.MS_BOOL:	return EvalResult{rt: mstype.MS_INT, val: boolint + utils.BoolToInt(right.val.(bool))}
		}
	case mstype.RT_STRING:
		switch right.rt {
		case mstype.MS_STRING:	return EvalResult{rt: mstype.MS_STRING, val: left.val.(string) + right.val.(string)}
		}
	}

	return evalErr(invalidBinop(left, right, "+"))
}

func evalSub(_, _ EvalResult) EvalResult {
	// This code should never happen as there is no
	// way in which the parser can generate a subtraction
	// instead it will generate an addition with a negative unary
	return evalErr("Subtraction should be replaced with addition and unary negation. If you see this error, the parser is broken.")
}