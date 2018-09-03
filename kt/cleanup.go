package kt

import (
	"context"
	"fmt"
)

// DestroyDB cleans up the specified DB after tests run
func (c *Context) DestroyDB(name string) {
	err := Retry(func() error {
		return c.Admin.DestroyDB(context.Background(), name, c.Options("db"))
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to clean up db '%s': %s", name, err))
	}
}
