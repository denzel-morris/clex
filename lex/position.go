package lex

import "fmt"

type Position struct {
	Line, Column int
}

func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}
