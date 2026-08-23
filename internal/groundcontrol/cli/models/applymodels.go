package models

import "gopkg.in/yaml.v3"

type Action string

const (
	ActionCreated    Action = "created"
	ActionConfigured Action = "configured"
	ActionUnchanged  Action = "unchanged"
	ActionError      Action = "error"
	ActionNil        Action = "nil"
)

type ApplyResult struct {
	Kind   string `json:"kind"`
	Name   string `json:"name"`
	Action Action `json:"action"`
	Error  error  `json:"error,omitempty"`
}

type Envelope struct {
	APIVersion string    `yaml:"apiVersion"`
	Kind       string    `yaml:"kind"`
	Metadata   Metadata  `yaml:"metadata"`
	Spec       yaml.Node `yaml:"spec"`
}

type Metadata struct {
	Name   string            `yaml:"name"`
	Labels map[string]string `yaml:"labels,omitempty"`
}

type GroupSpec struct {
	RegistryURL string   `yaml:"registryURL"`
	Projects    []string `yaml:"projects"`
}

type ConfigSpec struct {
	RegistryURL string         `yaml:"registryURL"`
	Config      map[string]any `yaml:"config"`
}

type SatelliteSpec struct {
	Groups []string `yaml:"groups"`
	Config string   `yaml:"config,omitempty"`
}

type YAMLResource struct {
	APIVersion string
	Kind       string
	Metadata   Metadata
	Spec       any
}
