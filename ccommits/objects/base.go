package objects

import "github.com/gdamore/tcell/v2"

type Object interface {
	Display(tcell.Screen)
	IsColliding(int, int) bool
	SetFocus(bool)
	HasFocus() bool
	GetCursorPosition() (int, int)
	GetContent() string
	HandleEventKey(tcell.Screen, *tcell.EventKey)
	HandleEventMouse(tcell.Screen, *tcell.EventMouse)
}
