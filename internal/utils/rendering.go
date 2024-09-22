package utils

import "strings"

// TextEllipsis checks if the text is longer than the width of its container and
// adds a "..." as the last characters accordingly.
func TextEllipsis(text string, charactersLimit int) string {
	if len(text) >= charactersLimit {
		chopped := text[:charactersLimit]
		if len(chopped) < 3 {
			return "..."
		}
		
		return chopped[:len(chopped) - 3] + "..."
	}
	
	return text
}

// CenterText centers the text inside the widget
func CenterText(text string, width int, withBorders bool) string {
	if len(text) >= width {
		return text // If text is longer than the width, no need to center
	}
	
	// Subtract the borders if they exist.
	if withBorders {
		width = width - 2
	}
	
	// Calculate the padding needed on each side
	padding := (width - len(text)) / 2
	
	leftPadding := strings.Repeat(" ", padding)
	rightPadding := strings.Repeat(" ", width-len(text)-padding)

	// Return the padded text
	return leftPadding + text + rightPadding
}