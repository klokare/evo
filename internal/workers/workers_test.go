package workers

import (
	"errors"
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

	Do(tasks, func(wt Task) error {
		t0 := wt.(*task)
		time.Sleep(100 * time.Nanosecond) // delay a bit like we're doing something
		t0.called = true
		return nil
	})

	for i, wt := range tasks {
		if !wt.(*task).called {
			t.Errorf("task %d was not called", i)
		}
	}
}

func TestDoWithError(t *testing.T) {

	type task struct {
		called   bool
		hasError bool
	}

	tasks := make([]Task, runtime.NumCPU()*2) // so some have to wait
	for i := 0; i < len(tasks); i++ {
		tasks[i] = &task{hasError: i == 1}
	}

	err := Do(tasks, func(wt Task) error {
		t0 := wt.(*task)
		time.Sleep(100 * time.Nanosecond) // delay a bit like we're doing something
		t0.called = true
		if t0.hasError {
			return errors.New("mock task error")
		}
		return nil
	})

	if err == nil {
		t.Errorf("expected error not found")
	}

	all := true
	for _, wt := range tasks {
		all = all && wt.(*task).called
	}
	if all {
		t.Errorf("did not expect all tasks to be called")
	}
}
