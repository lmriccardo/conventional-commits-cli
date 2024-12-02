package objects

import "github.com/gdamore/tcell/v2"

type Object interface {
	HandleEvent(tcell.Screen, *tcell.EventKey)
}
