package stat

import (
	"sync"
)

// Counter maintains primer and message counts for tests by ID.
type Counter struct {
	mu     sync.Mutex
	counts map[string]*count
}

// NewCounter initializes and returns a new counter.
func NewCounter() *Counter {
	return &Counter{
		counts: make(map[string]*count),
	}
}

// Add will add the given primer and message counts to the counts with the given
// ID.
func (c *Counter) Add(id string, primers, msgs int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ct, ok := c.counts[id]
	if !ok {
		ct = &count{}
		c.counts[id] = ct
	}

	ct.primers += primers
	ct.messages += msgs

}

// Counts will return the primer and message counts for a given ID.
func (c *Counter) Counts(id string) (primers, messages int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ct, ok := c.counts[id]
	if !ok {
		return 0, 0
	}
	primers, messages = ct.primers, ct.messages

	return primers, messages
}

type count struct {
	primers  int
	messages int
}
