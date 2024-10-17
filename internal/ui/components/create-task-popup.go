package components

import (
	"github.com/gizak/termui/v3"
	"github.com/okira-e/gotasks/internal/domain"
	"github.com/okira-e/gotasks/internal/ui/types"
	cw "github.com/okira-e/gotasks/internal/ui/custom-widgets"
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
	component.titleInput = cw.NewTextInput()
	component.descInput = cw.NewTextInput()
	component.userConfig = config
	component.boardName = boardName

	component.focusedField = component.titleInput
	component.titleInput.GetDrawableWidget().Title = "Title"
	component.descInput.GetDrawableWidget().Title = "Description"
	
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

func (self *CreateTaskPopup) HandleKeyboardEvent(event termui.Event) {
	if event.ID ==  "<Escape>" {
		self.Hide()
	} else if event.ID == "<Tab>" {
		self.ToggleFocusOnNextField()
	} else if event.ID == "<Enter>" {
		if self.titleInput.GetText() == "" {
			return
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
			return
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
		// Re-rendering already happens after this function call.
	} else {
		// Disable new lines in the title.
		disableNewLines := false
		if self.focusedField == self.titleInput {
			disableNewLines = true
		}

		self.focusedField.HandleInput(event.ID, disableNewLines)
	}
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

	self.titleInput.GetDrawableWidget().SetRect(
		self.window.Width/4,
		self.window.Height/4,

		self.window.Width/4*3,
		y1+3,
	)

	self.descInput.GetDrawableWidget().SetRect(
		self.window.Width/4,
		self.titleInput.GetDrawableWidget().Max.Y, // 3 is the height of the Title widget above it.
		self.window.Width/4*3,
		self.window.Height/4*3,
	)
	// Set the border to be blue on the input field that is in focus.
	if self.focusedField != nil {
		self.focusedField.GetDrawableWidget().BorderStyle = termui.NewStyle(self.userConfig.PrimaryColor)
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