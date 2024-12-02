package display

import (
	"log"

	"example.com/ccommits/src/styles"
	"github.com/gdamore/tcell/v2"
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
	for i, char := range content {
		screen.SetContent(start_x+i, start_y, char, nil, style)
	}
}

func DisplayTitle(screen tcell.Screen, title string) {
	width, _ := screen.Size()     // Get the width and the height of the screen
	titlelen := len(title) / 2    // Get the length of the title
	start_x := width/2 - titlelen // The starting position of the title
	var start_y int = 3           // The starting y position of the title

	// Draw the title
	DrawString(screen, title, start_x, start_y, styles.TitleStyle)
}
