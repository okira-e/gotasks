package components

import (
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/okira-e/gotasks/internal/ui/types"
	"github.com/okira-e/gotasks/internal/utils"
)


type ConfirmationComponent struct {
	Visible bool
	Action 	func(bool)
	
	window	*types.Window
	widget 	*widgets.Paragraph
}

func NewConfirmationPopupComponent(window *types.Window) *ConfirmationComponent {
	ret := new(ConfirmationComponent)
	
	ret.widget = widgets.NewParagraph()
	ret.widget.Title = "Confirmation"
	ret.window = window
	
	return ret
}

func (self *ConfirmationComponent) SetMessageAndAction(message string, action func(bool)) {
	const widgetHeight = 4
	widgetWidth := len(message) + 4
	
	self.widget.SetRect(
		self.window.Width / 2 - widgetWidth / 2,
		self.window.Height / 2 - widgetHeight / 2,
		
		self.window.Width / 2 + widgetWidth / 2,
		self.window.Height / 2 + widgetHeight / 2,
	)
	
	self.widget.Text = utils.CenterText(message, widgetWidth, true)
	self.widget.Text += "\n" + utils.CenterText("y/N", widgetWidth, true)
	
	self.widget.Border = true
	
	self.Action = action
}

// HandleInput handles keyboard inputs sent to this component. It returns a boolean
// indicating if we should clear before we re-render.
func (self *ConfirmationComponent) HandleInput(event termui.Event) bool {
	switch event.ID {
	case "y", "Y": // Confirm
		self.Action(true)
	case "n", "N": // Cancel
		self.Action(false)
	}
	
	self.Hide()
	return true
}

func (self *ConfirmationComponent) Hide() {
	self.Visible = false
}

func (self *ConfirmationComponent) Show() {
	self.Visible = true
}

func (self *ConfirmationComponent) Draw() {
	termui.Render(
		self.widget,
	)
}