package ui

import (
	"errors"
	"os"

	"github.com/gizak/termui/v3"
	"github.com/okira-e/gotasks/internal/domain"
	"github.com/okira-e/gotasks/internal/ui/components"
	"github.com/okira-e/gotasks/internal/vars"
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
	confirmationPopup				*components.ConfirmationComponent
}

// NewApp creates a new instance of the App with initial configurations.
func NewApp(userConfig *domain.UserConfig, boardName string) (*App, error) {
	if err := termui.Init(); err != nil {
		return nil, err
	}

	width, height := termui.TerminalDimensions()
	
	boardOpt := userConfig.GetBoard(boardName)
	if boardOpt.IsNone() {
		return nil, errors.New("Couldn't find the board while trying to add a task")
	}

	board := boardOpt.Unwrap()

	theme := os.Getenv(vars.ThemeFlag)
	if theme == "" {
		theme = "dark"
	}
	
	app := new(App)
	
	app.userConfig = userConfig
	app.boardName = boardName
	app.width = width
	app.height = height
	app.theme = theme
	app.createTaskPopup = components.NewCreateTaskPopupComponent(width, height, userConfig, boardName)
	app.confirmationPopup = components.NewConfirmationPopupComponent(width, height)
	app.columnsHeadersView = components.NewColumnsHeaderComponent(width, height, board.Columns)
	app.tasksView = components.NewTasksViewComponent(width, height, board, userConfig)

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
