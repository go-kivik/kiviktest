//go:build !js
// +build !js

package client

import (
	kivik "github.com/go-kivik/kivik/v3"
	"github.com/go-kivik/kiviktest/v3/kt"
)

func replicationOptions(_ *kt.Context, _ *kivik.Client, _, _, _ string, in map[string]interface{}) map[string]interface{} {
	if in == nil {
		in = make(map[string]interface{})
	}
	return in
}
