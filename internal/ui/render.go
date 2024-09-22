package ui

import (
	"github.com/gizak/termui/v3"
)

// render handles the full re-render of the app based on the current state.
func (app *App) render() {
	termui.Clear()

	app.width, app.height = termui.TerminalDimensions()

	applyTheme(app.columnsHeadersView, app.theme)
	// if app.columnsHeadersView.NeedsRedraw {
		app.columnsHeadersView.Draw()
	// }
	
	applyTheme(app.tasksView, app.theme)
	// if app.tasksView.NeedsRedraw {
		app.tasksView.Draw()
	// }
	
	if app.createTaskPopup.Visible { // && app.createTaskPopup.NeedsRedraw {
		applyTheme(app.createTaskPopup, app.theme)
		app.createTaskPopup.Draw()
	}
}
