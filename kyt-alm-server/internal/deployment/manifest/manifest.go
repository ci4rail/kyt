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
	Name            string            `json:"name"`
	Image           string            `json:"image"`
	CreateOptions   string            `json:"createOptions"`
	ImagePullPolicy string            `json:"imagePullPolicy"`
	RestartPolicy   string            `json:"restartPolicy"`
	Status          string            `json:"status"`
	StartupOrder    int               `json:"startupOrder"`
	Envs            map[string]string `json:"envs"`
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
	"add": func(x, y int) int {
		return x + y
	},
}

const (
	layeredManifestTemplate = `
{
	"content": {
		"modulesContent": {
		{{ $tenant := .Tenant }}
		{{ $application := .Application }}
		{{ $n := len .Modules }}
			"$edgeAgent": {
				{{ range $i, $element := .Modules }}
				{{ $envsLen := len $element.Envs }}
				"properties.desired.modules.{{$application}}_{{$element.Name}}": {
					{{ if $element.Envs }}
					"env": {
						{{ $counter:=0 }}
						{{ range $k, $v := $element.Envs }}
							"{{$k}}": {"value": "{{$v}}"}{{ if lt $counter (minus1 $envsLen) }},{{ end }}
							{{ $counter = (add $counter 1) }}
						{{ end }}
					},
					{{ end }}
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
			},
			"$edgeHub": {
				"properties.desired.routes.upstream": {
					"priority": 0,
					"route": "FROM /messages/* INTO $upstream",
					"timeToLiveSecs": 7200
				}
			},
			{{ range $i, $element := .Modules }}

			"{{$application}}_{{$element.Name}}": {
				"properties.desired": {
					"tenantId": "{{$tenant}}"
				}
			}{{ if lt $i (minus1 $n) }},{{ end }}
			{{ end }}
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
	fmt.Println(result.String())
	if err != nil {
		return "", err
	}
	return result.String(), nil
}
