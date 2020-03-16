package event

import (
	_ "fmt"
	"ioutil"
	"reflect"
	"sync"
	"time"
)

type DirectoryWatcher struct {
	Dispatcher
	path string // path to directory on disk
}

func (w *DirectoryWatcher) Init(path string) {
	w.path = path
	w.Dispatcher.Init(w)
}

// Watch the given directory and dispatch new file events
func (w *DirectoryWatcher) Watch() error {
	for {
		files, _ := ioutil.ReadDir(w.path)
		for _, file := range files {
			if w.Unseen(file) {
				w.Dispatch(NewFile, file)
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// Sample Event Types
const (
	UnknownEvent EventType = iota
	CommitEvent
	MessageEvent
	FileOpenedEvent
)

// Names of event types
var eventTypeStrings = [...]string{
	"unknown", "entryCommitted", "messageReceived", "fileOpened",
}

//===========================================================================
// Event Types
//===========================================================================

// EventType is an enumeration of the kind of events that can occur.
type EventType uint16

// String returns the name of event types
func (t EventType) String() string {
	return eventTypeStrings[t]
}

// Callback is a function that can receive events.
type Callback func(Event) error

//===========================================================================
// Event Dispatcher Definition and Methods
//===========================================================================

// Dispatcher objects can register callbacks for specific events, then when
// those events occur, dispatch them to all callback functions.
type Dispatcher struct {
	sync.RWMutex
	source    interface{}
	callbacks map[EventType][]Callback
}

// Init a dispatcher with the source, creating the callbacks map.
func (d *Dispatcher) Init(source interface{}) {
	d.source = source
	d.callbacks = make(map[EventType][]Callback)
}

// Register a callback function for the specified event type.
func (d *Dispatcher) Register(etype EventType, callback Callback) {
	d.Lock()
	defer d.Unlock()
	d.callbacks[etype] = append(d.callbacks[etype], callback)
}

// Remove a callback function for the specified event type.
func (d *Dispatcher) Remove(etype EventType, callback Callback) {
	d.Lock()
	defer d.Unlock()

	// Grab a reference to the function pointer
	ptr := reflect.ValueOf(callback).Pointer()

	// Find callback by pointer and remove it
	callbacks := d.callbacks[etype]
	for idx, cb := range callbacks {
		if reflect.ValueOf(cb).Pointer() == ptr {
			d.callbacks[etype] = append(callbacks[:idx], callbacks[idx+1:]...)
		}
	}
}

// Dispatch an event, ensuring that the event is properly formatted.
// Currently this method simply warns if there is an error.
// TODO: return list of errors or do better error handling.
func (d *Dispatcher) Dispatch(etype EventType, value interface{}) error {
	d.RLock()
	defer d.RUnlock()

	// Create the event
	e := &event{
		etype:  etype,
		source: d.source,
		value:  value,
	}

	// Dispatch the event to all callbacks
	for _, cb := range d.callbacks[etype] {
		if err := cb(e); err != nil {
			return err
		}
	}

	return nil
}

//===========================================================================
// Event Definition and Methods
//===========================================================================

// Event represents actions that occur during consensus. Listeners can
// register callbacks with event handlers for specific event types.
type Event interface {
	Type() EventType
	Source() interface{}
	Value() interface{}
}

// event is an internal implementation of the Event interface.
type event struct {
	etype  EventType
	source interface{}
	value  interface{}
}

// Type returns the event type.
func (e *event) Type() EventType {
	return e.etype
}

// Source returns the entity that dispatched the event.
func (e *event) Source() interface{} {
	return e.source
}

// Value returns the current value associated with the event.
func (e *event) Value() interface{} {
	return e.value
}
