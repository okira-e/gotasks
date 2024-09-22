package components

import (
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/okira-e/gotasks/internal/utils"
)

type ColumnsHeaderComponent struct {
	NeedsRedraw  	bool
    columnBoxes []*widgets.Paragraph
}

func NewColumnsHeaderComponent(fullWidth int, fullHeight int, columnNames []string) *ColumnsHeaderComponent {
	ret := &ColumnsHeaderComponent{
		NeedsRedraw: true,
		columnBoxes: []*widgets.Paragraph{},
	}
	
	// @Cleanup: It maybe better to move this initialization of columns when there isn't any (on first startup)
	// to somewhere else outside the renderer. So mutating the board for an initial state is handled before
	// asking the renderer to render anything.
	// if len(columnNames) == 0 {
	// 	err := app.userConfig.AddColumnToBoard(board.Name, "Todo")
	// 	if err != nil {
	// 		log.Fatalf("Failed to add a %s column. %s", "Todo", err)
	// 	}
	// 	err = app.userConfig.AddColumnToBoard(board.Name, "In Progress")
	// 	if err != nil {
	// 		log.Fatalf("Failed to add a %s column. %s", "In Progress", err)
	// 	}
	// 	err = app.userConfig.AddColumnToBoard(board.Name, "Done")
	// 	if err != nil {
	// 		log.Fatalf("Failed to add a %s column. %s", "Done", err)
	// 	}
	// }

	widgetWidth := fullWidth / len(columnNames)

	for i, columnName := range columnNames {
		widget := widgets.NewParagraph()
		widget.Border = true

		x1 := i * widgetWidth
		x2 := x1 + widgetWidth
		y1 := 0
		y2 := 3

		widget.SetRect(x1, y1, x2, y2)

		widget.Text = utils.CenterText(columnName, widgetWidth, true)
		
		widget.WrapText = true

		ret.columnBoxes = append(ret.columnBoxes, widget)
	}

	return ret
}

func (self *ColumnsHeaderComponent) GetAllDrawableWidgets() []termui.Drawable {
	ret := []termui.Drawable{}
	
	for _, w := range self.columnBoxes {
		ret = append(ret, w)
	}
	
	return ret
}


func (self *ColumnsHeaderComponent) Draw() {
	self.NeedsRedraw = false
	
	termui.Render(
		self.GetAllDrawableWidgets()...
	)
}