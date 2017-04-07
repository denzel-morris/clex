package lex

import "io"

type Reader interface {
	PeekRune() rune
	ReadRune() rune
	UnreadRune()
	Err() error
}

type reader struct {
	scanner io.RuneScanner
	err     error
}

const (
	runeEOF   = -1
	runeError = -2
)

func NewReader(scanner io.RuneScanner) Reader {
	return &reader{scanner: scanner}
}

func (rd *reader) PeekRune() rune {
	r := rd.ReadRune()
	rd.UnreadRune()
	return r

}

func (rd *reader) ReadRune() rune {
	if rd.err != nil {
		return runeError
	}

	r, _, err := rd.scanner.ReadRune()

	switch {
	case err == io.EOF:
		return runeEOF
	case err != nil:
		rd.err = err
		return runeError
	}

	return r
}

func (rd *reader) UnreadRune() {
	rd.scanner.UnreadRune()
}

func (rd *reader) Err() error {
	err := rd.err
	rd.err = nil
	return err
}
