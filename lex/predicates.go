package lex

func startsEOF(r rune) bool { return r == runeEOF }

func startsIdentifier(r rune) bool {
	return isNonDigit(r) || startsEscape(r)
}

func startsEscape(r rune) bool { return r == '\\' }

func startsNumericConstant(r rune) bool {
	return isDecimalDigit(r) || isDecimalPoint(r)
}

func startsDecimalConstant(r rune) bool { return r >= '1' && r <= '9' }
func startsOctalConstant(r rune) bool   { return r == '0' }
func startsHexConstant(r rune) bool     { return r == '0' }

func startsCharLiteral(r rune) bool   { return r == '\'' }
func startsStringLiteral(r rune) bool { return r == '"' }
func startsWideLiteral(r rune) bool {
	switch r {
	case 'L', 'U', 'u':
		return true
	default:
		return false
	}
}

func startsComment(r rune) bool       { return r == '/' }
func commentIsSingleLine(r rune) bool { return r == '/' }
func commentIsMultiLine(r rune) bool  { return r == '*' }

func startsWhiteSpace(r rune) bool { return isWhitespace(r) }

func startsIntegerSuffix(r rune) bool {
	switch r {
	case 'u', 'U', 'l', 'L':
		return true
	default:
		return false
	}
}

func startsFloatingSuffix(r rune) bool {
	return isDecimalPoint(r) || startsExponentPart(r)
}

func startsExponentPart(r rune) bool {
	return r == 'e' || r == 'E' || r == 'p' || r == 'P'
}
