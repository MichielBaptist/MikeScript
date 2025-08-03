package parser

import (
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
	// 4. array types 'type[]'

	// Case 1: basic types
	switch _, tok := p.match(token.SimpleTypeKeywords...) ; tok.Type {
	case token.INT_TYPE:		return mstype.MS_INT, nil
	case token.FLOAT_TYPE:		return mstype.MS_FLOAT, nil
	case token.BOOLEAN_TYPE:	return mstype.MS_BOOL, nil
	case token.STRING_TYPE:		return mstype.MS_STRING, nil
	case token.NOTHING_TYPE:	return mstype.MS_NOTHING, nil
	}

	// Array type
	if ok, _ := p.match(token.LEFT_SQUARE) ; ok{
		return p.parseArrayType()
	}

	// Composite or function type
	if ok, _ := p.match(token.LEFT_PAREN) ; ok{
		return p.parseCompositeOrFunctionType()
	}

	// cannot match current token
	err := p.unexpectedToken(p.peek(),  token.LEFT_PAREN, token.LEFT_SQUARE)
	return mstype.MS_NOTHING, err

}

func (p *MSParser) parseArrayType() (mstype.MSType, error) {

	// expect a ']'
	if ok, tok := p.match(token.RIGHT_SQUARE) ; !ok {
		return mstype.MS_NOTHING, p.unexpectedToken(tok, token.RIGHT_SQUARE)
	}

	// parse type
	base, err := p.parseType()

	if err != nil {
		return mstype.MS_NOTHING, err
	}

	return &mstype.MSArrayType{Type: base}, err

}

func (p *MSParser) parseCompositeOrFunctionType() (mstype.MSType, error) {
	// Check if we have a closing ')' immediately or a '->'
	// we have an empty typelist. In this case we don't need
	// to parse typelist
	var types []mstype.MSType
	var err error

	if ok, _ := p.lookahead(token.RIGHT_PAREN, token.MINUS_GREAT) ; ok {
		types = []mstype.MSType{}
	} else {
		types, err = p.parseTypeList()
	}

	// Potential parse error
	if err != nil {
		return mstype.MS_NOTHING, err
	}

	// Founc composite type
	if ok, _ := p.match(token.RIGHT_PAREN) ; ok {
		return &mstype.MSCompositeTypeS{Types: types}, nil
	}

	// Found operation type
	if ok, _ := p.match(token.MINUS_GREAT) ; ok {

		// Check if the returntype is empty by checking if we have
		// a ')' immediately after the '->' 
		if ok, _ := p.match(token.RIGHT_PAREN) ; ok {
			return &mstype.MSOperationTypeS{Left: types, Right: mstype.MS_NOTHING}, err
		}

		// No ')' so expect a type
		returnType, err := p.parseType()

		if err != nil {
			return mstype.MS_NOTHING, err
		}

		// expect a ')'
		if ok, tok := p.match(token.RIGHT_PAREN) ; !ok {
			return mstype.MS_NOTHING, p.unexpectedToken(tok, token.RIGHT_PAREN)
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