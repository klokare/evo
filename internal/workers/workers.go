package workers

import (
	"runtime"
	"sync"
)

// Task can be any unit of work
type Task interface{}

// Do performs the tasks in parallel
func Do(tasks []Task, action func(Task)) {

	// Feed in the tasks
	ch := make(chan Task)
	go func(ch chan Task) {
		for _, t := range tasks {
			ch <- t
		}
		close(ch)
	}(ch)

	// Spin up some workers
	var wg sync.WaitGroup
	for w := 0; w < runtime.NumCPU(); w++ {
		wg.Add(1)
		go func(ch <-chan Task) {
			defer wg.Done()
			for t := range ch {
				action(t)
			}
		}(ch)
	}
	wg.Wait()
}
