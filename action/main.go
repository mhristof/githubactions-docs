package action

import (
	"bytes"
	"text/template"

	"gopkg.in/yaml.v2"
)

type Input struct {
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
	Default     string `yaml:"default"`
}

type Output struct {
	Description string `yaml:"description"`
	Value       string `yaml:"value"`
}

type Config struct {
	Description string            `yaml:"description"`
	Inputs      map[string]Input  `yaml:"inputs"`
	Name        string            `yaml:"name"`
	Outputs     map[string]Output `yaml:"outputs"`
	Runs        struct {
		Steps []struct {
			ID    string `yaml:"id"`
			Run   string `yaml:"run"`
			Shell string `yaml:"shell"`
		} `yaml:"steps"`
		Using string `yaml:"using"`
	} `yaml:"runs"`
}

func Load(yamlData []byte) (*Config, error) {
	var ret Config
	err := yaml.Unmarshal(yamlData, &ret)

	if err != nil {
		return nil, err
	}

	return &ret, nil
}

var md = `# {{ .Name }}

{{ .Description }}

{{ if .Inputs -}}
## Inputs

| Name | Description | Default | Required |
| ---- | ----------- | ------- | -------- |
{{ range $key, $value  := .Inputs -}}
| {{ $key }} | {{$value.Description}} | {{$value.Default }} | {{ $value.Required }}|
{{ end }}
{{- end }}

{{ if .Outputs -}}
## Outputs

| Name | Description |
| ---- | ----------- |
{{ range $key, $value  := .Outputs -}}
| {{ $key }} | {{$value.Description}} |
{{ end }}
{{- end -}}`

func (c *Config) Markdown() string {
	tmpl, err := template.New("markdown").Parse(md)
	if err != nil {
		panic(err)
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, struct {
		Name        string
		Description string
		Inputs      map[string]Input
		Outputs     map[string]Output
	}{
		Name:        c.Name,
		Description: c.Description,
		Inputs:      c.Inputs,
		Outputs:     c.Outputs,
	})

	return tpl.String()
}
