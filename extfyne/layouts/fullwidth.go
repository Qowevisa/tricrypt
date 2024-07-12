package layouts

import (
	"fyne.io/fyne/v2"
)

// Implements: fyne.Layout
type FullScale struct {
	size fyne.Size
}

// Implements: fyne.Layout 1
func (p *FullScale) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return p.size
}

// Implements: fyne.Layout 2
func (p *FullScale) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) != 1 {
		return
	}
	if p.size.Width < size.Width {
		p.size.Width = size.Width
	}
	if p.size.Height < size.Height {
		p.size.Height = size.Height
	}

	objects[0].Resize(fyne.NewSize(p.size.Width, p.size.Height))
	objects[0].Move(fyne.NewPos(0, 0))
}

func NewFullWidth() fyne.Layout {
	return &FullScale{}
}

func NewFullWidthWithSize(s fyne.Size) fyne.Layout {
	return &FullScale{
		size: s,
	}
}
