package manifest

import (
	"gopkg.in/yaml.v2"
)

type Manifest struct {
	ID          string              `yaml:"id"`
	Name        string              `yaml:"name"`
	Description string              `yaml:"description"`
	Variables   map[string]Variable `yaml:"variables"`
	Files       []File              `yaml:"files"`
}

type Variable struct {
	Type    string   `yaml:"type"`
	Options []Option `yaml:"options,omitempty"`
	Default string   `yaml:"default,omitempty"`
}

type Option struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type File struct {
	Path    string `yaml:"path"`
	Content string `yaml:"content,omitempty"`
	Source  string `yaml:"source,omitempty"`
}

func Parse(data []byte) (*Manifest, error) {
	var m Manifest
	err := yaml.Unmarshal(data, &m)
	return &m, err
}
