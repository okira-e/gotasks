package components

import (
	"log"
	"math"
	"strings"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/okira-e/gotasks/internal/domain"
	"github.com/okira-e/gotasks/internal/opt"
	"github.com/okira-e/gotasks/internal/utils"
)

type TasksViewComponent struct {
	// ID of the task that should be in focus.
	TaskInFocus		*domain.Task
	
	// Represents a card on the board. Wherever it is.
	tasksWidgets 	[]*widgets.Paragraph
	width        	int
	height       	int
	board        	*domain.Board
	userConfig      *domain.UserConfig
	filter 			opt.Option[string]
	scroll			int
}

func NewTasksViewComponent(fullWidth int, fullHeight int, board *domain.Board, userConfig *domain.UserConfig) *TasksViewComponent {
	ret := new(TasksViewComponent)
	
	ret.width = fullWidth
	ret.height = fullHeight
	ret.board = board
	ret.userConfig = userConfig
	ret.tasksWidgets = []*widgets.Paragraph{}
	
	ret.tasksWidgets = ret.drawTasks()
	
	return ret
}

// HandleKeymap changes the reference in self.taskInFocus
func (self *TasksViewComponent) HandleKeymap(key string) {
	if self.TaskInFocus == nil {
		self.SetDefaultFocusedWidget()
	}
	
	colName, _ := self.board.GetColumnForTask(self.TaskInFocus)
	tasks := self.getFilteredTasks(colName)
	
	// @Speed: Movement now is an O(n) operation on every key stroke because we use a simple dynamic array
	// to store tasks for each column. A more sophesticated DS like a Linked List would benefit vertical 
	// movemnet here for example. But n here is so small that it isn't worth it to waste a second optimizing this.
	switch key {
	case "j", "<Down>", "k", "<Up>":
		if key == "k" || key == "<Up>" {
			for i := range tasks {
				if tasks[i].Id == self.TaskInFocus.Id {
					if i - 1 >= 0 {
						self.TaskInFocus = tasks[i - 1]
					}
					
					break
				}
			}
		} else if key == "j" || key == "<Down>" {
			for i := range tasks {
				if tasks[i].Id == self.TaskInFocus.Id {
					if i + 1 <= len(tasks) - 1 {
						self.TaskInFocus = tasks[i + 1]
					}
					
					break
				}
			}
		}
		
	case "h", "<Left>", "l", "<Right>":
		if len(self.board.Columns) == 0 {
			return
		}
		
		self.scroll = 0 // Reset the scroll to be ontop
		
		_, columnIndexForTask := self.board.GetColumnForTask(self.TaskInFocus)
		var columnToMoveTo string
		
		if key == "l" || key == "<Right>" {
			nextColumnIndex := columnIndexForTask + 1
			
			if nextColumnIndex >= len(self.board.Columns) {
				return
			}
			
			nextColumnName := self.board.Columns[nextColumnIndex]
			tasksInNextColumn := self.getFilteredTasks(nextColumnName)
			
			for len(tasksInNextColumn) == 0 {
				nextColumnIndex += 1
				
				if nextColumnIndex > len(self.board.Columns) - 1 {
					return
				}
				
				nextColumnName = self.board.Columns[nextColumnIndex]
				tasksInNextColumn = self.getFilteredTasks(nextColumnName)
			}
			
			columnToMoveTo = nextColumnName
		} else {
			prevColumnIndex := columnIndexForTask - 1
			
			if prevColumnIndex < 0 {
				return
			}
			
			prevColumnName := self.board.Columns[prevColumnIndex]
			tasksInPrevColumn := self.getFilteredTasks(prevColumnName)
			
			for len(tasksInPrevColumn) == 0 {
				prevColumnIndex -= 1
				
				if prevColumnIndex < 0 {
					return
				}
				
				prevColumnName = self.board.Columns[prevColumnIndex]
				tasksInPrevColumn = self.getFilteredTasks(prevColumnName)
			}
			
			columnToMoveTo = prevColumnName
		}
		
		tasksToMoveTo := self.getFilteredTasks(columnToMoveTo)
		
		self.TaskInFocus = tasksToMoveTo[0]

	case "n":
		// Scroll infinitely for now.
		self.scroll += 1
		self.SetDefaultFocusedWidget()
		
	case "p":
		newScroll := self.scroll - 1
		
		if newScroll >= 0 {
			self.scroll = newScroll
		}
		
	case "g":
		column, i := self.board.GetColumnForTask(self.TaskInFocus)
		if i == -1 {
			log.Fatalf("Failed to find the column for task on scrolling to top.")
		}
		
		self.scroll = 0
		self.setFocusOnTopTask(column)
		
	case "G":
		column, i := self.board.GetColumnForTask(self.TaskInFocus)
		if i == -1 {
			log.Fatalf("Failed to find the column for task on scrolling to bottom.")
		}
		
		self.setFocusOnBottonTask(column)
		
	case "]":
		err := self.userConfig.MoveTaskRight(self.board, self.TaskInFocus)
		if err != nil {
			utils.SaveLog(
				utils.Error, 
				"Failed to move task to the right. " + err.Error(), 
				map[string]any{
					"task": self.TaskInFocus.Title,
				},
			)
		}
		
	case "[":
		err := self.userConfig.MoveTaskLeft(self.board, self.TaskInFocus)
		if err != nil {
			utils.SaveLog(
				utils.Error, 
				"Failed to move task to the left. " + err.Error(), 
				map[string]any{
					"task": self.TaskInFocus.Title,
				},
			)
		}
	}
	
	self.UpdateTasks()
}

func (self *TasksViewComponent) UpdateTasks() {
	self.tasksWidgets = self.drawTasks()
}

func (self *TasksViewComponent) SetDefaultFocusedWidget() {
	// Set the task in focus to be the first task you encounter (doesn't necessarily mean the first column.)
	found := false
	for _, columnName := range self.board.Columns {
		if found {
			break
		}
		
		if _, ok := self.board.Tasks[columnName]; !ok {
			continue
		}
		
		tasks := self.getFilteredTasks(columnName)
		
		for _, task := range tasks {
			self.TaskInFocus = task
			found = true
			break
		}
	}
}

// SetTextFilter applies a searching phase to the state.
func (self *TasksViewComponent) SetTextFilter(filter string) {
	self.filter = utils.Cond(filter == "", opt.None[string](), opt.Some(filter))
	self.TaskInFocus = nil
}

// getFilteredTasks returns all the tasks for a column but in reverse.
func (self *TasksViewComponent) getFilteredTasks(columnName string) []*domain.Task {
	ret := []*domain.Task{}

	for i := (len(self.board.Tasks[columnName]) - 1 - self.scroll); i >= 0; i -= 1 {
		task := self.board.Tasks[columnName][i]
		
		// If a filter is porvided, make sure to only draw the tasks that match the searched for phrase
		// by skipping the ones that don't.
		if self.filter.IsSome() {
			title := strings.ToLower(task.Title)
			desc := strings.ToLower(task.Description)
			
			if !utils.IncludesFuzzy(title, self.filter.Unwrap()) && !utils.IncludesFuzzy(desc, self.filter.Unwrap()) {
				continue
			}
		}
		
		ret = append(ret, task)
	}
		
	return ret
}

func (self *TasksViewComponent) setFocusOnTopTask(columnName string) {
	if _, ok := self.board.Tasks[columnName]; !ok {
		return
	}
	
	len := len(self.board.Tasks[columnName])
	
	if len == 0 {
		return
	}
	
	self.TaskInFocus = self.board.Tasks[columnName][len - 1]
}

func (self *TasksViewComponent) setFocusOnBottonTask(columnName string) {
	if _, ok := self.board.Tasks[columnName]; !ok {
		return
	}
	
	if len(self.board.Tasks[columnName]) == 0 {
		return
	}
	
	self.TaskInFocus = self.board.Tasks[columnName][0]
}

func (self *TasksViewComponent) drawTasks() []*widgets.Paragraph {
	ret := []*widgets.Paragraph{}
	
	widgetWidth := self.width / len(self.board.Columns)
	const widthPadding = 4

	if self.TaskInFocus == nil {
		self.SetDefaultFocusedWidget()
	}
	
	for columnIndex, columnName := range self.board.Columns {
		differentWidgetsLengths := []int{}
		
		tasks := self.getFilteredTasks(columnName)
		
		for _, task := range tasks {
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

			// Set a minimum length size for every task.
			if widgetLength < 6 {
				widgetLength = 6
			}

			widget := widgets.NewParagraph()
			widget.Border = true
			
			
			// if task.Id == "74ac4c49-e5f6-4bb1-86a7-a050adb6295d" {
			if task == self.TaskInFocus{
				widget.BorderStyle = termui.NewStyle(self.userConfig.PrimaryColor)
			}
			
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
	
	// for columnIndex, columnName := range self.board.Columns {
	// 	// We sum them up instead of doing `rowIndex * widgetLength` because each widget has a different length.
	// 	differentWidgetsLengths := []int{}

	// 	// Make the paragraph widgets and append them but in reverse. So last task in self.board.Tasks["Todo"] is rendered ontop.
	// 	for i := (len(self.board.Tasks[columnName]) - 1 - self.scroll); i >= 0; i -= 1 {
	// 		task := self.board.Tasks[columnName][i]
			
	// 		widgetLength := 2 // Border lines.

	// 		widgetLength += int(math.Ceil(
	// 			float64(len(task.Title)) / float64(widgetWidth-2),
	// 		))

	// 		if task.Description != "" {
	// 			widgetLength += 1 // The separator line "-------" between the title and the description
	// 			widgetLength += int(math.Ceil(
	// 				float64(len(task.Description)) / float64(widgetWidth-2), // 2 here is for border lines
	// 			))
	// 		}

	// 		// Set a minimum length size for every task.
	// 		if widgetLength < 6 {
	// 			widgetLength = 6
	// 		}

	// 		widget := widgets.NewParagraph()
	// 		widget.Border = true
			
			
	// 		// if task.Id == "74ac4c49-e5f6-4bb1-86a7-a050adb6295d" {
	// 		if task == self.TaskInFocus{
	// 			widget.BorderStyle = termui.NewStyle(self.userConfig.PrimaryColor)
	// 		}
			
	// 		widget.WrapText = true
	// 		// widget.Text = TextEllipsis(ticket.Title, (widgetWidth - widthPadding))
	// 		widget.Text = task.Title
	// 		widget.Text += "\n"
	// 		widget.Text += strings.Repeat("-", widgetWidth-widthPadding)
	// 		widget.Text += "\n"

	// 		if task.Description != "" {
	// 			widget.Text += task.Description
	// 		} else {
	// 			// See how much the title has taken up. If it took only one line, add a new line to the description
	// 			// because it looks better.
	// 			if math.Ceil(
	// 				float64(len(task.Title))/float64(widgetWidth),
	// 			) == 1 {
	// 				widget.Text += "\n"
	// 			}

	// 			widget.Text += utils.CenterText("No description found.", widgetWidth, true)
	// 		}

	// 		widget.PaddingLeft = 1
	// 		widget.PaddingRight = 1

	// 		x1 := columnIndex * widgetWidth
	// 		x2 := x1 + widgetWidth

	// 		sumOfPreviousWidgetsLengths := 0
	// 		for _, length := range differentWidgetsLengths {
	// 			sumOfPreviousWidgetsLengths += length
	// 		}

	// 		y1 := sumOfPreviousWidgetsLengths + 3 // 3 here is the y length of the header.
	// 		y2 := y1 + widgetLength

	// 		widget.SetRect(x1, y1, x2, y2)

	// 		differentWidgetsLengths = append(differentWidgetsLengths, widgetLength)
			
	// 		ret = append(ret, widget)
	// 	}
	// }

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
	self.UpdateTasks()
	
	termui.Render(
		self.GetAllDrawableWidgets()...
	)
}