package components

import (
	"log"
	"math"
	"strings"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/okira-e/gotasks/internal/domain"
	"github.com/okira-e/gotasks/internal/opt"
	"github.com/okira-e/gotasks/internal/ui/types"
	"github.com/okira-e/gotasks/internal/utils"
)

type TasksViewComponent struct {
	// ID of the task that should be in focus.
	TaskInFocus *domain.Task

	// Represents a card on the board. Wherever it is.
	tasksWidgets          []*widgets.Paragraph
	window                *types.Window
	board                 *domain.Board
	userConfig            *domain.UserConfig
	filter                opt.Option[string]
	scroll                int
	// goToFirstTaskInColumn Tells the draw function to set the 
	// the task in focus to be the pointer to the first task
	// in the column list that is set.
	// Its value is the name of a column.
	goToFirstTaskInColumn opt.Option[string]
}

func NewTasksViewComponent(window *types.Window, board *domain.Board, userConfig *domain.UserConfig) *TasksViewComponent {
	ret := new(TasksViewComponent)
	
	ret.window = window
	ret.board = board
	ret.userConfig = userConfig
	ret.tasksWidgets = []*widgets.Paragraph{}
	
	ret.tasksWidgets = ret.drawTasks()
	
	return ret
}

// HandleKeymap changes the reference in self.taskInFocus
// Returns a flag saying if the renderer should clear the old view.
func (self *TasksViewComponent) HandleKeymap(key string) bool {
	shouldClear := false
	
	if self.TaskInFocus == nil {
		self.SetDefaultFocusedWidget()
	}
	
	colName, _ := self.board.GetColumnForTask(self.TaskInFocus)
	tasks := self.getFilteredTasks(colName)
	
	// @Speed: Movement now is an O(n) operation on every key stroke because we use a simple dynamic array
	// to store tasks for each column. A more sophesticated DS like a Linked List would benefit vertical 
	// movement here for example.
	switch key {
	case "j", "<Down>", "k", "<Up>":
		if key == "k" || key == "<Up>" {
			for i := range tasks {
				if tasks[i].Id == self.TaskInFocus.Id {
					// Scroll up one task if you're on a the first task in the view but not in the list.
					if i == 0 && self.scroll > 0 {
						self.scroll -= 1
						
						self.goToFirstTaskInColumn = opt.Some(colName)
					} else if i - 1 >= 0 { // If we're not the first task in the list.
						self.TaskInFocus = tasks[i - 1]
					}
					
					break
				}
			}
		} else if key == "j" || key == "<Down>" {
			for i := range tasks {
				if tasks[i].Id == self.TaskInFocus.Id {
					if i + 1 <= len(tasks) - 1 { // We're not the last task in the list.
						self.TaskInFocus = tasks[i + 1]
					}
					
					break
				}
			}
		}
		
	case "h", "<Left>", "l", "<Right>":
		if len(self.board.Columns) == 0 {
			return shouldClear
		}
		
		self.scroll = 0 // Reset the scroll to be ontop
		
		_, columnIndexForTask := self.board.GetColumnForTask(self.TaskInFocus)
		var columnToMoveTo string
		
		if key == "l" || key == "<Right>" {
			nextColumnIndex := columnIndexForTask + 1
			
			if nextColumnIndex >= len(self.board.Columns) {
				return shouldClear
			}
			
			nextColumnName := self.board.Columns[nextColumnIndex]
			tasksInNextColumn := self.getFilteredTasks(nextColumnName)
			
			for len(tasksInNextColumn) == 0 {
				nextColumnIndex += 1
				
				if nextColumnIndex > len(self.board.Columns) - 1 {
					return shouldClear
				}
				
				nextColumnName = self.board.Columns[nextColumnIndex]
				tasksInNextColumn = self.getFilteredTasks(nextColumnName)
			}
			
			columnToMoveTo = nextColumnName
		} else {
			prevColumnIndex := columnIndexForTask - 1
			
			if prevColumnIndex < 0 {
				return shouldClear
			}
			
			prevColumnName := self.board.Columns[prevColumnIndex]
			tasksInPrevColumn := self.getFilteredTasks(prevColumnName)
			
			for len(tasksInPrevColumn) == 0 {
				prevColumnIndex -= 1
				
				if prevColumnIndex < 0 {
					return shouldClear
				}
				
				prevColumnName = self.board.Columns[prevColumnIndex]
				tasksInPrevColumn = self.getFilteredTasks(prevColumnName)
			}
			
			columnToMoveTo = prevColumnName
		}
		
		tasksToMoveTo := self.getFilteredTasks(columnToMoveTo)
		
		self.TaskInFocus = tasksToMoveTo[0]

	case "n":
		// @Todo: This scrolls infinitely for now.
		self.scroll += 1
		
		// Move the task in focus to the first of its list if scrolled beyond it
		for i := range tasks {
			if tasks[i] == self.TaskInFocus {
				if self.scroll > i {
					self.goToFirstTaskInColumn = opt.Some(colName)
				}
			}
		}
		
		shouldClear = true
		
	case "p":
		newScroll := self.scroll - 1
		
		if newScroll >= 0 {
			self.scroll = newScroll
		}

		shouldClear = true
		
	case "g":
		columnName, i := self.board.GetColumnForTask(self.TaskInFocus)
		if i == -1 {
			log.Fatalf("Failed to find the columnName for task on scrolling to top.")
		}
		
		self.setFocusOnTopTask(columnName)
		self.scroll = 0
		shouldClear = true
		
	case "G":
		columnName, i := self.board.GetColumnForTask(self.TaskInFocus)
		if i == -1 {
			log.Fatalf("Failed to find the columnName for task on scrolling to bottom.")
		}
		
		self.setFocusOnBottomTask(columnName)
		self.scroll = len(self.board.Tasks[columnName]) - 1
			
		shouldClear = true
		
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
		shouldClear = true
		
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
		shouldClear = true
		
	}
	
	self.UpdateTasks()
	
	return shouldClear
}

func (self *TasksViewComponent) UpdateTasks() {
	self.tasksWidgets = self.drawTasks()
}

func (self *TasksViewComponent) SetDefaultFocusedWidget() {
	if self.board.IsEmpty() {
		self.TaskInFocus = nil
		return
	}
	
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

// getFilteredTasks returns all the tasks for a column but in reverse accounting 
// for the current scroll value (therefore tasks that we scrolled beyond aren't even 
// accounted for, or rendered) as well as filtered if a filter is in effect. 
// They are filtered because in the board we should show the last added task first.
func (self *TasksViewComponent) getFilteredTasks(columnName string) []*domain.Task {
	ret := []*domain.Task{}

	for i := (len(self.board.Tasks[columnName]) - 1 - self.scroll); i >= 0; i -= 1 {
		task := self.board.Tasks[columnName][i]
		
		// If a filter is provided, make sure to only draw the tasks that match the searched for phrase
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

func (self *TasksViewComponent) setFocusOnBottomTask(columnName string) {
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
	
	widgetWidth := self.window.Width / len(self.board.Columns)
	const widthPadding = 4

	if self.TaskInFocus == nil {
		self.SetDefaultFocusedWidget()
	}
	
	// Set the task in focus to be the first one for its column if the flag is set
	// from the last render.
	if self.goToFirstTaskInColumn.IsSome() {
		// No need to check for an empty board because logically we shouldn't be calling this
		// in that case.
		

		// Set the task in focus to be the first task you encounter (doesn't necessarily mean the first column.)
		found := false
		for _, columnName := range self.board.Columns {
			if columnName != self.goToFirstTaskInColumn.Unwrap() {
				continue
			}
			
			if found {
				break
			}
			
			tasks := self.getFilteredTasks(columnName)
			
			for _, task := range tasks {
				self.TaskInFocus = task
				found = true
				break
			}
		}
	}
	self.goToFirstTaskInColumn = opt.None[string]()
	
	
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