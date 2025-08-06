package ast

import (
	"mikescript/src/mstype"
	"mikescript/src/token"
)

type AssignmentNodeS struct {
	Identifier *VariableExpNodeS
	Exp        ExpNodeI
}

type DeclAssignNodeS struct {
	Identifier 	*VariableExpNodeS
	Exp 		ExpNodeI
}

// exp, exp, ... >> exp
type FuncAppNodeS struct {
	Args 	[]ExpNodeI
	Fun		ExpNodeI
}

// = exp
type FuncCallNodeS struct {
	Fun ExpNodeI
	Op token.Token
}

type BinaryExpNodeS struct {
	Left  ExpNodeI
	Op    token.Token		// +, -, /, ...
	Right ExpNodeI
}

type LogicalExpNodeS struct {
	Left  ExpNodeI
	Op    token.Token		// &&, ||
	Right ExpNodeI
}

type UnaryExpNodeS struct {
	Op   token.Token
	Node ExpNodeI
}

type LiteralExpNodeS struct {
	Tk token.Token
}

type VariableExpNodeS struct {
	Name token.Token
}

type GroupExpNodeS struct {
	Node       ExpNodeI
	TokenLeft  token.Token
	TokenRight token.Token
}

type ArrayIndexNodeS struct {
	Target ExpNodeI
	Index ExpNodeI
}

type ArrayConstructorNodeS struct {
	// '['']' type '{' {expression ','} * '}' 
	Type mstype.MSType	
	Vals []ExpNodeI
	N ExpNodeI
}

type ArrayAssignmentNodeS struct {
	Target ExpNodeI
	Index ExpNodeI
	Value ExpNodeI
}

// forces possible structs for ExpNode
// pointer to these structs implement expression
func (*AssignmentNodeS) expressionPlaceholder() {}
func (*DeclAssignNodeS) expressionPlaceholder() {}
func (*FuncAppNodeS) expressionPlaceholder() {}
func (*FuncCallNodeS) expressionPlaceholder() {}
func (*BinaryExpNodeS) expressionPlaceholder() {}
func (*UnaryExpNodeS) expressionPlaceholder() {}
func (*LiteralExpNodeS) expressionPlaceholder() {}
func (*GroupExpNodeS) expressionPlaceholder() {}
func (*VariableExpNodeS) expressionPlaceholder() {}
func (*LogicalExpNodeS) expressionPlaceholder() {}
func (*ArrayIndexNodeS) expressionPlaceholder() {}
func (*ArrayConstructorNodeS) expressionPlaceholder() {}
func (*ArrayAssignmentNodeS) expressionPlaceholder() {}

func (ve *VariableExpNodeS) VarName() string {

	// 'anonymous' functions have no name
	if ve == nil {
		return ""
	}

	return ve.Name.Lexeme
}