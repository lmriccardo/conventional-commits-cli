package display

import (
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

func InitializeScreen() tcell.Screen {
	screen, err := tcell.NewScreen() // Initialize a screen
	if err != nil {
		// Error logging when creating the screen
		log.Fatalf("Error creating screen: %v", err)
	}

	if err := screen.Init(); err != nil {
		// Log error when initializing the screen
		log.Fatalf("Error initializing screen: %v", err)
	}

	screen.Clear() // Clearup the screen

	// Set the cursor style for the screen
	screen.SetCursorStyle(tcell.CursorStyleBlinkingUnderline)
	screen.HideCursor()
	screen.EnableMouse()

	return screen
}

func DrawString(screen tcell.Screen, content string, start_x, start_y int, style tcell.Style) {
	prev_width := 0
	for _, char := range content {
		width := runewidth.RuneWidth(char)
		screen.SetContent(start_x+prev_width, start_y, char, nil, style)
		prev_width += width
	}
}

func CenterString(total_len int, content string) string {
	if total_len <= 1 {
		return content
	}

	// Compute the actual char dimension of the string
	actual_dim := runewidth.StringWidth(content)

	// Compute the padding at each side
	left_padding := (total_len-1)/2 - actual_dim/2
	right_padding := total_len - 2 - left_padding

	// Build the string
	return strings.Repeat(" ", left_padding) + content + strings.Repeat(" ", right_padding)
}
