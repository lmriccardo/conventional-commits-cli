package objects

import (
	"example.com/ccommits/src/display"
	"example.com/ccommits/src/styles"
	"github.com/gdamore/tcell/v2"
)

var DIRECTIONS map[tcell.Key]int = map[tcell.Key]int{
	tcell.KeyUp:   -1,
	tcell.KeyDown: 1,
}

const TRIANGLE_DOWN string = "🔻"

type MultiOptionBox struct {
	rec      Rectangle         // The rectangle containing the box
	content  map[string]string // The vector with all the content
	keys     []string          // The keys of the content map
	curr_idx int               // The current index of the selected content
	title    string            // The title of the box
	focus    bool              // If the current focus is on this object
	view     int               // The current view of the content
}

func MultiOptionBox_new(title string, x, y, size_w, size_h int, content map[string]string) *MultiOptionBox {
	// First, creates the rectangle
	rect := new(Rectangle)
	rect.Start_x = x
	rect.Start_y = y
	rect.Width = size_w
	rect.Height = size_h + 1

	keys := make([]string, 0, len(content))
	for key := range content {
		keys = append(keys, key)
	}

	return &MultiOptionBox{*rect, content, keys, 0, title, false, 0}
}

func (mob *MultiOptionBox) getMaxNofLines() int {
	return mob.rec.Height - 5
}

func (mob *MultiOptionBox) getMaxRowSize() int {
	return mob.rec.Width - 4
}

func (mob *MultiOptionBox) getMaxNofViews() int {
	return len(mob.content) / mob.getMaxNofLines()
}

func (mob *MultiOptionBox) getStringContent(idx int) string {
	return mob.keys[idx] + " - " + mob.content[mob.keys[idx]]
}

func (mob *MultiOptionBox) getCurrentStringContent() string {
	return mob.getStringContent(mob.curr_idx)
}

func (mob *MultiOptionBox) drawContent(screen tcell.Screen, start_line, end_line int) {
	max_lines := mob.getMaxNofLines()
	end_line = min(end_line, start_line+max_lines)
	start_x := mob.rec.Start_x + 2

	// Draw all the specified lines in the screen
	for idx := start_line; idx < end_line; idx++ {
		style := styles.SimpleStyle
		if idx == mob.curr_idx {
			style = styles.SelectStyle
		}

		str_content := mob.getStringContent(idx)
		str_content = str_content[0:min(mob.getMaxRowSize(), len(str_content))]
		start_y := mob.rec.Start_y + 2 + (idx - start_line)
		display.DrawString(screen, str_content, start_x, start_y, style)
	}

	// If the current view is not the last one then we can add at the botton
	// an upside-down triangle indicating that there is more content below
	if mob.view < mob.getMaxNofViews() {
		arrow_down_str := display.CenterString(mob.getMaxRowSize(), TRIANGLE_DOWN)
		start_y := mob.rec.Start_y + mob.rec.Height - 2
		display.DrawString(screen, arrow_down_str, start_x,
			start_y, styles.ArrowDown)
	}
}

func (mob *MultiOptionBox) clearContent(screen tcell.Screen) {
	for idx := 0; idx < mob.getMaxNofLines()+2; idx++ {
		start_y := mob.rec.Start_y + idx + 2

		for row_idx := 0; row_idx < mob.getMaxRowSize(); row_idx++ {
			start_x := mob.rec.Start_x + row_idx + 2
			display.DrawString(screen, " ", start_x, start_y, styles.SimpleStyle)
		}
	}
}

func (mob *MultiOptionBox) getView(curr_idx int) int {
	var next_start int // Initialize a temporary variable
	start_index, view_index := 0, 0

	for start_index < len(mob.content) {
		next_start = min(len(mob.content), start_index+mob.getMaxNofLines())
		if curr_idx >= start_index && curr_idx < next_start {
			return view_index
		}

		view_index++             // Increase the view index
		start_index = next_start // Increase also the next start idx
	}

	return -1
}

func (mob *MultiOptionBox) handleArrowPressed(screen tcell.Screen, direction int) {
	next_view := mob.getView(mob.curr_idx + direction) // Get the next possible view
	if next_view < 0 {
		// If the next view is not possible returns
		return
	}

	prev_curr_idx := mob.curr_idx           // Take the previous index
	mob.curr_idx = mob.curr_idx + direction // Update the current index

	// If the next view is different from the previous one
	// then we need to update the content displayed
	if next_view != mob.view {
		mob.clearContent(screen) // Clear the previous content

		// Then we need to display the new one
		start_line := next_view * mob.getMaxNofLines()
		stop_line := min((next_view+1)*mob.getMaxNofLines(), len(mob.content))
		mob.view = next_view // Update the current view with the new one
		mob.drawContent(screen, start_line, stop_line)
		return
	}

	// On the other hand, the view remains the same no new content should be
	// displayed, we just need to change which line is highlighted and which not
	prev_pos_x := mob.rec.Start_x + 2
	prev_pos_y := mob.rec.Start_y + 2 + prev_curr_idx - mob.view*mob.getMaxNofLines()
	prev_content := mob.getStringContent(prev_curr_idx)
	end_content := min(mob.getMaxRowSize(), len(prev_content))
	prev_content = prev_content[0:end_content]
	display.DrawString(screen, prev_content, prev_pos_x, prev_pos_y, styles.SimpleStyle)

	pos_x := mob.rec.Start_x + 2
	pos_y := mob.rec.Start_y + 2 + mob.curr_idx - mob.view*mob.getMaxNofLines()
	content := mob.getCurrentStringContent()
	end_content = min(mob.getMaxRowSize(), len(content))
	content = content[0:end_content]
	display.DrawString(screen, content, pos_x, pos_y, styles.SelectStyle)
}

func (mob *MultiOptionBox) Display(screen tcell.Screen) {
	mob.rec.DrawRectangle(screen) // Draw the rectangle for the text box

	// Display the Title
	display.DrawString(screen, mob.title, mob.rec.Start_x+3, mob.rec.Start_y, styles.TextBoxTitle)

	// Draw the contents
	mob.clearContent(screen)
	mob.drawContent(screen, 0, len(mob.content))
}

func (mob *MultiOptionBox) IsColliding(x, y int) bool {
	return ((x >= mob.rec.Start_x && x <= mob.rec.Start_x+mob.rec.Width) &&
		(y >= mob.rec.Start_y && y <= mob.rec.Start_y+mob.rec.Height))
}

func (mob *MultiOptionBox) SetFocus(value bool) {
	mob.focus = value
}

func (mob *MultiOptionBox) HasFocus() bool {
	return mob.focus
}

// Returns the current cursor position relative to the object
func (mob *MultiOptionBox) GetCursorPosition() (int, int) {
	cursor_pos_x := mob.rec.Start_x + 2
	cursor_pos_y := mob.rec.Start_y + 2 + (mob.curr_idx - mob.view*mob.getMaxNofLines())
	return cursor_pos_x, cursor_pos_y
}

func (mob *MultiOptionBox) GetContent() string {
	return mob.keys[mob.curr_idx]
}

func (mob *MultiOptionBox) HandleEventKey(screen tcell.Screen, event *tcell.EventKey) {
	if !mob.focus {
		return
	}

	switch event.Key() {
	case tcell.KeyEscape:
		// When Escape is pressed it removes the focus
		// from the current object
		mob.focus = false
		screen.HideCursor()

	case tcell.KeyUp, tcell.KeyDown:
		// Handle pressing arrow up or arrow down. In both cases
		// the current selected lines must change and so the current
		// view in the screen in case.
		mob.handleArrowPressed(screen, DIRECTIONS[event.Key()])
		screen.ShowCursor(mob.GetCursorPosition())

	default:
		return
	}
}

func (mob *MultiOptionBox) HandleEventMouse(screen tcell.Screen, event *tcell.EventMouse) {
	if !mob.focus {
		mob.focus = true
		screen.ShowCursor(mob.GetCursorPosition())
	}
}
