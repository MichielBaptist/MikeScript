package parser

import (
	"fmt"
	"mikescript/src/mstype"
	"mikescript/src/token"
	"slices"
)

func (p *MSParser) parseType() (mstype.MSType, error) {
	// arg: type of declartion;
	// Need to consider 3 cases:
	// 1. basic types 'int', 'bool', 'float', 'string'
	// 2. composite types (type, type), ()
	// 3. function types (type, type, type -> type), ( -> type), (->)

	// Case 1: basic types
	switch _, tok := p.match(token.SimpleTypeKeywords...) ; tok.Type {
	case token.INT_TYPE:		return mstype.MS_INT, nil
	case token.FLOAT_TYPE:		return mstype.MS_FLOAT, nil
	case token.BOOLEAN_TYPE:	return mstype.MS_BOOL, nil
	case token.STRING_TYPE:		return mstype.MS_STRING, nil
	case token.NOTHING_TYPE:	return mstype.MS_NOTHING, nil
	}

	// Case 2 & 3: we expect a '('
	if ok, _ := p.expect(token.LEFT_PAREN) ; !ok {
		return mstype.MS_NOTHING, nil
		// return mstype.MS_NOTHING, p.unexpectedToken(tok, token.LEFT_PAREN)
	}

	// Check if we have a closing ')' immediately or a '->'
	// we have an empty typelist. In this case we don't need
	// to parse typelist
	var types []mstype.MSType
	var err error

	if ok, _ := p.lookahead(token.RIGHT_PAREN, token.MINUS_GREAT) ; ok {
		types = []mstype.MSType{}
		fmt.Println("Found no type list...")
	} else {
		types, err = p.parseTypeList()
	}

	// Potential parse error
	if err != nil {
		return mstype.MS_NOTHING, err
	}

	if ok, _ := p.match(token.RIGHT_PAREN) ; ok {
		// composite type
		return &mstype.MSCompositeTypeS{Types: types}, nil
	}

	if ok, _ := p.match(token.MINUS_GREAT) ; ok {

		// operation type, parse return type
		returnType, err := p.parseType()

		if err != nil {
			return mstype.MS_NOTHING, err
		}

		if ok, tk := p.match(token.RIGHT_PAREN); !ok {
			return mstype.MS_NOTHING, p.unexpectedToken(tk, token.RIGHT_PAREN)
		}

		return &mstype.MSOperationTypeS{Left: types, Right: returnType}, nil
	}

	return mstype.MS_NOTHING, p.unexpectedToken(p.peek(), token.RIGHT_PAREN, token.MINUS_GREAT)
	
}

func (p *MSParser) parseTypeList() ([]mstype.MSType, error) {

	typelist := []mstype.MSType{}
	for {
		
		// parse type
		t, err := p.parseType()

		if err != nil {
			return typelist, err
		}

		typelist = append(typelist, t)

		// When not encountering a comma, we break
		// from the parse loop
		if ok, _ := p.match(token.COMMA) ; !ok {
			break
		}
	}

	return typelist, nil
}

func IsSimpleTypeToken(t token.Token) bool {
	return slices.Contains(token.SimpleTypeKeywords, t.Type)
}