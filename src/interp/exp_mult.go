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
		case MSFloat:		return MSFloat{Val: float64(l.Val) * r.Val}, err
		case MSBool:		return MSInt{Val: l.Val * utils.BoolToInt(r.Val)}, err
		}
	case MSFloat:
		switch r := rval.(type){
		case MSInt:		return MSFloat{Val: l.Val * float64(r.Val)}, err
		case MSFloat:		return MSFloat{Val: l.Val * r.Val}, err
		case MSBool:		return MSFloat{Val: l.Val * utils.BoolToFloat(r.Val)}, err
		}
	case MSBool:
		// cast to int
		lint := utils.BoolToInt(l.Val)

		switch r := rval.(type){
		case MSInt:		return MSInt{Val: lint * r.Val}, err
		case MSFloat:		return MSFloat{Val: float64(lint) * r.Val}, err
		case MSBool:		return MSInt{Val: lint * utils.BoolToInt(r.Val)}, err
		}
	case MSString:
		switch r := rval.(type){
		case MSInt:	return MSString{Val: strings.Repeat(l.Val, r.Val)}, err
		}
	}

	return MSNothing{}, &EvalError{invalidBinop(lval, rval, "*")}
}

// func evalMultSimple(left, right EvalResult, lt, rt *mstype.MSSimpleTypeS) EvalResult{

// 	switch lt.Rt {
// 	case mstype.RT_STRING:
// 		switch rt.Rt {
// 		case mstype.RT_INT:		return EvalResult{Rt: mstype.MS_STRING, Val: utils.RepeatString(left.Val.(string), right.Val.(int))}
// 		}
// 	case mstype.RT_INT:
// 		switch rt.Rt {
// 		case mstype.RT_INT:		return EvalResult{Rt: mstype.MS_INT, Val: left.Val.(int) * right.Val.(int)}
// 		case mstype.RT_FLOAT:	return EvalResult{Rt: mstype.MS_FLOAT, Val: float64(left.Val.(int)) * right.Val.(float64)}
// 		case mstype.RT_BOOL:	return EvalResult{Rt: mstype.MS_INT, Val: left.Val.(int) * utils.BoolToInt(right.Val.(bool))}
// 		case mstype.RT_STRING: 	return EvalResult{Rt: mstype.MS_STRING, Val: utils.RepeatString(right.Val.(string), left.Val.(int))}
// 		}
// 	case mstype.RT_FLOAT:
// 		switch right.Rt {
// 		case mstype.MS_INT:		return EvalResult{Rt: mstype.MS_FLOAT, Val: left.Val.(float64) * float64(right.Val.(int))}
// 		case mstype.MS_FLOAT:	return EvalResult{Rt: mstype.MS_FLOAT, Val: left.Val.(float64) * right.Val.(float64)}
// 		case mstype.MS_BOOL:	return EvalResult{Rt: mstype.MS_FLOAT, Val: left.Val.(float64) * float64(utils.BoolToInt(right.Val.(bool)))}
// 		}
// 	case mstype.RT_BOOL:
// 		switch rt.Rt {
// 		case mstype.RT_INT:		return EvalResult{Rt: mstype.MS_INT, Val: utils.BoolToInt(left.Val.(bool)) * right.Val.(int)}
// 		case mstype.RT_FLOAT:	return EvalResult{Rt: mstype.MS_FLOAT, Val: float64(utils.BoolToInt(left.Val.(bool))) * right.Val.(float64)}
// 		case mstype.RT_BOOL:	return EvalResult{Rt: mstype.MS_INT, Val: utils.BoolToInt(left.Val.(bool)) * utils.BoolToInt(right.Val.(bool))}
// 		}
// 	}
// 	return evalErr(invalidBinop(left, right, "*"))
// }
