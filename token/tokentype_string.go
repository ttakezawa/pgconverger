// Code generated by "stringer -type=TokenType"; DO NOT EDIT.

package token

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Illegal-0]
	_ = x[EOF-1]
	_ = x[Space-2]
	_ = x[Comment-3]
	_ = x[CommentBlock-4]
	_ = x[CommentLine-5]
	_ = x[Identifier-6]
	_ = x[String-7]
	_ = x[Number-8]
	_ = x[Semicolon-9]
	_ = x[Comma-10]
	_ = x[LParen-11]
	_ = x[RParen-12]
	_ = x[Plus-13]
	_ = x[Minus-14]
	_ = x[Asterisk-15]
	_ = x[Slash-16]
	_ = x[Typecast-17]
	_ = x[Add-18]
	_ = x[Alter-19]
	_ = x[By-20]
	_ = x[Cache-21]
	_ = x[Column-22]
	_ = x[Concurrently-23]
	_ = x[Constraint-24]
	_ = x[Create-25]
	_ = x[Default-26]
	_ = x[Exists-27]
	_ = x[Extension-28]
	_ = x[False-29]
	_ = x[Grant-30]
	_ = x[If-31]
	_ = x[Increment-32]
	_ = x[Index-33]
	_ = x[Insert-34]
	_ = x[Maxvalue-35]
	_ = x[Minvalue-36]
	_ = x[No-37]
	_ = x[Not-38]
	_ = x[Null-39]
	_ = x[On-40]
	_ = x[Only-41]
	_ = x[Owner-42]
	_ = x[Schema-43]
	_ = x[Select-44]
	_ = x[Sequence-45]
	_ = x[Set-46]
	_ = x[Start-47]
	_ = x[Table-48]
	_ = x[To-49]
	_ = x[True-50]
	_ = x[Unique-51]
	_ = x[Update-52]
	_ = x[Using-53]
	_ = x[Varying-54]
	_ = x[View-55]
	_ = x[With-56]
	_ = x[Zone-57]
	_ = x[Bigint-58]
	_ = x[Bigserial-59]
	_ = x[Boolean-60]
	_ = x[Bytea-61]
	_ = x[Character-62]
	_ = x[Date-63]
	_ = x[Integer-64]
	_ = x[Jsonb-65]
	_ = x[Numeric-66]
	_ = x[Serial-67]
	_ = x[Text-68]
	_ = x[Timestamp-69]
	_ = x[Time-70]
	_ = x[Tsvector-71]
}

const _TokenType_name = "IllegalEOFSpaceCommentCommentBlockCommentLineIdentifierStringNumberSemicolonCommaLParenRParenPlusMinusAsteriskSlashTypecastAddAlterByCacheColumnConcurrentlyConstraintCreateDefaultExistsExtensionFalseGrantIfIncrementIndexInsertMaxvalueMinvalueNoNotNullOnOnlyOwnerSchemaSelectSequenceSetStartTableToTrueUniqueUpdateUsingVaryingViewWithZoneBigintBigserialBooleanByteaCharacterDateIntegerJsonbNumericSerialTextTimestampTimeTsvector"

var _TokenType_index = [...]uint16{0, 7, 10, 15, 22, 34, 45, 55, 61, 67, 76, 81, 87, 93, 97, 102, 110, 115, 123, 126, 131, 133, 138, 144, 156, 166, 172, 179, 185, 194, 199, 204, 206, 215, 220, 226, 234, 242, 244, 247, 251, 253, 257, 262, 268, 274, 282, 285, 290, 295, 297, 301, 307, 313, 318, 325, 329, 333, 337, 343, 352, 359, 364, 373, 377, 384, 389, 396, 402, 406, 415, 419, 427}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
