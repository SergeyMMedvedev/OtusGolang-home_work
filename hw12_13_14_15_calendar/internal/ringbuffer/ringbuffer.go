package ringbuffer

import (
	"fmt"
)

type RingBuffer struct {
	buffer []interface{}
	size   int
	head   int
	tail   int
	count  int
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		buffer: make([]interface{}, size),
		size:   size,
	}
}

func (rb *RingBuffer) Enqueue(item interface{}) {
	if rb.count == rb.size {
		rb.head = (rb.head + 1) % rb.size
	} else {
		rb.count++
	}
	rb.buffer[rb.tail] = item
	rb.tail = (rb.tail + 1) % rb.size
}

func (rb *RingBuffer) Dequeue() (interface{}, error) {
	if rb.count == 0 {
		return nil, fmt.Errorf("ring buffer is empty")
	}
	item := rb.buffer[rb.head]
	rb.head = (rb.head + 1) % rb.size
	rb.count--
	return item, nil
}

func (rb *RingBuffer) IsEmpty() bool {
	return rb.count == 0
}

func (rb *RingBuffer) IsFull() bool {
	return rb.count == rb.size
}

func (rb *RingBuffer) Size() int {
	return rb.count
}
