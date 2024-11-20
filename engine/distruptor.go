package engine

import (
	"errors"
	"fmt"
	"sync"
)

type HandlerFunc[T any] interface {
	Process(data T) error
}

type Distruptor interface {
	Start(wg *sync.WaitGroup)
	Stop() error
	IsFull() bool
	IsEmpty() bool
	Publish(int) error
	Consume() error
	GetBufferSize() int
	GetPendingCount() int
	HandleException(err error) // Handle exceptions with actual error messages
}

type DistruptorEngine[T any] struct {
	ringbuffer *RingBuffer[T]
	handler    HandlerFunc[T]
}

func (de *DistruptorEngine[T]) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := de.Consume(); err != nil {
				break
			}
		}
	}()
}

func (de *DistruptorEngine[T]) Stop() error {
	return nil
}

func (de *DistruptorEngine[T]) IsFull() bool {
	return de.ringbuffer.IsFull()
}

func (de *DistruptorEngine[T]) IsEmpty() bool {
	return de.ringbuffer.IsEmpty()
}

func (de *DistruptorEngine[T]) Publish(data T) error {
	if de.IsFull() {
		return errors.New("unable to publish, buffer is full")
	}
	de.ringbuffer.Write(data)
	return nil
}

func (de *DistruptorEngine[T]) Consume() error {
	if de.IsEmpty() {
		return errors.New("unable to consume, buffer is empty")
	}
	data := de.ringbuffer.Read()
	if err := de.handler.Process(*data); err != nil {
		de.HandleException(err)
		return err
	}
	return nil
}

func (de *DistruptorEngine[T]) GetBufferSize() int {
	return de.ringbuffer.bufferSize
}

func (de *DistruptorEngine[T]) GetPendingCount() int {
	return de.ringbuffer.SpaceForReading()
}

func (de *DistruptorEngine[T]) HandleException(err error) {
	fmt.Printf("error occurred: %v\n", err)
}

func NewDistruptor[T any](bufferSize int, handler HandlerFunc[T]) *DistruptorEngine[T] {
	return &DistruptorEngine[T]{
		ringbuffer: NewRingBuffer[T](bufferSize),
		handler:    handler,
	}
}
