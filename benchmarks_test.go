package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/gokul656/go-lmax-distruptor/engine"
)

func BenchmarkDisruptor(b *testing.B) {
	count := 1000
	handler := &EventHandler{name: "EventHandler"}
	engine := engine.NewDistruptor(10, handler)

	b.ResetTimer()

	go func() {
		for i := 0; i < count; i++ {
			if err := engine.Publish(i); err != nil {
				fmt.Println("Error:", err)
			}
		}
	}()

	go func() {
		for i := 0; i < count; i++ {
			if err := engine.Consume(); err != nil {
				fmt.Println("Error:", err)
			}
		}
	}()

	// Start the engine and benchmark the performance
	b.StartTimer()
	engine.Start(&sync.WaitGroup{})
	b.StopTimer()

	time.Sleep(1 * time.Second) // Give time for goroutines to finish
}

func BenchmarkChannelGoroutines(b *testing.B) {
	count := 1000
	messageChannel := make(chan int, count)
	handler := &EventHandler{name: "EventHandler"}

	b.ResetTimer()

	go func() {
		for i := 0; i < count; i++ {
			messageChannel <- i
		}
		close(messageChannel)
	}()

	go func() {
		timeout := time.After(5 * time.Second)
		for {
			select {
			case <-timeout:
				return
			case data, ok := <-messageChannel:
				if !ok {
					return
				}
				handler.Process(data)                 // Process message
				timeout = time.After(3 * time.Second) // Reset the timeout
			}
		}
	}()

	b.StartTimer()
	time.Sleep(1 * time.Second)
	b.StopTimer()
}
