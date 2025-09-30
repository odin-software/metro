package observation

import (
	"sync"
	"time"
)

// Collector receives and buffers events from all trains
type Collector struct {
	eventChannel <-chan interface{}
	buffer       []interface{}
	mu           sync.RWMutex
	lastCollect  time.Time
}

// NewCollector creates a new event collector
func NewCollector(eventChannel <-chan interface{}) *Collector {
	return &Collector{
		eventChannel: eventChannel,
		buffer:       make([]interface{}, 0, 100),
		lastCollect:  time.Now(),
	}
}

// Collect gathers all events that have arrived since last collection
// Returns the collected events and clears the buffer
func (c *Collector) Collect() []interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Copy buffer and reset
	collected := make([]interface{}, len(c.buffer))
	copy(collected, c.buffer)
	c.buffer = c.buffer[:0]
	c.lastCollect = time.Now()

	return collected
}

// Start begins collecting events from the channel
// Should be run in its own goroutine
func (c *Collector) Start() {
	for event := range c.eventChannel {
		c.mu.Lock()
		c.buffer = append(c.buffer, event)
		c.mu.Unlock()
	}
}

// GetBufferSize returns current number of events in buffer
func (c *Collector) GetBufferSize() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.buffer)
}

// GetLastCollectTime returns when the last collection happened
func (c *Collector) GetLastCollectTime() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastCollect
}
