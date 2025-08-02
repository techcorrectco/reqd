package types

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Project represents a collection of requirements for a Product Requirements Document
type Project struct {
	Name         string        `yaml:"name"`
	IDPrefix     string        `yaml:"id_prefix"`
	Requirements []Requirement `yaml:"requirements,omitempty"`
}

// LoadProject loads a project from requirements.yaml file
func LoadProject() (*Project, error) {
	const filename = "requirements.yaml"

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var project Project
	if err := yaml.Unmarshal(data, &project); err != nil {
		return nil, err
	}

	return &project, nil
}

// Save saves the project to requirements.yaml file
func (p *Project) Save() error {
	const filename = "requirements.yaml"

	data, err := yaml.Marshal(p)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// Requirement represents a single requirement in a Product Requirements Document
type Requirement struct {
	ID          string        `yaml:"id"`
	Title       string        `yaml:"title"`
	Keyword     string        `yaml:"keyword"`
	Description string        `yaml:"description"`
	Children    []Requirement `yaml:"children,omitempty"`
}