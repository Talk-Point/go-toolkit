// Package signal implements a simple event dispatch system that allows
// components within an application to communicate with each other in a
// loosely coupled manner. It provides a SignalDispatcher which maintains
// a mapping of signals (events) to callbacks that are to be executed when
// a signal is emitted.
//
// A Signal in this context is a string identifier that represents a specific
// type of event. Callback functions registered to a signal are called with
// the signal and any accompanying data when that signal is emitted.
//
// Basic Usage:
//
// To use this package, start by creating a SignalDispatcher:
//
//	dispatcher := signal.NewSignalDispatcher()
//
// Then, register one or more callbacks to listen for specific signals:
//
//	dispatcher.Connect("order-created", func(signal signal.Signal, data interface{}) {
//	    order, ok := data.(*models.Order)
//	    if !ok {
//	        fmt.Println("Invalid order data")
//	        return
//	    }
//	    fmt.Printf("Received order created signal for order ID: %v\n", order.ID)
//	})
//
// You can emit a signal anywhere in your application with any related data:
//
//	dispatcher.Emit("order-created", newOrder)
//
// This package ensures that all registered callbacks for a signal are executed
// in parallel, which can help in improving the performance of operations
// that are independent and can be processed concurrently.
//
// Note: This package uses goroutines to handle signals concurrently and it's
// important to ensure that callbacks are thread-safe.
//
// This package is thread-safe and can be used concurrently across multiple
// goroutines.
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
