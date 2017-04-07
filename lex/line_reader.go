package lex

import "bytes"

type LineReader interface {
	Reader
	Position() Position
	Line() string
}

type lineReader struct {
	reader                 Reader
	position, lastPosition Position
	lineBuf                *bytes.Buffer
}

func NewLineReader(rd Reader) LineReader {
	return &lineReader{
		reader:       rd,
		position:     Position{Line: 1, Column: 1},
		lastPosition: Position{Line: 1, Column: 1},
		lineBuf:      new(bytes.Buffer),
	}
}

func (rd *lineReader) Position() Position {
	return rd.position
}

func (rd *lineReader) Line() string {
	return rd.lineBuf.String()
}

func (rd *lineReader) PeekRune() rune {
	return rd.reader.PeekRune()
}

func (rd *lineReader) ReadRune() rune {
	rd.savePosition()
	r := rd.reader.ReadRune()
	rd.updateLineAndPosition(r)
	return r
}

func (rd *lineReader) UnreadRune() {
	rd.restorePosition()
	rd.unreadRune()
	rd.reader.UnreadRune()
}

func (rd *lineReader) Err() error {
	return rd.reader.Err()
}

func (rd *lineReader) savePosition() {
	rd.lastPosition = rd.position
}

func (rd *lineReader) restorePosition() {
	rd.position = rd.lastPosition
}

func (rd *lineReader) updateLineAndPosition(r rune) {
	rd.updateLine(r)
	rd.updatePosition(r)
}

func (rd *lineReader) updateLine(r rune) {
	switch {
	case r > 0:
		rd.lineBuf.WriteRune(r)
	case r == '\n':
		rd.lineBuf.Reset()
	}
}

func (rd *lineReader) updatePosition(r rune) {
	switch {
	case r > 0:
		rd.position.Column++
	case r == '\n':
		rd.position.Column = 1
		rd.position.Line++
	}
}

func (rd *lineReader) unreadRune() {
	rd.lineBuf.Truncate(rd.lineBuf.Len() - 1)
}
