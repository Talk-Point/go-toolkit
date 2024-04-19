package signal

import (
	"sync"
)

// Signal type for demonstration
type Signal string

// Callback function type
type Callback func(signal Signal, data interface{})

// SignalDispatcher to hold registered callbacks
type SignalDispatcher struct {
	listeners map[Signal][]Callback
	lock      sync.Mutex
}

// NewSignalDispatcher creates a new instance of SignalDispatcher
func NewSignalDispatcher() *SignalDispatcher {
	return &SignalDispatcher{
		listeners: make(map[Signal][]Callback),
	}
}

// Connect registers a callback for a given signal
func (d *SignalDispatcher) Connect(signal Signal, callback Callback) {
	d.lock.Lock()
	defer d.lock.Unlock()

	if _, exists := d.listeners[signal]; !exists {
		d.listeners[signal] = []Callback{}
	}
	d.listeners[signal] = append(d.listeners[signal], callback)
}

// Send emits a signal to all registered callbacks, executing them in parallel
func (d *SignalDispatcher) Emit(signal Signal, data interface{}) {
	d.lock.Lock()
	callbacks, exists := d.listeners[signal]
	d.lock.Unlock() // Unlock as soon as possible, before invoking callbacks

	if exists {
		var wg sync.WaitGroup
		for _, callback := range callbacks {
			wg.Add(1)
			go func(cb Callback) {
				defer wg.Done()
				cb(signal, data)
			}(callback)
		}
		wg.Wait()
	}
}
