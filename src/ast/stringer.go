package ast

import (
	"fmt"
	"mikescript/src/token"
	"mikescript/src/utils"
	"strings"
)

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

func (node VarDeclNodeS) String() string {
	return fmt.Sprintf("%v %v", node.Vartype, node.Identifier)
}

func (node AssignmentNodeS) String() string {
	return fmt.Sprintf("(%v -> %v)", node.Exp, node.Identifier)
}

func (node FuncAppNodeS) String() string {
	
	// map array to strings
	args := utils.MapArrayString(node.Args)

	// join the strings
	argsStr := strings.Join(args, ", ")

	return fmt.Sprintf("(%v)>>%v", argsStr, node.Fun)
}

func (node VariableExpNodeS) String() string {
	return node.Name.Lexeme
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

func (node ReturnNodeS) String() string {
	return fmt.Sprintf("return %s", node.Node)
}

func (node FuncDeclNodeS) String() string {

	// Get args string
	argss := []string{}
	for _, arg := range node.Params {
		argss = append(argss, arg.String())
	}

	// format
	return fmt.Sprintf(
		"%s %s %s %s %s",
		strings.Join(argss, ", "),		// int x, int y
		token.GREATER_GREATER.String(), // >>
		node.Fname.String(),			// fname
		token.MINUS_GREAT.String(),		// ->
		node.Rt.String(),				// int
	)
}


func (fa FuncParamS) String() string {
	typeS := fa.Type.String()
	nameS := fa.Iden.Name.Lexeme
	return fmt.Sprintf("%s %s", typeS, nameS)
}