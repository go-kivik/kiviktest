package kt

import (
	"context"
	"net/url"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/flimzy/kivik"
)

const envPrefix = "KIVIK_TEST_DSN"

// getDSNs takes the result of os.Environ() and returns a list of values for
// keys in sorted order. KIVIK_TEST_DSN's value is returned first, followed by
// any other KIVIK_TEST_DSN_* keys, in reverse sorted order.
func getDSN(env []string) string {
	dsnMap := make(map[string]string)
	for _, kv := range env {
		if !strings.HasPrefix(kv, envPrefix) {
			continue
		}
		parts := strings.SplitN(kv, "=", 2)
		if parts[0] == envPrefix {
			return parts[1]
		}
		dsnMap[parts[0]] = parts[1]
	}
	if len(dsnMap) == 0 {
		return ""
	}
	keys := make([]string, 0, len(dsnMap))
	for key := range dsnMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return dsnMap[keys[len(keys)-1]]
}

// DSN returns a testing DSN from the environment.
func DSN(t *testing.T) string {
	if dsn := getDSN(os.Environ()); dsn != "" {
		return dsn
	}
	t.Skip("DSN not set")
	return ""
}

// NoAuthDSN returns a testing DSN with credentials stripped.
func NoAuthDSN(t *testing.T) string {
	dsn := DSN(t)
	parsed, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("invalid DSN: %s", err)
	}
	parsed.User = nil
	return parsed.String()
}

func connect(dsn string, t *testing.T) *kivik.Client {
	client, err := kivik.New(context.Background(), "couch", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to '%s': %s", dsn, err)
	}
	return client
}

// GetClient returns a connection to a CouchDB client, for testing.
func GetClient(t *testing.T) *kivik.Client {
	return connect(DSN(t), t)
}

// GetNoAuthClient returns an unauthenticated connection to a CouchDB client, for testing.
func GetNoAuthClient(t *testing.T) *kivik.Client {
	return connect(NoAuthDSN(t), t)
}
