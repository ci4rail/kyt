package manifest

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
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

type customerManifestWithTenant struct {
	Application string       `json:"application"`
	Modules     []ModuleType `json:"modules"`
	Tenant      string       `json:"tenandId"`
	Now         string       `json:"now"`
}

var fns = template.FuncMap{
	"minus1": func(x int) int {
		return x - 1
	},
}

const (
	layeredManifestTemplate = `
    {
        "content": {
            "modulesContent": {
                "$edgeAgent": {
					{{ $tenant := .Tenant }}
					{{ $application := .Application }}
                    {{ $n := len .Modules }}
                    {{ range $i, $element := .Modules }}
                    "properties.desired.modules.{{$tenant}}_{{$application}}_{{$element.Name}}": {
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
func CreateLayeredManifest(c *CustomerManifest, tenantID string) (string, error) {
	ct := &customerManifestWithTenant{
		Application: strings.ToLower(c.Application),
		Modules:     c.Modules,
		Tenant:      tenantID,
	}
	t, err := template.New("LayeredDeployment").Funcs(fns).Parse(layeredManifestTemplate)
	if err != nil {
		return "", err
	}
	var result bytes.Buffer
	err = t.Execute(&result, ct)
	if err != nil {
		return "", err
	}
	fmt.Println(result.String())
	return result.String(), nil
}
