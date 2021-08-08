package main

import "sync"

type Queue struct {
	data  [][]byte
	mutex sync.Mutex
}

func (q *Queue) Push(b []byte) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.data = append(q.data, b)
}

func (q *Queue) Pull() []byte {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	result := q.data[0]
	q.data = q.data[1:]

	return result
}
