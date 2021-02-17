package manifest

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateLayeredManifest(t *testing.T) {
	assert := assert.New(t)

	c := CustomerManifest{
		Application: "myApplication",
		Modules: []ModuleType{
			{
				Name:            "module1",
				Image:           "image1",
				CreateOptions:   "{createOptions1}",
				ImagePullPolicy: "image_pull_1",
				RestartPolicy:   "policy1",
				Status:          "status1",
				StartupOrder:    1,
			},
			{
				Name:            "module2",
				Image:           "image2",
				CreateOptions:   "{createOptions2}",
				ImagePullPolicy: "image_pull_2",
				RestartPolicy:   "policy2",
				Status:          "status2",
				StartupOrder:    2,
			},
		},
	}
	layeredManifest, err := CreateLayeredManifest(&c, "myTenantId")

	assert.Nil(err)
	fmt.Println(layeredManifest)

	objs := make(map[string]interface{})
	if err := json.Unmarshal([]byte(layeredManifest), &objs); err != nil {
		panic(err)
	}
	fmt.Println(objs)
	now := fmt.Sprintf("%d", time.Now().Unix())
	content := objs["content"].(map[string]interface{})
	modulesContent := content["modulesContent"].(map[string]interface{})
	edgeAgent := modulesContent["$edgeAgent"].(map[string]interface{})
	module1 := edgeAgent["properties.desired.modules.myTenantId.myApplication.module1."+now].(map[string]interface{})
	assert.Equal(module1["type"], "docker")
	assert.Equal(module1["status"], "status1")
	assert.Equal(module1["restartPolicy"], "policy1")
	// In reality this isn't float. But because we are unmarshaling without a struct, it can't know.
	assert.Equal(module1["startupOrder"], float64(1))
	assert.Equal(module1["imagePullPolicy"], "image_pull_1")
	settings1 := module1["settings"].(map[string]interface{})
	assert.Equal(settings1["image"], "image1")
	assert.Equal(settings1["createOptions"], "{createOptions1}")

	module2 := edgeAgent["properties.desired.modules.myTenantId.myApplication.module2."+now].(map[string]interface{})
	assert.Equal(module2["type"], "docker")
	assert.Equal(module2["status"], "status2")
	assert.Equal(module2["restartPolicy"], "policy2")
	// In reality this isn't float. But because we are unmarshaling without a struct, it can't know.
	assert.Equal(module2["startupOrder"], float64(2))
	assert.Equal(module2["imagePullPolicy"], "image_pull_2")
	settings2 := module2["settings"].(map[string]interface{})
	assert.Equal(settings2["image"], "image2")
	assert.Equal(settings2["createOptions"], "{createOptions2}")
}
