// +build go1.8

package kt

import (
	"strings"
	"testing"
)

func tName(t *testing.T) string {
	return t.Name()
}

func tSuite(t *testing.T) string {
	parts := strings.SplitN(t.Name(), "/", 2)
	return parts[0]
}
