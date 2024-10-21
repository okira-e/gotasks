package ui

import (
	"os"

	"github.com/gizak/termui/v3"
)

// handleEvent processes user input and other events.
func (app *App) handleEvent(event termui.Event) {
	shouldClear := false
	
	if event.Type == termui.ResizeEvent {
		app.window.Width, app.window.Height = termui.TerminalDimensions()
		shouldClear = true
		
	} else if event.Type == termui.KeyboardEvent {
		shouldClear = app.handleKeymap(event)
	}
	
	app.render(shouldClear)
}

// handleKeymap handles every keystroke given. One of the things it handles
// is, if a text input widget is in focus, it sends the characters to it instead
// of handling the global keymap as an example. So 'q' could conditionaly write "q"
// on a widget or it could exit the app.
// It returns a flag indicating if we should clear before the next render.
func (app *App) handleKeymap(event termui.Event) bool {
	shouldClear := false

	if app.createTaskPopup.Visible {
		shouldClear = app.createTaskPopup.HandleKeyboardEvent(event)
		
	} else if app.confirmationPopup.Visible {
		shouldClear = app.confirmationPopup.HandleInput(event)
		
	} else if app.searchDialogPopup.Visible {
		shouldClear = app.searchDialogPopup.HandleInput(event)
		
	} else { // Default view is the tasks-view (the board itself)
		switch event.ID {
		case "?":
			// Handle help key logic here
			
		case "/", "s":
			if !app.searchDialogPopup.Visible {
				app.searchDialogPopup.Show()
			}
		
		case "q", "<C-c>":
			termui.Close()
			os.Exit(0)
			
		case "c":
			if !app.createTaskPopup.Visible {
				app.createTaskPopup.Show()
			}
			
		case "e":
			if app.tasksView.TaskInFocus != nil {
				app.createTaskPopup.SetEditEditingTask(app.tasksView.TaskInFocus)
				app.createTaskPopup.Show()
			}
			
		case "d":
			if !app.confirmationPopup.Visible {
				action := func(choice bool) {
					if choice == false {
						return
					}
					
					app.userConfig.DeleteTask(app.boardName, app.tasksView.TaskInFocus)
					app.tasksView.SetDefaultFocusedWidget()
				}
				
				app.confirmationPopup.SetMessageAndAction("Are you sure you want to delete this task?", action)
				app.confirmationPopup.Show()
			}
			
		default: // Handles the movements/action in the board view itself
			shouldClear = app.tasksView.HandleKeymap(event.ID)
			
		}
	}

	return shouldClear
}

