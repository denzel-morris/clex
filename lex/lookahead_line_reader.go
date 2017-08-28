package lex

import (
	"bytes"

	"github.com/denzel-morris/clex/lex/container"
)

type lookaheadLineReader struct {
	position                       Position
	readPositions, unreadPositions *container.RingBuffer
	lineBuf                        *bytes.Buffer
	reader                         Reader
}

func NewLookaheadLineReader(rd Reader, lookahead uint64) LineReader {
	return &lookaheadLineReader{
		reader:          rd,
		position:        Position{Line: 1, Column: 1},
		lineBuf:         new(bytes.Buffer),
		readPositions:   container.NewRingBuffer(lookahead),
		unreadPositions: container.NewRingBuffer(lookahead),
	}
}

func (rd *lookaheadLineReader) Position() Position {
	return rd.position
}

func (rd *lookaheadLineReader) Line() string {
	return rd.lineBuf.String()
}

func (rd *lookaheadLineReader) PeekRune() rune {
	return rd.reader.PeekRune()
}

func (rd *lookaheadLineReader) ReadRune() rune {
	pv, err := rd.unreadPositions.Pop()
	if err == nil {
		rd.readPositions.Push(pv)
		r := rd.reader.ReadRune()
		rd.lineBuf.WriteRune(r)
		return r
	}

	rd.readPositions.Push(rd.position)
	r := rd.reader.ReadRune()
	switch {
	case r == '\n':
		rd.lineBuf.Reset()
		rd.position.Column = 1
		rd.position.Line++
	case r > 0:
		rd.lineBuf.WriteRune(r)
		rd.position.Column++
	}
	return r
}

func (rd *lookaheadLineReader) UnreadRune() {
	if rd.lineBuf.Len() == 0 {
		return
	}
	p, err := rd.readPositions.Pop()
	if err != nil {
		return
	}

	rd.reader.UnreadRune()
	rd.lineBuf.Truncate(rd.lineBuf.Len() - 1)

	rd.unreadPositions.Push(p)
}

func (rd *lookaheadLineReader) Err() error {
	return rd.reader.Err()
}
