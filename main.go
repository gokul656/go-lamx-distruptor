package main

import "fmt"

func main() {
	ringBuffer := NewRingBuffer(10)
	ringBuffer.Write(1)
	ringBuffer.Write(2)
	ringBuffer.Write(3)

	fmt.Printf("read: %v\n", *ringBuffer.Read()) // 1
	fmt.Printf("read: %v\n", *ringBuffer.Read()) // 2
	fmt.Printf("read: %v\n", *ringBuffer.Read()) // 3

	fmt.Printf("write: %v\n", ringBuffer.Write(4))
	fmt.Printf("write: %v\n", ringBuffer.Write(5))

	fmt.Printf("read: %v\n", *ringBuffer.Read()) // 4
	fmt.Printf("read: %v\n", *ringBuffer.Read()) // 5
	fmt.Printf("read: %v\n", ringBuffer.Read())  // nil
}
