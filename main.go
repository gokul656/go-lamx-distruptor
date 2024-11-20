package main

import (
	"fmt"
	"sync"

	"github.com/gokul656/go-lmax-distruptor/engine"
)

func main() {
	wg := &sync.WaitGroup{}
	handler := &EventHandler{name: "EventHandler"}
	engine := engine.NewDistruptor(8, handler)

	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 8; i++ {
			if err := engine.Publish(i); err != nil {
				fmt.Printf("error publishing data: %v\n", err)
				break // Exit if there's an error
			}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 8; i++ {
			if err := engine.Consume(); err != nil {
				fmt.Printf("error consuming data: %v\n", err)
				break // Exit if there's an error
			}
		}
	}()

	engine.Start(wg)
	wg.Wait()
	engine.Stop()
}

type EventHandler struct {
	name string
}

func (e *EventHandler) Process(data int) error {
	fmt.Printf("Processed event from [%s]: %d\n", e.name, data)
	return nil
}
