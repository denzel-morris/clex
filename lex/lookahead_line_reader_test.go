package lex

import (
	"strings"
	"testing"
)

func TestLookaheadLineReaderPeekRune(t *testing.T) {
	rd := NewLookaheadLineReader(NewLookaheadReader(strings.NewReader("ab"), 4), 4)
	r := rd.PeekRune()
	if r != 'a' {
		t.Error("Expected 'a', got", r)
	}

	r = rd.PeekRune()
	if r != 'a' {
		t.Error("Expected 'a', got", r)
	}

	rd.ReadRune()

	r = rd.PeekRune()
	if r != 'b' {
		t.Error("Expected 'b', got", r)
	}

	r = rd.PeekRune()
	if r != 'b' {
		t.Error("Expected 'b', got", r)
	}
	r = rd.PeekRune()
	if r != 'b' {
		t.Error("Expected 'b', got", r)
	}

	rd.ReadRune()
	rd.UnreadRune()
	rd.ReadRune()
	rd.UnreadRune()
	rd.ReadRune()
	rd.UnreadRune()
	rd.ReadRune()
	rd.UnreadRune()
	rd.ReadRune()
	rd.UnreadRune()

	r = rd.PeekRune()
	if r != 'b' {
		t.Error("Expected 'b', got", r)
	}

	rd.ReadRune()

	r = rd.PeekRune()
	if r != runeEOF {
		t.Error("Expected runeEOF, got", r)
	}
}
