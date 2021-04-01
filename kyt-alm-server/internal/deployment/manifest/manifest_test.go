package manifest

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateLayeredManifest(t *testing.T) {
	assert := assert.New(t)

	c := CustomerManifest{
		Application: "myapplication",
		Modules: []ModuleType{
			{
				Name:            "module1",
				Image:           "image1",
				CreateOptions:   "{}",
				ImagePullPolicy: "image_pull_1",
				RestartPolicy:   "policy1",
				Status:          "status1",
				StartupOrder:    1,
				Envs: map[string]string{
					"ENV1": "value1",
					"ENV2": "value2",
				},
			},
			{
				Name:            "module2",
				Image:           "image2",
				CreateOptions:   "{}",
				ImagePullPolicy: "image_pull_2",
				RestartPolicy:   "policy2",
				Status:          "status2",
				StartupOrder:    2,
			},
		},
	}
	layeredManifest, err := CreateLayeredManifest(c, "mytenantid")

	assert.Nil(err)
	fmt.Println(layeredManifest)

	objs := make(map[string]interface{})
	if err := json.Unmarshal([]byte(layeredManifest), &objs); err != nil {
		panic(err)
	}
	content := objs["content"].(map[string]interface{})
	modulesContent := content["modulesContent"].(map[string]interface{})
	edgeAgent := modulesContent["$edgeAgent"].(map[string]interface{})
	module1 := edgeAgent["properties.desired.modules.myapplication_module1"].(map[string]interface{})
	assert.Equal(module1["type"], "docker")
	assert.Equal(module1["status"], "status1")
	assert.Equal(module1["restartPolicy"], "policy1")
	// In reality this isn't float. But because we are unmarshaling without a struct, it can't know.
	assert.Equal(module1["startupOrder"], float64(1))
	assert.Equal(module1["imagePullPolicy"], "image_pull_1")
	settings1 := module1["settings"].(map[string]interface{})
	assert.Equal(settings1["image"], "image1")
	assert.Contains(settings1["createOptions"], fmt.Sprintf("host.docker.internal:%s", dockerHostIP))
	env1 := module1["env"].(map[string]interface{})
	env1env1 := env1["ENV1"].(map[string]interface{})
	env1env2 := env1["ENV2"].(map[string]interface{})
	assert.Equal(env1env1["value"], "value1")
	assert.Equal(env1env2["value"], "value2")

	module2 := edgeAgent["properties.desired.modules.myapplication_module2"].(map[string]interface{})
	assert.Equal(module2["type"], "docker")
	assert.Equal(module2["status"], "status2")
	assert.Equal(module2["restartPolicy"], "policy2")
	// In reality this isn't float. But because we are unmarshaling without a struct, it can't know.
	assert.Equal(module2["startupOrder"], float64(2))
	assert.Equal(module2["imagePullPolicy"], "image_pull_2")
	settings2 := module2["settings"].(map[string]interface{})
	assert.Equal(settings2["image"], "image2")
	assert.Contains(settings1["createOptions"], fmt.Sprintf("host.docker.internal:%s", dockerHostIP))
}

func TestAppendToCreateOptions(t *testing.T) {
	assert := assert.New(t)
	CreateOptions := "{\"Tty\":true,\"HostConfig\":{\"PortBindings\":{\"22/tcp\":[{\"HostPort\":\"11022\"}]}}}"
	str, err := appendExtraHosts(CreateOptions)
	assert.Nil(err)
	fmt.Println(str)
	assert.Contains(str, "ExtraHosts")
	assert.Contains(str, "PortBindings")
	assert.Contains(str, fmt.Sprintf("host.docker.internal:%s", dockerHostIP))
}

func TestAppendToCreateOptionsExtraHostExisting(t *testing.T) {
	assert := assert.New(t)
	CreateOptions := "{\"HostConfig\":{\"ExtraHosts\":[\"myfancyhost:unicorn\"]}}"
	str, err := appendExtraHosts(CreateOptions)
	assert.Nil(err)
	fmt.Println(str)
	assert.Contains(str, "ExtraHosts")
	assert.Contains(str, "myfancyhost:unicorn")
	assert.Contains(str, fmt.Sprintf("host.docker.internal:%s", dockerHostIP))
}

func TestAppendToCreateOptionsNoHostConfig(t *testing.T) {
	assert := assert.New(t)
	CreateOptions := "{}"
	str, err := appendExtraHosts(CreateOptions)
	assert.Nil(err)
	fmt.Println(str)
	assert.Contains(str, "ExtraHosts")
	assert.Contains(str, fmt.Sprintf("host.docker.internal:%s", dockerHostIP))
}
