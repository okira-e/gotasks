package ui

import (
	"log"
	"math"
	"strings"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/okira-e/gotasks/internal/domain"
)

// render re-renders the UI based on the current state of the app.
func (app *App) render(boardName string) {
	widgetsToRender := []termui.Drawable{}

	app.Width, app.Height = termui.TerminalDimensions()

	// App has the entire userConfig object that includes all the boards.
	// Here we are extracting the pointer for this specific board that
	// ther renderer was asked to render and are passing it down.
	boardOpt := app.UserConfig.GetBoard(boardName)
	board := boardOpt.Expect("Board found to be null. boardName: " + boardName)

	app.drawColumnHeaders(board)
	for _, columnWidget := range app.ColumnHeadersWidgets {
		widgetsToRender = append(widgetsToRender, columnWidget)
	}

	app.drawTickets(board)
	for _, ticket := range app.Tickets {
		widgetsToRender = append(widgetsToRender, ticket)
	}

	termui.Clear()
	termui.Render(widgetsToRender...)
}

func (app *App) drawColumnHeaders(board *domain.Board) {
	// @Cleanup: It maybe better to move this initialization of columns when there isn't any (on first startup)
	// to somewhere else outside the renderer. So mutating the board for an initial state is handled before
	// asking the renderer to render anything.
	if len(board.Columns) == 0 {
		err := app.UserConfig.AddColumnToBoard(board.Name, "Todo")
		if err != nil {
			log.Fatalf("Failed to add a %s column. %s", "Todo", err)
		}
		err = app.UserConfig.AddColumnToBoard(board.Name, "In Progress")
		if err != nil {
			log.Fatalf("Failed to add a %s column. %s", "In Progress", err)
		}
		err = app.UserConfig.AddColumnToBoard(board.Name, "Done")
		if err != nil {
			log.Fatalf("Failed to add a %s column. %s", "Done", err)
		}
	}

	boardOpt := app.UserConfig.GetBoard(board.Name)
	board = boardOpt.Unwrap()

	widgetWidth := app.Width / len(board.Columns)

	for i, columnName := range board.Columns {
		widget := widgets.NewParagraph()
		widget.Border = true

		x1 := i * widgetWidth
		x2 := x1 + widgetWidth
		y1 := 0
		y2 := 3

		widget.SetRect(x1, y1, x2, y2)

		// Center the text
		widget.Text = centerText(columnName, widgetWidth, true)
		widget.WrapText = true

		app.ColumnHeadersWidgets = append(app.ColumnHeadersWidgets, widget)
	}
}

func (app *App) drawTickets(board *domain.Board) {
	widgetWidth := app.Width / len(board.Columns)
	widthPadding := 4

	for columnIndex, columnName := range board.Columns {
		// We sum them up instead of doing `rowIndex * widgetLength` because each widget has a different length.
		differentWidgetsLengths := []int{}

		for _, ticket := range board.Tickets[columnName] {
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
			// widget.Text = textEllipsis(ticket.Title, (widgetWidth - widthPadding))
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

				widget.Text += centerText("No description found.", widgetWidth, true)
			}

			widget.PaddingLeft = 1
			widget.PaddingRight = 1

			x1 := columnIndex * widgetWidth
			x2 := x1 + widgetWidth

			sumOfPreviousWidgetsLengths := 0
			for _, length := range differentWidgetsLengths {
				sumOfPreviousWidgetsLengths += length
			}

			y1 := sumOfPreviousWidgetsLengths + app.ColumnHeadersWidgets[0].Dy()
			y2 := y1 + widgetLength

			widget.SetRect(x1, y1, x2, y2)

			app.Tickets = append(app.Tickets, widget)

			differentWidgetsLengths = append(differentWidgetsLengths, widgetLength)
		}
	}
}
