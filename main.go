package main

import (
	"fmt"
	"sync"
	"time"
)

type SharedBuffer struct {
	buffer []byte
	mutex  sync.Mutex
}

func main() {
	var M, N int

	fmt.Print("Enter the number of reader goroutines (M): ")
	fmt.Scan(&M)

	fmt.Print("Enter the number of writer goroutines (N): ")
	fmt.Scan(&N)

	// sum of reader and writer is taken as the buffersize here
	bufferSize := M + N 

	buffer := SharedBuffer{buffer: make([]byte, bufferSize)}
	readCh := make(chan byte, bufferSize)
	writeCh := make(chan byte, bufferSize)

	// Start N writer goroutines
	for i := 0; i < N; i++ {
		go writer(i, &buffer, writeCh)
	}

	// Start M reader goroutines
	for i := 0; i < M; i++ {
		go reader(i, &buffer, readCh)
	}

	// to read from both channels
	go func() {
		for {
			select {
			case <-readCh:
			case <-writeCh:
			}
		}
	}()

	select {}
}

func writer(id int, buffer *SharedBuffer, writeCh chan<- byte) {
	for {

		value := byte(id + 1)

		buffer.mutex.Lock()
		select {
		case writeCh <- value:
			buffer.buffer = append(buffer.buffer, value)
			fmt.Printf("Writer %d wrote: %d\n", id, value)
		default:
		}
		buffer.mutex.Unlock()

		// adding some delay
		time.Sleep(time.Millisecond * 100)

	}
}

func reader(id int, buffer *SharedBuffer, readCh chan<- byte) {
	for {
		buffer.mutex.Lock()
		if len(buffer.buffer) > 0 {
			value := buffer.buffer[0]
			buffer.buffer = buffer.buffer[1:]
			fmt.Printf("Reader %d read: %d\n", id, value)
			readCh <- value
		}
		buffer.mutex.Unlock()

		// adding some delay
		time.Sleep(time.Millisecond * 100)

	}
}


