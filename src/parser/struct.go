package parser

import (
	"mikescript/src/ast"
	"mikescript/src/mstype"
	token "mikescript/src/token"
)


func (p *MSParser) parseStructDeclaration() (*ast.StructDeclarationNodeS, error) {

	var fields map[*ast.VariableExpNodeS]mstype.MSType
	var typ mstype.MSType
	var sname *ast.VariableExpNodeS		// struct name placeholder
	var fname *ast.VariableExpNodeS		// field name placeholder
	var err error

	fields = make(map[*ast.VariableExpNodeS]mstype.MSType)

	sname, err = p.parseIdentifier()

	if err != nil {
		return nil, err
	}

	if ok, tok := p.expect(token.LEFT_BRACE) ; !ok {
		return nil, p.unexpectedToken(tok, token.RIGHT_BRACE)
	}

	for {

		if ok, _ := p.lookahead(token.RIGHT_BRACE) ; ok {
			break
		}

		// parse type and id
		typ, err = p.parseType()

		if err != nil {
			return nil, err
		}

		fname, err = p.parseIdentifier()

		if err != nil {
			return nil, err
		}

		fields[fname] = typ

		// break ok no ';'
		if ok, _ := p.match(token.SEMICOLON) ; !ok {
			break
		}
	}

	// we want closing brace
	if ok, tok := p.match(token.RIGHT_BRACE) ; !ok {
		return nil, p.unexpectedToken(tok, token.RIGHT_BRACE)
	}

	return &ast.StructDeclarationNodeS{Name: sname, Fields: fields}, nil
}

func (p *MSParser) parseStructConstructor() (ast.ExpNodeI, error) {
	// IDENTIFIER
	// IDENTIFIER '{' exp? {',' exp} '}'

	var err error
	var name *ast.VariableExpNodeS
	// var exp ast.ExpNodeI
	// var fields map[*ast.VariableExpNodeS]ast.ExpNodeI

	// parse identifier
	name, err = p.parseIdentifier()

	if err != nil {
		return nil, err
	}

	return name, err

	// if ok, _ := p.match(token.LEFT_BRACE) ; !ok {
	// 	// just identifier
	// 	return name, err
	// }

	// // struct constructor 100 % because we got '{'

	// fields = make(map[*ast.VariableExpNodeS]ast.ExpNodeI)

	// for {

	// 	if ok, _ := p.lookahead(token.RIGHT_BRACE) ; ok {
	// 		break
	// 	}

	// 	exp, err = p.parseExpression()

	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	assignment, ok := exp.(*ast.AssignmentNodeS)

	// 	if !ok {
	// 		msg := fmt.Sprintf("Expected assignment expression, got '%v'", exp)
	// 		return nil, p.error(msg, 0, 0)
	// 	}

	// 	// add to fields
	// 	fields[assignment.Identifier] = assignment.Exp

	// 	if ok, _ := p.match(token.SEMICOLON) ; !ok {
	// 		break
	// 	}
	// }

	// // expect '}'
	// if ok, tok := p.expect(token.RIGHT_BRACE) ; !ok {
	// 	return nil, p.unexpectedToken(tok, token.RIGHT_BRACE)
	// }

	// structConstructor := &ast.StructConstructorNodeS{
	// 	Name: &mstype.MSNamedTypeS{Name: name.VarName()},	// convert varaible to named type
	// 	Fields: fields,
	// }

	// return structConstructor, nil
}