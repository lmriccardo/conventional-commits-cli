package objects

import "github.com/gdamore/tcell/v2"

type Object interface {
	Display(tcell.Screen)
	IsColliding(int, int) bool
	SetFocus(bool)
	GetCursorPosition() (int, int)
	SetCursorPosition(int, int)
	CheckNextCursorPosition(int, int) bool
	SetNextCursorPosition(int, int)
	HandleEventKey(tcell.Screen, *tcell.EventKey)
	HandleEventMouse(tcell.Screen, *tcell.EventMouse)
}
