package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gokul656/go-lmax-distruptor/engine"
)

const (
	count         = 3000
	consumerCount = 100
)

func main() {
	now := time.Now()
	fmt.Print("Using Channels:")
	chans()
	fmt.Println(time.Since(now))

	now = time.Now()
	fmt.Print("\nUsing Disruptor:")
	optimized()
	fmt.Println(time.Since(now))

}

func optimized() {
	wg := &sync.WaitGroup{}
	handler := &EventHandler{name: "EventHandler"}
	engine := engine.NewDistruptor(10, handler)

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(consumerCount + 1)

	// Producer Goroutine
	go func() {
		defer wg.Done()
		for i := 0; i < count; i++ {
			for {
				if err := engine.Publish(i); err != nil {
					time.Sleep(time.Millisecond * 10) // Small sleep to retry
					fmt.Printf("waiting for data to be published: %v\n", err)
				} else {
					break
				}
			}
		}

		// After publishing all data, we cancel the context to signal consumers to stop
		cancel()
	}()

	// Multiple Consumer Goroutines
	for i := 0; i < consumerCount; i++ {
		go func(consumerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done(): // Context is cancelled, stop consuming
					fmt.Printf("Consumer %d: Context cancelled, exiting.\n", consumerID)
					return
				default:
					if err := engine.Consume(); err != nil {
						time.Sleep(time.Millisecond * 10) // Small sleep on error
					}
				}
			}
		}(i)
	}

	engine.Start(wg)
	wg.Wait()
	engine.Stop()
}

func chans() {
	wg := &sync.WaitGroup{}
	producerChannel := make(chan Event, count)
	consumerChannel := make(chan Event, count)

	wg.Add(1)

	// Producer Goroutine
	go func() {
		defer wg.Done()
		for i := 0; i < count; i++ {
			producerChannel <- Event{Index: i, Data: fmt.Sprintf("Event %d", i)}
			time.Sleep(time.Millisecond * 10) // Simulate delay
		}
		close(producerChannel)
	}()

	// Consumer Goroutines
	consumerWg := &sync.WaitGroup{}
	consumerWg.Add(consumerCount)
	for i := 0; i < consumerCount; i++ {
		go func(consumerID int) {
			defer consumerWg.Done()
			for event := range producerChannel {
				// Simulate processing
				time.Sleep(time.Millisecond * 50)
				fmt.Printf("Consumer %d processed: %d -> %s\n", consumerID, event.Index, event.Data)

				// Forward processed event to the final ordering stage
				consumerChannel <- event
			}
		}(i)
	}

	// Final Ordering Stage
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Wait for all consumers to finish
		go func() {
			consumerWg.Wait()
			close(consumerChannel)
		}()

		// Map to temporarily store out-of-order events
		eventMap := make(map[int]Event)
		expectedIndex := 0

		for event := range consumerChannel {
			eventMap[event.Index] = event

			// Deliver events in order
			for {
				if evt, exists := eventMap[expectedIndex]; exists {
					fmt.Printf("Processed in Order: %d -> %s\n", evt.Index, evt.Data)
					delete(eventMap, expectedIndex)
					expectedIndex++
				} else {
					break
				}
			}
		}
	}()

	wg.Wait()
}

// Event struct to maintain order
type Event struct {
	Index int
	Data  string
}

type EventHandler struct {
	name string
}

func (e *EventHandler) Process(data int) error {
	// Simulate processing
	time.Sleep(time.Millisecond * 100)
	fmt.Printf("Processed event from [%s]: %d\n", e.name, data)
	return nil
}
