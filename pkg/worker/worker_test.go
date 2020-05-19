package worker

import (
	"context"
	"testing"
	"time"
)

func TestStart(t *testing.T) {

	rChan := make(chan bool, 10)
	tt := testTask{result: rChan}

	ctx, cancel := context.WithCancel(context.Background())

	go Start(&tt, ctx)

	if result := <-rChan; !result {
		t.Error("!result")
	}
	cancel()
}

type testTask struct {
	result chan bool
}

func (t *testTask) Exec() {
	t.result <- true
}

func (t *testTask) GetInterval() time.Duration {
	return time.Millisecond * 500
}
