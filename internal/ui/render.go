package ui

import (
	"log"
	"math"
	"strings"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/okira-e/gotasks/internal/domain"
	"github.com/okira-e/gotasks/internal/utils"
)

// render re-renders the UI based on the current state of the app.
func (app *App) render() {
	app.widgetsToRender = []termui.Drawable{}
	defer termui.Clear()

	app.width, app.height = termui.TerminalDimensions()

	// App has the entire userConfig object that includes all the boards.
	// Here we are extracting the pointer for this specific board that
	// ther renderer was asked to render and are passing it down.
	boardOpt := app.userConfig.GetBoard(app.boardName)
	board := boardOpt.Expect("Board found to be null. boardName: " + app.boardName)

	app.drawColumnHeaders(board)
	app.drawTasks(board)

	if app.shouldRenderCreateTaskPopup {
		app.drawCreateTaskPopup()
	}

	termui.Render(app.widgetsToRender...)
}

func (app *App) drawColumnHeaders(board *domain.Board) {
	// @Cleanup: It maybe better to move this initialization of columns when there isn't any (on first startup)
	// to somewhere else outside the renderer. So mutating the board for an initial state is handled before
	// asking the renderer to render anything.
	if len(board.Columns) == 0 {
		err := app.userConfig.AddColumnToBoard(board.Name, "Todo")
		if err != nil {
			log.Fatalf("Failed to add a %s column. %s", "Todo", err)
		}
		err = app.userConfig.AddColumnToBoard(board.Name, "In Progress")
		if err != nil {
			log.Fatalf("Failed to add a %s column. %s", "In Progress", err)
		}
		err = app.userConfig.AddColumnToBoard(board.Name, "Done")
		if err != nil {
			log.Fatalf("Failed to add a %s column. %s", "Done", err)
		}
	}

	boardOpt := app.userConfig.GetBoard(board.Name)
	board = boardOpt.Unwrap()

	widgetWidth := app.width / len(board.Columns)

	for i, columnName := range board.Columns {
		widget := widgets.NewParagraph()
		app.applyTheme(widget)
		widget.Border = true

		x1 := i * widgetWidth
		x2 := x1 + widgetWidth
		y1 := 0
		y2 := 3

		widget.SetRect(x1, y1, x2, y2)

		// Center the text
		widget.Text = centerText(columnName, widgetWidth, true)
		// log.Fatalf("THEME: %s", app.theme)
		if app.theme.IsSome() {
			widget.TextStyle = utils.Cond(
				app.theme.Unwrap() == "dark",
				termui.NewStyle(termui.ColorWhite),
				termui.NewStyle(termui.ColorBlack),
			)
		}
		widget.WrapText = true

		app.widgetsToRender = append(app.widgetsToRender, widget)
	}
}

func (app *App) drawTasks(board *domain.Board) {
	widgetWidth := app.width / len(board.Columns)
	widthPadding := 4

	for columnIndex, columnName := range board.Columns {
		// We sum them up instead of doing `rowIndex * widgetLength` because each widget has a different length.
		differentWidgetsLengths := []int{}

		for _, ticket := range board.Tasks[columnName] {
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
			app.applyTheme(widget)
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

			y1 := sumOfPreviousWidgetsLengths + 3 // 3 here is the y length of the header.
			y2 := y1 + widgetLength

			widget.SetRect(x1, y1, x2, y2)

			app.widgetsToRender = append(app.widgetsToRender, widget)

			differentWidgetsLengths = append(differentWidgetsLengths, widgetLength)
		}
	}
}

func (app *App) drawCreateTaskPopup() {
	boxWidget := termui.NewBlock()
	app.applyTheme(boxWidget)
	boxWidget.SetRect(
		app.width/4,
		app.height/4,
		app.width/4*3,
		app.height/4*3,
	)

	app.widgetsToRender = append(app.widgetsToRender, boxWidget)
}

// applyTheme applies the theme set on the App object to any widget given.
func (app *App) applyTheme(widget termui.Drawable) {
	if app.theme.IsSome() {
		color := utils.Cond(
			app.theme.Unwrap() == "dark",
			termui.NewStyle(termui.ColorWhite),
			termui.NewStyle(termui.ColorBlack),
		)

		switch widget := widget.(type) {
		case *widgets.Paragraph:
			{
				widget.TextStyle = color
				widget.BorderStyle = color
				widget.TitleStyle = color
			}
		case *widgets.List:
			{
				widget.TextStyle = color
				widget.BorderStyle = color
				widget.TitleStyle = color
			}
		case *widgets.Table:
			{
				widget.BorderStyle = color
				widget.TitleStyle = color
			}

		case *widgets.BarChart:
			{
				widget.BorderStyle = color
				widget.TitleStyle = color
			}
		case *widgets.Gauge:
			{
				widget.BorderStyle = color
				widget.TitleStyle = color
			}
		case *widgets.PieChart:
			{
				widget.BorderStyle = color
				widget.TitleStyle = color
			}
		case *widgets.SparklineGroup:
			{
				widget.BorderStyle = color
				widget.TitleStyle = color
			}
		case *termui.Block:
			{
				widget.BorderStyle = color
				widget.TitleStyle = color
			}
		case *termui.Canvas:
			{
				widget.BorderStyle = color
				widget.TitleStyle = color
			}
		case *termui.Grid:
			{
				widget.BorderStyle = color
				widget.TitleStyle = color
			}
		}
	}
}
