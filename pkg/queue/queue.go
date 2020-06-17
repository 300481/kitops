package queue

import (
	"container/list"
	"sync"
)

type Consumer interface {
	Process(q *Queue)
}

type Status string

const (
	Init       Status = "Init"
	InProgress Status = "InProgress"
	Failed     Status = "Failed"
	Successful Status = "Successful"
)

type Queue struct {
	list           *list.List
	current        interface{}
	previous       interface{}
	lastSuccessful interface{}
	status         Status
	mux            *sync.Mutex
	consumer       Consumer
}

// New creates a new Queue and returns *Queue.
func New(consumer Consumer) *Queue {
	return &Queue{
		list:           list.New(),
		current:        nil,
		previous:       nil,
		lastSuccessful: nil,
		status:         Init,
		mux:            &sync.Mutex{},
		consumer:       consumer,
	}
}

// Add adds a new commit string as element to the Queue and initiate a threaded processing.
func (q *Queue) Add(v interface{}) {
	q.list.PushBack(v)
	go q.consumer.Process(q)
}

// StartNext returns the next element in the queue and clean the queue by one.
// It also Locks a Mutex.
// Must be finalized with Finish() to unlock the mutex.
func (q *Queue) StartNext() interface{} {
	if q.list.Len() == 0 {
		return nil
	}

	q.mux.Lock()
	q.previous = q.current
	q.current = q.list.Front().Value
	q.list.Remove(q.list.Front())
	q.status = InProgress

	return q.current
}

// Finish finishes the current element processing of the queue and unlock the mutex.
func (q *Queue) Finish(success bool) {
	if success {
		q.status = Successful
		q.lastSuccessful = q.current
	} else {
		q.status = Failed
	}
	q.mux.Unlock()
}

// Current returns the current element
func (q *Queue) Current() interface{} {
	return q.current
}

// Previous returns the previous element
func (q *Queue) Previous() interface{} {
	return q.previous
}

// LastSuccessful returns the last successful processed element
func (q *Queue) LastSuccessful() interface{} {
	return q.lastSuccessful
}

// Status returns the status of the current processed element
func (q *Queue) Status() Status {
	return q.status
}
