package ui

import (
	"os"

	"github.com/gizak/termui/v3"
	"github.com/okira-e/gotasks/internal/domain"
	"github.com/okira-e/gotasks/internal/opt"
	"github.com/okira-e/gotasks/internal/utils"
)

// App represents the entire UI entity that has both the state and behavior of
// rendering anything to the screen.
type App struct {
	userConfig                  *domain.UserConfig
	boardName                   string
	width                       int
	height                      int
	// theme could be "dark" or "light". Is set through an environment variable.
	theme                       opt.Option[string]
	shouldRenderCreateTaskPopup bool
	// widgetsToRender is the final slice of widget pointers to render.
	widgetsToRender []termui.Drawable
}

// NewApp creates a new instance of the App with initial configurations.
func NewApp(userConfig *domain.UserConfig) (*App, error) {
	if err := termui.Init(); err != nil {
		return nil, err
	}

	width, height := termui.TerminalDimensions()

	theme := os.Getenv("GOTASKS_THEME")
	
	// Initialize the app struct with state
	app := &App{
		userConfig: userConfig,
		width:      width,
		height:     height,
		theme:      utils.Cond(theme != "", opt.Some(theme), opt.None[string]()),
	}

	return app, nil
}

// Run starts the main event loop and rendering
func (app *App) Run(boardName string) {
	defer termui.Close()

	app.boardName = boardName

	app.render()

	for event := range termui.PollEvents() {
		app.handleEvent(boardName, event)
	}
}

// Quit exits the application gracefully
func (app *App) Quit() {
	termui.Close()
	os.Exit(0)
}
