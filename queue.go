package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Queue struct {
	sync.Mutex
	awaiting []*websocket.Conn
}

func NewQueue() *Queue {
	return &Queue{
		awaiting: []*websocket.Conn{},
	}
}

func (queue *Queue) Append(conn *websocket.Conn) {
	queue.Lock()
	queue.awaiting = append(queue.awaiting, conn)
	queue.Unlock()
}

func (queue *Queue) DoMatchmaking(
	fn func(a, b *websocket.Conn),
	chunkSize int,
) {
	queue.Lock()
	defer queue.Unlock()

	chunks := chunksOf(queue.awaiting, chunkSize)
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

	next := []*websocket.Conn{}
	for conn := range withoutMatch {
		next = append(next, conn)
	}
	queue.awaiting = next

}
