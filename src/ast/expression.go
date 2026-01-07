package ast

import (
	"mikescript/src/mstype"
	"mikescript/src/token"
)

type AssignmentNodeS struct {
	// exp '->' IDENTIFIER ';' 
	Identifier *VariableExpNodeS
	Exp        ExpNodeI
}

type DeclAssignNodeS struct {
	// exp '=>' IDENTIFIER ';'
	Identifier 	*VariableExpNodeS
	Exp 		ExpNodeI
}

// exp, exp, ... >> exp
type FuncAppNodeS struct {
	Args 	[]ExpNodeI
	Fun		ExpNodeI
}

// exp, exp, ... .>> exp
type IterableFuncAppNodeS struct {
	Args 	ExpNodeI
	Fun		ExpNodeI
}

// exp, exp, ... .>>= exp
type IterableFuncAppAndCallNodeS struct {
	Args 	ExpNodeI
	Fun		ExpNodeI
}

// '=' exp
type FuncCallNodeS struct {
	Fun ExpNodeI
	Op token.Token
}

// .= exp
type IterableFuncCallNodeS struct {
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

type TupleNodeS struct {
	Expressions []ExpNodeI
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

type RangeConstructorNodeS struct {
	From ExpNodeI // expects to be evaluate to int
	To ExpNodeI // expects to be evaluate to int
}

type ArrayAssignmentNodeS struct {
	Target ExpNodeI
	Index ExpNodeI
	Value ExpNodeI
}

type StructConstructorNodeS struct {
	Name *mstype.MSNamedTypeS
	Fields map[*VariableExpNodeS]ExpNodeI
}

type FieldAccessNodeS struct {
	Target ExpNodeI
	Field *VariableExpNodeS
}

type FieldAssignmentNode struct {
	Target ExpNodeI
	Field *VariableExpNodeS
	Value ExpNodeI
}

type StarredExpNodeS struct {
	Node ExpNodeI
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
func (*TupleNodeS) expressionPlaceholder() {}
func (*StructConstructorNodeS) expressionPlaceholder() {}
func (*FieldAccessNodeS) expressionPlaceholder() {}
func (*FieldAssignmentNode) expressionPlaceholder() {}
func (*IterableFuncCallNodeS) expressionPlaceholder() {}
func (*IterableFuncAppNodeS) expressionPlaceholder() {}
func (*IterableFuncAppAndCallNodeS) expressionPlaceholder() {}
func (*RangeConstructorNodeS) expressionPlaceholder() {}
func (*StarredExpNodeS) expressionPlaceholder() {}

func (ve *VariableExpNodeS) VarName() string {

	// 'anonymous' functions have no name
	if ve == nil {
		return ""
	}

	return ve.Name.Lexeme
}