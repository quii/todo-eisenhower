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

// Adaptive color palette - automatically adjusts for light/dark backgrounds
var (
	// Text colors
	TextPrimary = lipgloss.AdaptiveColor{
		Light: "#000000", // Black text on light background
		Dark:  "#FFFFFF", // White text on dark background
	}
	TextSecondary = lipgloss.AdaptiveColor{
		Light: "#666666", // Dark gray on light background
		Dark:  "#888888", // Light gray on dark background
	}
	TextMuted = lipgloss.AdaptiveColor{
		Light: "#999999", // Medium gray on light background
		Dark:  "#666666", // Darker gray on dark background
	}

	// Border and UI element colors
	BorderColor = lipgloss.AdaptiveColor{
		Light: "#CCCCCC", // Light gray border on light background
		Dark:  "#444444", // Dark gray border on dark background
	}
	BorderAccent = lipgloss.AdaptiveColor{
		Light: "#999999", // Darker border on light background
		Dark:  "#666666", // Lighter border on dark background
	}

	// Background colors for selections and highlights
	SelectionBg = lipgloss.AdaptiveColor{
		Light: "#E0E0E0", // Light gray selection on light background
		Dark:  "#444444", // Dark gray selection on dark background
	}

	// Completed todo color
	CompletedColor = lipgloss.AdaptiveColor{
		Light: "#999999", // Medium gray on light background
		Dark:  "#808080", // Light gray on dark background
	}

	// Empty state color
	EmptyColor = lipgloss.AdaptiveColor{
		Light: "#AAAAAA", // Medium gray on light background
		Dark:  "#666666", // Darker gray on dark background
	}

	// Due date colors (these are bright enough to work in both modes)
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
		
		// Use adaptive text color: black on light backgrounds, white on dark
		style := lipgloss.NewStyle().
			Background(bgColor).
			Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"})
		result.WriteString(style.Render(string(ch)))
	}
	
	return result.String()
}
