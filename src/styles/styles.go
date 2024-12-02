/*
This file contains all the styles for the content of the ccommits screen manager.
*/
package styles

import "github.com/gdamore/tcell/v2"

var (
	SimpleStyle  = tcell.StyleDefault.Foreground(tcell.ColorWhite)
	TitleStyle   = tcell.StyleDefault.Foreground(tcell.ColorFloralWhite).Bold(true).Underline(true)
	TextBoxTitle = tcell.StyleDefault.Foreground(tcell.ColorCadetBlue).Italic(true).Underline(true)
)
