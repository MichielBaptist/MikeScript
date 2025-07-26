package interp

import (
	"mikescript/src/mstype"
	"mikescript/src/utils"
)


func evalAdd(left, right EvalResult) EvalResult {

	// Based on type we need to perform add diff
	leftT := left.Rt
	rightT := right.Rt

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
		case mstype.RT_INT:		return EvalResult{Rt: mstype.MS_INT, Val: left.Val.(int) + right.Val.(int)}
		case mstype.RT_FLOAT:	return EvalResult{Rt: mstype.MS_FLOAT, Val: float64(left.Val.(int)) + right.Val.(float64)}
		case mstype.RT_BOOL:	return EvalResult{Rt: mstype.MS_INT, Val: left.Val.(int) + utils.BoolToInt(right.Val.(bool))}
		}
	case mstype.RT_FLOAT:
		switch right.Rt {
		case mstype.MS_INT:		return EvalResult{Rt: mstype.MS_FLOAT, Val: left.Val.(float64) + float64(right.Val.(int))}
		case mstype.MS_FLOAT:	return EvalResult{Rt: mstype.MS_FLOAT, Val: left.Val.(float64) + right.Val.(float64)}
		case mstype.MS_BOOL:	return EvalResult{Rt: mstype.MS_FLOAT, Val: left.Val.(float64) + float64(utils.BoolToInt(right.Val.(bool)))}
		}
	case mstype.RT_BOOL:

		// Cast left value to int
		boolint := utils.BoolToInt(left.Val.(bool))

		// Check right side of add
		switch right.Rt {
		case mstype.MS_INT:		return EvalResult{Rt: mstype.MS_INT, Val: boolint + right.Val.(int)}
		case mstype.MS_FLOAT:	return EvalResult{Rt: mstype.MS_FLOAT, Val: float64(boolint) + right.Val.(float64)}
		case mstype.MS_BOOL:	return EvalResult{Rt: mstype.MS_INT, Val: boolint + utils.BoolToInt(right.Val.(bool))}
		}
	case mstype.RT_STRING:
		switch right.Rt {
		case mstype.MS_STRING:	return EvalResult{Rt: mstype.MS_STRING, Val: left.Val.(string) + right.Val.(string)}
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