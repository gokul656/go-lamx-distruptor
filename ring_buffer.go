package main

type RingBuffer struct {
	rPointer   int   // Read pointer
	wPointer   int   // Write pointer
	bufferSize int   // Total size of the buffer
	buffer     []int // Underlying array for the buffer
}

// Read data from the ring buffer
func (rb *RingBuffer) Read() *int {
	if rb.IsEmpty() {
		return nil // Return nil if the buffer is empty
	}

	// Read data from the current read pointer
	data := rb.buffer[rb.rPointer%rb.bufferSize]
	rb.rPointer++ // Advance the read pointer

	return &data
}

// Write data to the ring buffer
func (rb *RingBuffer) Write(data int) bool {
	if rb.IsFull() {
		return false // Return false if the buffer is full
	}

	// Write data at the current write pointer
	rb.buffer[rb.wPointer%rb.bufferSize] = data
	rb.wPointer++ // Advance the write pointer

	return true
}

// Calculate the number of items available for reading
func (rb *RingBuffer) SpaceForReading() int {
	return rb.wPointer - rb.rPointer
}

// Check if the ring buffer is empty
func (rb *RingBuffer) IsEmpty() bool {
	return rb.rPointer == rb.wPointer
}

// Check if the ring buffer is full
func (rb *RingBuffer) IsFull() bool {
	return rb.SpaceForReading() == rb.bufferSize
}

// Initialize a new ring buffer
func NewRingBuffer(bufferSize int) *RingBuffer {
	return &RingBuffer{
		rPointer:   0,
		wPointer:   0,
		bufferSize: bufferSize,
		buffer:     make([]int, bufferSize),
	}
}
