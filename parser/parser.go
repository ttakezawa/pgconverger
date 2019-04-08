package parser

import (
	"fmt"

	"github.com/ttakezawa/pgconverger/ast"
	"github.com/ttakezawa/pgconverger/lexer"
	"github.com/ttakezawa/pgconverger/token"
)

const (
	_ int = iota
	precedenceLowest
	precedenceSum
	precedenceProduct
	precedenceTypecast
	precedencePrefix // -x
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l         *lexer.Lexer
	token     token.Token
	peekToken token.Token
	errors    []error

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.String, p.parseStringLiteral)
	p.registerPrefix(token.Number, p.parseNumberLiteral)
	p.registerPrefix(token.Identifier, p.parseIdentifierAsExpression)
	// p.registerPrefix(token.Minus, p.parsePrefixExpression)
	// p.registerPrefix(token.Plus, p.parsePrefixExpression)
	p.registerPrefix(token.True, p.parseBoolean)
	p.registerPrefix(token.False, p.parseBoolean)
	// p.registerPrefix(token.Null, p.parseNull)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.Plus, p.parseInfixExpression)
	p.registerInfix(token.Minus, p.parseInfixExpression)
	p.registerInfix(token.Slash, p.parseInfixExpression)
	p.registerInfix(token.Asterisk, p.parseInfixExpression)
	p.registerInfix(token.Typecast, p.parseInfixExpression)

	p.advance()
	p.advance()
	return p
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) Errors() []error {
	return p.errors
}

var precedences = map[token.TokenType]int{
	token.Plus:     precedenceSum,
	token.Minus:    precedenceSum,
	token.Slash:    precedenceProduct,
	token.Asterisk: precedenceProduct,
	token.Typecast: precedenceTypecast,
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return precedenceLowest
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.token.Type]; ok {
		return p
	}
	return precedenceLowest
}

func (p *Parser) advance() {
	p.token = p.peekToken
	p.peekToken = p.l.NextToken()
}

type parseError struct {
	error
	line int
}

func (e *parseError) Error() string {
	return fmt.Sprintf("<input>:%d: %s", e.line, e.error.Error())
}

func (p *Parser) errorf(line int, format string, a ...interface{}) {
	p.errors = append(p.errors,
		&parseError{
			fmt.Errorf(format, a...),
			line,
		},
	)
}

func (p *Parser) expectPeek(typ token.TokenType) bool {
	if p.peekToken.Type != typ {
		p.errorf(p.peekToken.Line, "expected %s, found %s", typ, p.peekToken.Literal)
		return false
	}
	p.advance()
	return true
}

func (p *Parser) expect(typ token.TokenType) (token.Token, bool) {
	token := p.token
	if p.token.Type != typ {
		p.errorf(p.token.Line, "expected %s, found %s", typ, p.token.Literal)
		return token, false
	}
	p.advance()
	return token, true
}

func (p *Parser) ParseDataDefinition() *ast.DataDefinition {
	statements := p.parseStatementList()
	return &ast.DataDefinition{StatementList: statements}
}

func (p *Parser) parseStatementList() (list []ast.Statement) {
	for p.token.Type != token.EOF {
		if statement := p.parseStatement(); statement != nil {
			list = append(list, statement)
		}

	skipRest:
		for {
			switch p.token.Type {
			case token.EOF:
				break skipRest
			case token.Semicolon:
				p.advance()
				break skipRest
			default:
				p.advance()
			}
		}
	}

	return list
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.token.Type {
	case token.Create:
		switch p.peekToken.Type {
		case token.Schema:
			return p.parseCreateSchemaStatement()
		case token.Table:
			return p.parseCreateTableStatement()
		case token.Extension, token.Sequence:
			// Not yet implemented
			return nil
		case token.View:
			// Not yet implemented
			return nil
		default:
			p.errorf(p.peekToken.Line, "unknown token: CREATE %s", p.peekToken.Literal)
		}
	case token.Alter:
		switch p.peekToken.Type {
		case token.Schema:
			// Not yet implemented
			return nil
		}
	case token.Set, token.Comment:
		// Not yet implemented
		return nil
	default:
		p.errorf(p.token.Line, "unknown token: %s", p.token.Literal)
	}
	return nil
}

func (p *Parser) parseCreateSchemaStatement() ast.Statement {
	var createSchemaStatement ast.CreateSchemaStatement
	if !p.expectPeek(token.Schema) {
		return nil
	}
	p.advance()
	identifier := p.parseIdentifier()
	if identifier == nil {
		return nil
	}
	createSchemaStatement.Name = identifier
	return &createSchemaStatement
}

func (p *Parser) parseCreateTableStatement() ast.Statement {
	createTableStatement := &ast.CreateTableStatement{}

	if !p.expectPeek(token.Table) {
		return nil
	}

	p.advance()
	identifier := p.parseIdentifier()
	if identifier == nil {
		return nil
	}
	createTableStatement.TableName = identifier

	if !p.expectPeek(token.LParen) {
		return nil
	}
	p.advance()
	columnDefinitionList := p.parseColumnDefinitionList()
	if !p.expectPeek(token.RParen) {
		return nil
	}
	createTableStatement.ColumnDefinitionList = columnDefinitionList

	switch p.peekToken.Type {
	case token.Semicolon:
		p.advance()
	case token.EOF:
	default:
		p.errorf(p.peekToken.Line, "expected %s, found %s", token.Semicolon, p.peekToken.Literal)
	}

	return createTableStatement
}

// { column_definition, ... }
func (p *Parser) parseColumnDefinitionList() (defs []*ast.ColumnDefinition) {
	if p.peekToken.Type == token.RParen {
		return
	}
	if def := p.parseColumnDefinition(); def != nil {
		defs = append(defs, def)
	}
	for p.peekToken.Type == token.Comma {
		p.advance()
		p.advance()
		if def := p.parseColumnDefinition(); def != nil {
			defs = append(defs, def)
		}
	}
	return
}

func (p *Parser) parseIdentifierAsExpression() ast.Expression {
	return p.parseIdentifier()
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	switch {
	case p.token.Type == token.Identifier:
		identifier := ast.Identifier{Token: p.token}
		switch {
		case len(p.token.Literal) == 0:
			identifier.Value = ""
		case p.token.Literal[0] == '"':
			identifier.Value = p.token.Literal[1 : len(p.token.Literal)-1]
		default:
			identifier.Value = p.token.Literal
		}
		return &identifier
	case !p.token.IsReserved():
		return &ast.Identifier{
			Token: p.token,
			Value: p.token.Literal,
		}
	default:
		p.errorf(p.token.Line, "expected identifier, found %s", p.token.Literal)
		return nil
	}
}

// column_name data_type [ COLLATE collation ] [ column_constraint [ ... ] ]
func (p *Parser) parseColumnDefinition() *ast.ColumnDefinition {
	var def ast.ColumnDefinition
	identifier := p.parseIdentifier()
	if identifier == nil {
		return nil
	}
	def.Name = identifier
	p.advance()

	dataType := p.parseDataType()
	if dataType == nil {
		return nil
	}
	def.Type = dataType

	switch p.peekToken.Type {
	case token.Comma, token.RParen:
		// Do nothing
	default:
		p.advance()
		// Parse constraint
		columnConstraintList := p.parseColumnConstraintList()
		if columnConstraintList == nil {
			return nil
		}
		def.ConstraintList = columnConstraintList
	}
	return &def
}

func (p *Parser) parseDataType() ast.DataType {
	switch p.token.Type {
	case token.Integer:
		return &ast.DataTypeInteger{p.token}
	case token.Bigint:
		return &ast.DataTypeBigint{p.token}
	case token.Boolean:
		return &ast.DataTypeBoolean{}
	case token.Numeric:
		return &ast.DataTypeNumeric{}
	case token.Character:
		var dataTypeCharacter ast.DataTypeCharacter
		if p.peekToken.Type == token.Varying {
			p.advance()
			dataTypeCharacter.Varying = true
		}
		if p.peekToken.Type == token.LParen {
			p.advance()
			optionLength := p.parseDataTypeOptionLength()
			if optionLength == nil {
				return nil
			}
			dataTypeCharacter.OptionLength = optionLength
		}
		return &dataTypeCharacter
	case token.Date:
		return &ast.DataTypeDate{}
	case token.Timestamp:
		var dataTypeTimestamp ast.DataTypeTimestamp
		// timestamp with time zone
		if p.peekToken.Type == token.With {
			p.advance()
			if ok := p.expectPeek(token.Time) && p.expectPeek(token.Zone); !ok {
				return nil
			}
			dataTypeTimestamp.WithTimeZone = true
		}
		return &dataTypeTimestamp
	case token.Text:
		return &ast.DataTypeText{}
	case token.Jsonb:
		return &ast.DataTypeJsonb{}
	case token.Bytea:
		return &ast.DataTypeBytea{}
	case token.Tsvector:
		return &ast.DataTypeTsvector{}
	default:
		switch p.token.Literal {
		case `"date"`:
			return &ast.DataTypeDate{}
		case `"text"`:
			return &ast.DataTypeText{}
		case `"jsonb"`:
			return &ast.DataTypeJsonb{}
		case `"bytea"`:
			return &ast.DataTypeBytea{}
		case `"tsvector"`:
			return &ast.DataTypeTsvector{}
		}

		p.errorf(p.token.Line, "expected DataType, found %s", p.token.Literal)
		return nil
	}
}

// Parse: ( n )
func (p *Parser) parseDataTypeOptionLength() *ast.DataTypeOptionLength {
	if ok := p.expectPeek(token.Number); !ok {
		return nil
	}
	tok := p.token
	if ok := p.expectPeek(token.RParen); !ok {
		return nil
	}
	return &ast.DataTypeOptionLength{tok}
}

func (p *Parser) parseColumnConstraintList() (constraints []ast.ColumnConstraint) {
	for {
		constraint := p.parseColumnConstraint()
		if constraint == nil {
			break
		}
		constraints = append(constraints, constraint)

		switch p.peekToken.Type {
		case token.Not, token.Null, token.Default:
			p.advance()
			continue
		}
		break
	}
	return constraints
}

func (p *Parser) parseColumnConstraint() ast.ColumnConstraint {
	switch p.token.Type {
	case token.Not:
		// NOT NULL
		if !p.expectPeek(token.Null) {
			return nil
		}
		return &ast.ColumnConstraintNotNull{}
	case token.Null:
		// NULL
		return &ast.ColumnConstraintNull{}
	case token.Default:
		// DEFAULT expr
		p.advance()
		expr := p.parseExpression(precedenceLowest)
		if expr == nil {
			return nil
		}
		return &ast.ColumnConstraintDefault{
			Expr: expr,
		}
	default:
		return nil
	}
}

func (p *Parser) noPrefixParseFnError(t token.Token) {
	p.errorf(t.Line, "no prefix parse function for %s found", t.Type)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.token.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.token)
		return nil
	}
	leftExp := prefix()

	for precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.advance()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Operator: p.token,
		Left:     left,
	}

	precedence := p.currentPrecedence()
	p.advance()
	right := p.parseExpression(precedence)
	if right == nil {
		return nil
	}
	expression.Right = right

	return expression
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.token}
}

func (p *Parser) parseNumberLiteral() ast.Expression {
	return &ast.NumberLiteral{Token: p.token}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.BooleanLiteral{Token: p.token}
}
