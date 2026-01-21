package ui

import (
	"testing"

	"github.com/matryer/is"
)

func TestTagInventory_IsHighWIP(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		name     string
		count    int
		expected bool
	}{
		{
			name:     "Zero items is not high WIP",
			count:    0,
			expected: false,
		},
		{
			name:     "One item is not high WIP",
			count:    1,
			expected: false,
		},
		{
			name:     "Five items is not high WIP (at threshold)",
			count:    5,
			expected: false,
		},
		{
			name:     "Six items is high WIP (exceeds threshold)",
			count:    6,
			expected: true,
		},
		{
			name:     "Ten items is high WIP",
			count:    10,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			inventory := TagInventory{
				Tag:   "test",
				Count: tt.count,
			}
			is.Equal(inventory.IsHighWIP(), tt.expected)
		})
	}
}
