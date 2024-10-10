package customwidgets

import (
	"strings"

	"github.com/gizak/termui/v3/widgets"
	"github.com/jinzhu/copier"
)

type TextInput struct {
	textBuilder    strings.Builder
	drawableWidget *widgets.Paragraph
	cursorPosition int
}

func NewTextInput() *TextInput {
	var builder strings.Builder

	ret := new(TextInput)
	
	ret.textBuilder = builder
	ret.drawableWidget = widgets.NewParagraph()
	
	return ret
}

// GetDrawableWidget returns the drawable widget for the text input in its drawable ready state.
func (self *TextInput) GetDrawableWidget() *widgets.Paragraph {
	self.drawableWidget.Text = self.textBuilder.String()// + "❚"
	
	// Insert cursor (❚) at the cursor position
	cursorChar := "❚"
	if self.cursorPosition >= len(self.drawableWidget.Text) {
		// If cursor is at the end of the text
		textWithCursor := self.drawableWidget.Text + cursorChar
		self.drawableWidget.Text = textWithCursor
	} else {
		// Split the text and insert cursor at the cursor position
		textWithCursor := self.drawableWidget.Text[:self.cursorPosition] + cursorChar + self.drawableWidget.Text[self.cursorPosition:]
		self.drawableWidget.Text = textWithCursor
	}
		
	copier.Copy(&self.drawableWidget, &self)
	
	return self.drawableWidget
}

func (self *TextInput) SetText(text string) {
	self.textBuilder.WriteString(text)
	self.cursorPosition = self.textBuilder.Len()
}

func (self *TextInput) GetText() string {
	return self.textBuilder.String()
}

func (self *TextInput) MoveCursorRight() {
	if self.cursorPosition < self.textBuilder.Len() {
		self.cursorPosition += 1
	}
}

func (self *TextInput) MoveCursorLeft() {
	if self.cursorPosition > 0 {
		self.cursorPosition -= 1
	}
}

func (self *TextInput) MoveCursorRightOneWord() {
	var indexForTextSinceNextWord int
	for i := self.textBuilder.Len() - 1; i > self.cursorPosition; i -= 1 {
		char := self.textBuilder.String()[i]
		
		if string(char) == " " {
			indexForTextSinceNextWord = i
		}
	}
	
	if indexForTextSinceNextWord != 0 {
		self.moveCursorToPosition(indexForTextSinceNextWord)
	} else { // Cursor is probably at the last word
		self.moveCursorToPosition(self.textBuilder.Len())
	}
}

func (self *TextInput) MoveCursorLeftOneWord() {
	textTillCursor := self.textBuilder.String()[:self.cursorPosition]

	var indexForTextSinceLastWord int
	for i, char := range textTillCursor {
		if string(char) == " " {
			indexForTextSinceLastWord = i
		}
	}
	
	self.moveCursorToPosition(indexForTextSinceLastWord)
}

// Flush clears the text input.
func (self *TextInput) Flush() {
	self.textBuilder.Reset()
	self.cursorPosition = 0
}

// PopChar removes one character to the left of the cursor.
func (self *TextInput) PopChar() {
	if self.textBuilder.Len() == 0 {
		return
	}
	
	if self.cursorPosition == 0 {
		return
	}
	
	text := self.textBuilder.String()
	self.textBuilder.Reset()
	
	newText := text[:self.cursorPosition - 1] + text[self.cursorPosition:]
	
	self.textBuilder.WriteString(newText)
	self.MoveCursorLeft()
}

// PopWord removes one word to the left of the cursor.
func (self *TextInput) PopWord() {
	if self.textBuilder.Len() == 0 {
		return
	}
	
	if self.cursorPosition == 0 {
		return
	}
	
	textTillCursor := self.textBuilder.String()[:self.cursorPosition]

	var indexForTextSinceLastWord int
	for i, char := range textTillCursor {
		if string(char) == " " {
			indexForTextSinceLastWord = i
		}
	}	
	
	text := self.textBuilder.String()
	self.textBuilder.Reset()
	
	newText := text[:indexForTextSinceLastWord] + text[self.cursorPosition:]
	
	self.moveCursorToPosition(indexForTextSinceLastWord)
	self.textBuilder.WriteString(newText)
}

// AppendText appends the text event to the text input.
func (self *TextInput) AppendText(textEvent string) {
	if len(textEvent) == 1 {
		text := self.textBuilder.String()
		self.textBuilder.Reset()
		
		newText := text[:self.cursorPosition] + textEvent + text[self.cursorPosition:]
		
		self.textBuilder.WriteString(newText)
		self.MoveCursorRight()
	}
}

func (self *TextInput) moveCursorToPosition(position int) {
	self.cursorPosition = position
}
