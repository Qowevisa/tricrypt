package tui

import (
	"bufio"
	"crypto/tls"
	"os"
	"sync"

	"golang.org/x/term"
)

type mySignal struct {
	Type int
}

var UI TUI

type dataT struct {
	raw  string
	rawP *[]rune
	op1  string
	op2  string
	ops  []string
}

type dataProcessHandler func(t *TUI, data dataT) error

type tuiPointPair struct {
	startCol int
	startRow int
	endCol   int
	endRow   int
}

type widget struct {
	row           int
	col           int
	MinWidth      int
	MinHeight     int
	Width         int
	Height        int
	Title         string
	Input         *[]rune
	Handler       dataProcessHandler
	Data          dataT
	Next          *widgetConfig
	Finale        dataProcessHandler
	FinaleData    dataT
	percentPair   tuiPointPair
	snappedPair   tuiPointPair
	startupConfig widgetConfig
}

type notifier struct {
	Row    int
	Col    int
	Width  int
	Height int
	Buf    string
}

type widgetDraw struct {
	Buf string
	Row int
	Col int
}

type cursorPosConfigValue int
type widgetPosConfigValue int

type widgetConfig struct {
	MinWidth        int
	MinHeight       int
	Input           *[]rune
	CursorPosConfig cursorPosConfigValue
	Title           string
	WidgetPosConfig widgetPosConfigValue
	HasBorder       bool
	DataHandler     dataProcessHandler
	Data            dataT
	Next            *widgetConfig
	Finale          dataProcessHandler
	FinaleData      dataT
}

type TUI struct {
	width            int
	height           int
	cursorPosRow     int
	cursorPosCol     int
	writeMu          sync.Mutex
	sizeMutex        sync.Mutex
	oldState         *term.State
	input            chan (rune)
	printRunes       chan (rune)
	mySignals        chan (mySignal)
	osSignals        chan (os.Signal)
	errors           chan (error)
	readInputState   chan (bool)
	readEnterState   chan (bool)
	stateChannel     chan (string)
	messageChannel   chan ([]byte)
	widgets          []*widget
	widgetsMutext    sync.Mutex
	writer           *bufio.Writer
	isConnected      bool
	selectedWidget   *widget
	selectedNotifier *notifier
	storage          map[string]string
	tlsConnection    *tls.Conn
}
