package objects

import (
	"example.com/ccommits/src/display"
	"example.com/ccommits/src/styles"
	"github.com/gdamore/tcell/v2"
)

type TextBox struct {
	rec         Rectangle // The rectangle of the text box
	title       string    // The title of the text box (Optional)
	start_pos_x int       // Start position x of the first character
	start_pos_y int       // Start position y of the first character
	content     string    // The content of the text box
	curr_pos    int       // The current position of the cursor in the text
}

// Create a simple text box using a rectangle and leaving
// one single space for writing text into
func TextBox_new(title string, x, y, size_w, size_h int) *TextBox {
	// First, creates the rectangle
	rect := new(Rectangle)
	rect.Start_x = x
	rect.Start_y = y
	rect.Width = size_w
	rect.Height = size_h + 1

	// Compute the starting position of the text
	startpos_y := rect.Start_y + 2
	startpos_x := rect.Start_x + 2

	return &TextBox{rec: *rect, title: title, start_pos_x: startpos_x,
		start_pos_y: startpos_y, content: "", curr_pos: 0}
}

// Display the text box
func (tb TextBox) Display(screen tcell.Screen, style tcell.Style) {
	tb.rec.DrawRectangle(screen) // Draw the rectangle for the text box

	// Display the Title
	display.DrawString(screen, tb.title, tb.rec.Start_x+3,
		tb.rec.Start_y, style)

	// Display the content
	display.DrawString(screen, tb.content, tb.start_pos_x,
		tb.start_pos_y, styles.SimpleStyle)
}

// Check if the textbox collides with input coordinates
func (tb TextBox) IsColliding(x, y int) bool {
	return ((x > tb.rec.Start_x && x < tb.rec.Start_x+tb.rec.Width) ||
		(y > tb.rec.Start_y && y < tb.rec.Start_y+tb.rec.Height))
}

func (tb TextBox) HandleEvent(screen tcell.Screen, event *tcell.EventKey) {

}
