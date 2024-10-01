package components

import (
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/okira-e/gotasks/internal/utils"
)

type ColumnsHeaderComponent struct {
    columnBoxes []*widgets.Paragraph
}

func NewColumnsHeaderComponent(fullWidth int, fullHeight int, columnNames []string) *ColumnsHeaderComponent {
	component := new(ColumnsHeaderComponent)
	
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

		component.columnBoxes = append(component.columnBoxes, widget)
	}

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
	termui.Render(
		self.GetAllDrawableWidgets()...
	)
}