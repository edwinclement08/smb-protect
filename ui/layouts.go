package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

type ExpandedVBoxLayout struct {
}

func (d *ExpandedVBoxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	for _, o := range objects {
		childSize := o.MinSize()

		if w < childSize.Width {
			w = childSize.Width
		}
		h += childSize.Height
	}
	return fyne.NewSize(w, h)
}

func (d *ExpandedVBoxLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, 0)
	for _, o := range objects {
		size := o.MinSize()
		newSize := fyne.NewSize(containerSize.Width, size.Height)
		o.Resize(newSize)
		o.Move(pos)
		pos = pos.Add(fyne.NewPos(0, size.Height))
	}
}

func NewExpandedVBoxLayout(objects ...fyne.CanvasObject) fyne.CanvasObject {
	return container.New(&ExpandedVBoxLayout{}, objects...)
}

type VSplitLayout struct {
}

func (d *VSplitLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	padding := theme.Padding()
	w += padding * 3
	for i, o := range objects {
		if i > 1 {
			break
		}
		childSize := o.MinSize()

		w += childSize.Width
		if h < childSize.Height {
			h = childSize.Height
		}
	}
	return fyne.NewSize(w, h)
}

func (d *VSplitLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	halfwayPoint := containerSize.Width / 2
	minHeight := d.MinSize(objects).Height
	padding := theme.Padding()

	newSize := fyne.NewSize(halfwayPoint-padding*(1+0.5), minHeight)

	pos0 := fyne.NewPos(padding, 0)
	pos1 := fyne.NewPos(halfwayPoint+padding/2, 0)
	objects[0].Resize(newSize)
	objects[1].Resize(newSize)

	objects[0].Move(pos0)
	objects[1].Move(pos1)
}

func NewVSplitLayout(object0, object1 fyne.CanvasObject) fyne.CanvasObject {
	return container.New(&VSplitLayout{}, object0, object1)
}
