package evo

import "github.com/klokare/evo/internal/workers"

// Event is key used with a Callback
type Event byte

// Events associated with the Experiment
const (
	Started Event = iota + 1
	Decoded
	Evaluated
	Advanced
	Completed
)

// Callback functions are called when the event to which they are subscribed occurs. The final flag
// is true when the experiment is solved or on its final iteration
type Callback func(pop Population) error

// Subscription pairs a listener with its event
type Subscription struct {
	Event
	Callback
}

// SubscriptionProvider informs the caller of waiting subscriptions
type SubscriptionProvider interface {
	Subscriptions() []Subscription
}

// Publish an event to the listeners. Callbacks will be called concurrently so there is no
// guarantee the order in which they are called.
func publish(listeners map[Event][]Callback, event Event, pop Population) (err error) {
	callbacks := listeners[event]
	if len(callbacks) == 0 {
		return
	}

	tasks := make([]workers.Task, len(callbacks))
	for i, cb := range callbacks {
		tasks[i] = cb
	}

	err = workers.Do(tasks, func(wt workers.Task) error {
		cb := wt.(Callback)
		return cb(pop)
	})

	return
}
