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
	_ = x[LBracket-14]
	_ = x[RBracket-15]
	_ = x[Equal-16]
	_ = x[Plus-17]
	_ = x[Minus-18]
	_ = x[Asterisk-19]
	_ = x[Slash-20]
	_ = x[Typecast-21]
	_ = x[Add-22]
	_ = x[Alter-23]
	_ = x[Asc-24]
	_ = x[BackslashConnect-25]
	_ = x[By-26]
	_ = x[Cache-27]
	_ = x[Column-28]
	_ = x[Concurrently-29]
	_ = x[Constraint-30]
	_ = x[Create-31]
	_ = x[Database-32]
	_ = x[Default-33]
	_ = x[Desc-34]
	_ = x[Exists-35]
	_ = x[Extension-36]
	_ = x[False-37]
	_ = x[Function-38]
	_ = x[Grant-39]
	_ = x[If-40]
	_ = x[Increment-41]
	_ = x[Index-42]
	_ = x[Insert-43]
	_ = x[Is-44]
	_ = x[Key-45]
	_ = x[Maxvalue-46]
	_ = x[Minvalue-47]
	_ = x[No-48]
	_ = x[Not-49]
	_ = x[Null-50]
	_ = x[On-51]
	_ = x[Only-52]
	_ = x[Operator-53]
	_ = x[Owned-54]
	_ = x[Owner-55]
	_ = x[Primary-56]
	_ = x[Revoke-57]
	_ = x[Role-58]
	_ = x[Schema-59]
	_ = x[Select-60]
	_ = x[Sequence-61]
	_ = x[Set-62]
	_ = x[Start-63]
	_ = x[Table-64]
	_ = x[To-65]
	_ = x[Trigger-66]
	_ = x[True-67]
	_ = x[Unique-68]
	_ = x[Update-69]
	_ = x[Using-70]
	_ = x[Varying-71]
	_ = x[View-72]
	_ = x[With-73]
	_ = x[Without-74]
	_ = x[Zone-75]
	_ = x[Bigint-76]
	_ = x[Smallint-77]
	_ = x[Bigserial-78]
	_ = x[Boolean-79]
	_ = x[Bytea-80]
	_ = x[Character-81]
	_ = x[Date-82]
	_ = x[Integer-83]
	_ = x[Jsonb-84]
	_ = x[Numeric-85]
	_ = x[Serial-86]
	_ = x[Text-87]
	_ = x[Timestamp-88]
	_ = x[Time-89]
	_ = x[Tsvector-90]
	_ = x[Uuid-91]
}

const _TokenType_name = "IllegalEOFSpaceCommentCommentBlockCommentLineIdentifierStringNumberDotSemicolonCommaLParenRParenLBracketRBracketEqualPlusMinusAsteriskSlashTypecastAddAlterAscBackslashConnectByCacheColumnConcurrentlyConstraintCreateDatabaseDefaultDescExistsExtensionFalseFunctionGrantIfIncrementIndexInsertIsKeyMaxvalueMinvalueNoNotNullOnOnlyOperatorOwnedOwnerPrimaryRevokeRoleSchemaSelectSequenceSetStartTableToTriggerTrueUniqueUpdateUsingVaryingViewWithWithoutZoneBigintSmallintBigserialBooleanByteaCharacterDateIntegerJsonbNumericSerialTextTimestampTimeTsvectorUuid"

var _TokenType_index = [...]uint16{0, 7, 10, 15, 22, 34, 45, 55, 61, 67, 70, 79, 84, 90, 96, 104, 112, 117, 121, 126, 134, 139, 147, 150, 155, 158, 174, 176, 181, 187, 199, 209, 215, 223, 230, 234, 240, 249, 254, 262, 267, 269, 278, 283, 289, 291, 294, 302, 310, 312, 315, 319, 321, 325, 333, 338, 343, 350, 356, 360, 366, 372, 380, 383, 388, 393, 395, 402, 406, 412, 418, 423, 430, 434, 438, 445, 449, 455, 463, 472, 479, 484, 493, 497, 504, 509, 516, 522, 526, 535, 539, 547, 551}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
