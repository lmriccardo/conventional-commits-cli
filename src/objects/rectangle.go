package objects

import (
	"example.com/ccommits/src/styles"
	"github.com/gdamore/tcell/v2"
)

// A simple rectangle
type Rectangle struct {
	Width   int // The width of the rectangle
	Height  int // The height of the rectangle
	Start_x int // The starting x position
	Start_y int // The starting y position
}

// Draw a simple rectangle
func (rec Rectangle) DrawRectangle(screen tcell.Screen) {
	// Draw the top border
	screen.SetContent(rec.Start_x, rec.Start_y, '┌', nil, styles.BorderStyle)
	screen.SetContent(rec.Start_x+rec.Width-1, rec.Start_y, '┐', nil, styles.BorderStyle)
	for i := 1; i < rec.Width-1; i++ {
		screen.SetContent(rec.Start_x+i, rec.Start_y, '─', nil, styles.BorderStyle)
	}

	// Draw the sides
	for i := 1; i < rec.Height-1; i++ {
		screen.SetContent(rec.Start_x, rec.Start_y+i, '│', nil, styles.BorderStyle)
		screen.SetContent(rec.Start_x+rec.Width-1, rec.Start_y+i, '│', nil, styles.BorderStyle)
	}

	// Draw the bottom border
	screen.SetContent(rec.Start_x, rec.Start_y+rec.Height-1, '└', nil, styles.BorderStyle)
	screen.SetContent(rec.Start_x+rec.Width-1, rec.Start_y+rec.Height-1, '┘', nil, styles.BorderStyle)
	for i := 1; i < rec.Width-1; i++ {
		screen.SetContent(rec.Start_x+i, rec.Start_y+rec.Height-1, '─', nil, styles.BorderStyle)
	}
}
