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
	if p.peekToken.Typ != typ {
		p.errorf(p.peekToken.Line, "expected %s, found %s", typ, p.peekToken.Val)
		return false
	}
	p.advance()
	return true
}

func (p *Parser) ParseDataDefinition() *ast.DataDefinition {
	statements := p.parseStatementList()
	return &ast.DataDefinition{StatementList: statements}
}

func (p *Parser) parseStatementList() (list []ast.Statement) {
	for p.token.Typ != token.EOF {
		if statement := p.parseStatement(); statement != nil {
			list = append(list, statement)
		}

	skipRest:
		for {
			switch p.token.Typ {
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
	createTableStatement.TableName = p.token.Val

	if !p.expectPeek(token.LParen) {
		return nil
	}

	if !p.expectPeek(token.RParen) {
		return nil
	}

	switch p.peekToken.Typ {
	case token.Semicolon:
		p.advance()
	case token.EOF:
	default:
		p.errorf(p.peekToken.Line, "expected %s, found %s", token.Semicolon, p.peekToken.Val)
	}

	return createTableStatement
}
