package parser

import (
	"fmt"

	"github.com/ttakezawa/pgconverger/ast"
	"github.com/ttakezawa/pgconverger/lexer"
	"github.com/ttakezawa/pgconverger/token"
)

type Parser struct {
	l         *lexer.Lexer
	token     token.Token
	peekToken token.Token
	errors    []error
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
	}
	p.advance()
	p.advance()
	return p
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
	return p.parseCreateTableStatement()
}

func (p *Parser) parseCreateTableStatement() ast.Statement {
	createTableStatement := &ast.CreateTableStatement{}

	if !p.expectPeek(token.Table) {
		return nil
	}

	if !p.expectPeek(token.Identifier) {
		return nil
	}
	createTableStatement.TableName = p.token.Literal

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

// column_name data_type [ COLLATE collation ] [ column_constraint [ ... ] ]
func (p *Parser) parseColumnDefinition() *ast.ColumnDefinition {
	var def ast.ColumnDefinition
	tok, ok := p.expect(token.Identifier)
	if !ok {
		return nil
	}
	def.Name = tok

	dataType := p.parseDataType()
	if dataType == nil {
		return nil
	}
	def.Type = dataType

	p.advance()
	columnConstraintList := p.parseColumnConstraintList()
	def.ConstraintList = columnConstraintList

	return &def
}

func (p *Parser) parseDataType() ast.DataType {
	switch p.token.Type {
	case token.Bigint:
		return &ast.DataTypeBigint{p.token}
	case token.Character:
		var dataTypeCharacter ast.DataTypeCharacter
		if p.peekToken.Type == token.Varying {
			p.advance()
			dataTypeCharacter.Varying = true
		}
		if p.peekToken.Type == token.LParen {
			p.advance()
			// Parse: ( n )
			if ok := p.expectPeek(token.Number); !ok {
				return nil
			}
			tok := p.token
			dataTypeCharacter.Length = &tok
			if ok := p.expectPeek(token.RParen); !ok {
				return nil
			}
		}
		return &dataTypeCharacter
	default:
		p.errorf(p.token.Line, "expected DataType, found %s", p.token.Literal)
		return nil
	}
}

func (p *Parser) parseColumnConstraintList() (constraints []ast.ColumnConstraint) {
	for {
		constraint := p.parseColumnConstraint()
		if constraint == nil {
			break
		}
		constraints = append(constraints, constraint)

		switch p.peekToken.Type {
		case token.Not, token.Default:
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
	case token.Default:
		// DEFAULT default_expr
		return &ast.ColumnConstraintDefault{
			Expr: p.parseDefaultExpr(),
		}
	default:
		return nil
	}
}

func (p *Parser) parseDefaultExpr() ast.DefaultExpr {
	tok := p.token
	p.advance()
	return tok
}
