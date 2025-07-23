package ast

import "mikescript/src/token"

type AssignmentNodeS struct {
	Identifier VariableExpNodeS
	Exp        ExpNodeI
}

type DeclAssignNodeS struct {
	Identifier 	VariableExpNodeS
	Exp 		ExpNodeI
}

// exp >> exp
type FuncAppNodeS struct {
	Args 	[]ExpNodeI
	Fun		ExpNodeI
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

// forces possible structs for ExpNode
func (AssignmentNodeS) expressionPlaceholder() {}
func (DeclAssignNodeS) expressionPlaceholder() {}
func (FuncAppNodeS) expressionPlaceholder() {}
func (BinaryExpNodeS) expressionPlaceholder() {}
func (UnaryExpNodeS) expressionPlaceholder() {}
func (LiteralExpNodeS) expressionPlaceholder() {}
func (GroupExpNodeS) expressionPlaceholder() {}
func (VariableExpNodeS) expressionPlaceholder() {}
func (LogicalExpNodeS) expressionPlaceholder() {}