package interp

import (
	"fmt"
	"mikescript/src/ast"
	"mikescript/src/token"
	"mikescript/src/utils"
	"strconv"
)

////////////////////////////////////////////////////////////////////////
// Evaluate expressions
////////////////////////////////////////////////////////////////////////

func (evaluator *MSEvaluator) evaluateExpression(node *ast.ExpNodeI) EvalResult {
	switch node := (*node).(type) {
	case ast.BinaryExpNodeS:	return evaluator.evaluateBinaryExpression(&node)
	case ast.UnaryExpNodeS:		return evaluator.evaluateUnaryExpression(&node)
	case ast.LiteralExpNodeS:	return evaluator.evaluateLiteralExpression(&node)
	case ast.GroupExpNodeS:		return evaluator.evaluateGroupExpression(&node)
	case ast.AssignmentNodeS:	return evaluator.evaluateAssignmentExpression(&node)
	case ast.FuncAppNodeS:		return evaluator.evaluateFunctionApplication(&node)
	case ast.VariableExpNodeS:	return evaluator.evalVariable(&node)
	case ast.LogicalExpNodeS:	return evaluator.evaluateLogicalExpression(&node)
	default:					return evalErr(fmt.Sprintf("Unknown expression type: %v", node))
	}
}

func (evaluator *MSEvaluator) evaluateBinaryExpression(node *ast.BinaryExpNodeS) EvalResult {

	// Evaluate left and right side of the expression
	left := evaluator.evaluateExpression(&node.Left)
	right := evaluator.evaluateExpression(&node.Right)

	// If either of the evaluations failed, return the error
	if !left.Valid() || !right.Valid() {
		return EvalResult{err: append(left.err, right.err...)}
	}

	switch node.Op.Type {
	case token.PLUS: 				return evalAdd(left, right)
	case token.MINUS:				return evalSub(left, right)
	case token.MULT:				return evalMult(left, right)
	case token.SLASH:				return evalDiv(left, right)
	case token.GREATER:				return evalGreater(left, right, node.Op.Type)
	case token.LESS:				return evalGreater(right, left, node.Op.Type)
	case token.GREATER_EQ:			return evalGreaterEq(left, right, node.Op.Type)
	case token.LESS_EQ:				return evalGreaterEq(right, left, node.Op.Type)
	case token.EQ_EQ:				return evalEq(left, right, node.Op.Type)
	case token.EXCLAMATION_EQ:		return evalNeq(left, right, node.Op.Type)
	case token.GREATER_GREATER: 	return evalGrGr(left, right)
	case token.COMMA:				return evalTuple(left, right)
	case token.PERCENT:				return evalMod(left, right)
	default:						return evalErr(fmt.Sprintf("Unknown binary operator: %v", node.Op.Lexeme))
	}
}

func (evaluator *MSEvaluator) evaluateUnaryExpression(node *ast.UnaryExpNodeS) EvalResult {
	// evaluate the node
	res := evaluator.evaluateExpression(&node.Node)

	// check if the evaluation was.Valid()
	if !res.Valid() {
		return res
	}

	// handle unary operators
	switch node.Op.Type {
	case token.MINUS:
		switch res.rt {
		case RT_INT:	return EvalResult{rt: RT_INT, 	val: -res.val.(int)}
		case RT_FLOAT:	return EvalResult{rt: RT_FLOAT, val: -res.val.(float64)}
		}
	case token.EXCLAMATION:
		switch res.rt {
		case RT_BOOL:	return EvalResult{rt: RT_BOOL, val: !res.val.(bool)}
		}
	}
	return evalErr(fmt.Sprintf("Operator %v is not defined for type %v", node.Op.Lexeme, res.rt))
}

func (evaluator *MSEvaluator) evaluateGroupExpression(node *ast.GroupExpNodeS) EvalResult {
	return evaluator.evaluateExpression(&node.Node)
}

func (evaluator *MSEvaluator) evaluateLiteralExpression(node *ast.LiteralExpNodeS) EvalResult {
	switch node.Tk.Type {
	case token.NUMBER_INT:		return evalIntLiteral(node)
	case token.NUMBER_FLOAT:	return evalFloatLiteral(node)
	case token.STRING:			return evalStringLiteral(node)
	case token.TRUE:			return EvalResult{rt: RT_BOOL, val: true}
	case token.FALSE:			return EvalResult{rt: RT_BOOL, val: false}
	case token.IDENTIFIER:		return evalErr(fmt.Sprintf("Trying to evaluate identifier '%v' as a literal.", node.Tk.Lexeme))
	default:					return evalErr(fmt.Sprintf("Literal type '%v' is not defined.", node.Tk.Type))
	}
}

func (evaluator *MSEvaluator) evaluateAssignmentExpression(node *ast.AssignmentNodeS) EvalResult {
	// Evaluate the expression
	res := evaluator.evaluateExpression(&node.Exp)

	// set the variable in the environment
	err := evaluator.env.SetVar(node.Identifier.Name.Lexeme, res)

	if err != nil {
		return EvalResult{err: []error{err}}
	}

	return EvalResult{rt: res.rt, val: res.val}
}

func (evaluator *MSEvaluator) evalVariable(node *ast.VariableExpNodeS) EvalResult {

	// Get the value from the environment
	val, err := evaluator.env.GetVar(node.Name.Lexeme)

	if err != nil {
		return EvalResult{err: []error{err}}
	}

	return val
	
}

func (evaluator *MSEvaluator) evaluateLogicalExpression(node *ast.LogicalExpNodeS) EvalResult {

	// Evaluate the left and right side of the expression
	left := evaluator.evaluateExpression(&node.Left)

	// If either of the evaluations failed, return the error
	if !left.Valid() {
		return left
	}

	// Check if the left side is a boolean
	if left.rt != RT_BOOL {
		return evalErr(invalidLogop(left, node.Op.Lexeme))
	}

	// Check the truth value of the left side
	leftb := left.val.(bool)

	// short circuit evaluation
	switch  {
	case node.Op.Type == token.AMP_AMP && !leftb: return EvalResult{rt: RT_BOOL, val: false}
	case node.Op.Type == token.BAR_BAR &&  leftb: return EvalResult{rt: RT_BOOL, val: true}
	}

	// Means the first operand is inconclusive
	right := evaluator.evaluateExpression(&node.Right)

	// Check if the right side is a boolean
	if right.rt != RT_BOOL {
		return evalErr(invalidLogop(right, node.Op.Lexeme))
	}

	return right

}

func (evaluator *MSEvaluator) evaluateFunctionApplication(node *ast.FuncAppNodeS) EvalResult {
	// check if the function is 'print'
	// if it is, print the arguments

	switch fun := node.Fun.(type) {
	case ast.VariableExpNodeS:
		// Get the function name
		if fun.Name.Lexeme == "print" || fun.Name.Lexeme == "println" {
			// Evaluate the arguments
			args := make([]EvalResult, len(node.Args))
			for i, arg := range node.Args {
				args[i] = evaluator.evaluateExpression(&arg)
			}

			// Accumulate all errors into one
			errs := []error{}
			for _, arg := range args {
				errs = append(errs, arg.err...)
			}
			if len(errs) > 0 {
				return EvalResult{err: errs}
			}

			// Get all result values in a slice
			// and convert them to strings
			strs := make([]string, len(args))
			for i, arg := range args {
				strs[i] = fmt.Sprint(arg.val)
			}

			// Print the joined strings
			fmt.Println(utils.StrJoin(strs, ", "))
			return EvalResult{rt: RT_NONE}
		}
	}
	return evalErr("Function application is not implemented yet.")
}

////////////////////////////////////////////////////////////////////////
// Evaluation literals
////////////////////////////////////////////////////////////////////////

func evalIntLiteral(node *ast.LiteralExpNodeS) EvalResult {
	// convert the lexeme to an int
	val, err := strconv.Atoi(node.Tk.Lexeme)

	if err != nil {
		return evalErr(fmt.Sprintf("Could not convert '%v' to int.", node.Tk.Lexeme))
	}

	return EvalResult{rt: RT_INT, val: val}
}

func evalFloatLiteral(node *ast.LiteralExpNodeS) EvalResult {
	// convert the lexeme to a float
	val, err := strconv.ParseFloat(node.Tk.Lexeme, 64)

	if err != nil {
		return evalErr(fmt.Sprintf("Could not convert '%v' to float64.", node.Tk.Lexeme))
	}

	return EvalResult{rt: RT_FLOAT, val: val}
}

func evalStringLiteral(node *ast.LiteralExpNodeS) EvalResult {
	return EvalResult{rt: RT_STRING, val: node.Tk.Lexeme}
}

////////////////////////////////////////////////////////////////////////
// Evaluation Binary operations                 					  //
////////////////////////////////////////////////////////////////////////

func invalidBinop(left, right EvalResult, op string) string {
	return fmt.Sprintf("Operator '%v' is not defined for types '%v' and '%v'", op, left.rt, right.rt)
}

func invalidLogop(left EvalResult, op string) string {
	return fmt.Sprintf("Logical operator '%v' is not defined for type '%v'", op, left.rt)
}

func evalAdd(left, right EvalResult) EvalResult {

	// Prepare error message
	errRes := evalErr(invalidBinop(left, right, "+"))

	// Handle all numeric/boolean cases
	switch left.rt {
	case RT_INT:
		switch right.rt {
		case RT_INT:	return EvalResult{rt: RT_INT, val: left.val.(int) + right.val.(int)}
		case RT_FLOAT:	return EvalResult{rt: RT_FLOAT, val: float64(left.val.(int)) + right.val.(float64)}
		case RT_BOOL:	return EvalResult{rt: RT_INT, val: left.val.(int) + utils.BoolToInt(right.val.(bool))}
		}
	case RT_FLOAT:
		switch right.rt {
		case RT_INT:	return EvalResult{rt: RT_FLOAT, val: left.val.(float64) + float64(right.val.(int))}
		case RT_FLOAT:	return EvalResult{rt: RT_FLOAT, val: left.val.(float64) + right.val.(float64)}
		case RT_BOOL:	return EvalResult{rt: RT_FLOAT, val: left.val.(float64) + float64(utils.BoolToInt(right.val.(bool)))}
		}
	case RT_BOOL:
		// Cast left value to int
		boolint := utils.BoolToInt(left.val.(bool))
		switch right.rt {
		case RT_INT:	return EvalResult{rt: RT_INT, val: boolint + right.val.(int)}
		case RT_FLOAT:	return EvalResult{rt: RT_FLOAT, val: float64(boolint) + right.val.(float64)}
		case RT_BOOL:	return EvalResult{rt: RT_INT, val: boolint + utils.BoolToInt(right.val.(bool))}
		}
	case RT_STRING:
		switch right.rt {
		case RT_STRING:	return EvalResult{rt: RT_STRING, val: left.val.(string) + right.val.(string)}
		}
	case RT_TUPLE:
		switch right.rt{
		case RT_TUPLE:	return EvalResult{rt: RT_TUPLE, val: append(left.val.([]EvalResult), right.val.([]EvalResult)...)}
		}
	}
	return errRes
}

func evalSub(_, _ EvalResult) EvalResult {
	// This code should never happen as there is no
	// way in which the parser can generate a subtraction
	// instead it will generate an addition with a negative unary
	return evalErr("Subtraction should be replaced with addition and unary negation. If you see this error, the parser is broken.")
}

func evalMult(left, right EvalResult) EvalResult {

	switch left.rt {
	case RT_STRING:
		switch right.rt {
		case RT_INT:	return EvalResult{rt: RT_STRING, val: utils.RepeatString(left.val.(string), right.val.(int))}
		}
	case RT_INT:
		switch right.rt {
		case RT_INT:	return EvalResult{rt: RT_INT, val: left.val.(int) * right.val.(int)}
		case RT_FLOAT:	return EvalResult{rt: RT_FLOAT, val: float64(left.val.(int)) * right.val.(float64)}
		case RT_BOOL:	return EvalResult{rt: RT_INT, val: left.val.(int) * utils.BoolToInt(right.val.(bool))}
		case RT_STRING: return EvalResult{rt: RT_STRING, val: utils.RepeatString(right.val.(string), left.val.(int))}
		}
	case RT_FLOAT:
		switch right.rt {
		case RT_INT:	return EvalResult{rt: RT_FLOAT, val: left.val.(float64) * float64(right.val.(int))}
		case RT_FLOAT:	return EvalResult{rt: RT_FLOAT, val: left.val.(float64) * right.val.(float64)}
		case RT_BOOL:	return EvalResult{rt: RT_FLOAT, val: left.val.(float64) * float64(utils.BoolToInt(right.val.(bool)))}
		}
	case RT_BOOL:
		switch right.rt {
		case RT_INT:	return EvalResult{rt: RT_INT, val: utils.BoolToInt(left.val.(bool)) * right.val.(int)}
		case RT_FLOAT:	return EvalResult{rt: RT_FLOAT, val: float64(utils.BoolToInt(left.val.(bool))) * right.val.(float64)}
		case RT_BOOL:	return EvalResult{rt: RT_INT, val: utils.BoolToInt(left.val.(bool)) * utils.BoolToInt(right.val.(bool))}
		}
	}
	return evalErr(invalidBinop(left, right, "*"))
}

func evalDiv(left, right EvalResult) EvalResult {

	// prepare error message
	errRes := evalErr(invalidBinop(left, right, "/"))
	errDiv := evalErr("Division by zero.")

	// Div with dynamic type checking
	switch {
	case left.rt == RT_NONE || right.rt == RT_NONE:
		return errRes
	case left.rt == RT_STRING || right.rt == RT_STRING:
		return errRes
	}

	// Handle all numeric/boolean cases
	switch left.rt {
	case RT_INT:
		switch right.rt {
		case RT_INT:
			if right.val.(int) == 0 {
				return errDiv
			}
			return EvalResult{rt: RT_FLOAT, val: float64(left.val.(int)) / float64(right.val.(int))}
		case RT_FLOAT:
			if right.val.(float64) == 0 {
				return errDiv
			}
			return EvalResult{rt: RT_FLOAT, val: float64(left.val.(int)) / right.val.(float64)}
		case RT_BOOL:
			if utils.BoolToInt(right.val.(bool)) == 0 {
				return errDiv
			}
			return EvalResult{rt: RT_FLOAT, val: left.val.(int) / utils.BoolToInt(right.val.(bool))}
		default:
			return evalErr(fmt.Sprintf("Cannot divide '%v' with type '%v'", left.rt, right.rt))
		}
	case RT_FLOAT:
		switch right.rt {
		case RT_INT:
			if right.val.(int) == 0 {
				return errDiv
			}
			return EvalResult{rt: RT_FLOAT, val: left.val.(float64) / float64(right.val.(int))}
		case RT_FLOAT:
			if right.val.(float64) == 0 {
				return errDiv
			}
			return EvalResult{rt: RT_FLOAT, val: left.val.(float64) / right.val.(float64)}
		case RT_BOOL:
			if utils.BoolToInt(right.val.(bool)) == 0 {
				return errDiv
			}
			return EvalResult{rt: RT_FLOAT, val: left.val.(float64) / float64(utils.BoolToInt(right.val.(bool)))}
		}
	case RT_BOOL:
		switch right.rt {
		case RT_INT:
			if right.val.(int) == 0 {
				return evalErr("Division by zero.")
			}
			return EvalResult{rt: RT_FLOAT, val: utils.BoolToInt(left.val.(bool)) / right.val.(int)}
		case RT_FLOAT:
			if right.val.(float64) == 0 {
				return evalErr("Division by zero.")
			}
			return EvalResult{rt: RT_FLOAT, val: float64(utils.BoolToInt(left.val.(bool))) / right.val.(float64)}
		case RT_BOOL:
			if utils.BoolToInt(right.val.(bool)) == 0 {
				return evalErr("Division by zero.")
			}
			return EvalResult{rt: RT_FLOAT, val: utils.BoolToInt(left.val.(bool)) / utils.BoolToInt(right.val.(bool))}
		default:
			return evalErr(invalidBinop(left, right, "/"))
		}
	}
	return evalErr(invalidBinop(left, right, "/"))
}

func evalGreater(left, right EvalResult, op token.TokenType) EvalResult {

	switch left.rt {
	case RT_INT:
		switch right.rt {
		case RT_INT:	return EvalResult{rt: RT_BOOL, val: left.val.(int) > right.val.(int)}
		case RT_FLOAT:	return EvalResult{rt: RT_BOOL, val: float64(left.val.(int)) > right.val.(float64)}
		}
	case RT_FLOAT:
		switch right.rt {
		case RT_INT:	return EvalResult{rt: RT_BOOL, val: left.val.(float64) > float64(right.val.(int))}
		case RT_FLOAT:	return EvalResult{rt: RT_BOOL, val: left.val.(float64) > right.val.(float64)}
		}
	}
	return evalErr(invalidBinop(left, right, op.String()))
}

func evalGreaterEq(left, right EvalResult, op token.TokenType) EvalResult {

	switch left.rt {
	case RT_INT:
		switch right.rt {
		case RT_INT:	return EvalResult{rt: RT_BOOL, val: left.val.(int) >= right.val.(int)}
		case RT_FLOAT:	return EvalResult{rt: RT_BOOL, val: float64(left.val.(int)) >= right.val.(float64)}
		}
	case RT_FLOAT:
		switch right.rt {
		case RT_INT:	return EvalResult{rt: RT_BOOL, val: left.val.(float64) >= float64(right.val.(int))}
		case RT_FLOAT:	return EvalResult{rt: RT_BOOL, val: left.val.(float64) >= right.val.(float64)}
		}
	}
	return evalErr(invalidBinop(left, right, op.String()))
}

func evalEq(left, right EvalResult, op token.TokenType) EvalResult {
	switch {
	case left.rt == RT_NONE   || right.rt == RT_NONE:	return EvalResult{rt: RT_BOOL, val: false}
	case left.rt == RT_STRING && right.rt == RT_STRING:	return EvalResult{rt: RT_BOOL, val: left.val.(string) == right.val.(string)}
	case left.rt == RT_STRING || right.rt == RT_STRING:	return EvalResult{rt: RT_BOOL, val: false}
	case left.rt == RT_BOOL   && right.rt == RT_BOOL:	return EvalResult{rt: RT_BOOL, val: left.val.(bool) == right.val.(bool)}
	case left.rt == RT_BOOL   || right.rt == RT_BOOL:	return EvalResult{rt: RT_BOOL, val: false}
	case left.rt == RT_INT    && right.rt == RT_INT:	return EvalResult{rt: RT_BOOL, val: left.val.(int) == right.val.(int)}
	case left.rt == RT_FLOAT  && right.rt == RT_FLOAT:	return EvalResult{rt: RT_BOOL, val: left.val.(float64) == right.val.(float64)}
	case left.rt == RT_INT    && right.rt == RT_FLOAT:	return EvalResult{rt: RT_BOOL, val: float64(left.val.(int)) == right.val.(float64)}
	case left.rt == RT_FLOAT  && right.rt == RT_INT:	return EvalResult{rt: RT_BOOL, val: left.val.(float64) == float64(right.val.(int))}
	default:											return evalErr(invalidBinop(left, right, op.String()))
	}
}

func evalNeq(_, _ EvalResult, _ token.TokenType) EvalResult {
	return evalErr("Not equal should be replaced with equality and unary negation. If you see this error, the parser is broken.")
}

func evalGrGr(_, _ EvalResult) EvalResult {
	return evalErr("Trying to evaluate function call as binary operator. If you see this error, the parser is broken.")
}

func evalTuple(left, right EvalResult) EvalResult {
	return EvalResult{rt: RT_TUPLE, val: []EvalResult{left, right}}
}

func evalMod(left, right EvalResult) EvalResult {
	switch left.rt {
	case RT_INT:
		switch right.rt {
		case RT_INT:
			if right.val.(int) == 0 { return evalErr("Division by zero.")}
			return EvalResult{rt: RT_INT, val: left.val.(int) % right.val.(int)}
		}
	}
	return evalErr(invalidBinop(left, right, "%"))
}