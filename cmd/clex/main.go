package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/denzel-morris/clex/lex"
)

func main() {
	input, err := os.Open(os.Args[1])
	panicErr(err)
	defer input.Close()

	output, err := os.OpenFile(os.Args[2], os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	panicErr(err)
	defer output.Close()

	lexer := lex.NewLexer(
		lex.NewLineReader(lex.NewReader(bufio.NewReader(input))),
		&LogErrorPolicy{},
	)
	lexemelist, err := lexer.Lex()
	panicErr(err)

	for _, lexeme := range lexemelist {
		fmt.Fprintln(output, lexeme)
	}
}

type LogErrorPolicy struct{}

func (ep LogErrorPolicy) ReportError(message string, line string, position lex.Position) {
	log.Printf("error:%s: %s\n\t%s", position, message, line)
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
