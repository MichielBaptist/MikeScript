package ast

type ASTNodeI any

type ExpNodeI interface {
	expressionPlaceholder()
}
type StmtNodeI interface {
	statmentPlaceholder()
}
