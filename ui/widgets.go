package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func NewBorderStyle(objects ...fyne.CanvasObject) fyne.CanvasObject {
	return container.NewPadded(
		container.NewMax(
			&canvas.Rectangle{StrokeColor: color.Black, StrokeWidth: 1},
			container.NewPadded(objects...)))
}

func NewCustomBoldLabel(text string, color color.Color, textSize float32) fyne.CanvasObject {
	return container.NewPadded(
		&canvas.Text{
			Text:      text,
			TextSize:  textSize,
			Color:     color,
			TextStyle: fyne.TextStyle{Bold: true},
		},
	)
}
