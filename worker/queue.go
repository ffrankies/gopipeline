package worker

import (
	"fmt"
	"sync"
)

// Element struct
type Element struct {
	Data interface{}
}

// MessageString returns the string value of the element data
func (n *Element) MessageString() string {
	return fmt.Sprint(n.Data)
}

// Queue Struct
type Queue struct {
	elements  []interface{}
	length    int
	semaphore chan int
	mutex     *sync.Mutex
}

// makeQueue makes a new queue and returns it along with its size.
func makeQueue() *Queue {
	return &Queue{
		elements:  make([]interface{}, 0),
		semaphore: make(chan int, 9999), // Make a buffered channel with a high buffer
		mutex:     &sync.Mutex{},
	}
}

// Push adds an element to the queue
func (q *Queue) Push(element interface{}) {
	q.mutex.Lock()
	q.elements = append(q.elements, element)
	q.semaphore <- 1 // Write to channel (semaphore++)
	q.mutex.Unlock()
}

// Pop removes an element from the queue
func (q *Queue) Pop() interface{} {
	<-q.semaphore // Read from channel (semaphore--)
	q.mutex.Lock()
	result := q.elements[0]
	q.elements = q.elements[1:]
	q.mutex.Unlock()
	return result
}

// GetLength gives out the length of the queue
func (q *Queue) GetLength() int {
	q.mutex.Lock()
	length := len(q.elements)
	q.mutex.Unlock()
	return length
}
