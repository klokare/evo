package workers

import (
	"runtime"
	"sync"
)

// Task can be any unit of work
type Task interface{}

// Do performs the tasks in parallel
func Do(tasks []Task, action func(Task) error) (err error) {

	// Feed in the tasks
	ch := make(chan Task, len(tasks))
	go func(ch chan Task) {
		defer close(ch)
		for _, t := range tasks {
			ch <- t
		}
	}(ch)

	// Listen for error
	ec := make(chan error, len(tasks))
	done := make(chan struct{})
	go func(ec <-chan error, abort chan struct{}) {
		defer close(done)
		err = <-ec
	}(ec, done)

	// Spin up some workers
	var wg sync.WaitGroup
	for w := 0; w < runtime.NumCPU(); w++ {
		wg.Add(1)
		go func(ch <-chan Task, ec chan<- error, done <-chan struct{}) {
			defer wg.Done()
			for {
				select {
				case t := <-ch:
					if t == nil {
						return
					}
					if err := action(t); err != nil {
						ec <- err
					}
				case <-done:
					return
				}
			}
		}(ch, ec, done)
	}
	wg.Wait()
	close(ec)
	<-done
	return err
}
