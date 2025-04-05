package components

import (
	"github.com/gizak/termui/v3"
	"github.com/okira-e/gotasks/internal/domain"
	cw "github.com/okira-e/gotasks/internal/ui/custom-widgets"
	"github.com/okira-e/gotasks/internal/ui/types"
	"github.com/okira-e/gotasks/internal/utils"
)

type CreateTaskPopup struct {
	Visible 		bool
	// EditingTask if this is set, the widget becomes an edit popup that shows & edits existing data.
	EditingTask		*domain.Task
	
	window			*types.Window
	titleInput   	*cw.TextInput
	descInput    	*cw.TextInput
	focusedField 	*cw.TextInput
	userConfig		*domain.UserConfig
	boardName		string
}

// NewCreateTaskPopupComponent initializes a new popup.
func NewCreateTaskPopupComponent(window *types.Window, config *domain.UserConfig, boardName string) *CreateTaskPopup {
	component := new(CreateTaskPopup)
	
	component.Visible = false
	component.window = window
	component.userConfig = config
	component.boardName = boardName
	component.titleInput = cw.NewTextInput()
	component.descInput = cw.NewTextInput()

	component.focusedField = component.titleInput
	
	return component
}

func (self *CreateTaskPopup) SetEditEditingTask(task *domain.Task) {
	self.EditingTask = task
	
	self.titleInput.SetText(task.Title)
	self.descInput.SetText(task.Description)
}

func (self *CreateTaskPopup) GetAllDrawableWidgets() []termui.Drawable {
	return []termui.Drawable{
		self.titleInput.GetDrawableWidget(),
		self.descInput.GetDrawableWidget(),
	}
}

// HandleKeyboardEvent handles every event for this widget. It returns a flag
// indicating if the next render should clear the view.
func (self *CreateTaskPopup) HandleKeyboardEvent(event termui.Event) bool {
	if event.ID ==  "<C-c>" {
		self.Hide()
		return true
		
	} else if event.ID == "<Tab>" {
		self.ToggleFocusOnNextField()
		
	} else if event.ID == "<Enter>" {
		if self.titleInput.GetText() == "" {
			return false
		}
		
		// Save the task.
		boardOpt := self.userConfig.GetBoard(self.boardName)
		board := boardOpt.Expect("Board was found to be null while handling <Enter> on task creation.")
		
		if len(board.Columns) == 0 {
			utils.SaveLog(
				utils.Error, 
				"Board has no columns. Cannot create a task without a column.",
				nil,
			)
			
			return false
		}
		
		// If we are not in edit mode, create a new task. Otherwise, just simple edit the pointer
		// to the task we're editing.
		
		if self.EditingTask == nil {
			task := domain.NewTask(
				self.titleInput.GetText(), 
				self.descInput.GetText(),
			)
			
			err := self.userConfig.AddTask(self.boardName, task)
			if err != nil {
				utils.SaveLog(utils.Error, err.Error(), map[string]any{"boardName": self.boardName, "task": task})
			}
		} else {
			self.EditingTask.Title = self.titleInput.GetText()
			self.EditingTask.Description = self.descInput.GetText()
			
			self.userConfig.UpdateBoard(board)
		}
		
		self.Hide()
		return true
		
	} else {
		// Disable new lines in the title.
		disableNewLines := false
		if self.focusedField == self.titleInput {
			disableNewLines = true
		}

		self.focusedField.HandleInput(event.ID, disableNewLines)
	}
	
	return false
}

func (self *CreateTaskPopup) Show() {
	self.Visible = true
}

func (self *CreateTaskPopup) Hide() {
	self.Visible = false
	self.reset()
}

// Toggles the focus onto the next input field.
func (self *CreateTaskPopup) ToggleFocusOnNextField() {
	if self.focusedField != nil {
		self.focusedField = utils.Cond(
			self.focusedField == self.titleInput,
			self.descInput,
			self.titleInput,
		)
	} else { // Shouldn't happen
		self.focusedField = self.titleInput
	}
}

// Draw renders the popup if visible.
func (self *CreateTaskPopup) Draw() {
	self.Visible = true
	
	y1 := self.window.Height/4

	self.titleInput.GetDrawableWidget().Title = "Title"
	
	self.titleInput.GetDrawableWidget().SetRect(
		self.window.Width/4,
		self.window.Height/4,

		self.window.Width/4*3,
		y1+3,
	)

	self.descInput.GetDrawableWidget().Title = "Description"
	self.descInput.GetDrawableWidget().SetRect(
		self.window.Width/4,
		self.titleInput.GetDrawableWidget().Max.Y, // 3 is the height of the Title widget above it.
		self.window.Width/4*3,
		self.window.Height/4*3,
	)
	
	// Set the border to be the primary color on the input field that is in focus.
	self.focusedField.GetDrawableWidget().BorderStyle = termui.NewStyle(self.userConfig.PrimaryColor)
	// Remove the primary color from the non-focused
	if self.focusedField == self.titleInput {
		self.descInput.GetDrawableWidget().BorderStyle = termui.NewStyle(termui.ColorClear)
	} else {
		self.titleInput.GetDrawableWidget().BorderStyle = termui.NewStyle(termui.ColorClear)
	}
	
	termui.Render(
		self.GetAllDrawableWidgets()...
	)
}

func (self *CreateTaskPopup) reset() {
	self.focusedField = self.titleInput
	self.EditingTask = nil
	self.titleInput.Flush()
	self.descInput.Flush()
}