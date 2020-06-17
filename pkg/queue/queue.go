package queue

import (
	"container/list"
	"sync"
)

type Processor interface {
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
	current        string
	previous       string
	lastSuccessful string
	status         Status
	mux            *sync.Mutex
	processor      Processor
}

// New creates a new Queue and returns *Queue.
func New(processor Processor) *Queue {
	return &Queue{
		list:           list.New(),
		current:        "",
		previous:       "",
		lastSuccessful: "",
		status:         Init,
		mux:            &sync.Mutex{},
		processor:      processor,
	}
}

// Add adds a new commit string as element to the Queue and initiate a threaded processing.
func (q *Queue) Add(commit string) {
	q.list.PushBack(commit)
	go q.processor.Process(q)
}

// Next returns the next element in the queue and clean the queue by one.
// Must be finalized with Finish() to unlock the mutex.
func (q *Queue) Next() string {
	if q.list.Len() == 0 {
		return ""
	}

	q.mux.Lock()
	q.previous = q.current
	q.current = q.list.Front().Value.(string)
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
