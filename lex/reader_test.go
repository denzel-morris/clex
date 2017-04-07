package lex

import (
	"errors"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestPeekRuneBeforeReadRune(t *testing.T) {
	str := "abc"
	rd := newReader(str)
	expected, _ := utf8.DecodeRuneInString(str[1:])

	rd.ReadRune()

	r := rd.PeekRune()
	if r != expected {
		t.Error("Expected", expected, "on PeekRune, got", r)
	}

	r = rd.ReadRune()
	if r != expected {
		t.Error("Expected", expected, "on ReadRune, got", r)
	}
}

func TestPeekOnEOF(t *testing.T) {
	rd := newReader("a")

	rd.ReadRune()

	r := rd.PeekRune()
	expectEOFRune(t, r, "peek")

	r = rd.PeekRune()
	expectEOFRune(t, r, "repeated peek")

	r = rd.ReadRune()
	expectEOFRune(t, r, "read")
}

func TestErrorOnUnderlyingScannerError(t *testing.T) {
	rd := NewReader(mockErrorRuneScanner{})

	r := rd.ReadRune()
	expectErrorRune(t, r, "read")

	r = rd.PeekRune()
	expectErrorRune(t, r, "peek")

	err := rd.Err()
	expectError(t, err, "after underlying scanner error")
}

type mockErrorRuneScanner struct{}

func (rs mockErrorRuneScanner) ReadRune() (r rune, size int, err error) {
	return 0, 0, errors.New("")
}
func (rs mockErrorRuneScanner) UnreadRune() error {
	return errors.New("")
}

func newReader(contents string) Reader {
	return NewReader(strings.NewReader(contents))
}

func expectError(t *testing.T, err error, on string) {
	if err == nil {
		t.Error("Expected error", on)
	}
}

func expectEOFRune(t *testing.T, r rune, on string) {
	if r != runeEOF {
		t.Error("Expected runeEOF on", on, "got", r)
	}
}

func expectErrorRune(t *testing.T, r rune, on string) {
	if r != runeError {
		t.Error("Expected runeError on", on, "got", r)
	}
}
