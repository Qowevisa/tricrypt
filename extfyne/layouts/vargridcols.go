package layouts

import "fyne.io/fyne/v2"

// Implements: fyne.Layout
type VarGridCols struct {
	cols int
	cfg  []int
}

func NewVariableGridWithColumns(cols int, cfg []int) *VarGridCols {
	return &VarGridCols{
		cols: cols,
		cfg:  cfg,
	}
}

// Implements: fyne.Layout : 1
func (v *VarGridCols) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	i := 0
	for _, o := range objects {
		if !o.Visible() {
			continue
		}
		childSize := o.MinSize()

		// NOTE: can be out-of-bounds error if not checked at creation
		takenCols := v.cfg[i]
		if takenCols > 1 {
			w += childSize.Width * float32(takenCols)
		} else {
			w += childSize.Width
		}
		if h+childSize.Height > h {
			h += childSize.Height
		}
		i++
	}
	return fyne.NewSize(w, h)
}

// Implements: fyne.Layout : 2
func (v *VarGridCols) Layout(objects []fyne.CanvasObject, contSize fyne.Size) {
	// pos := fyne.NewPos(0, contSize.Height-v.MinSize(objects).Height)
	pos := fyne.NewPos(0, 0)
	cellWidthSingular := float64(contSize.Width) / float64(v.cols)
	cellHeight := contSize.Height

	for i, o := range objects {
		numCols := v.cfg[i]
		w := cellWidthSingular * float64(numCols)
		o.Move(pos)
		o.Resize(fyne.NewSize(float32(w), cellHeight))

		pos = pos.Add(fyne.NewPos(float32(w), 0))
	}
}
