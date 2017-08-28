package container

import "testing"

func TestRingBufferWrapsProperly(t *testing.T) {
	rb := NewRingBuffer(2)
	rb.Push('a')
	rb.Push('b')
	rb.Push('c')

	r, _ := rb.Pop()
	if r != 'c' {
		t.Error("Expected to pop 'c' after pushing ['a', 'b', 'c'], got", r)
	}

	r, _ = rb.Pop()
	if r != 'b' {
		t.Error("Expected to pop 'b' after popping 'c', got", r)
	}
}

func TestRingBufferErrorsWhenEmpty(t *testing.T) {
	rb := NewRingBuffer(2)
	_, err := rb.Pop()
	if err != ErrRingBufferEmpty {
		t.Error("Expected error when popping while empty, got", err)
	}
}
