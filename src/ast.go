package main

import "fmt"

// Add statements and replace "ASTNode" with "Expression"
// Add "ASTNode" interface

////////////////////////////////////////$
// Node types
////////////////////////////////////////

type ASTNodeI interface {}
type ExpNodeI interface {
	expressionPlaceholder()
}
type StmtNodeI interface {
	statmentPlaceholder()
}

////////////////////////////////////////
// Statements
////////////////////////////////////////

type Program struct {
	Statements []StmtNodeI
}

type DeclarationNodeS struct {
	Identifier 	VariableExpNodeS
	Vartype 	Token
}

type ExStmtNodeS struct {
	Ex ExpNodeI
}

// forces possible structs for StmtNode
func (Program) statmentPlaceholder() {}
func (DeclarationNodeS) statmentPlaceholder() {}
func (ExStmtNodeS) statmentPlaceholder() {}

////////////////////////////////////////
// Expressios
////////////////////////////////////////

type AssignmentNodeS struct {
	Identifier VariableExpNodeS
	Exp        ExpNodeI
}

type FuncAppNodeS struct {
	Args 	[]ExpNodeI
	fun		ExpNodeI
}

type BinaryExpNodeS struct {
	Left  ExpNodeI
	Op    Token
	Right ExpNodeI
}

type UnaryExpNodeS struct {
	Op   Token
	Node ExpNodeI
}

type LiteralExpNodeS struct {
	Tk Token
}

type VariableExpNodeS struct {
	Name Token
}

type GroupExpNodeS struct {
	Node       ExpNodeI
	TokenLeft  Token
	TokenRight Token
}

// forces possible structs for ExpNode
func (AssignmentNodeS) expressionPlaceholder() {}
func (FuncAppNodeS) expressionPlaceholder() {}
func (BinaryExpNodeS) expressionPlaceholder() {}
func (UnaryExpNodeS) expressionPlaceholder() {}
func (LiteralExpNodeS) expressionPlaceholder() {}
func (GroupExpNodeS) expressionPlaceholder() {}
func (VariableExpNodeS) expressionPlaceholder() {}


////////////////////////////////////////
// Stringer
////////////////////////////////////////

func (node Program) String() string {
	s := ""
	for i, stmt := range node.Statements {
		s += fmt.Sprintf("[%v] %v\n", i, stmt)
	}
	return s
}


func (node ExStmtNodeS) String() string {
	return fmt.Sprintf("%v;", node.Ex)
}

func (node BinaryExpNodeS) String() string {
	return fmt.Sprintf("(%v %v %v)", node.Left, node.Op.Lexeme, node.Right)
}

func (node UnaryExpNodeS) String() string {
	return fmt.Sprintf("(%v %v)", node.Op.Lexeme, node.Node)
}

func (node LiteralExpNodeS) String() string {
	return node.Tk.Lexeme
}

func (node GroupExpNodeS) String() string {
	return fmt.Sprintf("(%v)", node.Node)
}

func (node DeclarationNodeS) String() string {
	return fmt.Sprintf("%v %v", node.Vartype.Lexeme, node.Identifier)
}

func (node AssignmentNodeS) String() string {
	return fmt.Sprintf("(%v -> %v)", node.Exp, node.Identifier)
}

func (node FuncAppNodeS) String() string {
	
	// map array to strings
	args := mapArrayString(node.Args)

	// join the strings
	argsStr := strJoin(args, ", ")
	
	for _, arg := range node.Args {
		fmt.Println(arg)
	}

	return fmt.Sprintf("%v(%v)", node.fun, argsStr)
}

func (node VariableExpNodeS) String() string {
	return "Var: " + node.Name.Lexeme
}

////////////////////////////////////////
// Helper functions
////////////////////////////////////////
func fmap[T any, F any](a []T, f func(T) F) []F {
	fs := make([]F, len(a))
	for i, v := range a {
		fs[i] = f(v)
	}
	return fs
}

func mapArrayString[T any](a []T) []string {
	return fmap[T, string](a, func(v T) string { return fmt.Sprint(v) })
}

func strJoin(a []string, sep string) string {
	out := ""
	for i, s := range a {
		if i > 0 {
			out += sep
		}
		out += s
	}
	return out
}

