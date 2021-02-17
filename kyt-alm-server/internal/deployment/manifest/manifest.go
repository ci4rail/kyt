package manifest

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
)

// CustomerManifest describes the format the customer defines modules
type CustomerManifest struct {
	Application string       `json:"application"`
	Modules     []ModuleType `json:"modules"`
}

// ModuleType contains single module options
type ModuleType struct {
	Name            string `json:"name"`
	Image           string `json:"image"`
	CreateOptions   string `json:"createOptions"`
	ImagePullPolicy string `json:"imagePullPolicy"`
	RestartPolicy   string `json:"restartPolicy"`
	Status          string `json:"status"`
	StartupOrder    int    `json:"startupOrder"`
}

type CustomerManifestWithTenant struct {
	Application string       `json:"application"`
	Modules     []ModuleType `json:"modules"`
	Tenant      string       `json:"tenand_id"`
	Now         string       `json:"now"`
}

var fns = template.FuncMap{
	"minus1": func(x int) int {
		return x - 1
	},
}

const (
	LayeredManifestTemplate = `
    {    
        "content": {
            "modulesContent": {
                "$edgeAgent": {
					{{ $tenant := .Tenant }}
					{{ $now := .Now }}
					{{ $application := .Application }}
                    {{ $n := len .Modules }}
                    {{ range $i, $element := .Modules }}
                    "properties.desired.modules.{{$tenant}}.{{$application}}.{{$element.Name}}.{{$now}}": {
                        "settings": {
                            "image": "{{ $element.Image }}",
                            "createOptions": "{{ $element.CreateOptions }}"
                        },
                        "type": "docker",
                        "imagePullPolicy": "{{ $element.ImagePullPolicy }}",
                        "status": "{{ $element.Status }}",
                        "restartPolicy": "{{ $element.RestartPolicy }}",
                        "startupOrder": {{ $element.StartupOrder }}
                    }{{ if lt $i (minus1 $n) }},{{ end }}
                    {{ end }}
                }
            }
        }
    }`
)

// CreateLayeredManifest creates a new LayeredManifest from a CustomerManifest
func CreateLayeredManifest(c *CustomerManifest, tenantId string) (string, error) {
	ct := &CustomerManifestWithTenant{
		Application: c.Application,
		Modules:     c.Modules,
		Tenant:      tenantId,
		Now:         fmt.Sprintf("%d", time.Now().Unix()),
	}
	t, err := template.New("LayeredDeployment").Funcs(fns).Parse(LayeredManifestTemplate)
	if err != nil {
		return "", err
	}
	var result bytes.Buffer
	err = t.Execute(&result, ct)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}
