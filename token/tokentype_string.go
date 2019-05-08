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
	_ = x[Dot-9]
	_ = x[Semicolon-10]
	_ = x[Comma-11]
	_ = x[LParen-12]
	_ = x[RParen-13]
	_ = x[Equal-14]
	_ = x[Plus-15]
	_ = x[Minus-16]
	_ = x[Asterisk-17]
	_ = x[Slash-18]
	_ = x[Typecast-19]
	_ = x[Add-20]
	_ = x[Alter-21]
	_ = x[Asc-22]
	_ = x[BackslashConnect-23]
	_ = x[By-24]
	_ = x[Cache-25]
	_ = x[Column-26]
	_ = x[Concurrently-27]
	_ = x[Constraint-28]
	_ = x[Create-29]
	_ = x[Database-30]
	_ = x[Default-31]
	_ = x[Desc-32]
	_ = x[Exists-33]
	_ = x[Extension-34]
	_ = x[False-35]
	_ = x[Function-36]
	_ = x[Grant-37]
	_ = x[If-38]
	_ = x[Increment-39]
	_ = x[Index-40]
	_ = x[Insert-41]
	_ = x[Is-42]
	_ = x[Maxvalue-43]
	_ = x[Minvalue-44]
	_ = x[No-45]
	_ = x[Not-46]
	_ = x[Null-47]
	_ = x[On-48]
	_ = x[Only-49]
	_ = x[Operator-50]
	_ = x[Owner-51]
	_ = x[Revoke-52]
	_ = x[Role-53]
	_ = x[Schema-54]
	_ = x[Select-55]
	_ = x[Sequence-56]
	_ = x[Set-57]
	_ = x[Start-58]
	_ = x[Table-59]
	_ = x[To-60]
	_ = x[Trigger-61]
	_ = x[True-62]
	_ = x[Unique-63]
	_ = x[Update-64]
	_ = x[Using-65]
	_ = x[Varying-66]
	_ = x[View-67]
	_ = x[With-68]
	_ = x[Zone-69]
	_ = x[Bigint-70]
	_ = x[Bigserial-71]
	_ = x[Boolean-72]
	_ = x[Bytea-73]
	_ = x[Character-74]
	_ = x[Date-75]
	_ = x[Integer-76]
	_ = x[Jsonb-77]
	_ = x[Numeric-78]
	_ = x[Serial-79]
	_ = x[Text-80]
	_ = x[Timestamp-81]
	_ = x[Time-82]
	_ = x[Tsvector-83]
}

const _TokenType_name = "IllegalEOFSpaceCommentCommentBlockCommentLineIdentifierStringNumberDotSemicolonCommaLParenRParenEqualPlusMinusAsteriskSlashTypecastAddAlterAscBackslashConnectByCacheColumnConcurrentlyConstraintCreateDatabaseDefaultDescExistsExtensionFalseFunctionGrantIfIncrementIndexInsertIsMaxvalueMinvalueNoNotNullOnOnlyOperatorOwnerRevokeRoleSchemaSelectSequenceSetStartTableToTriggerTrueUniqueUpdateUsingVaryingViewWithZoneBigintBigserialBooleanByteaCharacterDateIntegerJsonbNumericSerialTextTimestampTimeTsvector"

var _TokenType_index = [...]uint16{0, 7, 10, 15, 22, 34, 45, 55, 61, 67, 70, 79, 84, 90, 96, 101, 105, 110, 118, 123, 131, 134, 139, 142, 158, 160, 165, 171, 183, 193, 199, 207, 214, 218, 224, 233, 238, 246, 251, 253, 262, 267, 273, 275, 283, 291, 293, 296, 300, 302, 306, 314, 319, 325, 329, 335, 341, 349, 352, 357, 362, 364, 371, 375, 381, 387, 392, 399, 403, 407, 411, 417, 426, 433, 438, 447, 451, 458, 463, 470, 476, 480, 489, 493, 501}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
