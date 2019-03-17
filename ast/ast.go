package ast

import (
	"fmt"
	"strings"
)

type Node interface {
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type DataDefinition struct {
	StatementList []Statement
}

func (dataDefinition *DataDefinition) String() string {
	var b strings.Builder
	for _, statement := range dataDefinition.StatementList {
		b.WriteString(statement.String())
	}
	return b.String()
}

// CREATE [ [ GLOBAL | LOCAL ] { TEMPORARY | TEMP } | UNLOGGED ] TABLE [ IF NOT EXISTS ] table_name ( [
//   { column_name data_type [ COLLATE collation ] [ column_constraint [ ... ] ]
//     | table_constraint
//     | LIKE source_table [ like_option ... ] }
//     [, ... ]
// ] )
// [ INHERITS ( parent_table [, ... ] ) ]
// [ PARTITION BY { RANGE | LIST } ( { column_name | ( expression ) } [ COLLATE collation ] [ opclass ] [, ... ] ) ]
// [ WITH ( storage_parameter [= value] [, ... ] ) | WITH OIDS | WITHOUT OIDS ]
// [ ON COMMIT { PRESERVE ROWS | DELETE ROWS | DROP } ]
// [ TABLESPACE tablespace_name ]
type CreateTableStatement struct {
	TableName string
}

func (*CreateTableStatement) statementNode() {}
func (createTableStatement *CreateTableStatement) String() string {
	var b strings.Builder
	_, _ = fmt.Fprintf(&b, "CREATE TABLE '%s' ();", createTableStatement.TableName)
	return b.String()
}
