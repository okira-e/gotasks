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
	NeedsRedraw  	bool
	// Represents a card on the board. Wherever it is.
	tasksWidgets []*widgets.Paragraph
}

func NewTasksViewComponent(fullWidth int, fullHeight int, columnNames []string, tasks map[string][]domain.Task) *TasksViewComponent {
	ret := &TasksViewComponent{
		NeedsRedraw: true,
		tasksWidgets: []*widgets.Paragraph{},
	}
	
	widgetWidth := fullWidth / len(columnNames)
	const widthPadding = 4

	for columnIndex, columnName := range columnNames {
		// We sum them up instead of doing `rowIndex * widgetLength` because each widget has a different length.
		differentWidgetsLengths := []int{}

		for _, ticket := range tasks[columnName] {
			widgetLength := 2 // Border lines.

			widgetLength += int(math.Ceil(
				float64(len(ticket.Title)) / float64(widgetWidth-2),
			))

			if ticket.Description != "" {
				widgetLength += 1 // The separator line "-------" between the title and the description
				widgetLength += int(math.Ceil(
					float64(len(ticket.Description)) / float64(widgetWidth-2), // 2 here is for border lines
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
			widget.Text = ticket.Title
			widget.Text += "\n"
			widget.Text += strings.Repeat("-", widgetWidth-widthPadding)
			widget.Text += "\n"

			if ticket.Description != "" {
				widget.Text += ticket.Description
			} else {
				// See how much the title has taken up. If it took only one line, add a new line to the description
				// because it looks better.
				if math.Ceil(
					float64(len(ticket.Title))/float64(widgetWidth),
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
			
			ret.tasksWidgets = append(ret.tasksWidgets, widget)
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
	
	termui.Render(
		self.GetAllDrawableWidgets()...
	)
}