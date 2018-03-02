package workers

import (
	"runtime"
	"testing"
	"time"
)

func TestDo(t *testing.T) {

	type task struct {
		called bool
	}

	tasks := make([]Task, runtime.NumCPU()*2) // so some have to wait
	for i := 0; i < len(tasks); i++ {
		tasks[i] = new(task)
	}

	Do(tasks, func(wt Task) {
		t0 := wt.(*task)
		time.Sleep(100 * time.Nanosecond) // delay a bit like we're doing something
		t0.called = true
	})

	for i, wt := range tasks {
		if !wt.(*task).called {
			t.Errorf("task %d was not called", i)
		}
	}
}
