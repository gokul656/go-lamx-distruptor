package engine

import (
	"errors"
	"fmt"
	"sync"
)

type HandlerFunc interface {
	Process(data int) error
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

type DistruptorEngine struct {
	ringbuffer *RingBuffer
	handler    HandlerFunc
}

func (de *DistruptorEngine) Start(wg *sync.WaitGroup) {
	fmt.Println("starting disruptor engine")
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

func (de *DistruptorEngine) Stop() error {
	fmt.Println("stopping disruptor engine...")
	return nil
}

func (de *DistruptorEngine) IsFull() bool {
	return de.ringbuffer.IsFull()
}

func (de *DistruptorEngine) IsEmpty() bool {
	return de.ringbuffer.IsEmpty()
}

func (de *DistruptorEngine) Publish(data int) error {
	if de.IsFull() {
		return errors.New("unable to publish, buffer is full")
	}
	de.ringbuffer.Write(data)
	return nil
}

func (de *DistruptorEngine) Consume() error {
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

func (de *DistruptorEngine) GetBufferSize() int {
	return de.ringbuffer.bufferSize
}

func (de *DistruptorEngine) GetPendingCount() int {
	return de.ringbuffer.SpaceForReading()
}

func (de *DistruptorEngine) HandleException(err error) {
	fmt.Printf("error occurred: %v\n", err)
}

func NewDistruptor(bufferSize int, handler HandlerFunc) *DistruptorEngine {
	return &DistruptorEngine{
		ringbuffer: NewRingBuffer(bufferSize),
		handler:    handler,
	}
}
