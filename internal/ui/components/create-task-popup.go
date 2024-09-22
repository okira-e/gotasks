package components

import (
	"github.com/gizak/termui/v3"
	"github.com/okira-e/gotasks/internal/domain"
	cw "github.com/okira-e/gotasks/internal/ui/custom-widgets"
	"github.com/okira-e/gotasks/internal/utils"
)

type CreateTaskPopup struct {
	Visible 		bool
	NeedsRedraw  	bool
	
	titleInput   	*cw.TextInput
	descInput    	*cw.TextInput
	focusedField 	*cw.TextInput
	userConfig		*domain.UserConfig
	boardName		string
}

// NewCreateTaskPopup initializes a new popup.
func NewCreateTaskPopup(fullWidth int, fullHeight int, config *domain.UserConfig, boardName string) *CreateTaskPopup {
	ret := &CreateTaskPopup{
		Visible:    false,
		titleInput: cw.NewTextInput(),
		descInput:  cw.NewTextInput(),
		userConfig: config,
		boardName:  boardName,
	}

	y1 := fullHeight/4

	ret.titleInput.GetDrawableWidget().Title = "Title"
	ret.titleInput.GetDrawableWidget().SetRect(
		fullWidth/4,
		fullHeight/4,

		fullWidth/4*3,
		y1+3,
	)

	ret.focusedField = ret.titleInput

	ret.descInput.GetDrawableWidget().Title = "Description"
	ret.descInput.GetDrawableWidget().SetRect(
		fullWidth/4,
		ret.titleInput.GetDrawableWidget().Max.Y, // 3 is the height of the Title widget above it.
		fullWidth/4*3,
		fullHeight/4*3,
	)

	return ret
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
		// Save the task.
		task := domain.NewTask(
			self.titleInput.GetText(), 
			self.descInput.GetText(),
		)
		
		err := self.userConfig.AddTask(self.boardName, task)
		if err != nil {
			utils.SaveLog(utils.Error, err.Error(), map[string]any{"boardName": self.boardName, "task": task})
		}
		
		self.Hide()
		// Re-rendering already happens after this function call.
	} else if event.ID == "<Backspace>" {
		self.focusedField.Pop()
	} else if parsedString := utils.ParseEventId(event.ID); parsedString != "" {
		// Disable new lines in the title.
		if self.focusedField == self.titleInput && parsedString == "\n" {
			return
		}
		
		self.focusedField.AppendText(parsedString)
	}
}

func (self *CreateTaskPopup) Hide() {
	self.Visible = false
	self.focusedField = self.titleInput
	self.titleInput.Flush()
	self.descInput.Flush()
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
	self.NeedsRedraw = false
	
	// Set the border to be blue on the input field that is in focus.
	if self.focusedField != nil {
		self.focusedField.GetDrawableWidget().BorderStyle = termui.NewStyle(termui.ColorBlue)
	}
	
	termui.Render(
		self.GetAllDrawableWidgets()...
	)
}
