package lex

import (
	"fmt"

	"github.com/DenzelMorris/clex/lex/lexemes"
)

type Lexeme struct {
	Type  lexemes.Type
	Value string
}

func (l Lexeme) String() string {
	return fmt.Sprintf("%s{%s}", l.Type, l.Value)
}

func makeLexeme(typ lexemes.Type, value string) Lexeme {
	return Lexeme{
		Type:  typ,
		Value: value,
	}
}

func (l Lexeme) IsNot(typ lexemes.Type) bool {
	return !l.Is(typ)
}

func (l Lexeme) Is(typ lexemes.Type) bool {
	return l.Type == typ
}
