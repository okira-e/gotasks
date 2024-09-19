package ui

import (
	"os"

	"github.com/gizak/termui/v3"
)

// handleEvent processes user input and other events.
func (app *App) handleEvent(boardName string, event termui.Event) {

	if event.Type == termui.ResizeEvent {
		app.render()
	}

	if event.Type == termui.KeyboardEvent {
		app.handleKeymap(event)
	}
}

func (app *App) handleKeymap(event termui.Event) {
	switch event.ID {
	case "<Escape>":
		{
			if app.shouldRenderCreateTaskPopup {
				app.shouldRenderCreateTaskPopup = false
			}

			app.render()
		}
	case "h", "?":
		{
			// Handle help key logic here
		}
	case "q", "<C-c>":
		{
			termui.Close()
			os.Exit(0)
		}
	case "c":
		{
			if !app.shouldRenderCreateTaskPopup {
				app.shouldRenderCreateTaskPopup = true
				app.render()
			}
		}

	// Optionally, you can handle the default case if no matching keys are found
	default:
		// Handle other keys if necessary
	}
}
