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
	_ = x[Equal-13]
	_ = x[Plus-14]
	_ = x[Minus-15]
	_ = x[Asterisk-16]
	_ = x[Slash-17]
	_ = x[Typecast-18]
	_ = x[Add-19]
	_ = x[Alter-20]
	_ = x[By-21]
	_ = x[Cache-22]
	_ = x[Column-23]
	_ = x[Concurrently-24]
	_ = x[Constraint-25]
	_ = x[Create-26]
	_ = x[Default-27]
	_ = x[Exists-28]
	_ = x[Extension-29]
	_ = x[False-30]
	_ = x[Grant-31]
	_ = x[If-32]
	_ = x[Increment-33]
	_ = x[Index-34]
	_ = x[Insert-35]
	_ = x[Is-36]
	_ = x[Maxvalue-37]
	_ = x[Minvalue-38]
	_ = x[No-39]
	_ = x[Not-40]
	_ = x[Null-41]
	_ = x[On-42]
	_ = x[Only-43]
	_ = x[Owner-44]
	_ = x[Schema-45]
	_ = x[Select-46]
	_ = x[Sequence-47]
	_ = x[Set-48]
	_ = x[Start-49]
	_ = x[Table-50]
	_ = x[To-51]
	_ = x[True-52]
	_ = x[Unique-53]
	_ = x[Update-54]
	_ = x[Using-55]
	_ = x[Varying-56]
	_ = x[View-57]
	_ = x[With-58]
	_ = x[Zone-59]
	_ = x[Bigint-60]
	_ = x[Bigserial-61]
	_ = x[Boolean-62]
	_ = x[Bytea-63]
	_ = x[Character-64]
	_ = x[Date-65]
	_ = x[Integer-66]
	_ = x[Jsonb-67]
	_ = x[Numeric-68]
	_ = x[Serial-69]
	_ = x[Text-70]
	_ = x[Timestamp-71]
	_ = x[Time-72]
	_ = x[Tsvector-73]
}

const _TokenType_name = "IllegalEOFSpaceCommentCommentBlockCommentLineIdentifierStringNumberSemicolonCommaLParenRParenEqualPlusMinusAsteriskSlashTypecastAddAlterByCacheColumnConcurrentlyConstraintCreateDefaultExistsExtensionFalseGrantIfIncrementIndexInsertIsMaxvalueMinvalueNoNotNullOnOnlyOwnerSchemaSelectSequenceSetStartTableToTrueUniqueUpdateUsingVaryingViewWithZoneBigintBigserialBooleanByteaCharacterDateIntegerJsonbNumericSerialTextTimestampTimeTsvector"

var _TokenType_index = [...]uint16{0, 7, 10, 15, 22, 34, 45, 55, 61, 67, 76, 81, 87, 93, 98, 102, 107, 115, 120, 128, 131, 136, 138, 143, 149, 161, 171, 177, 184, 190, 199, 204, 209, 211, 220, 225, 231, 233, 241, 249, 251, 254, 258, 260, 264, 269, 275, 281, 289, 292, 297, 302, 304, 308, 314, 320, 325, 332, 336, 340, 344, 350, 359, 366, 371, 380, 384, 391, 396, 403, 409, 413, 422, 426, 434}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
