package engine

type RingBuffer struct {
	rPointer   int   // Read pointer
	wPointer   int   // Write pointer
	bufferSize int   // Total size of the buffer [must be powers of 2]
	buffer     []int // Underlying array for the buffer
}

func (rb *RingBuffer) Read() *int {
	if rb.IsEmpty() {
		return nil
	}

	data := rb.buffer[rb.rPointer%rb.bufferSize]
	rb.rPointer++

	return &data
}

func (rb *RingBuffer) Write(data int) bool {
	if rb.IsFull() {
		return false
	}

	rb.buffer[rb.wPointer%rb.bufferSize] = data
	rb.wPointer++

	return true
}

// Calculate the number of items available for reading
func (rb *RingBuffer) SpaceForReading() int {
	return rb.wPointer - rb.rPointer
}

func (rb *RingBuffer) IsEmpty() bool {
	return rb.rPointer == rb.wPointer
}

func (rb *RingBuffer) IsFull() bool {
	return rb.SpaceForReading() == rb.bufferSize
}

func NewRingBuffer(bufferSize int) *RingBuffer {
	return &RingBuffer{
		rPointer:   0,
		wPointer:   0,
		bufferSize: bufferSize,
		buffer:     make([]int, bufferSize),
	}
}
