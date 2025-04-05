package customwidgets

import (
	"strings"

	"github.com/gizak/termui/v3/widgets"
	"github.com/jinzhu/copier"
	"github.com/okira-e/gotasks/internal/utils"
)

type TextInput struct {
	textBuilder    strings.Builder
	drawableWidget *widgets.Paragraph
	cursorPosition int
	showCursor     bool
}

func NewTextInput() *TextInput {
	var builder strings.Builder

	ret := new(TextInput)

	ret.textBuilder = builder
	ret.drawableWidget = widgets.NewParagraph()
	ret.showCursor = true
	
	return ret
}

// GetDrawableWidget returns the drawable widget for the text input in its drawable ready state.
func (self *TextInput) GetDrawableWidget() *widgets.Paragraph {
	// Split text at cursor position
	beforeCursor := self.textBuilder.String()[:self.cursorPosition]
	afterCursor := ""
	if self.cursorPosition < self.textBuilder.Len() {
		afterCursor = self.textBuilder.String()[self.cursorPosition:]
	}

	// Set the character at cursor position to be highlighted (reversed)
	if self.cursorPosition < self.textBuilder.Len() {
		// Cursor is on a character - highlight it
		cursorChar := string(self.textBuilder.String()[self.cursorPosition])
		self.drawableWidget.Text = beforeCursor

		// Store normal style to restore after cursor
		self.drawableWidget.WrapText = false

		// Add the cursor and remaining text
		if self.showCursor {
			self.drawableWidget.Text += "[" + cursorChar + "](fg:black,bg:white)" + afterCursor[1:]
		} else {
			self.drawableWidget.Text += cursorChar + afterCursor[1:]
		}
		
	} else {
		// Cursor is at end of text - show a block
		if self.showCursor {
			self.drawableWidget.Text = beforeCursor + "[█](fg:white,bg:black)"
		}
	}

	copier.Copy(&self.drawableWidget, &self)

	return self.drawableWidget
}

func (self *TextInput) SetText(text string) {
	self.textBuilder.Reset()
	self.textBuilder.WriteString(text)
	self.cursorPosition = self.textBuilder.Len()
}

func (self *TextInput) GetText() string {
	return self.textBuilder.String()
}

// Flush clears the text input.
func (self *TextInput) Flush() {
	self.textBuilder.Reset()
	self.cursorPosition = 0
}

func (self *TextInput) HandleInput(char string, disableNewLines bool) {
	if char == "<Backspace>" {
		self.popChar()

	} else if char == "<C-u>" {
		self.popWord()

	} else if char == "<Right>" {
		self.moveCursorRight()

	} else if char == "<C-e>" {
		self.moveCursorEnd()

	} else if char == "<C-a>" {
		self.moveCursorStart()

	} else if char == "<Left>" {
		self.moveCursorLeft()

	} else if char == "<C-right>" || char == "<M-right>" || char == "<A-right>" {
		self.moveCursorRightOneWord()

	} else if char == "<C-left>" || char == "<M-left>" || char == "<A-left>" {
		self.moveCursorLeftOneWord()

	} else if char == "<Home>" {
		self.moveCursorStart()

	} else if char == "<End>" {
		self.moveCursorEnd()

	} else if char == "<Delete>" {
		self.deleteChar()

	} else {
		char = utils.ParseEventId(char)

		if disableNewLines && char == "\n" {
			return
		}

		self.appendText(char)
	}
}

func (self *TextInput) moveCursorRight() {
	if self.cursorPosition < self.textBuilder.Len() {
		self.cursorPosition += 1
	}
}

func (self *TextInput) moveCursorLeft() {
	if self.cursorPosition > 0 {
		self.cursorPosition -= 1
	}
}

func (self *TextInput) moveCursorStart() {
	self.cursorPosition = 0
}

func (self *TextInput) moveCursorEnd() {
	self.cursorPosition = self.textBuilder.Len()
}

func (self *TextInput) moveCursorRightOneWord() {
	text := self.textBuilder.String()

	// Start from cursor position and move right
	for i := self.cursorPosition; i < len(text); i += 1 {
		// If we're on a space, find the next non-space
		if string(text[i]) == " " {
			// Find next non-space
			for j := i + 1; j < len(text); j++ {
				if string(text[j]) != " " {
					self.moveCursorToPosition(j)
					return
				}
			}
			// If we only found spaces, move to end
			self.moveCursorToPosition(len(text))
			return
		} else if i == self.cursorPosition {
			// If we're on a non-space, find the next space
			for j := i + 1; j < len(text); j++ {
				if string(text[j]) == " " {
					self.moveCursorToPosition(j)
					return
				}
			}
			// If we found no space, move to end
			self.moveCursorToPosition(len(text))
			return
		}
	}

	// Default to end of text
	self.moveCursorToPosition(len(text))
}

func (self *TextInput) moveCursorLeftOneWord() {
	textTillCursor := self.textBuilder.String()[:self.cursorPosition]

	// If we're at start, do nothing
	if len(textTillCursor) == 0 {
		return
	}

	// Start from cursor position and move left
	foundNonSpace := false
	lastSpacePos := 0

	// Start from one character before cursor
	for i := len(textTillCursor) - 1; i >= 0; i -= 1 {
		char := string(textTillCursor[i])

		// If we find a space after finding non-space, we found a word boundary
		if char == " " && foundNonSpace {
			// Move to position after the space
			self.moveCursorToPosition(i + 1)
			return
		}

		// Track the last space we saw
		if char == " " {
			lastSpacePos = i
		} else {
			foundNonSpace = true
		}

		// If we reach the start, move there
		if i == 0 {
			self.moveCursorToPosition(0)
			return
		}
	}

	// If we only saw spaces, move to the last space position
	if !foundNonSpace && lastSpacePos > 0 {
		self.moveCursorToPosition(lastSpacePos)
	} else {
		// Default to beginning of text if we couldn't find a word boundary
		self.moveCursorToPosition(0)
	}
}

// PopChar removes one character to the left of the cursor.
func (self *TextInput) popChar() {
	if self.textBuilder.Len() == 0 || self.cursorPosition == 0 {
		return
	}

	text := self.textBuilder.String()
	self.textBuilder.Reset()

	newText := text[:self.cursorPosition-1] + text[self.cursorPosition:]

	self.textBuilder.WriteString(newText)
	self.moveCursorLeft()
}

// DeleteChar removes one character at the cursor position.
func (self *TextInput) deleteChar() {
	if self.textBuilder.Len() == 0 || self.cursorPosition >= self.textBuilder.Len() {
		return
	}

	text := self.textBuilder.String()
	self.textBuilder.Reset()

	newText := text[:self.cursorPosition] + text[self.cursorPosition+1:]

	self.textBuilder.WriteString(newText)
	// Cursor position stays the same
}

// PopWord removes one word to the left of the cursor.
func (self *TextInput) popWord() {
	if self.textBuilder.Len() == 0 || self.cursorPosition == 0 {
		return
	}

	text := self.textBuilder.String()
	textTillCursor := text[:self.cursorPosition]

	// Start from cursor position and move left
	deletePos := 0
	foundNonSpace := false

	// Start from one character before cursor
	for i := len(textTillCursor) - 1; i >= 0; i -= 1 {
		char := string(textTillCursor[i])

		// If we find a space after finding non-space, we found a word boundary
		if char == " " && foundNonSpace {
			deletePos = i + 1
			break
		}

		if char != " " {
			foundNonSpace = true
		}

		// If we reach the start, delete from beginning
		if i == 0 {
			deletePos = 0
		}
	}

	// Delete from the determined position to the cursor
	self.textBuilder.Reset()
	newText := text[:deletePos] + text[self.cursorPosition:]
	self.textBuilder.WriteString(newText)
	self.moveCursorToPosition(deletePos)
}

// AppendText appends the text event to the text input.
func (self *TextInput) appendText(textEvent string) {
	if textEvent == "" {
		return
	}

	text := self.textBuilder.String()
	self.textBuilder.Reset()

	newText := text[:self.cursorPosition] + textEvent + text[self.cursorPosition:]

	self.textBuilder.WriteString(newText)
	self.cursorPosition += len(textEvent)
}

func (self *TextInput) moveCursorToPosition(position int) {
	if position < 0 {
		position = 0
	} else if position > self.textBuilder.Len() {
		position = self.textBuilder.Len()
	}
	self.cursorPosition = position
}
