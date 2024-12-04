package src

import (
	"fmt"
	"strings"

	"example.com/ccommits/src/display"
	"example.com/ccommits/src/objects"
	"example.com/ccommits/src/styles"
	"example.com/ccommits/src/util"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

var DIRECTIONS map[tcell.Key]int = map[tcell.Key]int{
	tcell.KeyLeft:  -1,
	tcell.KeyRight: 1,
}

type CCommitWindow struct {
	screen   tcell.Screen            // The main screen of tcell
	tb_desc1 *objects.TextBox        // The textbox for the main description
	tb_desc2 *objects.TextBox        // The textbox for the longer description
	mb_slct1 *objects.MultiOptionBox // The box for selecting the commit type
	mb_slct2 *objects.MultiOptionBox // The box for selecting the gitmoji
	gitinfo  *util.GitInfo           // Git Information of the current repo

	size_w int // The total size in width of the screen
	size_h int // The total size in height of the screen

	cursor_x int // The x position of the cursor
	cursor_y int // The y position of the cursor

	prev_focus_obj objects.Object // Previously focused object
	prev_focus_idx int            // Previous focused object index
}

func CCommitWindow_new(gitinfo *util.GitInfo) *CCommitWindow {
	// First of all, let's create the screen of the main app
	win := new(CCommitWindow)
	win.screen = display.InitializeScreen()
	win.size_w, win.size_h = win.screen.Size()

	// Creates the textbox for the main description
	tbd1_x := win.size_w/2 + 32
	tbd1_y := 9
	tbd1_size := win.size_w - 3 - tbd1_x
	win.tb_desc1 = objects.TextBox_new(MAIN_DESC, tbd1_x, tbd1_y, tbd1_size, 5)

	// Creates the textbox for the longer description
	tbd2_x := tbd1_x
	tbd2_y := tbd1_y + 7
	tbd2_size_w := win.size_w - 3 - tbd1_x
	tbd2_size_h := win.size_h - 3 - tbd2_y
	win.tb_desc2 = objects.TextBox_new(LONG_DESC, tbd2_x, tbd2_y, tbd2_size_w, tbd2_size_h)

	// Creates the first Multi option selection box
	mob1_x, mob1_y := 5, 9
	mob1_size_w := win.size_w/4 - mob1_x + 8
	mob1_size_h := win.size_h - 3 - mob1_y
	win.mb_slct1 = objects.MultiOptionBox_new(TYPE, mob1_x, mob1_y,
		mob1_size_w, mob1_size_h, CHANGE_TYPE)

	// Creates the second Multi option selection box
	mob2_x, mob2_y := mob1_x+mob1_size_w+3, 9
	mob2_size_w := tbd1_x - 3 - mob2_x
	mob2_size_h := win.size_h - 3 - mob1_y
	win.mb_slct2 = objects.MultiOptionBox_new(GITMOJI, mob2_x, mob2_y,
		mob2_size_w, mob2_size_h, GITMOJI_ARRAY)

	// Sets the cursor position
	win.cursor_x = 0
	win.cursor_y = 0
	win.prev_focus_obj = nil
	win.prev_focus_idx = -1
	win.gitinfo = gitinfo

	return win
}

func (win *CCommitWindow) handleArrowPressed(key tcell.Key) {
	direction := DIRECTIONS[key]
	next_focus_idx := win.prev_focus_idx + direction
	if next_focus_idx < 0 {
		return
	}

	objs := []objects.Object{win.mb_slct1, win.mb_slct2, win.tb_desc1, win.tb_desc2}
	next_focus_obj := objs[next_focus_idx]
	next_focus_obj.HandleEventMouse(win.screen, nil)
	win.cursor_x, win.cursor_y = next_focus_obj.GetCursorPosition()
	win.screen.Sync()
	win.prev_focus_obj = next_focus_obj
	win.prev_focus_idx = next_focus_idx
}

func (win *CCommitWindow) displayTitle(title string) {
	width, _ := win.screen.Size() // Get the width and the height of the screen
	titlelen := len(title) / 2    // Get the length of the title
	start_x := width/2 - titlelen // The starting position of the title

	// Draw the title
	display.DrawString(win.screen, title, start_x, TITLE_Y, styles.TitleStyle)
}

func (win *CCommitWindow) displaySubTitle(sub_title string) {
	width, _ := win.screen.Size()  // Get the width and the height of the screen
	titlelen := len(sub_title) / 2 // Get the length of the title
	start_x := width/2 - titlelen  // The starting position of the title
	var start_y int = TITLE_Y + 2  // The starting y position of the title

	// Draw the title
	display.DrawString(win.screen, sub_title, start_x, start_y, styles.SubTitleStyle)
}

func (win *CCommitWindow) displayGitInfo() {
	width, _ := win.screen.Size() // Get the width and the height of the screen
	var start_y int = TITLE_Y + 4 // The starting y position of the title

	// Prepare the strings that needs to be displayed
	repo_name := strings.Join([]string{REPO, win.gitinfo.Reponame}, " ")
	branch_name := strings.Join([]string{BRANCH, win.gitinfo.Curr_branch}, " ")
	remote_name := strings.Join([]string{REMOTE, win.gitinfo.Curr_remote}, " ")

	// Compute the total length and respective starting x positions
	repo_name_len := runewidth.StringWidth(repo_name)
	branch_name_len := runewidth.StringWidth(branch_name)
	remote_name_len := runewidth.StringWidth(remote_name)
	tot_len := repo_name_len + branch_name_len + 4 + remote_name_len

	// Display the repository name
	repo_start_x := width/2 - tot_len/2
	display.DrawString(win.screen, repo_name, repo_start_x, start_y, styles.GitInfoStyle)

	// Display the branch name
	branch_start_x := repo_start_x + repo_name_len + 4
	display.DrawString(win.screen, branch_name, branch_start_x, start_y, styles.GitInfoStyle)

	// Display the branch name
	remote_start_x := branch_start_x + branch_name_len + 4
	display.DrawString(win.screen, remote_name, remote_start_x, start_y, styles.GitInfoStyle)
}

func (win *CCommitWindow) Display() {
	// Display the title
	win.displayTitle(TITLE)
	win.displaySubTitle(VERSION)
	win.displayGitInfo()

	// Draw the text boxes
	win.tb_desc1.Display(win.screen)
	win.tb_desc2.Display(win.screen)
	win.mb_slct1.Display(win.screen)
	win.mb_slct2.Display(win.screen)

	// Show the screen
	win.screen.Show()
}

func (win *CCommitWindow) getColliding(x, y int, focus bool) (int, objects.Object) {
	objs := []objects.Object{win.mb_slct1, win.mb_slct2, win.tb_desc1, win.tb_desc2}

	for idx, element := range objs {
		if element.IsColliding(x, y) && (!focus || (focus && element.HasFocus())) {
			return idx, element
		}
	}

	return -1, nil
}

func (win *CCommitWindow) Run() string {
	defer win.screen.Fini()

	win.Display()

	// Wait for a key event
	for {
		event := win.screen.PollEvent()
		switch ev := event.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyCtrlC {
				// Fetch all the results
				short_desc := win.tb_desc1.GetContent()
				long_desc := win.tb_desc2.GetContent()
				commit_type := strings.ToLower(win.mb_slct1.GetContent())
				commit_emoji := win.mb_slct2.GetContent()

				// Check that both descriptions have at least one char
				if !(len(short_desc) > 1 && len(long_desc) > 1) {
					return ""
				}

				// Otherwise, we can returns the formatted commit
				return fmt.Sprintf("%s: %s %s\n\n%s", commit_type, commit_emoji,
					short_desc, long_desc)
			}

			_, obj := win.getColliding(win.cursor_x, win.cursor_y, true)
			if obj == nil { // Check that the returned object is not null
				if ev.Key() == tcell.KeyLeft || ev.Key() == tcell.KeyRight {
					win.handleArrowPressed(ev.Key())
				}
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

			idx, obj := win.getColliding(mouse_x, mouse_y, false)
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
			win.prev_focus_idx = idx

		case *tcell.EventResize:
			win.screen.Sync()
		}
	}
}
