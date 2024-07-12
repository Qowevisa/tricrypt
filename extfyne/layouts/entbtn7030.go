package layouts

import "fyne.io/fyne/v2"

// Implements: fyne.Layout
type EntryBtn7030 struct{}

// Implements: fyne.Layout 1
func (p *EntryBtn7030) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minWidth := float32(0)
	minHeight := float32(0)
	for _, obj := range objects {
		minSize := obj.MinSize()
		minWidth += minSize.Width
		if minSize.Height > minHeight {
			minHeight = minSize.Height
		}
	}
	return fyne.NewSize(minWidth, minHeight)
}

// Implements: fyne.Layout 2
func (p *EntryBtn7030) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) != 2 {
		return
	}

	entryWidth := size.Width * 0.7
	buttonWidth := size.Width * 0.3

	objects[0].Resize(fyne.NewSize(entryWidth, size.Height))
	objects[1].Resize(fyne.NewSize(buttonWidth, size.Height))

	objects[0].Move(fyne.NewPos(0, 0))
	objects[1].Move(fyne.NewPos(entryWidth, 0))
}

func NewEntryBtn7030() fyne.Layout {
	return &EntryBtn7030{}
}
