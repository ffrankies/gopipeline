package worker

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/semaphore"
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
	semaphore *semaphore.Weighted
	mutex     *sync.Mutex
	context   context.Context
}

// makeQueue makes a new queue and returns it along with its size.
func makeQueue() *Queue {
	return &Queue{
		elements:  make([]interface{}, 0),
		semaphore: semaphore.NewWeighted(0),
		mutex:     &sync.Mutex{},
		context:   context.TODO(),
	}
}

// Push adds an element to the queue
func (q *Queue) Push(element interface{}) {
	q.mutex.Lock()
	q.elements = append(q.elements, element)
	q.semaphore.Release(1)
	q.mutex.Unlock()
}

// Pop removes an element from the queue
func (q *Queue) Pop() interface{} {
	q.semaphore.Acquire(q.context, 1)
	q.mutex.Lock()
	result := q.elements[0]
	q.elements = q.elements[1:]
	q.mutex.Unlock()
	return result
}

// GetLength gives out the length of the queue
func (q *Queue) GetLength() int {
	return q.length
}
