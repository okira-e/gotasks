package ui

import (
	"os"

	"github.com/gizak/termui/v3"
)

// handleEvent processes user input and other events.
func (app *App) handleEvent(event termui.Event) {

	if event.Type == termui.ResizeEvent {
		app.render()
	}

	if event.Type == termui.KeyboardEvent {
		app.width, app.height = termui.TerminalDimensions()
			
		app.handleKeymap(event)
	}
}

// handleKeymap handles every keystroke given. One of the things it handles
// is, if a text input widget is in focus, it sends the characters to it instead
// of handling the global keymap as an example. So 'q' could conditionaly write "q"
// on a widget or it could exit the app.
func (app *App) handleKeymap(event termui.Event) {
	if app.createTaskPopup.Visible {
		app.createTaskPopup.HandleKeyboardEvent(event)
		app.render()
		
		return
	}
	
	if app.confirmationPopup.Visible {
		app.confirmationPopup.HandleInput(event)
		app.render()
		
		return
	}
	
	switch event.ID {
	case "?":
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
			if !app.createTaskPopup.Visible {
				app.createTaskPopup.Show()
			}
		}
	case "d":
		{
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
		}
	default: // Handles the movements/action in the board view itself
		{
			app.tasksView.HandleKeymap(event.ID)
		}
	}
	
	app.render()
}
