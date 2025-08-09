package types

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Project represents a collection of requirements for a Product Requirements Document
type Project struct {
	Name         string        `yaml:"name"`
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
	ID       string        `yaml:"id"`
	Text     string        `yaml:"text"`
	Children []Requirement `yaml:"children,omitempty"`
}

// FindRequirement finds a requirement by ID in the project's requirement tree
func (p *Project) FindRequirement(id string) *Requirement {
	return findRequirement(p.Requirements, id)
}

// findRequirement recursively searches for a requirement by ID
func findRequirement(requirements []Requirement, id string) *Requirement {
	for i := range requirements {
		req := requirements[i]
		if req.ID == id {
			return &requirements[i]
		}
		if found := findRequirement(req.Children, id); found != nil {
			return found
		}
	}
	return nil
}

// GetBranches returns only the requirements that have children (branches, not leaves)
func (p *Project) GetBranches() []Requirement {
	return getBranches(p.Requirements)
}

// getBranches recursively collects requirements that have children
func getBranches(requirements []Requirement) []Requirement {
	var branches []Requirement
	for _, req := range requirements {
		if len(req.Children) > 0 {
			branches = append(branches, req)
			branches = append(branches, getBranches(req.Children)...)
		}
	}
	return branches
}
