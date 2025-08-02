package types

// Project represents a collection of requirements for a Product Requirements Document
type Project struct {
	Name         string        `yaml:"name"`
	IDPrefix     string        `yaml:"id_prefix"`
	Requirements []Requirement `yaml:"requirements,omitempty"`
}

// Requirement represents a single requirement in a Product Requirements Document
type Requirement struct {
	ID          string        `yaml:"id"`
	Title       string        `yaml:"title"`
	Keyword     string        `yaml:"keyword"`
	Description string        `yaml:"description"`
	Children    []Requirement `yaml:"children,omitempty"`
}