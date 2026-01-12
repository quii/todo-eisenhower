package ui

import (
	"hash/fnv"

	"github.com/charmbracelet/lipgloss"
)

// Bright color palette for tags
var tagColors = []string{
	"#FF6B6B", // Red
	"#4ECDC4", // Teal
	"#FFE66D", // Yellow
	"#95E1D3", // Light green
	"#FF8C94", // Pink
	"#A8E6CF", // Mint
	"#FFD93D", // Golden
	"#6BCF7F", // Green
	"#C7CEEA", // Lavender
	"#FFDAC1", // Peach
	"#B4F8C8", // Pale green
	"#FBE7C6", // Cream
	"#A0C4FF", // Light blue
	"#CAFFBF", // Light lime
	"#FDFFB6", // Pale yellow
	"#FFD6A5", // Light orange
}

// HashColor generates a consistent color for a given tag name
func HashColor(tag string) lipgloss.Color {
	h := fnv.New32a()
	h.Write([]byte(tag))
	hash := h.Sum32()
	colorIndex := int(hash) % len(tagColors)
	return lipgloss.Color(tagColors[colorIndex])
}
