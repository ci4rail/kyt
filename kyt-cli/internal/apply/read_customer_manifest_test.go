package apply

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	manifest string = `---
application: influxdb
modules:
  - name: influxdb
    image: influxdb:latest
    createOptions: '{\"HostConfig\":{\"PortBindings\":{\"8086/tcp\":[{\"HostPort\":\"8086\"}]}}}'
    imagePullPolicy: on-create
    restartPolicy: always
    status: running
    startupOrder: 1
    envs:
      ENV1: influx1
      ENV2: influx2
  - name: nginx
    image: nginx:1.19.6
    createOptions: '{\"HostConfig\":{\"PortBindings\":{\"80/tcp\":[{\"HostPort\":\"12010\"}]}}}'
    imagePullPolicy: on-create
    restartPolicy: always
    status: running
    startupOrder: 1
    envs:
      ENV1: nginx1`
)

func prepare() *os.File {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "prefix-")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}
	text := []byte(manifest)
	if _, err = tmpFile.Write(text); err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}

	return tmpFile
}

func cleanup(file *os.File) {
	// Close the file
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
	os.Remove(file.Name())
}

func TestReadCustomerManifestt(t *testing.T) {
	assert := assert.New(t)
	file := prepare()
	manifest := ReadCustomerManifest(file.Name())
	assert.Equal(manifest.Application, "influxdb")
	assert.Equal(manifest.Modules[0].Name, "influxdb")
	assert.Equal(manifest.Modules[0].Image, "influxdb:latest")
	fmt.Println(*manifest.Modules[0].CreateOptions)
	assert.Equal((*manifest.Modules[0].ImagePullPolicy), "on-create")
	assert.Equal((*manifest.Modules[0].RestartPolicy), "always")
	assert.Equal((*manifest.Modules[0].Status), "running")
	assert.Equal((*manifest.Modules[0].StartupOrder), int32(1))
	assert.Equal((*manifest.Modules[0].Envs)["ENV1"], "influx1")
	assert.Equal((*manifest.Modules[0].Envs)["ENV2"], "influx2")
	assert.Equal((*manifest.Modules[1].Envs)["ENV1"], "nginx1")
	cleanup(file)
}
