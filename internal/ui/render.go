package ui

import (
	"github.com/gizak/termui/v3"
)

// render handles the full re-render of the app based on the current state.
func (app *App) render() {
	termui.Clear()

	app.width, app.height = termui.TerminalDimensions()

	applyTheme(app.columnsHeadersView, app.theme)
	app.columnsHeadersView.Draw()
	
	applyTheme(app.tasksView, app.theme)
	app.tasksView.Draw()
	
	if app.createTaskPopup.Visible {
		app.createTaskPopup.Draw()

	} else if app.confirmationPopup.Visible {
		app.confirmationPopup.Draw()
		
	} else if app.searchDialogPopup.Visible {
		app.searchDialogPopup.Draw()
		
	}
}
