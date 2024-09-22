package components

import (
	"math"
	"strings"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/okira-e/gotasks/internal/domain"
	"github.com/okira-e/gotasks/internal/utils"
)

type TasksViewComponent struct {
	NeedsRedraw bool
	// Represents a card on the board. Wherever it is.
	tasksWidgets []*widgets.Paragraph
	width        int
	height       int
	board        *domain.Board
}

func NewTasksViewComponent(fullWidth int, fullHeight int, board *domain.Board) *TasksViewComponent {
	ret := &TasksViewComponent{
		width:   		fullWidth,
		height:   		fullHeight,
		board:			board,
		NeedsRedraw: 	true,
		tasksWidgets: 	[]*widgets.Paragraph{},
	}
	
	ret.tasksWidgets = ret.drawTasks()

	return ret
}

func (self *TasksViewComponent) UpdateTasks() {
	self.tasksWidgets = self.drawTasks()
}

func (self *TasksViewComponent) drawTasks() []*widgets.Paragraph {
	ret := []*widgets.Paragraph{}
	
	widgetWidth := self.width / len(self.board.Columns)
	const widthPadding = 4

	// @Todo: Right now we don't have any sort of scrolling for overflowing tasks.
	for columnIndex, columnName := range self.board.Columns {
		// We sum them up instead of doing `rowIndex * widgetLength` because each widget has a different length.
		differentWidgetsLengths := []int{}

		// Make the paragraph widgets and append them but in reverse. So last task in self.board.Tasks["Todo"] is rendered ontop.
		for i := len(self.board.Tasks[columnName]) - 1; i > 0; i -= 1 {
			task := self.board.Tasks[columnName][i]
			
			widgetLength := 2 // Border lines.

			widgetLength += int(math.Ceil(
				float64(len(task.Title)) / float64(widgetWidth-2),
			))

			if task.Description != "" {
				widgetLength += 1 // The separator line "-------" between the title and the description
				widgetLength += int(math.Ceil(
					float64(len(task.Description)) / float64(widgetWidth-2), // 2 here is for border lines
				))
			}

			// Set a minimum length size for every ticket.
			if widgetLength < 6 {
				widgetLength = 6
			}

			widget := widgets.NewParagraph()
			widget.Border = true
			widget.WrapText = true
			// widget.Text = TextEllipsis(ticket.Title, (widgetWidth - widthPadding))
			widget.Text = task.Title
			widget.Text += "\n"
			widget.Text += strings.Repeat("-", widgetWidth-widthPadding)
			widget.Text += "\n"

			if task.Description != "" {
				widget.Text += task.Description
			} else {
				// See how much the title has taken up. If it took only one line, add a new line to the description
				// because it looks better.
				if math.Ceil(
					float64(len(task.Title))/float64(widgetWidth),
				) == 1 {
					widget.Text += "\n"
				}

				widget.Text += utils.CenterText("No description found.", widgetWidth, true)
			}

			widget.PaddingLeft = 1
			widget.PaddingRight = 1

			x1 := columnIndex * widgetWidth
			x2 := x1 + widgetWidth

			sumOfPreviousWidgetsLengths := 0
			for _, length := range differentWidgetsLengths {
				sumOfPreviousWidgetsLengths += length
			}

			y1 := sumOfPreviousWidgetsLengths + 3 // 3 here is the y length of the header.
			y2 := y1 + widgetLength

			widget.SetRect(x1, y1, x2, y2)

			differentWidgetsLengths = append(differentWidgetsLengths, widgetLength)
			
			ret = append(ret, widget)
		}
	}

	return ret
}

func (self *TasksViewComponent) GetAllDrawableWidgets() []termui.Drawable {
	ret := []termui.Drawable{}
	
	for _, w := range self.tasksWidgets {
		ret = append(ret, w)
	}
	
	return ret
}


func (self *TasksViewComponent) Draw() {
	self.NeedsRedraw = false
	
	self.UpdateTasks()
	
	termui.Render(
		self.GetAllDrawableWidgets()...
	)
}