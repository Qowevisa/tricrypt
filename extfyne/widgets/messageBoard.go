package widgets

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type MBoardMessage struct {
	LeftAlign bool
	Data      *widget.Label
}

type MessageBoard struct {
	widget.BaseWidget
	Messages []MBoardMessage
}

func NewEmptyMessageBoard() *MessageBoard {
	board := &MessageBoard{
		Messages: make([]MBoardMessage, 0),
	}
	board.ExtendBaseWidget(board)
	return board
}

func NewMessageBoard(msgs []MBoardMessage) *MessageBoard {
	board := &MessageBoard{
		Messages: msgs,
	}
	board.ExtendBaseWidget(board)
	return board
}

func (m *MessageBoard) Add(msg MBoardMessage) {
	m.Messages = append(m.Messages, msg)
	m.Refresh()
}

func (m *MessageBoard) CreateRenderer() fyne.WidgetRenderer {
	content := container.NewVBox()
	for _, msg := range m.Messages {
		hbox := container.NewHBox()
		background := canvas.NewRectangle(&color.RGBA{R: 0, G: 0, B: 255, A: 128})
		background.SetMinSize(msg.Data.MinSize())

		if msg.LeftAlign {
			hbox.Add(container.NewStack(background, msg.Data))
		} else {
			hbox.Add(layout.NewSpacer())
			hbox.Add(container.NewStack(background, msg.Data))
		}
		content.Add(hbox)
	}

	scrollContainer := container.NewScroll(content)

	return &messageBoardRenderer{
		board:   m,
		content: scrollContainer,
	}
}

type messageBoardRenderer struct {
	board   *MessageBoard
	content *container.Scroll
}

func (r *messageBoardRenderer) Layout(size fyne.Size) {
	r.content.Resize(size)
}

func (r *messageBoardRenderer) MinSize() fyne.Size {
	return r.content.MinSize()
}

func (r *messageBoardRenderer) Refresh() {
	r.content.Content.(*fyne.Container).Objects = nil
	for _, msg := range r.board.Messages {
		hbox := container.NewHBox()
		if msg.LeftAlign {
			background := canvas.NewRectangle(&color.RGBA{R: 0, G: 0, B: 255, A: 128})
			background.SetMinSize(msg.Data.MinSize())
			hbox.Add(container.NewStack(background, msg.Data))
		} else {
			background := canvas.NewRectangle(&color.RGBA{R: 0, G: 255, B: 0, A: 128})
			background.SetMinSize(msg.Data.MinSize())
			hbox.Add(layout.NewSpacer())
			hbox.Add(container.NewStack(background, msg.Data))
		}
		r.content.Content.(*fyne.Container).Add(hbox)
	}
	r.content.Refresh()
}

func (r *messageBoardRenderer) Destroy() {
	// Clean up resources if needed
}

func (r *messageBoardRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.content}
}
