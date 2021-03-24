package action

import (
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/stretchr/testify/assert"
)

var githubExample = heredoc.Doc(`
# example from https://docs.github.com/en/free-pro-team@latest/actions/creating-actions/creating-a-composite-run-steps-action#creating-an-action-metadata-file
name: 'Hello World'
description: 'Greet someone'
inputs:
  who-to-greet:  # id of input
    description: 'Who to greet'
    required: true
    default: 'World'
outputs:
  random-number: 
    description: "Random number"
    value: ${{ steps.random-number-generator.outputs.random-id }}
runs:
  using: "composite"
  steps: 
  - run: echo Hello ${{ inputs.who-to-greet }}.
    shell: bash
  - id: random-number-generator
    run: echo "::set-output name=random-id::$(echo $RANDOM)"
    shell: bash
  - run: ${{ github.action_path }}/goodbye.sh
    shell: bash
`)

func TestLoadl(t *testing.T) {
	var cases = []struct {
		name string
		yaml string
	}{
		{
			name: "simple yaml config",
			yaml: githubExample,
		},
	}

	for _, test := range cases {
		_, err := Load([]byte(test.yaml))
		assert.Equal(t, nil, err, test.name)
	}
}

func TestMarkdown(t *testing.T) {
	var cases = []struct {
		name     string
		yaml     string
		expected string
	}{
		{
			name: "sample yaml config",
			yaml: githubExample,
			expected: heredoc.Doc(`
				# Hello World
				
				Greet someone
				
				## Inputs
				
				| Name | Description | Default | Required |
				| ---- | ----------- | ------- | -------- |
				| who-to-greet | Who to greet | World | true|
				
				## Outputs
				
				| Name | Description |
				| ---- | ----------- |
				| random-number | Random number |
			`),
		},
		{
			name: "no outputs",
			yaml: heredoc.Doc(`
				name: 'Hello World'
				description: 'Greet someone'
				inputs:
				  who-to-greet:  # id of input
				    description: 'Who to greet'
				    required: true
				    default: 'World'
			`),
			expected: heredoc.Doc(`
				# Hello World
				
				Greet someone
				
				## Inputs
				
				| Name | Description | Default | Required |
				| ---- | ----------- | ------- | -------- |
				| who-to-greet | Who to greet | World | true|
			`),
		},
		{
			name: "no inputs",
			yaml: heredoc.Doc(`
				name: 'Hello World'
				description: 'Greet someone'
				outputs:
				  random-number: 
				    description: "Random number"
				    value: ${{ steps.random-number-generator.outputs.random-id }}
			`),
			expected: heredoc.Doc(`
				# Hello World
				
				Greet someone

				## Outputs
				
				| Name | Description |
				| ---- | ----------- |
				| random-number | Random number |
			`),
		},
	}

	for _, test := range cases {
		cfg, err := Load([]byte(test.yaml))
		assert.Equal(t, nil, err, test.name)
		assert.Equal(t, test.expected, cfg.Markdown(), test.name)
	}
}
