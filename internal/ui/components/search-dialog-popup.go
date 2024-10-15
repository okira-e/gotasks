package components

import (
	"strings"

	"github.com/gizak/termui/v3"
	cw "github.com/okira-e/gotasks/internal/ui/custom-widgets"
)


type SearchDialogPopupComponent struct {
	Visible bool
	
	fullWidth 	int
	fullHeight 	int
	action 		func(string)
	widget 		*cw.TextInput
}

func NewSearchDialogPopupComponent(fullWidth int, fullHeight int, action func(string)) *SearchDialogPopupComponent {
	ret := new(SearchDialogPopupComponent)
	
	ret.fullWidth 	= fullWidth
	ret.fullHeight 	= fullHeight
	ret.action 		= action
	
	ret.drawWidget()
	
	return ret
}

func (self *SearchDialogPopupComponent) drawWidget() {
	self.widget = cw.NewTextInput()
	self.widget.GetDrawableWidget().Title = "Search For"
	
	const widgetHeight = 3
	const widgetWidth = 60 // Just an arbitrary number.
	
	self.widget.GetDrawableWidget().SetRect(
		self.fullWidth / 2 - widgetWidth / 2,
		self.fullHeight / 2 - 1,
		
		self.fullWidth / 2 + widgetWidth / 2,
		self.fullHeight / 2 + 2,
	)
	
	self.widget.GetDrawableWidget().Border = true
}

func (self *SearchDialogPopupComponent) HandleInput(event termui.Event) {
	if event.ID ==  "<Escape>" {
		self.Hide()

	} else if event.ID == "<Enter>" {
		toSearchFor := self.widget.GetText()
		toSearchFor = strings.ToLower(toSearchFor)
		
		self.action(toSearchFor)
		
		self.Hide()
	} else {
		self.widget.HandleInput(event.ID, true)
	}
}

func (self *SearchDialogPopupComponent) Hide() {
	self.Visible = false
	self.widget.Flush()
}

func (self *SearchDialogPopupComponent) Show() {
	self.Visible = true
}

func (self *SearchDialogPopupComponent) Draw() {
	termui.Render(
		self.widget.GetDrawableWidget(),
	)
}