package types

import (
	"reflect"
	"testing"
)

func Test_getBranches(t *testing.T) {
	tests := []struct {
		name         string
		requirements []Requirement
		expected     []Requirement
	}{
		{
			name:         "empty slice",
			requirements: []Requirement{},
			expected:     []Requirement{},
		},
		{
			name: "only leaf requirements",
			requirements: []Requirement{
				{ID: "1", Text: "Leaf 1"},
				{ID: "2", Text: "Leaf 2"},
			},
			expected: []Requirement{},
		},
		{
			name: "single branch with children",
			requirements: []Requirement{
				{
					ID:   "1",
					Text: "Branch 1",
					Children: []Requirement{
						{ID: "1.1", Text: "Child 1"},
						{ID: "1.2", Text: "Child 2"},
					},
				},
			},
			expected: []Requirement{
				{
					ID:   "1",
					Text: "Branch 1",
					Children: []Requirement{
						{ID: "1.1", Text: "Child 1"},
						{ID: "1.2", Text: "Child 2"},
					},
				},
			},
		},
		{
			name: "mixed branches and leaves",
			requirements: []Requirement{
				{ID: "1", Text: "Leaf 1"},
				{
					ID:   "2",
					Text: "Branch 1",
					Children: []Requirement{
						{ID: "2.1", Text: "Child 1"},
						{ID: "2.2", Text: "Child 2"},
					},
				},
				{ID: "3", Text: "Leaf 2"},
			},
			expected: []Requirement{
				{
					ID:   "2",
					Text: "Branch 1",
					Children: []Requirement{
						{ID: "2.1", Text: "Child 1"},
						{ID: "2.2", Text: "Child 2"},
					},
				},
			},
		},
		{
			name: "nested branches",
			requirements: []Requirement{
				{
					ID:   "1",
					Text: "Root Branch",
					Children: []Requirement{
						{ID: "1.1", Text: "Leaf Child"},
						{
							ID:   "1.2",
							Text: "Nested Branch",
							Children: []Requirement{
								{ID: "1.2.1", Text: "Deep Child"},
							},
						},
					},
				},
			},
			expected: []Requirement{
				{
					ID:   "1",
					Text: "Root Branch",
					Children: []Requirement{
						{ID: "1.1", Text: "Leaf Child"},
						{
							ID:   "1.2",
							Text: "Nested Branch",
							Children: []Requirement{
								{ID: "1.2.1", Text: "Deep Child"},
							},
						},
					},
				},
				{
					ID:   "1.2",
					Text: "Nested Branch",
					Children: []Requirement{
						{ID: "1.2.1", Text: "Deep Child"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getBranches(tt.requirements)
			
			// Handle nil vs empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return
			}
			
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("getBranches() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func Test_findRequirement(t *testing.T) {
	requirements := []Requirement{
		{ID: "1", Text: "Root 1"},
		{
			ID:   "2",
			Text: "Root 2",
			Children: []Requirement{
				{ID: "2.1", Text: "Child 2.1"},
				{
					ID:   "2.2",
					Text: "Child 2.2",
					Children: []Requirement{
						{ID: "2.2.1", Text: "Grandchild 2.2.1"},
					},
				},
			},
		},
		{ID: "3", Text: "Root 3"},
	}

	tests := []struct {
		name         string
		requirements []Requirement
		id           string
		expected     *Requirement
	}{
		{
			name:         "empty requirements",
			requirements: []Requirement{},
			id:           "1",
			expected:     nil,
		},
		{
			name:         "find root level requirement",
			requirements: requirements,
			id:           "1",
			expected:     &Requirement{ID: "1", Text: "Root 1"},
		},
		{
			name:         "find child requirement",
			requirements: requirements,
			id:           "2.1",
			expected:     &Requirement{ID: "2.1", Text: "Child 2.1"},
		},
		{
			name:         "find nested child requirement",
			requirements: requirements,
			id:           "2.2.1",
			expected:     &Requirement{ID: "2.2.1", Text: "Grandchild 2.2.1"},
		},
		{
			name:         "find parent with children",
			requirements: requirements,
			id:           "2.2",
			expected: &Requirement{
				ID:   "2.2",
				Text: "Child 2.2",
				Children: []Requirement{
					{ID: "2.2.1", Text: "Grandchild 2.2.1"},
				},
			},
		},
		{
			name:         "non-existent id",
			requirements: requirements,
			id:           "999",
			expected:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findRequirement(tt.requirements, tt.id)
			
			if tt.expected == nil && result == nil {
				return
			}
			
			if tt.expected == nil && result != nil {
				t.Errorf("findRequirement() = %v, want nil", result)
				return
			}
			
			if tt.expected != nil && result == nil {
				t.Errorf("findRequirement() = nil, want %v", tt.expected)
				return
			}
			
			if !reflect.DeepEqual(*result, *tt.expected) {
				t.Errorf("findRequirement() = %v, want %v", *result, *tt.expected)
			}
		})
	}
}