package objects

import (
	"github.com/gdamore/tcell/v2"
	"github.com/lmriccardo/conventional-commits-cli/ccommits/display"
	"github.com/lmriccardo/conventional-commits-cli/ccommits/styles"
)

// Mapping event keys arrow to direction (x, y)
var ARROW_MAPPING map[tcell.Key]Vec2 = map[tcell.Key]Vec2{
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
}

func TextBox_new(title string, x, y, size_w, size_h int) *TextBox {
	// First, creates the rectangle
	rect := Rectangle{size_w, size_h + 1, x, y}

	// Compute the starting position of the text
	startpos_y := rect.Start_y + 2
	startpos_x := rect.Start_x + 2

	// Creates the TextBox and returns it
	tb := new(TextBox)
	tb.rec = rect
	tb.title = title
	tb.start_pos_x = startpos_x
	tb.start_pos_y = startpos_y
	tb.content = ""
	tb.curr_pos_x = startpos_x
	tb.curr_pos_y = startpos_y
	tb.focus = false
	tb.nof_lines = 0
	tb.curr_line = 0

	return tb
}

func (tb *TextBox) displayContentPortion(screen tcell.Screen, start_x, start_y int, content string) {
	rel_pos_x := start_x - tb.start_pos_x
	var nof_rows int = (rel_pos_x+len(content))/tb.getMaxRowSize() + 1

	content_arr := []rune(content)
	curr_nof_rows := 0
	in_content_start_pos := 0

	for curr_nof_rows < nof_rows {
		// Compute the start and stop indices for taking the substring
		in_content_stop_pos := min(len(content), in_content_start_pos+(tb.getMaxRowSize()-rel_pos_x))
		subcontent := content_arr[in_content_start_pos:in_content_stop_pos]

		// Display the substring
		display.DrawString(screen, string(subcontent), start_x, start_y, styles.SimpleStyle)

		// Update positions
		curr_nof_rows++
		start_x = tb.start_pos_x
		start_y = start_y + curr_nof_rows
		rel_pos_x = start_x - tb.start_pos_x
		in_content_start_pos = in_content_stop_pos
	}
}

func (tb *TextBox) displayContent(screen tcell.Screen) {
	// Display all the content into the screen
	row_size := tb.getMaxRowSize()      // Get the maximum row size
	content_array := []rune(tb.content) // From string to rule array
	content_len := len(tb.content)      // Take the length of the string

	for row_idx := 0; row_idx < tb.nof_lines; row_idx++ {
		// Get the correct start and stop indexes
		start_idx := row_idx * row_size
		stop_idx := min(content_len, (row_idx+1)*row_size)
		substr := string(content_array[start_idx:stop_idx])

		// Draw the substring content into the screen
		tb.displayContentPortion(screen, tb.start_pos_x, tb.start_pos_y+row_idx, substr)
	}
}

func (tb *TextBox) getMaxRowSize() int {
	return tb.rec.Width - 2*(tb.start_pos_x-tb.rec.Start_x)
}

func (tb *TextBox) getMaxRows() int {
	return tb.rec.Height - 2*(tb.start_pos_y-tb.rec.Start_y)
}

func (tb *TextBox) getPositionInString(x, y int) int {
	rel_pos_x := x - tb.start_pos_x
	rel_pos_y := y - tb.start_pos_y
	return rel_pos_x + rel_pos_y*tb.getMaxRowSize()
}

func (tb *TextBox) getCurrentPositionInString() int {
	return tb.getPositionInString(tb.curr_pos_x, tb.curr_pos_y)
}

func (tb *TextBox) getLastAvailableStringPosition() (int, int) {
	content_len := len(tb.content) - 1
	rel_pos_x := content_len % tb.getMaxRowSize()
	rel_pos_y := content_len / tb.getMaxRowSize()

	return rel_pos_x, rel_pos_y
}

func (tb *TextBox) addCharacter(screen tcell.Screen, char rune) {
	max_row_size := tb.getMaxRowSize() // Get the maximum size of a single row
	max_nof_line := tb.getMaxRows()    // Get the maximum number of rows for the textbox

	// Check if new characters are allowed
	if len(tb.content)+1 >= max_nof_line*max_row_size {
		return
	}

	prev_pos_x := tb.curr_pos_x
	prev_pos_y := tb.curr_pos_y
	curr_pos := tb.getCurrentPositionInString()

	// Check if the new character can be placed in the same row
	// or a new row is required to be displayed
	if len(tb.content[:curr_pos])+1 >= (tb.curr_line+1)*max_row_size {
		tb.nof_lines++
		tb.curr_line++
		tb.curr_pos_x = tb.start_pos_x
		tb.curr_pos_y = tb.start_pos_y + tb.curr_line
	} else {
		tb.curr_pos_x++
	}

	// We need to detect if the add character is an insert
	// operation or append operation.
	if curr_pos < len(tb.content) {
		tb.content = tb.content[:curr_pos] + string(char) + tb.content[curr_pos:]

		// Now, only the portion of the content that has been shifted
		// needs to be re-displayed as a whole
		tb.displayContentPortion(screen, prev_pos_x, prev_pos_y, tb.content[curr_pos:])
	} else {
		tb.content += string(char)
		display.DrawString(screen, string(char), prev_pos_x, prev_pos_y, styles.SimpleStyle)
	}

	screen.ShowCursor(tb.curr_pos_x, tb.curr_pos_y)
}

func (tb *TextBox) handleArrowPressed(screen tcell.Screen, direction Vec2) {
	abs_pos_x := tb.curr_pos_x - tb.start_pos_x + direction.X
	new_pos_x := (abs_pos_x) % tb.getMaxRowSize()
	new_pos_y := tb.curr_pos_y - tb.start_pos_y + direction.Y

	bool2int := map[bool]int{false: 0, true: 1}

	// Check if going left or right we need to switch up or down row
	if new_pos_x < 0 || abs_pos_x >= tb.getMaxRowSize() {
		cond1_i := bool2int[new_pos_x < 0]
		cond2_i := bool2int[abs_pos_x >= tb.getMaxRowSize()]
		new_pos_y += -1*cond1_i + 1*cond2_i
		new_pos_x = (tb.getMaxRowSize()-1)*cond1_i + 0*(1-cond1_i)
	}

	// Check out of bound for new y position in the textbox
	if new_pos_y > tb.getMaxRows() || new_pos_y < 0 {
		return
	}

	// If the new position into the string is grater than the length of the
	// string, we need to replace the cursor at the last avaiable position
	new_str_pos := tb.getPositionInString(new_pos_x+tb.start_pos_x, new_pos_y+tb.start_pos_y)
	if new_str_pos > len(tb.content) {
		new_pos_x, new_pos_y = tb.getLastAvailableStringPosition()
		new_pos_x = (new_pos_x + 1) % tb.getMaxRowSize()

		if new_pos_x == 0 {
			new_pos_y++
		}
	}

	tb.curr_pos_x = new_pos_x + tb.start_pos_x
	tb.curr_pos_y = new_pos_y + tb.start_pos_y
	tb.curr_line = tb.curr_pos_y - tb.start_pos_y

	screen.ShowCursor(tb.curr_pos_x, tb.curr_pos_y)
}

func (tb *TextBox) handleBackspace(screen tcell.Screen) {
	// If there is no content return
	if len(tb.content) < 1 {
		return
	}

	str_position := tb.getCurrentPositionInString() // Take the position into the string
	content_array := []rune(tb.content)             // From string to rule array
	content_len := len(tb.content)                  // Take the length of the string

	// Check how many characters are left in the sx substr
	substr_sx := content_array[0:str_position]
	if len(substr_sx) < 1 {
		return
	}

	// Otherwise, we need to remove the element
	str_pos_x, str_pos_y := tb.getLastAvailableStringPosition()
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

	tb.curr_pos_x = next_pos_x
	tb.curr_pos_y = next_pos_y
	screen.ShowCursor(tb.curr_pos_x, tb.curr_pos_y)

	// Replace the previous character with a void one
	display.DrawString(screen, " ", tb.curr_pos_x, tb.curr_pos_y, styles.SimpleStyle)
	tb.displayContentPortion(screen, tb.curr_pos_x, tb.curr_pos_y, remaining_dx)

	// Reset the pixed that previously was holding the last character of the string
	abs_str_pos_x := str_pos_x + tb.start_pos_x
	abs_str_pos_y := str_pos_y + tb.start_pos_y
	display.DrawString(screen, " ", abs_str_pos_x, abs_str_pos_y, styles.SimpleStyle)
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

func (tb *TextBox) GetContent() string {
	return tb.content
}

func (tb *TextBox) Display(screen tcell.Screen) {
	tb.rec.DrawRectangle(screen) // Draw the rectangle for the text box
	display.DrawString(screen, tb.title, tb.rec.Start_x+3, tb.rec.Start_y, styles.TextBoxTitle)
	tb.displayContent(screen) // Display the string content
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

	case tcell.KeyEnter: // Not supported yet
		return

	case tcell.KeyBackspace, tcell.KeyBackspace2:
		// When backspace is pressed deletes the character where
		// the cursor is positioned
		tb.handleBackspace(screen)

	case tcell.KeyUp, tcell.KeyDown, tcell.KeyLeft, tcell.KeyRight:
		// When the arrow is pressed we need to select the correct
		// move and update the cursor position if it is possible
		direction := ARROW_MAPPING[event.Key()]
		tb.handleArrowPressed(screen, direction)

	default:
		// Otherwise, check if the key pressed is a letter
		if event.Rune() == 0 {
			return
		}

		// Append the pressed letter to the content
		tb.addCharacter(screen, event.Rune())
	}
}

func (tb *TextBox) HandleEventMouse(screen tcell.Screen, event *tcell.EventMouse) {
	if !tb.focus {
		tb.focus = true
		screen.ShowCursor(tb.curr_pos_x, tb.curr_pos_y)
	}
}
