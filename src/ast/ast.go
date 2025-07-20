package ast

import (
	"fmt"
	"mikescript/src/token"
	"mikescript/src/utils"
	"strings"
)

// Add statements and replace "ASTNode" with "Expression"
// Add "ASTNode" interface

////////////////////////////////////////$
// Node types
////////////////////////////////////////

type ASTNodeI any
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

type BlockNodeS struct {
	Statements []StmtNodeI
}

type DeclarationNodeS struct {
	Identifier 	VariableExpNodeS
	Vartype 	token.Token
}

type ExStmtNodeS struct {
	Ex ExpNodeI
}

type IfNodeS struct {
	Condition 	ExpNodeI
	ThenStmt 	StmtNodeI
	ElseStmt 	StmtNodeI
}

type WhileNodeS struct {
	Condition 	ExpNodeI
	Body 		BlockNodeS
}

type ContinueNodeS struct {
	Tk token.Token
}
type BreakNodeS struct {
	Tk token.Token
}

type FuncDeclNodeS struct {
	Fname VariableExpNodeS			// name, should be VariableExpNodeS???
	Args []FuncArgS 				// Arguments
	Rt token.Token					// Return type
	Body BlockNodeS					// Body of function
}

type FuncArgS struct {
	Type token.Token		// Type token
	Iden VariableExpNodeS 	// Var name
}

// forces possible structs for StmtNode
func (Program) statmentPlaceholder() {}
func (BlockNodeS) statmentPlaceholder() {}
func (DeclarationNodeS) statmentPlaceholder() {}
func (ExStmtNodeS) statmentPlaceholder() {}
func (IfNodeS) statmentPlaceholder() {}
func (WhileNodeS) statmentPlaceholder() {}
func (ContinueNodeS) statmentPlaceholder() {}
func (BreakNodeS) statmentPlaceholder() {}
func (FuncDeclNodeS) statmentPlaceholder() {}

////////////////////////////////////////
// Expressios
////////////////////////////////////////

type AssignmentNodeS struct {
	Identifier VariableExpNodeS
	Exp        ExpNodeI
}

type FuncAppNodeS struct {
	Args 	[]ExpNodeI
	Fun		ExpNodeI
}

type BinaryExpNodeS struct {
	Left  ExpNodeI
	Op    token.Token
	Right ExpNodeI
}

type LogicalExpNodeS struct {
	Left  ExpNodeI
	Op    token.Token
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
func (FuncAppNodeS) expressionPlaceholder() {}
func (BinaryExpNodeS) expressionPlaceholder() {}
func (UnaryExpNodeS) expressionPlaceholder() {}
func (LiteralExpNodeS) expressionPlaceholder() {}
func (GroupExpNodeS) expressionPlaceholder() {}
func (VariableExpNodeS) expressionPlaceholder() {}
func (LogicalExpNodeS) expressionPlaceholder() {}


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
	switch node.Op.Type {
	case token.GREATER_GREATER:	return fmt.Sprintf("%v(%v)", node.Right, node.Left)
	case token.COMMA:			return fmt.Sprintf("(%v, %v)", node.Left, node.Right)
	default: 					return fmt.Sprintf("(%v %v %v)", node.Left, node.Op.Lexeme, node.Right)
	}
}

func (node UnaryExpNodeS) String() string {
	return fmt.Sprintf("(%v %v)", node.Op.Lexeme, node.Node)
}

func (node LiteralExpNodeS) String() string {
	switch node.Tk.Type{
	case token.STRING:			return fmt.Sprintf("\"%v\"", node.Tk.Lexeme)
	case token.NUMBER_FLOAT:	return fmt.Sprintf("%v", node.Tk.Lexeme)
	case token.NUMBER_INT:		return fmt.Sprintf("%v", node.Tk.Lexeme)
	default:					return node.Tk.Lexeme
	}
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
	args := utils.MapArrayString(node.Args)

	// join the strings
	argsStr := utils.StrJoin(args, ", ")

	return fmt.Sprintf("%v(%v)", node.Fun, argsStr)
}

func (node VariableExpNodeS) String() string {
	return "$" + node.Name.Lexeme
}

func (node BlockNodeS) String() string {
	s := "{\n"
	for _, stmt := range node.Statements {
		s += fmt.Sprintf("%v\n", stmt)
	}
	s += "}"
	return s
}

func (node IfNodeS) String() string {
	s := fmt.Sprintf("if %v %v", node.Condition, node.ThenStmt)
	if node.ElseStmt != nil {
		s += fmt.Sprintf(" else %v", node.ElseStmt)
	}
	return s
}

func (node LogicalExpNodeS) String() string {
	return fmt.Sprintf("(%v %v %v)", node.Left, node.Op.Lexeme, node.Right)
}

func (node WhileNodeS) String() string {
	return fmt.Sprintf("while %v %v", node.Condition, node.Body)
}

func (node ContinueNodeS) String() string {
	return "continue"
}

func (node BreakNodeS) String() string {
	return "break"
}

func (node FuncDeclNodeS) String() string {

	// Get args string
	argss := []string{}
	for _, arg := range node.Args {
		argss = append(argss, arg.String())
	}

	// format
	return fmt.Sprintf(
		"%s %s %s %s %s",
		strings.Join(argss, ", "),		// int x, int y
		token.GREATER_GREATER.String(), // >>
		node.Fname.String(),			// fname
		token.MINUS_GREAT.String(),		// ->
		node.Rt.Lexeme,					// int
	)
}


func (fa FuncArgS) String() string {
	typeS := fa.Type.Lexeme
	nameS := fa.Iden.Name.Lexeme
	return fmt.Sprintf("%s %s", typeS, nameS)
}