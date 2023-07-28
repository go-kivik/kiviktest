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

package kt

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	servicePort   = "5984"
	adminUsername = "admin"
	adminPassword = "abc123"
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
// and returns the DSN to connect to the instance.
func StartContainer(version string) (string, error) {
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

	return addr, nil
}
