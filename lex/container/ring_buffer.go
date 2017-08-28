package container

import "errors"

type RingBuffer struct {
	buffer      []interface{}
	start, used uint64
}

var ErrRingBufferEmpty = errors.New("ring_buffer: empty")

func NewRingBuffer(size uint64) *RingBuffer {
	if size == 0 {
		panic("ring_buffer: must have a size > 0")
	}

	return &RingBuffer{
		buffer: make([]interface{}, size),
		start:  0,
		used:   0,
	}
}

func (rb *RingBuffer) Top() (interface{}, error) {
	if rb.empty() {
		return nil, ErrRingBufferEmpty
	}

	index := (rb.start + rb.used - 1) % uint64(len(rb.buffer))
	return rb.buffer[index], nil
}

func (rb *RingBuffer) Push(r interface{}) {
	if rb.full() {
		rb.buffer[rb.start] = r
		rb.start = (rb.start + 1) % uint64(len(rb.buffer))
	} else {
		index := (rb.start + rb.used) % uint64(len(rb.buffer))
		rb.buffer[index] = r
		rb.used++
	}
}

func (rb *RingBuffer) Pop() (interface{}, error) {
	r, err := rb.Top()
	if err == nil {
		rb.used--
	}
	return r, err
}

func (rb *RingBuffer) full() bool {
	return rb.used == uint64(len(rb.buffer))
}

func (rb *RingBuffer) empty() bool {
	return rb.used == 0
}
