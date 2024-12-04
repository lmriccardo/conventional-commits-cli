package objects

import (
	"example.com/ccommits/ccommits/display"
	"example.com/ccommits/ccommits/styles"
	"github.com/gdamore/tcell/v2"
)

type Point struct {
	X int // X Position
	Y int // Y Position
}

// Mapping event keys arrow to direction (x, y)
var ARROW_MAPPING map[tcell.Key]Point = map[tcell.Key]Point{
	tcell.KeyUp:    {0, -1},
	tcell.KeyDown:  {0, +1},
	tcell.KeyLeft:  {-1, 0},
	tcell.KeyRight: {+1, 0},
}

type TextBox struct {
	rec         Rectangle // The rectangle of the text box
	title       string    // The title of the text box (Optional)
	start_pos_x int       // Start position x of the first character
	start_pos_y int       // Start position y of the first character
	content     string    // The content of the text box
	curr_pos_x  int       // The current X position of the cursor in the text
	curr_pos_y  int       // The current Y position of the cursor in the text
	focus       bool      // If the current focus is on this object
	nof_lines   int       // Total number of lines
	curr_line   int       // Current line at which the cursor is positioned
	enters      []bool    // Number of spaces added for each line
	enter_idx   int       // Current enter pressed idx
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
		start_pos_y: startpos_y, content: "", curr_pos_x: startpos_x,
		curr_pos_y: startpos_y, nof_lines: 1, curr_line: 0,
		enters: make([]bool, 0), enter_idx: 0}
}

func (tb *TextBox) getMaxRowSize() int {
	return tb.rec.Width - 2*(tb.start_pos_x-tb.rec.Start_x)
}

func (tb *TextBox) getStringPosition(x, y int) int {
	curr_line := y - tb.start_pos_y
	row_size := tb.getMaxRowSize()                    // Max row size
	rel_cursor_x := x - tb.start_pos_x                // Relative X cursor position
	str_position := curr_line*row_size + rel_cursor_x // Position into the string
	return str_position
}

// func (tb *TextBox) getLastYStringPosition() int {
// 	row_size := tb.getMaxRowSize() // Max row size

// }

func (tb *TextBox) addEnterValue(value bool) {
	if tb.enter_idx+1 > len(tb.enters) {
		tb.enters = append(tb.enters, value)
		tb.enter_idx++
		return
	}

	tb.enters[tb.enter_idx] = value
	tb.enter_idx++
}

func (tb *TextBox) drawContent(screen tcell.Screen) {
	row_size := tb.getMaxRowSize()      // Get the maximum row size
	content_array := []rune(tb.content) // From string to rule array
	content_len := len(tb.content)      // Take the length of the string

	for row_idx := 0; row_idx < tb.nof_lines; row_idx++ {
		// Get the correct start and stop indexes
		start_idx := row_idx * row_size
		stop_idx := min(content_len, (row_idx+1)*row_size)
		substr := string(content_array[start_idx:stop_idx])

		// Draw the substring content into the screen
		display.DrawString(screen, substr, tb.start_pos_x,
			tb.start_pos_y+row_idx, styles.SimpleStyle)
	}
}

func (tb *TextBox) handleEnterKey(screen tcell.Screen) {
	prev_pos_y := tb.curr_pos_y // Save the previous y position

	// Update the next cursor position
	tb.SetNextCursorPosition(tb.start_pos_x, tb.curr_pos_y+1)
	screen.ShowCursor(tb.curr_pos_x, tb.curr_pos_y)

	// The following code must be evaluated only if the
	// new y position of the cursor have changed
	if prev_pos_y != tb.curr_pos_y {
		// We need to save the total number of spaces added
		// when jumping into the next line
		tb.addEnterValue(true)
		tb.nof_lines++
		tb.curr_line++
	}
}

func (tb *TextBox) appendCharacter(screen tcell.Screen, char rune) {
	// First we need to check if the input char can be added
	// to the content string. This is possible only if the next
	// position of the cursor is at the edge of the rectangle
	next_x := tb.curr_pos_x + 1
	next_y := tb.curr_pos_y

	if next_x == tb.rec.Start_x+tb.rec.Width-1 {
		next_y = next_y + 1
		next_x = tb.start_pos_x + 1
	}

	// Check the new cursor position
	if !tb.CheckNextCursorPosition(next_x, next_y) {
		return
	}

	if next_y > tb.curr_pos_y {
		tb.nof_lines++
		tb.curr_line++
	}

	// Set the new values for the cursor positions
	tb.SetNextCursorPosition(next_x, next_y)
	screen.ShowCursor(tb.curr_pos_x, tb.curr_pos_y)

	// Append the character at the end of the string
	tb.content += string(char)
	display.DrawString(screen, string(char), next_x-1, next_y, styles.SimpleStyle)
	tb.addEnterValue(false)
}

func (tb *TextBox) handleEnterPressed(screen tcell.Screen) {
	tb.enter_idx--

	// If the previous cell is still true then
	if tb.enter_idx == 0 || tb.enters[tb.enter_idx-1] {
		next_pos_x := tb.start_pos_x
		next_pos_y := tb.curr_pos_y - 1

		// Set the new cursor position
		tb.SetNextCursorPosition(next_pos_x, next_pos_y)
		screen.ShowCursor(tb.curr_pos_x, tb.curr_pos_y)
		tb.enters[tb.enter_idx] = false
		tb.curr_line--
		tb.nof_lines--
	}

	// Otherwise, in the previous cell has been appended a value

}

func (tb *TextBox) handleBackspace(screen tcell.Screen) {

	// Check that the index is not negative
	if tb.enter_idx-1 < 0 {
		return
	}

	// If previously, I pressed enter
	if tb.enters[tb.enter_idx-1] {
		tb.handleEnterPressed(screen)
		return
	}

	// If there is no content return
	if len(tb.content) < 1 {
		return
	}

	str_position := tb.getCurrentStringPosition() // Take the position into the string
	content_array := []rune(tb.content)           // From string to rule array
	content_len := len(tb.content)                // Take the length of the string

	// Check how many characters are left in the sx substr
	substr_sx := content_array[0:str_position]
	if len(substr_sx) < 1 {
		return
	}

	// Otherwise, we need to remove the element
	remaining_sx := string(substr_sx[0 : str_position-1])
	remaining_dx := string(substr_sx[str_position:content_len])
	tb.content = remaining_sx + remaining_dx

	// We need to move the cursor in the correct position
	next_pos_x := tb.curr_pos_x - 1
	next_pos_y := tb.curr_pos_y

	if tb.curr_pos_x == tb.start_pos_x {
		// If the relative position is currently 0
		next_pos_x = tb.start_pos_x + tb.getMaxRowSize() - 1
		next_pos_y = next_pos_y - 1
		tb.curr_line--
		tb.nof_lines--
	}

	// Set the new cursor position
	tb.SetNextCursorPosition(next_pos_x, next_pos_y)
	screen.ShowCursor(tb.curr_pos_x, tb.curr_pos_y)

	// Replace the previous character with a void one
	display.DrawString(screen, " ", tb.curr_pos_x, tb.curr_pos_y, styles.SimpleStyle)

	// Redraw the entire string
	tb.drawContent(screen)
}

func (tb *TextBox) handleArrowPressed(screen tcell.Screen, direction_x, direction_y int) {
	next_pos_x := tb.curr_pos_x + direction_x
	next_pos_y := tb.curr_pos_y + direction_y
	next_str_pos := tb.getStringPosition(next_pos_x, next_pos_y)

	// Handle next or previous line
	// next_pos_x-tb.start_pos_x < 0 || next_pos_y-tb.start_pos_y < 0
	if next_pos_x < tb.start_pos_x {
		next_pos_x = tb.start_pos_x + tb.getMaxRowSize()
		next_pos_y--
	} else if next_pos_x >= tb.start_pos_x+tb.getMaxRowSize() {
		next_pos_x = tb.start_pos_x
		next_pos_y++
	}

	// Check overflow for Y positions
	y_diff := tb.rec.Start_y + tb.rec.Height - (tb.start_pos_y - tb.rec.Start_y)
	if next_pos_y < tb.start_pos_y || next_pos_y > y_diff {
		return
	}

	// Check if the content is empty or not, or the new position
	// is bigger than the dimension of the string
	if len(tb.content) < 1 || next_str_pos > len(tb.content) {
		return
	}

	// Set the new cursor position
	tb.SetNextCursorPosition(next_pos_x, next_pos_y)
	screen.ShowCursor(tb.curr_pos_x, tb.curr_pos_y)
}

func (tb *TextBox) getCurrentStringPosition() int {
	return tb.getStringPosition(tb.curr_pos_x, tb.curr_pos_y)
}

// Display the text box
func (tb *TextBox) Display(screen tcell.Screen) {
	tb.rec.DrawRectangle(screen) // Draw the rectangle for the text box

	// Display the Title
	display.DrawString(screen, tb.title, tb.rec.Start_x+3, tb.rec.Start_y, styles.TextBoxTitle)

	// Display the content
	display.DrawString(screen, tb.content, tb.start_pos_x, tb.start_pos_y, styles.SimpleStyle)
}

// Check if the textbox collides with input coordinates
func (tb *TextBox) IsColliding(x, y int) bool {
	y_diff := tb.start_pos_y - tb.rec.Start_y - 1
	x_diff := tb.start_pos_x - tb.rec.Start_x - 1

	return ((x >= tb.rec.Start_x+x_diff && x <= tb.rec.Start_x+tb.rec.Width-x_diff) &&
		(y >= tb.rec.Start_y+y_diff && y <= tb.rec.Start_y+tb.rec.Height-y_diff))
}

func (tb *TextBox) SetFocus(value bool) {
	tb.focus = value
}

func (tb *TextBox) HasFocus() bool {
	return tb.focus
}

// Returns the current cursor position relative to the object
func (tb *TextBox) GetCursorPosition() (int, int) {
	return tb.curr_pos_x, tb.curr_pos_y
}

func (tb *TextBox) SetCursorPosition(x, y int) {
	tb.curr_pos_x = x
	tb.curr_pos_y = y
}

func (tb *TextBox) CheckNextCursorPosition(next_x, next_y int) bool {
	// Initialize the returns values as the current ones
	return tb.IsColliding(next_x+1, next_y+1)
}

func (tb *TextBox) SetNextCursorPosition(x, y int) {
	// Check if the input values lies inside the boundaries
	result := tb.CheckNextCursorPosition(x, y)

	// If the next cursor position is not permitted then return
	if !result {
		return
	}

	// Otherwise updates all the corresponding values
	tb.SetCursorPosition(x, y)
}

func (tb *TextBox) GetContent() string {
	return tb.content
}

func (tb *TextBox) HandleEventKey(screen tcell.Screen, event *tcell.EventKey) {
	if !tb.focus {
		return
	}

	switch event.Key() {
	case tcell.KeyEscape:
		// When Escape is pressed it removes the focus
		// from the current object
		tb.focus = false
		screen.HideCursor()

	case tcell.KeyEnter:
		// When Enter is pressed it moves the cursor to the
		// next line, until the end of the textbox
		// tb.handleEnterKey(screen)

	case tcell.KeyDelete:
		// In the case the CANC key is pressed deletes the
		// characther next to the current one

	case tcell.KeyBackspace, tcell.KeyBackspace2:
		// When backspace is pressed deletes the caracter where
		// the cursor is positioned
		tb.handleBackspace(screen)

	case tcell.KeyUp, tcell.KeyDown, tcell.KeyLeft, tcell.KeyRight:
		// When the arrow is pressed we need to select the correct
		// move and update the cursor position if it is possible
		direction := ARROW_MAPPING[event.Key()]
		tb.handleArrowPressed(screen, direction.X, direction.Y)

	default:
		// Otherwise, check if the key pressed is a letter
		if event.Rune() == 0 {
			return
		}

		// Append the pressed letter to the content
		tb.appendCharacter(screen, event.Rune())
	}
}

func (tb *TextBox) HandleEventMouse(screen tcell.Screen, event *tcell.EventMouse) {
	if !tb.focus {
		tb.focus = true
		screen.ShowCursor(tb.curr_pos_x, tb.curr_pos_y)
	}
}
