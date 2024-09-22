package ui

import (
	"os"

	"github.com/gizak/termui/v3"
	"github.com/okira-e/gotasks/internal/domain"
	"github.com/okira-e/gotasks/internal/ui/components"
)


type WidgetsThatCaptureInput int

const (
	None WidgetsThatCaptureInput = iota
	CreateTaskTitle
	CreateTaskDescription
)

type Component interface {
	GetAllDrawableWidgets() []termui.Drawable
	Draw()
}

// App represents the entire UI entity that has both the state and behavior of
// rendering anything to the screen.
type App struct {
	userConfig                  	*domain.UserConfig
	boardName                   	string
	width                       	int
	height                      	int
	// theme could be "dark" or "light". Is set through an environment variable.
	theme                       	string
	createTaskPopup					*components.CreateTaskPopup
	columnsHeadersView				*components.ColumnsHeaderComponent
	tasksView						*components.TasksViewComponent
}

// NewApp creates a new instance of the App with initial configurations.
func NewApp(userConfig *domain.UserConfig, boardName string) (*App, error) {
	if err := termui.Init(); err != nil {
		return nil, err
	}

	width, height := termui.TerminalDimensions()
	
	// @Todo: Expect here exits the program without calling termui.Close()
	boardOpt := userConfig.GetBoard(boardName)
	board := boardOpt.Expect("Board found to be null. boardName: " + boardName)

	theme := os.Getenv("GOTASKS_THEME")
	if theme == "" {
		theme = "dark"
	}
	
	app := &App{
		userConfig: 			userConfig,
		width:      			width,
		height:     			height,
		theme: 					theme,
		createTaskPopup: 		components.NewCreateTaskPopup(width, height),
		columnsHeadersView: 	components.NewColumnsHeaderComponent(width, height, board.Columns),
		tasksView: 				components.NewTasksViewComponent(width, height, board.Columns, board.Tasks),
	}

	return app, nil
}

// Run starts the main event loop and rendering
func (app *App) Run() {
	defer termui.Close()

	app.render()

	for event := range termui.PollEvents() {
		app.handleEvent(event)
	}
}

// Quit exits the application gracefully
func (app *App) Quit() {
	termui.Close()
	os.Exit(0)
}

func applyTheme(c Component, theme string) {
	for _, widget := range c.GetAllDrawableWidgets() {
		ColorizeWidget(widget, theme)
	}
}
