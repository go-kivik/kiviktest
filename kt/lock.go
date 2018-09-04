package kt

import "sync"

// Lock establishes a lock on name.
func (c *Context) Lock(name string) {
	l := c.getMU(name)
	l.Lock()
	c.locks = append(c.locks, l.Unlock)
}

// RLock establishes a read lock on name.
func (c *Context) RLock(name string) {
	l := c.getMU(name)
	l.RLock()
	c.locks = append(c.locks, l.RUnlock)
}

// Unlock unlocks any open locks established at this test level.
func (c *Context) Unlock() {
	for _, l := range c.locks {
		l()
	}
}

func (c *Context) getMU(name string) *sync.RWMutex {
	lockName := tSuite(c.T) + "/" + name
	if c.mus == nil {
		c.mus = make(map[string]*sync.RWMutex)
	}
	mu, ok := c.mus[lockName]
	if !ok {
		mu = &sync.RWMutex{}
		c.mus[lockName] = mu
	}
	return mu
}
