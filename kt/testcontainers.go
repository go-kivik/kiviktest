// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy of
// the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations under
// the License.

//go:build go1.18
// +build go1.18

package kt

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	servicePort        = "5984"
	adminUsername      = "admin"
	adminPassword      = "abc123"
	defaultVersion     = "3.3"
	replicatorInterval = `"1000"`
)

type version struct {
	image string
}

var supportedVersions = map[string]version{
	"2.2": {image: "couchdb:2.2.0"},
	"2.3": {image: "apache/couchdb:2.3.1"},
	"3.0": {image: "couchdb:3.0.1"},
	"3.1": {image: "couchdb:3.1.2"},
	"3.2": {image: "couchdb:3.2.3"},
	"3.3": {image: "couchdb3.3.2"},
}

// StartContainer starts a container running the requested version of CouchDB,
// and returns the DSN to connect to the instance. If version is empty, the
// latest version is used.
func StartContainer(version string) (string, error) {
	if version == "" {
		version = defaultVersion
	}
	ver, ok := supportedVersions[version]
	if !ok {
		return "", fmt.Errorf("unknown CouchDB version %s", version)
	}

	ctx := context.Background()

	containerName := "couchdb" + version

	couchContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: ver.image,
			Env: map[string]string{
				"COUCHDB_USER":     adminUsername,
				"COUCHDB_PASSWORD": adminPassword,
			},
			ExposedPorts: []string{servicePort + "/tcp"},
			Name:         fmt.Sprintf("%s-%v", containerName, os.Getpid()),
			WaitingFor: wait.
				ForListeningPort(servicePort + "/tcp").
				WithStartupTimeout(120 * time.Second),
			HostConfigModifier: func(cf *container.HostConfig) {
				cf.AutoRemove = true
			},
		},
		Started: true,
		Reuse:   true,
	})
	if err != nil {
		return "", err
	}

	host, err := couchContainer.Host(ctx)
	if err != nil {
		return "", err
	}
	port, err := couchContainer.MappedPort(ctx, servicePort)
	if err != nil {
		return "", err
	}
	addr := fmt.Sprintf("http://%s:%s@%s:%s/", adminUsername, adminPassword, host, port.Port())

	c := &http.Client{}
	for _, dbName := range []string{"_users", "_replicator", "_global_changes"} {
		if err := createDB(ctx, c, addr, dbName); err != nil {
			return "", err
		}
	}
	if err := setReplicatorInterval(ctx, c, addr); err != nil {
		return "", err
	}

	return addr, nil
}

func createDB(ctx context.Context, c *http.Client, addr, dbName string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, addr+dbName, nil)
	if err != nil {
		return err
	}
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	switch {
	case res.StatusCode == http.StatusPreconditionFailed:
		return nil
	case res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusMultipleChoices:
		return nil
	}
	return fmt.Errorf("unexpected response code %d", res.StatusCode)
}

func setReplicatorInterval(ctx context.Context, c *http.Client, addr string) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		addr+"_node/nonode@nohost/_config/replicator/interval",
		strings.NewReader(replicatorInterval),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusMultipleChoices {
		return nil
	}
	return fmt.Errorf("unexpected response code %d", res.StatusCode)
}
