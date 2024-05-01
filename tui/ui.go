package tui

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"unicode"

	"git.qowevisa.me/Qowevisa/gotell/debug"
	"git.qowevisa.me/Qowevisa/gotell/errors"
	"golang.org/x/term"
)

func (t *TUI) init() error {
	var err error
	t.inputChannel = make(chan rune, 32)
	t.printRunes = make(chan rune, 32)
	t.widgets = make([]*widget, 8)
	t.errorsChannel = make(chan error, 4)
	t.mySignals = make(chan mySignal, 1)
	t.osSignals = make(chan os.Signal, 1)
	t.writer = bufio.NewWriter(os.Stdout)
	t.storage = make(map[string]string)
	t.readInputState = make(chan bool, 1)
	t.readEnterState = make(chan bool, 1)
	t.stateChannel = make(chan string, 1)
	t.messageChannel = make(chan []byte, 8)
	signal.Notify(t.osSignals, syscall.SIGWINCH)
	err = t.setSizes()
	if err != nil {
		return errors.WrapErr("t.getSizes", err)
	}
	err = t.setTermToRaw()
	if err != nil {
		return errors.WrapErr("t.setTermToRaw", err)
	}
	err = t.setRoutines()
	if err != nil {
		return errors.WrapErr("t.setRoutines", err)
	}
	err = t.readRoutines()
	if err != nil {
		return errors.WrapErr("t.readRoutines", err)
	}
	err = t.launchAllChannels()
	if err != nil {
		return errors.WrapErr("t.launchAllChannels", err)
	}
	return nil
}

func (t *TUI) exit() {
	if t.oldState != nil {
		term.Restore(int(os.Stdin.Fd()), t.oldState)
	}
}

func (t *TUI) Run() error {
	defer t.exit()
	var err error
	err = t.init()
	if err != nil {
		return errors.WrapErr("t.init", err)
	}
	//
	if t.mySignals == nil {
		return errors.WrapErr("t.mySignals", errors.NOT_INIT)
	}
	err = t.Draw()
	if err != nil {
		return errors.WrapErr("t.Draw", err)
	}
	for mySignal := range t.mySignals {
		log.Printf("Receive signal: %#v\n", mySignal)
		if mySignal.Type == MY_SIGNAL_EXIT {
			t.errorsChannel <- t.clearScreen()
			t.errorsChannel <- t.moveCursor(0, 0)
			break
		}
		switch mySignal.Type {
		case MY_SIGNAL_CONNECT:
			var host []rune
			var port []rune
			hostHandler, hostData := AddToStorageEasy("host", &host)
			portHandler, portData := AddToStorageEasy("port", &port)
			err := t.addWidget(widgetConfig{
				Input:           &host,
				Title:           "Host",
				MinWidth:        16,
				HasBorder:       true,
				WidgetPosConfig: widgetPosGeneralCenter,
				CursorPosConfig: cursorPosGeneralCenter,
				DataHandler:     hostHandler,
				Data:            hostData,
				Finale:          nil,
				Next: &widgetConfig{
					Input:           &port,
					Title:           "Port",
					MinWidth:        8,
					HasBorder:       true,
					WidgetPosConfig: widgetPosGeneralCenter,
					CursorPosConfig: cursorPosGeneralCenter,
					DataHandler:     portHandler,
					Data:            portData,
					Next:            nil,
					Finale:          FE_ConnectTLS,
					FinaleData:      dataT{},
				},
			})
			if err != nil {
				t.errorsChannel <- errors.WrapErr("t.addWidget", err)
			}
			err = t.drawSelectedWidget()
			if err != nil {
				t.errorsChannel <- errors.WrapErr("t.drawSelectedWidget", err)
			}

		case MY_SIGNAL_MESSAGE:
			if t.isConnected {
				t.SendMessageToServer("Message", 20)
			}

		case MY_SIGNAL_CLOSE:
			if t.isConnected {
				CloseConnection(t.tlsConnCloseData.wg, t.tlsConnCloseData.cancel)
				t.isConnected = false
				t.errorsChannel <- t.tlsConnection.Close()
				t.stateChannel <- "Disconnected"
			}
		default:
		}
	}
	//
	return nil
}

func (t *TUI) createNotification(text, title string) {
	notifier, err := createNotification(text, title)
	t.selectedNotifier = &notifier
	t.errorsChannel <- err
	t.errorsChannel <- t.write(notifier.Buf)
	t.readEnterState <- true
}

func (t *TUI) setRoutines() error {
	if t.inputChannel == nil {
		return errors.WrapErr("t.input", errors.NOT_INIT)
	}
	if t.oldState == nil {
		return errors.WrapErr("t.oldState", errors.NOT_INIT)
	}
	if t.readInputState == nil {
		return errors.WrapErr("t.readInputState", errors.NOT_INIT)
	}
	if t.readEnterState == nil {
		return errors.WrapErr("t.readEnterState", errors.NOT_INIT)
	}
	if t.stateChannel == nil {
		return errors.WrapErr("t.stateChannel", errors.NOT_INIT)
	}
	if t.messageChannel == nil {
		return errors.WrapErr("t.messageChannel", errors.NOT_INIT)
	}
	var readInputMu sync.Mutex
	var readEnterdMu sync.Mutex
	readInput := false
	readCommand := false
	readDebug := false
	readEnter := false
	go func() {
		for newState := range t.readInputState {
			readInputMu.Lock()
			readInput = newState
			readInputMu.Unlock()
		}
	}()
	go func() {
		for newState := range t.readEnterState {
			readEnterdMu.Lock()
			readEnter = newState
			readEnterdMu.Unlock()
		}
	}()
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			r, _, err := reader.ReadRune()
			if err != nil {
				panic(err)
			}
			log.Printf("Read %#v rune\n", r)
			if readEnter {
				if r == 13 {
					readEnter = false
					if t.selectedNotifier != nil {
						t.errorsChannel <- t.write(t.selectedNotifier.Clear())
					}
					continue
				}
			}
			if readDebug {
				log.Printf("Reading debug")
				switch r {
				case 'm':
					log.Printf("get m for debug")
					debug.LogMemUsage()
				}
				readDebug = false
				continue
			}
			if readCommand {
				switch r {
				case 'q':
					t.mySignals <- mySignal{
						Type: MY_SIGNAL_EXIT,
					}
				case 'c':
					if t.isConnected {
						t.mySignals <- mySignal{
							Type: MY_SIGNAL_CLOSE,
						}
					} else {
						t.mySignals <- mySignal{
							Type: MY_SIGNAL_CONNECT,
						}
						readInputMu.Lock()
						readInput = true
						readInputMu.Unlock()
					}
				case 'm':
					t.mySignals <- mySignal{
						Type: MY_SIGNAL_MESSAGE,
					}
					readInputMu.Lock()
					readInput = true
					readInputMu.Unlock()
				}

				readCommand = false
				continue
			}
			//
			if unicode.IsControl(r) {
				switch r {
				case CTRL_A:
					log.Printf("CTRL_A received!\n")
					readCommand = true
				case CTRL_D:
					log.Printf("CTRL_D received!\n")
					readDebug = true
				}
			} else {
				if readInput {
					log.Printf("Send %c | %d to t.input", r, r)
					t.inputChannel <- r
				}
			}
			readInputMu.Lock()
			if readInput {
				switch r {
				case 13:
					t.inputChannel <- r
				case 127:
					t.inputChannel <- r
				}
			}
			readInputMu.Unlock()
		}
	}()
	//
	return nil
}

func (t *TUI) readRoutines() error {
	return nil
}

func (t *TUI) Draw() error {
	err := t.clearScreen()
	if err != nil {
		return errors.WrapErr("t.clearScreen", err)
	}
	err = t.drawFooter()
	if err != nil {
		return errors.WrapErr("t.drawFooter", err)
	}
	return nil
}

func (t *TUI) drawFooter() error {
	t.writeMu.Lock()
	defer t.writeMu.Unlock()
	err := t.moveCursor(t.height, 0)
	if err != nil {
		return errors.WrapErr("t.moveCursor", err)
	}
	err = t.write(footerStart)
	if err != nil {
		return errors.WrapErr("t.write", err)
	}
	return nil
}

func (t *TUI) write(s string) error {
	_, err := t.writer.WriteString(s)
	if err != nil {
		return errors.WrapErr("t.writer.WriteString", err)
	}
	err = t.writer.Flush()
	if err != nil {
		return errors.WrapErr("t.writer.Flush", err)
	}
	t.cursorPosCol += len(s)
	if t.cursorPosCol > t.width {
		t.cursorPosCol %= t.width
		t.cursorPosRow++
	}
	return nil
}

func (t *TUI) writeRune(r rune) error {
	_, err := t.writer.WriteRune(r)
	if err != nil {
		return errors.WrapErr("t.writer.WriteRune", err)
	}
	err = t.writer.Flush()
	if err != nil {
		return errors.WrapErr("t.writer.Flush", err)
	}
	t.cursorPosCol++
	return nil
}

func (t *TUI) getCursorPos() (int, int) {
	return t.cursorPosRow, t.cursorPosCol
}

func (t *TUI) _clearLine() {
	t.errorsChannel <- t.write("\033[0K")
}

func (t *TUI) moveCursor(row, col int) error {
	t.sizeMutex.Lock()
	defer t.sizeMutex.Unlock()
	if row > t.height {
		return errors.WrapErr(fmt.Sprintf("row: %d; height: %d", row, t.height), errors.OUT_OF_BOUND)
	}
	if col > t.width {
		return errors.WrapErr(fmt.Sprintf("col: %d; width: %d", col, t.width), errors.OUT_OF_BOUND)
	}
	log.Printf("t.cursorPosRow: %d ; t.cursorPosCol: %d\n", t.cursorPosRow, t.cursorPosCol)
	log.Printf("trying to move to row: %d ; col %d\n", row, col)
	_, err := t.writer.WriteString(fmt.Sprintf("\033[%d;%dH", row, col))
	if err != nil {
		return errors.WrapErr("t.writer.WriteString", err)
	}
	err = t.writer.Flush()
	if err != nil {
		return errors.WrapErr("t.writer.Flush", err)
	}
	t.cursorPosCol = col
	log.Printf("t.cursorPosCol now is = %d\n", col)
	t.cursorPosRow = row
	log.Printf("t.cursorPosRow now is = %d\n", row)
	return nil
}

func (t *TUI) clearScreen() error {
	_, err := t.writer.WriteString("\033[2J")
	if err != nil {
		return errors.WrapErr("t.writer.WriteString", err)
	}
	err = t.writer.Flush()
	if err != nil {
		return errors.WrapErr("t.writer.Flush", err)
	}
	return nil
}

func (t *TUI) drawSelectedWidget() error {
	wDraw, err := t.selectedWidget.Draw()
	if err != nil {
		return errors.WrapErr("t.selectedWidget.Draw", err)
	}
	t.writeMu.Lock()
	err = t.write(wDraw.Buf)
	t.writeMu.Unlock()
	if err != nil {
		return errors.WrapErr("t.write", err)
	}
	t.cursorPosRow = wDraw.Row
	t.cursorPosCol = wDraw.Col
	return nil
}

// Creating and adding widget from config.
// Also sets created widget as selectedWidget
func (t *TUI) addWidget(config widgetConfig) error {
	widget := &widget{}
	err := widget.init(config)
	if err != nil {
		return errors.WrapErr("widget.init", err)
	}
	t.widgetsMutext.Lock()
	defer t.widgetsMutext.Unlock()
	t.widgets = append(t.widgets, widget)
	t.selectedWidget = widget
	return nil
}

func (t *TUI) redraw() error {
	var err error
	err = t.setSizes()
	if err != nil {
		return errors.WrapErr("t.getSizes", err)
	}
	err = t.redrawWidgets()
	if err != nil {
		return errors.WrapErr("t.redrawWidgets", err)
	}
	return nil
}

func (t *TUI) setSizes() error {
	w, h, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return errors.WrapErr("term.GetSize", err)
	}
	t.sizeMutex.Lock()
	t.width = w
	t.height = h
	t.sizeMutex.Unlock()
	return nil
}

func (t *TUI) getSizes() (int, int) {
	t.sizeMutex.Lock()
	defer t.sizeMutex.Unlock()
	return t.width, t.height
}

func (t *TUI) redrawWidgets() error {
	if t.widgets == nil {
		return errors.WrapErr("t.widgets", errors.NOT_INIT)
	}
	// TODO
	return nil
}

func (t *TUI) setTermToRaw() error {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return errors.WrapErr("term.MakeRaw", err)
	}
	t.oldState = oldState
	return nil
}
