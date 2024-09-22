package ui

import (
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/okira-e/gotasks/internal/utils"
)

// ColorizeWidget applies the theme set on the App object to any widget given.
func ColorizeWidget(widget termui.Drawable, theme string) {
	color := utils.Cond(
		theme == "dark",
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
