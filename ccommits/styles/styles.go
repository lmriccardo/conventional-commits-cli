/*
This file contains all the styles for the content of the ccommits screen manager.
*/
package styles

import "github.com/gdamore/tcell/v2"

var (
	SimpleStyle   = tcell.StyleDefault.Foreground(tcell.ColorWhite)
	SelectStyle   = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorGray).Underline(true)
	BorderStyle   = tcell.StyleDefault.Foreground(tcell.ColorDarkSlateGray)
	TextBoxTitle  = tcell.StyleDefault.Foreground(tcell.ColorCadetBlue).Italic(true).Underline(true)
	ArrowDown     = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorDarkSlateGray)
	TitleStyle    = tcell.StyleDefault.Foreground(tcell.ColorDarkOrange).Bold(true).Underline(true)
	SubTitleStyle = tcell.StyleDefault.Foreground(tcell.ColorDarkSlateBlue).Bold(true).Italic(true)
	GitInfoStyle  = tcell.StyleDefault.Foreground(tcell.ColorMediumVioletRed).Underline(true)
)
