package engine

type RingBuffer[T any] struct {
	rPointer   int // Read pointer
	wPointer   int // Write pointer
	bufferSize int // Total size of the buffer [must be powers of 2]
	buffer     []T // Underlying array for the buffer
}

func (rb *RingBuffer[T]) Read() *T {
	if rb.IsEmpty() {
		return nil
	}

	data := rb.buffer[rb.rPointer%rb.bufferSize]
	rb.rPointer++

	return &data
}

func (rb *RingBuffer[T]) Write(data T) bool {
	if rb.IsFull() {
		return false
	}

	rb.buffer[rb.wPointer%rb.bufferSize] = data
	rb.wPointer++

	return true
}

// Calculate the number of items available for reading
func (rb *RingBuffer[T]) SpaceForReading() int {
	return rb.wPointer - rb.rPointer
}

func (rb *RingBuffer[T]) IsEmpty() bool {
	return rb.rPointer == rb.wPointer
}

func (rb *RingBuffer[T]) IsFull() bool {
	return rb.SpaceForReading() == rb.bufferSize
}

func NewRingBuffer[T any](bufferSize int) *RingBuffer[T] {
	return &RingBuffer[T]{
		rPointer:   0,
		wPointer:   0,
		bufferSize: bufferSize,
		buffer:     make([]T, bufferSize),
	}
}
