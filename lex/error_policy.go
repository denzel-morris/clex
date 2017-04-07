package lex

type ErrorPolicy interface {
	ReportError(message string, line string, position Position)
}
