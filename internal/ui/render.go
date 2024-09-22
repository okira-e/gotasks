package ui

import (
	"github.com/gizak/termui/v3"
)

// render handles the full re-render of the app based on the current state.
func (app *App) render() {
	termui.Clear()

	// Update the board in the state with the latest one.
	
	// boardOpt := app.userConfig.GetBoard(app.boardName)
	// board := boardOpt.Expect("Screwed?")
	
	// utils.SaveLog(
	// 	utils.Debug,
	// 	"Board: ",
	// 	map[string]any{"BOARD: ": board.Tasks["Todo"]},
	// )

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
