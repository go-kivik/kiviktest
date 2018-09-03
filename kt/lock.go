package kt

import "sync"

// Lock establishes a lock on name.
func (c *Context) Lock(name string) {
	c.getMU(name).Lock()
}

// RLock establishes a read lock on name.
func (c *Context) RLock(name string) {
	c.getMU(name).RLock()
}

// Unlock releases a lock on name.
func (c *Context) Unlock(name string) {
	c.getMU(name).Unlock()
}

// RUnlock releases a read lock on name.
func (c *Context) RUnlock(name string) {
	c.getMU(name).RUnlock()
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
