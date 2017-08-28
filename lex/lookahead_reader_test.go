package lex

import (
	"strings"
	"testing"
)

func TestLookaheadReaderSupportsNRuneLookahead(t *testing.T) {
	str := "abc"
	rd := newLookaheadReader(str, 2)

	rd.ReadRune()
	rd.ReadRune()
	rd.UnreadRune()
	rd.UnreadRune()

	r := rd.ReadRune()
	if r != 'a' {
		t.Error("Expected a, got", r)
	}

	r = rd.ReadRune()
	if r != 'b' {
		t.Error("Expected a, got", r)
	}

	rd.ReadRune()
	rd.UnreadRune()
	rd.UnreadRune()

	r = rd.ReadRune()
	if r != 'b' {
		t.Error("Expected a, got", r)
	}

	r = rd.ReadRune()
	if r != 'c' {
		t.Error("Expected a, got", r)
	}
}

func newLookaheadReader(contents string, lookahead uint64) Reader {
	return NewLookaheadReader(strings.NewReader(contents), lookahead)
}
