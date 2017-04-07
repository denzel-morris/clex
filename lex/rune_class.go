package lex

import (
	"strings"
	"unicode"
)

type runeClass interface {
	has(rune) bool
}

type runeClassWithCount interface {
	runeClass
	count() int
}

type runeClassFunc func(rune) bool

func (f runeClassFunc) has(r rune) bool {
	return f(r)
}

type complementedRuneClass struct {
	inner runeClass
}

func complement(rc runeClass) runeClass {
	return complementedRuneClass{rc}
}

func (rc complementedRuneClass) has(r rune) bool {
	return !rc.inner.has(r)
}

type decimalDigits int
type hexDigits int
type octalDigits int

func (d decimalDigits) has(r rune) bool { return isDecimalDigit(r) }
func (d decimalDigits) count() int      { return int(d) }
func (d hexDigits) has(r rune) bool     { return isHexDigit(r) }
func (d hexDigits) count() int          { return int(d) }
func (d octalDigits) has(r rune) bool   { return isOctalDigit(r) }
func (d octalDigits) count() int        { return int(d) }

type oneOf string

func (rc oneOf) has(r rune) bool {
	return strings.ContainsRune(string(rc), r)
}

type oneRune rune

func (rc oneRune) has(r rune) bool {
	return rune(rc) == r
}

var (
	any            runeClassFunc = isAny
	decimalDigit   runeClassFunc = isDecimalDigit
	hexDigit       runeClassFunc = isHexDigit
	octalDigit     runeClassFunc = isOctalDigit
	whitespace     runeClassFunc = isWhitespace
	decimalPoint   runeClassFunc = isDecimalPoint
	simpleEscape   runeClassFunc = isSimpleEscape
	identifierChar runeClassFunc = isIdentifierChar
)

func isAny(r rune) bool { return true }

func isDecimalDigit(r rune) bool {
	switch {
	case r >= '0' && r <= '9':
		return true
	default:
		return false
	}
}

func isHexDigit(r rune) bool {
	switch {
	case isDecimalDigit(r), r >= 'A' && r <= 'F', r >= 'a' && r <= 'f':
		return true
	default:
		return false
	}
}

func isOctalDigit(r rune) bool {
	switch {
	case r >= '0' && r <= '7':
		return true
	default:
		return false
	}
}

func isNonDigit(r rune) bool {
	switch {
	case r == '_' || unicode.IsLetter(r):
		return true
	default:
		return false
	}
}

func isWhitespace(r rune) bool {
	switch r {
	case ' ', '\f', '\n', '\r', '\t', '\v':
		return true
	default:
		return false
	}
}

func isDecimalPoint(r rune) bool { return r == '.' }

func isSimpleEscape(r rune) bool {
	switch r {
	case '\'', '"', '?', '\\', 'a', 'b', 'f', 'n', 'r', 't', 'v':
		return true
	default:
		return false
	}
}

func isIdentifierChar(r rune) bool {
	switch {
	case isNonDigit(r), isDecimalDigit(r), r == '\\':
		return true
	default:
		return false
	}
}
