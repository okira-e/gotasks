package ui

import "github.com/gizak/termui/v3"

// render handles rendering all widgets of the app. Takes in a flag.
// indicating if it should re-render everything by clearing first.
func (app *App) render(shouldClear bool) {
	if shouldClear {
		termui.Clear()
	}
	
	app.columnsHeadersView.Draw()
	
	app.tasksView.Draw()
	
	if app.createTaskPopup.Visible {
		app.createTaskPopup.Draw()

	} else if app.confirmationPopup.Visible {
		app.confirmationPopup.Draw()
		
	} else if app.searchDialogPopup.Visible {
		app.searchDialogPopup.Draw()
		
	}
}
