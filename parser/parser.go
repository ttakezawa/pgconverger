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
	precedenceIs
	precedenceSum
	precedenceProduct
	precedenceTypecast
	precedencePrefix // -x
	precedenceCall
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
	p.registerPrefix(token.Text, p.parseIdentifierAsExpression)
	// p.registerPrefix(token.Minus, p.parsePrefixExpression)
	// p.registerPrefix(token.Plus, p.parsePrefixExpression)
	p.registerPrefix(token.True, p.parseBoolean)
	p.registerPrefix(token.False, p.parseBoolean)
	p.registerPrefix(token.Null, p.parseNull)
	p.registerPrefix(token.LParen, p.parseGroupedExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.Plus, p.parseInfixExpression)
	p.registerInfix(token.Minus, p.parseInfixExpression)
	p.registerInfix(token.Slash, p.parseInfixExpression)
	p.registerInfix(token.Asterisk, p.parseInfixExpression)
	p.registerInfix(token.Typecast, p.parseInfixExpression)
	p.registerInfix(token.Is, p.parseInfixExpression)
	p.registerInfix(token.LParen, p.parseCallExpression)

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

func (p *Parser) InputName() string {
	return p.l.InputName()
}

func (p *Parser) Errors() []error {
	return p.errors
}

// https://www.postgresql.org/docs/10/sql-syntax-lexical.html#SQL-PRECEDENCE
var precedences = map[token.TokenType]int{
	token.Is:       precedenceIs,
	token.Plus:     precedenceSum,
	token.Minus:    precedenceSum,
	token.Slash:    precedenceProduct,
	token.Asterisk: precedenceProduct,
	token.Typecast: precedenceTypecast,
	token.LParen:   precedenceCall,
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
	p *Parser
	error
	line int
}

func (e *parseError) Error() string {
	return fmt.Sprintf("%s:%d: %s", e.p.InputName(), e.line, e.error.Error())
}

func (p *Parser) errorf(line int, format string, a ...interface{}) {
	p.errors = append(p.errors,
		&parseError{
			p:     p,
			error: fmt.Errorf(format, a...),
			line:  line,
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
		case token.Database:
			// Not yet implemented
			return nil
		case token.Schema:
			return p.parseCreateSchemaStatement()
		case token.Table:
			return p.parseCreateTableStatement()
		case token.Unique, token.Index:
			return p.parseCreateIndexStatement()
		case token.Sequence:
			return p.parseCreateSequenceStatement()
		case token.Extension:
			// Not yet implemented
			return nil
		case token.Function:
			// Not yet implemented
			return nil
		case token.View:
			// Not yet implemented
			return nil
		case token.Operator:
			// Not yet implemented
			return nil
		case token.Trigger:
			// Not yet implemented
			return nil
		case token.Role:
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
		case token.Table:
			return p.parseAlterTableStatement()
		case token.Sequence:
			return p.parseAlterSequenceStatement()
		}
	case token.Grant:
		// Not yet implemented
		return nil
	case token.Revoke:
		// Not yet implemented
		return nil
	case token.Set:
		return p.parseSetStatement()
	case token.Select:
		// Not yet implemented
		return nil
	case token.Comment:
		// Not yet implemented
		return nil
	case token.BackslashConnect:
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
	tableName := p.parseTableName()
	if tableName == nil {
		return nil
	}
	createTableStatement.TableName = tableName

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

// "table_name" | "schema_name"."table_name"
func (p *Parser) parseTableName() *ast.TableName {
	var tableName ast.TableName

	identifier := p.parseIdentifier()
	if identifier == nil {
		return nil
	}
	if p.peekToken.Type != token.Dot {
		// Case: CREATE TABLE "table_name" ( ...
		tableName.TableIdentifier = identifier
	} else {
		// Case: CREATE TABLE "schema_name"."table_name" ( ...
		tableName.SchemaIdentifier = identifier
		p.advance()
		p.advance()
		identifier := p.parseIdentifier()
		if identifier == nil {
			return nil
		}
		tableName.TableIdentifier = identifier
	}
	return &tableName
}

// "sequence_name" | "schema_name"."sequence_name"
func (p *Parser) parseSequenceName() *ast.SequenceName {
	var sequenceName ast.SequenceName

	identifier := p.parseIdentifier()
	if identifier == nil {
		return nil
	}
	if p.peekToken.Type != token.Dot {
		// Case: CREATE SEQUENCE "sequence_name" ( ...
		sequenceName.SequenceIdentifier = identifier
	} else {
		// Case: CREATE SEQUENCE "schema_name"."sequence_name" ( ...
		sequenceName.SchemaIdentifier = identifier
		p.advance()
		p.advance()
		identifier := p.parseIdentifier()
		if identifier == nil {
			return nil
		}
		sequenceName.SequenceIdentifier = identifier
	}
	return &sequenceName
}

func (p *Parser) parseIdentifierAsExpression() ast.Expression {
	return p.parseIdentifier()
}

func (p *Parser) isIdentifier() bool {
	switch {
	case p.token.Type == token.Identifier:
		return true
	case p.token.IsKeyword() && !p.token.IsReserved():
		return true
	default:
		return false
	}
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
	case p.token.IsKeyword() && !p.token.IsReserved():
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
	case token.Smallint:
		return &ast.DataTypeSmallint{p.token}
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
		switch p.peekToken.Type {
		case token.With:
			// timestamp with time zone
			p.advance()
			if ok := p.expectPeek(token.Time) && p.expectPeek(token.Zone); !ok {
				return nil
			}
			dataTypeTimestamp.WithTimeZone = true
		case token.Without:
			p.advance()
			if ok := p.expectPeek(token.Time) && p.expectPeek(token.Zone); !ok {
				return nil
			}
			dataTypeTimestamp.WithTimeZone = false
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
	case token.Uuid:
		return &ast.DataTypeUuid{}
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
		case `"uuid"`:
			return &ast.DataTypeUuid{}
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

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expr := &ast.CallExpression{
		Token:     p.token,
		Function:  function,
		Arguments: p.parseCallArguments(),
	}
	return expr
}

func (p *Parser) parseCallArguments() []ast.Expression {
	var args []ast.Expression

	if p.peekToken.Type == token.RParen {
		p.advance()
		return args
	}

	p.advance()
	args = append(args, p.parseExpression(precedenceLowest))

	for p.peekToken.Type == token.Comma {
		p.advance()
		p.advance()
		args = append(args, p.parseExpression(precedenceLowest))
	}

	if !p.expectPeek(token.RParen) {
		return nil
	}

	return args
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

func (p *Parser) parseNull() ast.Expression {
	return &ast.NullLiteral{Token: p.token}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.advance()
	expr := p.parseExpression(precedenceLowest)
	if !p.expectPeek(token.RParen) {
		return nil
	}
	return &ast.GroupedExpression{Expression: expr}
}

func (p *Parser) parseCreateIndexStatement() ast.Statement {
	createIndexStatement := &ast.CreateIndexStatement{}

	if p.peekToken.Type == token.Unique {
		createIndexStatement.UniqueIndex = true
		p.advance()
	}

	if !p.expectPeek(token.Index) {
		return nil
	}

	if p.peekToken.Type == token.Concurrently {
		createIndexStatement.Concurrently = true
		p.advance()
	}

	if p.peekToken.Type == token.If {
		createIndexStatement.Concurrently = true
		p.advance()
		if !p.expectPeek(token.Not) {
			return nil
		}
		if !p.expectPeek(token.Exists) {
			return nil
		}
	}

	p.advance()
	identifier := p.parseIdentifier()
	if identifier == nil {
		return nil
	}
	createIndexStatement.Name = identifier

	if !p.expectPeek(token.On) {
		return nil
	}

	p.advance()
	tableName := p.parseTableName()
	if tableName == nil {
		return nil
	}
	createIndexStatement.TableName = tableName

	if p.peekToken.Type == token.Using {
		p.advance()
		p.advance()
		identifier := p.parseIdentifier()
		if identifier == nil {
			return nil
		}
		createIndexStatement.UsingMethod = identifier
	}

	p.advance()
	indexTargets := p.parseIndexTargets()
	if !p.expectPeek(token.RParen) {
		return nil
	}
	createIndexStatement.IndexTargets = indexTargets
	p.advance()

	return createIndexStatement
}

// CREATE SEQUENCE "users_id_seq"
//     START WITH 1
//     INCREMENT BY 1
//     NO MINVALUE
//     NO MAXVALUE
//     CACHE 1;
func (p *Parser) parseCreateSequenceStatement() ast.Statement {
	createSequenceStatement := &ast.CreateSequenceStatement{}

	if !p.expectPeek(token.Sequence) {
		return nil
	}

	p.advance()
	identifier := p.parseIdentifier()
	if identifier == nil {
		return nil
	}
	createSequenceStatement.Name = identifier

	for p.peekToken.Type != token.Semicolon && p.peekToken.Type != token.EOF {
		p.advance()
		p.parseCreateSequenceOption(createSequenceStatement)
	}

	return createSequenceStatement
}

func (p *Parser) parseCreateSequenceOption(createSequenceStatement *ast.CreateSequenceStatement) {
	switch p.token.Type {
	case token.Start:
		if p.peekToken.Type == token.With {
			p.advance()
		}
		p.advance()
		createSequenceStatement.StartWith = p.parseNumberLiteral()
	case token.Increment:
		if !p.expectPeek(token.By) {
			return
		}
		p.advance()
		createSequenceStatement.IncrementBy = p.parseNumberLiteral()
	case token.No:
		p.advance()
		switch p.token.Type {
		case token.Maxvalue:
			createSequenceStatement.NoMaxvalue = true
		case token.Minvalue:
			createSequenceStatement.NoMinvalue = true
		default:
			p.errorf(p.token.Line, "expected MAXVALUE or MINVALUE, found %s", p.token.Literal)
			return
		}
	case token.Cache:
		p.advance()
		createSequenceStatement.Cache = p.parseNumberLiteral()
	}
}

// ( target [ASC|DESC], ... )
func (p *Parser) parseIndexTargets() []*ast.IndexTarget {
	var indexTargets []*ast.IndexTarget
	p.advance()
	if p.token.Type == token.RParen {
		return indexTargets
	}
	for {
		if p.isIdentifier() {
			identifier := p.parseIdentifier()
			indexTarget := &ast.IndexTarget{
				Node: identifier,
			}
			switch p.peekToken.Type {
			case token.Asc:
				p.advance()
			case token.Desc:
				indexTarget.IsDesc = true
				p.advance()
			}
			indexTargets = append(indexTargets, indexTarget)
		} else {
			expr := p.parseExpression(precedenceLowest)
			indexTargets = append(indexTargets, &ast.IndexTarget{Node: expr})
		}
		if p.peekToken.Type == token.RParen {
			return indexTargets
		}
		if !p.expectPeek(token.Comma) {
			return indexTargets
		}
		p.advance()
	}
}

func (p *Parser) parseAlterSequenceStatement() ast.Statement {
	alterSequenceStatement := &ast.AlterSequenceStatement{}

	if !p.expectPeek(token.Sequence) {
		return nil
	}

	p.advance()
	sequenceName := p.parseSequenceName()
	if sequenceName == nil {
		return nil
	}
	alterSequenceStatement.Name = sequenceName

	if !p.expectPeek(token.Owned) {
		return nil
	}
	if !p.expectPeek(token.By) {
		return nil
	}

	p.advance()
	ownedByID1 := p.parseIdentifier()
	if ownedByID1 == nil {
		return nil
	}

	if !p.expectPeek(token.Dot) {
		return nil
	}

	p.advance()
	ownedByID2 := p.parseIdentifier()
	if ownedByID2 == nil {
		return nil
	}

	if p.peekToken.Type == token.Dot {
		p.advance()
		p.advance()
		ownedByColumn := p.parseIdentifier()
		if ownedByColumn == nil {
			return nil
		}
		alterSequenceStatement.OwnedByColumn = ownedByColumn
		alterSequenceStatement.OwnedByTable = &ast.TableName{
			SchemaIdentifier: ownedByID1,
			TableIdentifier:  ownedByID2,
		}
	} else {
		alterSequenceStatement.OwnedByColumn = ownedByID2
		alterSequenceStatement.OwnedByTable = &ast.TableName{
			TableIdentifier: ownedByID1,
		}
	}

	return alterSequenceStatement
}

// ALTER TABLE ONLY users ADD CONSTRAINT users_name_key UNIQUE (name);
// ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);
func (p *Parser) parseAlterTableStatement() ast.Statement {
	alterTableStatement := &ast.AlterTableStatement{}

	if !p.expectPeek(token.Table) {
		return nil
	}

	if p.peekToken.Type == token.Only {
		alterTableStatement.Only = true
		p.advance()
	}

	p.advance()
	tableName := p.parseTableName()
	if tableName == nil {
		return nil
	}
	alterTableStatement.Name = tableName

	switch p.peekToken.Type {
	case token.Add:
		p.advance()
		if p.peekToken.Type == token.Constraint {
			p.advance()
			tableConstraint := p.parseTableConstraint()
			if tableConstraint == nil {
				return nil
			}
			alterTableStatement.Actions = append(alterTableStatement.Actions, tableConstraint)
			return alterTableStatement
		} else {
			p.errorf(p.peekToken.Line, "expected %s, found %s", token.Constraint, p.peekToken.Literal)
			return nil
		}
	case token.Alter:
		p.advance()
		if p.peekToken.Type == token.Column {
			p.advance()
		}
		p.advance()

		column := p.parseIdentifier()
		if column == nil {
			return nil
		}
		p.advance()
		if p.token.Type == token.Set && p.peekToken.Type == token.Default {
			p.advance()
			p.advance()
			expr := p.parseExpression(precedenceLowest)
			if expr == nil {
				return nil
			}
			alterColumnSetDefault := &ast.AlterColumnSetDefault{
				Column: column,
				Expr:   expr,
			}
			alterTableStatement.Actions = append(alterTableStatement.Actions, alterColumnSetDefault)
			return alterTableStatement
		}
	}

	return alterTableStatement
}

// [ CONSTRAINT constraint_name ]
// { CHECK ( expression ) [ NO INHERIT ] |
//   UNIQUE ( column_name [, ... ] ) index_parameters |
//   PRIMARY KEY ( column_name [, ... ] ) index_parameters |
//   EXCLUDE [ USING index_method ] ( exclude_element WITH operator [, ... ] ) index_parameters [ WHERE ( predicate ) ] |
//   FOREIGN KEY ( column_name [, ... ] ) REFERENCES reftable [ ( refcolumn [, ... ] ) ]
//     [ MATCH FULL | MATCH PARTIAL | MATCH SIMPLE ] [ ON DELETE action ] [ ON UPDATE action ] }
// [ DEFERRABLE | NOT DEFERRABLE ] [ INITIALLY DEFERRED | INITIALLY IMMEDIATE ]
func (p *Parser) parseTableConstraint() ast.Node {
	tableConstraint := &ast.TableConstraint{}
	if p.token.Type == token.Constraint {
		p.advance()
		identifier := p.parseIdentifier()
		tableConstraint.Name = identifier
	}
	if p.peekToken.Type == token.Unique {
		tableConstraint.Unique = true
		p.advance()
		p.advance()
		columnList := p.parseColumnList()
		if columnList == nil {
			return nil
		}
		tableConstraint.ColumnList = columnList
	}
	if p.peekToken.Type == token.Primary {
		p.advance()
		if !p.expectPeek(token.Key) {
			return nil
		}
		tableConstraint.PrimaryKey = true
		p.advance()
		columnList := p.parseColumnList()
		if columnList == nil {
			return nil
		}
		tableConstraint.ColumnList = columnList
	}
	return tableConstraint
}

func (p *Parser) parseColumnList() *ast.ColumnList {
	columnList := &ast.ColumnList{}
	_, ok := p.expect(token.LParen)
	if !ok {
		return nil
	}
	for p.token.Type != token.RParen {
		column := p.parseIdentifier()
		if column != nil {
			columnList.ColumnNames = append(columnList.ColumnNames, column)
		}
		p.advance()
		if p.token.Type == token.Comma {
			p.advance()
		}
	}
	return columnList
}

// SET name = { value | 'value' | DEFAULT }
func (p *Parser) parseSetStatement() ast.Statement {
	setStatement := &ast.SetStatement{}

	if !p.expectPeek(token.Identifier) {
		return nil
	}
	setStatement.Name = p.parseIdentifier()

	if setStatement.Name.Value != "search_path" {
		// `SET search_path` is only implemented.
		return nil
	}

	if !p.expectPeek(token.Equal) {
		return nil
	}

	p.advance()
	for {
		expr := p.parseExpression(precedenceLowest)
		if expr == nil {
			return nil
		}
		setStatement.Values = append(setStatement.Values, expr)

		if p.peekToken.Type != token.Comma {
			break
		}
		p.advance()
		p.advance()
	}

	return setStatement
}
