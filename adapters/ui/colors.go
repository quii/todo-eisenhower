package ui

import (
	"hash/fnv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
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

// Due date colors
var (
	dueDateColor = lipgloss.Color("#00CED1") // DarkCyan - for upcoming due dates
	overdueColor = lipgloss.Color("#FF0000") // Red - for overdue dates
)

// HashColor generates a consistent color for a given tag name
func HashColor(tag string) lipgloss.Color {
	h := fnv.New32a()
	_, _ = h.Write([]byte(tag)) // hash.Hash.Write never errors
	hash := h.Sum32()
	colorIndex := int(hash) % len(tagColors)
	return lipgloss.Color(tagColors[colorIndex])
}

// lightenColor makes a color lighter by a factor (0.0 = original, 1.0 = white)
func lightenColor(c lipgloss.Color, factor float64) lipgloss.Color {
	col, _ := colorful.MakeColor(c)
	
	// Interpolate towards white
	white := colorful.Color{R: 1, G: 1, B: 1}
	lightened := col.BlendLab(white, factor)
	
	return lipgloss.Color(lightened.Hex())
}

// GradientBackground applies a horizontal gradient to text as a background
func GradientBackground(text string, startColor, endColor lipgloss.Color) string {
	if text == "" {
		return text
	}
	
	// Create gradient colors
	colors := gamut.Blends(startColor, endColor, len([]rune(text)))
	
	var result strings.Builder
	chars := []rune(text)
	
	for i, ch := range chars {
		// Convert gamut color to lipgloss color
		c, _ := colorful.MakeColor(colors[i])
		bgColor := lipgloss.Color(c.Hex())
		
		style := lipgloss.NewStyle().
			Background(bgColor).
			Foreground(lipgloss.Color("#FFFFFF"))
		result.WriteString(style.Render(string(ch)))
	}
	
	return result.String()
}
