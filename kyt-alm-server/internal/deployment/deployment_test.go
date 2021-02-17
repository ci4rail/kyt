/*
Copyright © 2021 Ci4Rail GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package deployment

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDeployment(t *testing.T) {
	assert := assert.New(t)
	if err := os.Setenv("IOTHUB_SERVICE_CONNECTION_STRING", "HostName=myHub.azure-devices.net;SharedAccessKeyName=myPolicy;SharedAccessKey=asdfasdfasdfasdfasdfasdfBasdfasdfasdfasdfas="); err != nil {
		t.Fatalf("Setenv: %v", err)
	}

	type testManifest struct {
		Key string `json:key`
	}
	m := &testManifest{
		Key: "myValue",
	}
	j, err := json.Marshal(m)
	if err != nil {
		t.Error(err)
	}
	d, err := NewDeployment(string(j), "myDeployment")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(d.connectionString, "HostName=myHub.azure-devices.net;SharedAccessKeyName=myPolicy;SharedAccessKey=asdfasdfasdfasdfasdfasdfBasdfasdfasdfasdfas=", "They should be equal")
	assert.Equal(d.name, "myDeployment", "they should be equal")
	var u testManifest
	err = json.Unmarshal([]byte(d.manifest), &u)
	assert.Nil(err, "should be nil")
	assert.Equal(u.Key, "myValue", "they should be equal")
}
