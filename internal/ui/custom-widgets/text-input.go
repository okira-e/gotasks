package customwidgets

import (
	"fmt"
	"strings"

	"github.com/gizak/termui/v3/widgets"
	"github.com/jinzhu/copier"
	"github.com/okira-e/gotasks/internal/utils"
)

type TextInput struct {
	textBuilder    strings.Builder
	drawableWidget *widgets.Paragraph
}

func NewTextInput() *TextInput {
	var builder strings.Builder

	return &TextInput{
		textBuilder:    builder,
		drawableWidget: widgets.NewParagraph(),
	}
}

// GetDrawableWidget returns the drawable widget for the text input in its drawable ready state.
func (self *TextInput) GetDrawableWidget() *widgets.Paragraph {
	self.drawableWidget.Text = self.textBuilder.String() + "❚"
	
	copier.Copy(&self.drawableWidget, &self)
	
	return self.drawableWidget
}

func (self *TextInput) GetText() string {
	return self.textBuilder.String()
}

// Flush clears the text input.
func (self *TextInput) Flush() {
	self.textBuilder.Reset()
}

// Pop removes the last character from the text input.
func (self *TextInput) Pop() {
	text := self.textBuilder.String()
	if len(text) > 0 {
		self.textBuilder.Reset()
		self.textBuilder.WriteString(text[:len(text) - 1])
	}
}

// AppendText appends the text event to the text input.
func (self *TextInput) AppendText(textEvent string) {
	if len(textEvent) == 1 {
		if _, err := self.textBuilder.WriteString(textEvent); err != nil {
			utils.ExitApp(
				fmt.Sprintf("Failed to write to input widget. %s", err),
			)
		}
	}
}
