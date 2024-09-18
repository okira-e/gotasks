package ui

import (
	"os"

	"github.com/gizak/termui/v3"
)

// handleEvent processes user input and other events.
func (app *App) handleEvent(boardName string, event termui.Event) {
	
	if event.Type == termui.ResizeEvent {
		app.render(boardName)
	} 
	
	if event.Type == termui.KeyboardEvent {
		handleKeymap(event)
	}
}

func handleKeymap(event termui.Event) {
	if event.ID == "<Escape>" {
	}

	if event.ID == "q" || event.ID == "<C-c>" {
		termui.Close()
		os.Exit(0)
	}

	if event.ID == "h" || event.ID == "?" {
	}
}