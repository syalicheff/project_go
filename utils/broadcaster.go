package utils

import (
	"Project_go/model"
	"sync"
)

// Broadcaster is a structure that allows events to be broadcasted to subscribers
type Broadcaster struct {
	subscribers []chan *model.Payment // modified type of channel to only accept Payment events
	mux         sync.RWMutex
}

var b *Broadcaster
var once sync.Once

// NewBroadcaster creates a new instance of Broadcaster
func NewBroadcaster() *Broadcaster {
	return &Broadcaster{}
}

// Subscribe adds a subscriber to the Broadcaster
func (b *Broadcaster) Subscribe() chan *model.Payment { // modified type of channel to only accept Payment events
	b.mux.Lock()
	defer b.mux.Unlock()

	ch := make(chan *model.Payment) // modified type of channel to only accept Payment events
	b.subscribers = append(b.subscribers, ch)
	return ch
}

// Unsubscribe removes a subscriber from the Broadcaster
func (b *Broadcaster) Unsubscribe(ch chan *model.Payment) { // modified type of channel to only accept Payment events
	b.mux.Lock()
	defer b.mux.Unlock()

	for i, c := range b.subscribers {
		if c == ch {
			b.subscribers = append(b.subscribers[:i], b.subscribers[i+1:]...)
			break
		}
	}
}
func (b *Broadcaster) Broadcast(evt *model.Payment) { // modified type of channel to only accept Payment events
	b.mux.RLock()
	defer b.mux.RUnlock()

	for _, ch := range b.subscribers {
		ch <- evt
	}
}

// Broadcast sends an event to all subscribers of the Broadcaster

// GetBroadcaster returns the singleton instance of Broadcaster
func GetBroadcaster() *Broadcaster {
	once.Do(func() {
		b = NewBroadcaster()
	})
	return b
}
