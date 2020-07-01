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
	_ = x[Key-43]
	_ = x[Maxvalue-44]
	_ = x[Minvalue-45]
	_ = x[No-46]
	_ = x[Not-47]
	_ = x[Null-48]
	_ = x[On-49]
	_ = x[Only-50]
	_ = x[Operator-51]
	_ = x[Owned-52]
	_ = x[Owner-53]
	_ = x[Primary-54]
	_ = x[Revoke-55]
	_ = x[Role-56]
	_ = x[Schema-57]
	_ = x[Select-58]
	_ = x[Sequence-59]
	_ = x[Set-60]
	_ = x[Start-61]
	_ = x[Table-62]
	_ = x[To-63]
	_ = x[Trigger-64]
	_ = x[True-65]
	_ = x[Unique-66]
	_ = x[Update-67]
	_ = x[Using-68]
	_ = x[Varying-69]
	_ = x[View-70]
	_ = x[With-71]
	_ = x[Without-72]
	_ = x[Zone-73]
	_ = x[Bigint-74]
	_ = x[Smallint-75]
	_ = x[Bigserial-76]
	_ = x[Boolean-77]
	_ = x[Bytea-78]
	_ = x[Character-79]
	_ = x[Date-80]
	_ = x[Integer-81]
	_ = x[Jsonb-82]
	_ = x[Numeric-83]
	_ = x[Serial-84]
	_ = x[Text-85]
	_ = x[Timestamp-86]
	_ = x[Time-87]
	_ = x[Tsvector-88]
	_ = x[Uuid-89]
}

const _TokenType_name = "IllegalEOFSpaceCommentCommentBlockCommentLineIdentifierStringNumberDotSemicolonCommaLParenRParenEqualPlusMinusAsteriskSlashTypecastAddAlterAscBackslashConnectByCacheColumnConcurrentlyConstraintCreateDatabaseDefaultDescExistsExtensionFalseFunctionGrantIfIncrementIndexInsertIsKeyMaxvalueMinvalueNoNotNullOnOnlyOperatorOwnedOwnerPrimaryRevokeRoleSchemaSelectSequenceSetStartTableToTriggerTrueUniqueUpdateUsingVaryingViewWithWithoutZoneBigintSmallintBigserialBooleanByteaCharacterDateIntegerJsonbNumericSerialTextTimestampTimeTsvectorUuid"

var _TokenType_index = [...]uint16{0, 7, 10, 15, 22, 34, 45, 55, 61, 67, 70, 79, 84, 90, 96, 101, 105, 110, 118, 123, 131, 134, 139, 142, 158, 160, 165, 171, 183, 193, 199, 207, 214, 218, 224, 233, 238, 246, 251, 253, 262, 267, 273, 275, 278, 286, 294, 296, 299, 303, 305, 309, 317, 322, 327, 334, 340, 344, 350, 356, 364, 367, 372, 377, 379, 386, 390, 396, 402, 407, 414, 418, 422, 429, 433, 439, 447, 456, 463, 468, 477, 481, 488, 493, 500, 506, 510, 519, 523, 531, 535}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
