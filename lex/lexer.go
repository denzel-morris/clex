package lex

import (
	"bytes"
	"sort"

	"github.com/denzel-morris/clex/lex/lexemes"
)

type Lexer interface {
	Lex() ([]Lexeme, error)
	Next() (Lexeme, error)
}

type lexer struct {
	stream LineReader
	buf    *bytes.Buffer
	errors ErrorPolicy
}

func NewLexer(rd LineReader, policy ErrorPolicy) Lexer {
	return &lexer{
		stream: rd,
		buf:    new(bytes.Buffer),
		errors: policy,
	}
}

func (l *lexer) Lex() ([]Lexeme, error) {
	var lexemelist []Lexeme
	lexeme, err := l.Next()
	for ; err == nil && lexeme.IsNot(lexemes.EOF); lexeme, err = l.Next() {
		lexemelist = append(lexemelist, lexeme)
	}
	return lexemelist, err
}

func (l *lexer) Next() (Lexeme, error) {
	lexeme := l.next()
	l.buf.Reset()
	return lexeme, l.stream.Err()
}

func (l *lexer) next() Lexeme {
	typ := l.lex()
	return l.makeLexeme(typ)
}

func (l *lexer) lex() lexemes.Type {
	switch r := l.peek(); {
	case startsIdentifier(r), startsWideLiteral(r):
		return l.maybeKeyword(l.lexIdentifierOrLiteral())
	case startsNumericConstant(r):
		return l.lexNumericConstant()
	case startsStringLiteral(r):
		return l.lexStringLiteral()
	case startsCharLiteral(r):
		return l.lexCharLiteral()
	case startsComment(r):
		return l.lexCommentOrPunctuator()
	case startsWhiteSpace(r):
		return l.lexWhitespace()
	case startsEOF(r):
		return l.lexEOF()
	default:
		return l.lexPunctuator()
	}
}

func (l *lexer) lexIdentifierOrLiteral() lexemes.Type {
	l.consume(oneOf("LUu"))
	switch l.peek() {
	case '"':
		return l.lexStringLiteral()
	case '\'':
		return l.lexCharLiteral()
	case '8':
		l.consume(oneRune('8'))
		if l.lexStringLiteral() != lexemes.Invalid {
			return lexemes.StringLiteral
		}
		fallthrough
	default:
		return l.lexIdentifier()
	}
}

func (l *lexer) lexIdentifier() lexemes.Type {
	ok := l.consumeWhileDo(identifierChar, l.lookForUnicodeEscape)
	if !ok {
		return lexemes.Invalid
	}
	return lexemes.Identifier
}

func (l *lexer) maybeKeyword(typ lexemes.Type) lexemes.Type {
	if isKeyword(l.value()) {
		return lexemes.Keyword
	}
	return typ
}

func isKeyword(str string) bool {
	idx := sort.SearchStrings(keywords, str)
	return idx < len(keywords) && keywords[idx] == str
}

func (l *lexer) lexNumericConstant() lexemes.Type {
	typ := lexemes.IntegerConstant

	switch r := l.peek(); {
	case startsOctalConstant(r), startsHexConstant(r):
		typ = l.lexOctalOrHexConstant()
	case startsDecimalConstant(r):
		typ = l.lexDecimalConstant()
	}

	if typ != lexemes.Invalid {
		typ = l.lexNumericConstantSuffix()
	}

	if l.value() == "." {
		return l.lexPunctuator()
	}

	return typ
}

func (l *lexer) lexOctalOrHexConstant() lexemes.Type {
	l.consume(oneRune('0'))
	switch r := l.peek(); {
	case isOctalDigit(r):
		return l.lexOctalConstant()
	case r == 'x', r == 'X':
		return l.lexHexConstant()
	default:
		return lexemes.IntegerConstant
	}
}

func (l *lexer) lexDecimalConstant() lexemes.Type {
	l.consumeWhile(decimalDigit)
	return lexemes.IntegerConstant
}

func (l *lexer) lexOctalConstant() lexemes.Type {
	l.consumeWhile(octalDigit)
	return lexemes.IntegerConstant
}

func (l *lexer) lexHexConstant() lexemes.Type {
	l.consume(oneOf("xX"))
	ok := l.consumeAtLeastOne(hexDigit)
	if !ok {
		l.reportError("Hexadecimal constant must contain at least one digit")
		return lexemes.Invalid
	}
	return lexemes.IntegerConstant
}

func (l *lexer) lexNumericConstantSuffix() lexemes.Type {
	switch r := l.peek(); {
	case startsIntegerSuffix(r):
		l.lexIntegerSuffix()
		return lexemes.IntegerConstant
	case startsFloatingSuffix(r):
		return l.lexFloatingSuffix()
	default:
		return lexemes.IntegerConstant
	}
}

func (l *lexer) lexIntegerSuffix() {
	switch r, _ := l.consume(oneOf("uUlL")); r {
	case 'u', 'U':
		l.lexLongOrLongLongSuffix()
	case 'l':
		l.consume(oneRune('l'))
		l.lexUnsignedSuffix()
	case 'L':
		l.consume(oneRune('L'))
		l.lexUnsignedSuffix()
	}
}

func (l *lexer) lexFloatingSuffix() lexemes.Type {
	switch r := l.peek(); {
	case isDecimalPoint(r):
		l.consume(decimalPoint)
		l.consumeWhile(decimalDigit)
		fallthrough
	case startsExponentPart(r):
		if l.lexExponentPart() == lexemes.Invalid {
			return lexemes.Invalid
		}
	}
	l.consume(oneOf("fFlL"))
	return lexemes.FloatingConstant
}

func (l *lexer) lexExponentPart() lexemes.Type {
	if _, ok := l.consume(oneOf("eEpP")); !ok {
		return lexemes.FloatingConstant
	}
	l.consume(oneOf("+-"))
	if l.consumeAtLeastOne(decimalDigit) {
		return lexemes.FloatingConstant
	}
	l.reportError("Exponent must have at least one digit")
	return lexemes.Invalid
}

func (l *lexer) lexLongOrLongLongSuffix() {
	switch r, _ := l.consume(oneOf("lL")); r {
	case 'l':
		l.consume(oneRune('l'))
	case 'L':
		l.consume(oneRune('L'))
	}
}

func (l *lexer) lexUnsignedSuffix() {
	l.consume(oneOf("uU"))
}

func (l *lexer) lexStringLiteral() lexemes.Type {
	_, ok := l.consume(oneRune('"'))
	if !ok {
		return lexemes.Invalid
	}

	ok = l.consumeUntilDo(oneOf("\"\n"), l.lookForEscape)
	if !ok {
		return lexemes.Invalid
	}

	_, ok = l.consume(oneRune('"'))
	if !ok {
		l.reportError("Expected `\"` to end string literal after newline")
		return lexemes.Invalid
	}
	return lexemes.StringLiteral
}

func (l *lexer) lexCharLiteral() lexemes.Type {
	l.consume(oneRune('\''))

	ok := l.consumeUntilDo(oneOf("'\n"), l.lookForEscape)
	if !ok {
		return lexemes.Invalid
	}
	_, ok = l.consume(oneRune('\''))
	if !ok {
		l.reportError("Expected `'` to end character literal after newline")
		return lexemes.Invalid
	}
	return lexemes.CharLiteral
}

func (l *lexer) lexCommentOrPunctuator() lexemes.Type {
	l.consume(oneRune('/'))
	switch r := l.peek(); {
	case commentIsSingleLine(r):
		return l.lexSingleLineComment()
	case commentIsMultiLine(r):
		return l.lexMultiLineComment()
	default:
		return l.lexPunctuator()
	}
}

func (l *lexer) lexSingleLineComment() lexemes.Type {
	l.consume(oneRune('/'))
	l.consumeUntil(oneRune('\n'))
	return lexemes.Comment
}

func (l *lexer) lexMultiLineComment() lexemes.Type {
	l.consume(oneRune('*'))
	l.consumeWhileDo(any, l.lookForMultiLineCommentEnd)
	return lexemes.Comment
}

func (l *lexer) lexPunctuator() lexemes.Type {
	tok := l.value() + string(l.peek())
	for _, present := punctuatorToType[tok]; present; _, present = punctuatorToType[tok] {
		l.consume(any)
		tok += string(l.peek())
	}
	typ := punctuatorToType[l.value()]
	if typ == lexemes.Invalid {
		l.reportError("Unregonized character `" + l.value() + "`")
	}
	return typ
}

func (l *lexer) lexWhitespace() lexemes.Type {
	l.consumeWhile(whitespace)
	return lexemes.Whitespace
}

func (l *lexer) lexEOF() lexemes.Type {
	l.consume(any)
	return lexemes.EOF
}

func (l *lexer) lookForUnicodeEscape(r rune) (cont bool) {
	if startsEscape(r) {
		return l.consumeUnicodeEscape()
	}
	return true
}

func (l *lexer) lookForEscape(r rune) (cont bool) {
	if startsEscape(r) {
		return l.consumeEscape()
	}
	return true
}

func (l *lexer) lookForMultiLineCommentEnd(r rune) (cont bool) {
	if r == '*' && l.peek() == '/' {
		l.consume(oneRune('/'))
		return false
	}
	return true
}

func (l *lexer) consumeEscape() (ok bool) {
	switch r := l.peek(); {
	case isSimpleEscape(r):
		ok = l.consumeSimpleEscape()
	case isOctalDigit(r):
		ok = l.consumeOctalEscape()
	case r == 'x':
		ok = l.consumeHexEscape()
	case r == 'u' || r == 'U':
		ok = l.consumeUnicodeEscape()
	default:
		l.reportError("Unknown character `" + string(r) + "` escaped")
	}
	return ok
}

func (l *lexer) consumeUnicodeEscape() (ok bool) {
	switch r, _ := l.consume(oneOf("uU")); r {
	case 'u':
		ok = l.consumeN(hexDigits(4)) == 4
		if !ok {
			l.reportError("Expected 4 hexadecimal characters for universal character name")
		}
	case 'U':
		ok = l.consumeN(hexDigits(8)) == 8
		if !ok {
			l.reportError("Expected 8 hexadecimal characters for universal character name")
		}
	default:
		l.reportError("Expected universal character name starting with \\u or \\U")
	}
	return ok
}

func (l *lexer) consumeSimpleEscape() (ok bool) {
	_, ok = l.consume(simpleEscape)
	return ok
}

func (l *lexer) consumeOctalEscape() (ok bool) {
	return l.consumeOneUpTo(octalDigits(3))
}

func (l *lexer) consumeHexEscape() (ok bool) {
	l.consume(oneRune('x'))
	ok = l.consumeAtLeastOne(hexDigit)
	if !ok {
		l.reportError("Must provide at least one digit for hexadecimal escape")
	}
	return ok
}

func (l *lexer) consumeUntil(rc runeClass) (ok bool) {
	return l.consumeWhile(complement(rc))
}

func (l *lexer) consumeUntilDo(rc runeClass, body func(rune) bool) (ok bool) {
	return l.consumeWhileDo(complement(rc), body)
}

func (l *lexer) consumeWhile(rc runeClass) (ok bool) {
	return l.consumeWhileDo(rc, nil)
}

func (l *lexer) consumeAtLeastOne(rc runeClass) (ok bool) {
	once := false
	l.consumeWhileDo(rc, func(rune) bool { once = true; return true })
	return once
}

func (l *lexer) consumeWhileDo(rc runeClass, body func(rune) bool) (ok bool) {
	for r, ok := l.consume(rc); ok; r, ok = l.consume(rc) {
		if body != nil && !body(r) {
			return false
		}
	}
	return true
}

func (l *lexer) consumeOneUpTo(rc runeClassWithCount) (ok bool) {
	once := false
	for i, ok := 0, true; i < rc.count() && ok; i++ {
		_, ok = l.consume(rc)
		if ok {
			once = true
		}
	}
	return once
}

func (l *lexer) consumeN(rc runeClassWithCount) (consumed int) {
	for i, ok := 0, true; i < rc.count() && ok; i++ {
		_, ok = l.consume(rc)
		if ok {
			consumed++
		}
	}
	return consumed
}

func (l *lexer) consume(rc runeClass) (r rune, ok bool) {
	r = l.stream.ReadRune()
	if r < 0 {
		return r, false
	}

	includes := rc.has(r)
	switch includes {
	case true:
		l.buf.WriteRune(r)
	case false:
		l.stream.UnreadRune()
	}
	return r, includes
}

func (l *lexer) peek() rune    { return l.stream.PeekRune() }
func (l *lexer) value() string { return l.buf.String() }

func (l *lexer) makeLexeme(typ lexemes.Type) Lexeme {
	return makeLexeme(typ, l.value())
}

func (l *lexer) reportError(message string) {
	line, position := l.stream.Line(), l.stream.Position()
	l.errors.ReportError(message, line, position)
}
