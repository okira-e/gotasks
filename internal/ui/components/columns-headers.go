package components

import (
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/okira-e/gotasks/internal/ui/types"
	"github.com/okira-e/gotasks/internal/utils"
)

type ColumnsHeaderComponent struct {
    columnBoxes []*widgets.Paragraph
    columnNames	[]string
    window		*types.Window
}

func NewColumnsHeaderComponent(window *types.Window, columnNames []string) *ColumnsHeaderComponent {
	component := new(ColumnsHeaderComponent)
	
	component.window = window
	component.columnNames = columnNames
	
	return component
}

func (self *ColumnsHeaderComponent) GetAllDrawableWidgets() []termui.Drawable {
	ret := []termui.Drawable{}
	
	for _, w := range self.columnBoxes {
		ret = append(ret, w)
	}
	
	return ret
}


func (self *ColumnsHeaderComponent) Draw() {
	widgetWidth := self.window.Width / len(self.columnNames)

	for i, columnName := range self.columnNames {
		widget := widgets.NewParagraph()
		widget.Border = true

		x1 := i * widgetWidth
		x2 := x1 + widgetWidth
		y1 := 0
		y2 := 3

		widget.SetRect(x1, y1, x2, y2)

		widget.Text = utils.CenterText(columnName, widgetWidth, true)
		
		widget.WrapText = true

		self.columnBoxes = append(self.columnBoxes, widget)
	}
	
	termui.Render(
		self.GetAllDrawableWidgets()...
	)
}