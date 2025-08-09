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