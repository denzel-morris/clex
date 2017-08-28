package lex

import (
	"io"

	"github.com/denzel-morris/clex/lex/container"
)

type lookaheadReader struct {
	reader
	readRunes, unreadRunes *container.RingBuffer
}

func NewLookaheadReader(scanner io.RuneScanner, lookahead uint64) Reader {
	return &lookaheadReader{
		reader:      reader{scanner: scanner},
		readRunes:   container.NewRingBuffer(lookahead),
		unreadRunes: container.NewRingBuffer(lookahead),
	}
}

func (rd *lookaheadReader) PeekRune() rune {
	r := rd.ReadRune()
	rd.UnreadRune()
	return r
}

func (rd *lookaheadReader) ReadRune() rune {
	rv, err := rd.unreadRunes.Pop()
	if err == nil {
		rd.readRunes.Push(rv)
		return rv.(rune)
	}

	r := rd.reader.ReadRune()
	if r != runeError {
		rd.readRunes.Push(r)
	}
	return r
}

func (rd *lookaheadReader) UnreadRune() {
	r, err := rd.readRunes.Pop()
	if err != nil {
		return
	}

	rd.unreadRunes.Push(r)
}
