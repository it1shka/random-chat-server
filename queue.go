package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Queue struct {
	sync.Mutex
	awaiting map[*websocket.Conn]bool
}

func NewQueue() *Queue {
	return &Queue{
		awaiting: make(map[*websocket.Conn]bool),
	}
}

func (queue *Queue) Append(conn *websocket.Conn) {
	queue.Lock()
	defer queue.Unlock()
	queue.awaiting[conn] = true
}

func (queue *Queue) Delete(conn *websocket.Conn) {
	queue.Lock()
	defer queue.Unlock()
	delete(queue.awaiting, conn)
}

func (queue *Queue) DoMatchmaking(
	fn func(a, b *websocket.Conn),
	chunkSize int,
) {
	queue.Lock()
	defer queue.Unlock()

	allWaiting := keysOf(queue.awaiting)
	chunks := chunksOf(allWaiting, chunkSize)
	withoutMatch := make(chan *websocket.Conn)

	var wg sync.WaitGroup
	wg.Add(len(chunks))

	for _, chunk := range chunks {

		go func(connections []*websocket.Conn) {
			defer wg.Done()
			for _, pair := range chunksOf(connections, 2) {
				if len(pair) < 2 {
					withoutMatch <- pair[0]
				} else {
					a, b := pair[0], pair[1]
					fn(a, b)
				}
			}
		}(chunk)

	}

	go func() {
		defer close(withoutMatch)
		wg.Wait()
	}()

	next := make(map[*websocket.Conn]bool)
	for conn := range withoutMatch {
		next[conn] = true
	}
	queue.awaiting = next

}
