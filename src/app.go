package src

import (
	"example.com/ccommits/src/display"
	"example.com/ccommits/src/objects"
	"github.com/gdamore/tcell/v2"
)

type CCommitWindow struct {
	screen   tcell.Screen     // The main screen of tcell
	tb_desc1 *objects.TextBox // The textbox for the main description
	tb_desc2 *objects.TextBox // The textbox for the longer description

	size_w int // The total size in width of the screen
	size_h int // The total size in height of the screen

	cursor_x int // The x position of the cursor
	cursor_y int // The y position of the cursor

	prev_focus_obj objects.Object // Previously focused object
}

func CCommitWindow_new() *CCommitWindow {
	// First of all, let's create the screen of the main app
	win := new(CCommitWindow)
	win.screen = display.InitializeScreen()
	win.size_w, win.size_h = win.screen.Size()

	// Creates the textbox for the main description
	tbd1_x := win.size_w/2 + 2
	tbd1_y := 9
	tbd1_size := win.size_w - 3 - tbd1_x
	win.tb_desc1 = objects.TextBox_new(MAIN_DESC, tbd1_x, tbd1_y, tbd1_size, 5)

	// Creates the textbox for the longer description
	tbd2_x := tbd1_x
	tbd2_y := tbd1_y + 7
	tbd2_size_w := win.size_w - 3 - tbd1_x
	tbd2_size_h := win.size_h - 3 - tbd2_y
	win.tb_desc2 = objects.TextBox_new(LONG_DESC, tbd2_x, tbd2_y, tbd2_size_w, tbd2_size_h)

	// Sets the cursor position
	win.cursor_x = 0
	win.cursor_y = 0
	win.prev_focus_obj = nil

	return win
}

func (win *CCommitWindow) Display() {
	// Display the title
	display.DisplayTitle(win.screen, TITLE)

	// Draw the text boxes
	win.tb_desc1.Display(win.screen)
	win.tb_desc2.Display(win.screen)

	// Show the screen
	win.screen.Show()
}

func (win *CCommitWindow) getColliding(x, y int) objects.Object {
	// Check if the colliding is for the main textbox
	if win.tb_desc1.IsColliding(x, y) {
		return win.tb_desc1
	}

	// Check for the longer description textbox
	if win.tb_desc2.IsColliding(x, y) {
		return win.tb_desc2
	}

	return nil
}

func (win *CCommitWindow) Run() {
	defer win.screen.Fini()

	win.Display()

	// Wait for a key event
	for {
		event := win.screen.PollEvent()
		switch ev := event.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyCtrlC {
				return
			}

			obj := win.getColliding(win.cursor_x, win.cursor_y)
			if obj == nil { // Check that the returned object is not null
				continue
			}

			obj.HandleEventKey(win.screen, ev) // Handle the event
			win.cursor_x, win.cursor_y = obj.GetCursorPosition()
			win.screen.Sync()

		case *tcell.EventMouse:
			mouse_x, mouse_y := ev.Position()

			// Check for actual mouse pressed
			if !(ev.Buttons() == tcell.Button1 || ev.Buttons() == tcell.Button2) {
				continue
			}

			obj := win.getColliding(mouse_x, mouse_y)
			if obj == nil { // Check that the returned object is not null
				continue
			}

			if win.prev_focus_obj != nil {
				win.prev_focus_obj.SetFocus(false) // Unset the focus of previous object
			}

			obj.HandleEventMouse(win.screen, ev) // Handle the event
			win.cursor_x, win.cursor_y = obj.GetCursorPosition()
			win.screen.Sync()
			win.prev_focus_obj = obj

		case *tcell.EventResize:
			win.screen.Sync()
		}
	}
}
