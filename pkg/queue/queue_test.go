package queue_test

// Must be improved to test the queue processing element by element

import (
	"testing"
	"time"

	"github.com/300481/kitops/pkg/queue"
)

const Value string = "Hello World"

var ReturnedValue string

type TestStruct struct{}

func (t *TestStruct) Process(q *queue.Queue) {
	ReturnedValue = q.StartNext().(string)
	q.Finish(true)
}

func TestNew(t *testing.T) {
	ts := &TestStruct{}
	q := queue.New(ts)
	q.Add(Value)
	time.Sleep(time.Second * 2)
	if Value != ReturnedValue {
		t.Errorf("got %s, want %s", ReturnedValue, Value)
	}
}
