package layouts

import "fyne.io/fyne/v2"

// Implements: fyne.Layout
type VarGridRows struct {
	rows int
	cfg  []int
}

func NewVariableGridWithRows(cols int, cfg []int) *VarGridRows {
	return &VarGridRows{
		rows: cols,
		cfg:  cfg,
	}
}

// Implements: fyne.Layout : 1
func (v *VarGridRows) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	for i, o := range objects {
		childSize := o.MinSize()

		// NOTE: can be out-of-bounds error if not checked at creation
		takenRows := v.cfg[i]
		if takenRows > 1 {
			h += childSize.Height * float32(takenRows)
		} else {
			h += childSize.Height
		}
		if h+childSize.Height > h {
			w += childSize.Width
		}
	}
	return fyne.NewSize(w, h)
}

// Implements: fyne.Layout : 2
func (v *VarGridRows) Layout(objects []fyne.CanvasObject, contSize fyne.Size) {
	// pos := fyne.NewPos(0, contSize.Height-v.MinSize(objects).Height)
	pos := fyne.NewPos(0, 0)
	cellHeightSingular := float64(contSize.Height) / float64(v.rows)
	cellWidth := contSize.Height

	for i, o := range objects {
		numRows := v.cfg[i]
		h := cellHeightSingular * float64(numRows)
		o.Move(pos)
		o.Resize(fyne.NewSize(cellWidth, float32(h)))

		pos = pos.Add(fyne.NewPos(0, float32(h)))
	}
}
