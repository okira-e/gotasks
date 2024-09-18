package ui

import (
	"os"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/okira-e/gotasks/internal/domain"
)

// App represents the entire UI entity that has both the state and behavior of
// rendering anything to the screen.
type App struct {
	UserConfig *domain.UserConfig
	Width      int
	Height     int
	// The names of the columns for this kanban (Todo, In Progress, etc).
	// ColumnHeadersWidgets is the widgets that represent the headers for each column.
	ColumnHeadersWidgets []*widgets.Paragraph
	// TicketsWidgets are a slice representing each ticket on the board.
	Tickets []*widgets.Paragraph
}

// NewApp creates a new instance of the App with initial configurations.
func NewApp(userConfig *domain.UserConfig) (*App, error) {
	if err := termui.Init(); err != nil {
		return nil, err
	}

	width, height := termui.TerminalDimensions()

	// Initialize the app struct with state
	app := &App{
		UserConfig:           userConfig,
		Width:                width,
		Height:               height,
		ColumnHeadersWidgets: []*widgets.Paragraph{},
		Tickets:              []*widgets.Paragraph{},
	}

	return app, nil
}

// Run starts the main event loop and rendering
func (app *App) Run(boardName string) {
	defer termui.Close()

	app.render(boardName)

	for event := range termui.PollEvents() {
		app.handleEvent(boardName, event)
	}
}

// Quit exits the application gracefully
func (app *App) Quit() {
	termui.Close()
	os.Exit(0)
}
