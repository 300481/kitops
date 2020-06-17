package queue_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/300481/kitops/pkg/queue"
)

const Value string = "Test"

var ReturnedValues []string

type TestStruct struct{}

func (t *TestStruct) Process(q *queue.Queue) {
	ReturnedValues = append(ReturnedValues, q.StartNext().(string))
	q.Finish(true)
}

func TestNew(t *testing.T) {
	ts := &TestStruct{}
	q := queue.New(ts)
	for i := 0; i < 10; i++ {
		q.Add(Value + strconv.Itoa(i))
	}
	time.Sleep(time.Second)
	for i := 0; i < 10; i++ {
		if Value+strconv.Itoa(i) != ReturnedValues[i] {
			t.Errorf("got %s, want %s", ReturnedValues[i], Value+strconv.Itoa(i))
		}
	}
}
