package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"ribbirc/client"
	"ribbirc/utils"
	"unicode"
)

type Application struct {
	screen   tcell.Screen
	width    int
	height   int
	listener chan int

	server     *client.Server
	channelTab string

	inputActive bool
	inputCursor int
	inputText   []rune
}

func New() (*Application, error) {
	// @todo: temporary
	listener := make(chan int)
	server := client.New(listener, "irc.libera.chat", 6697, "ribbirc")
	err := server.Connect()
	if err != nil {
		return nil, err
	}

	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	return &Application{
		screen:   screen,
		listener: listener,
		server:   server,
	}, nil
}

func (a *Application) Run() error {
	err := a.screen.Init()
	if err != nil {
		return err
	}

	go a.listenToChannel()

	for {
		ev := a.screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			a.width, a.height = ev.Size()
			a.screen.Sync()
		case *tcell.EventKey:
			a.handleKeyEvent(ev)
		}

		a.draw()
	}
}

func (a *Application) Stop() {
	a.screen.Fini()
}

func (a *Application) listenToChannel() {
	for {
		<-a.listener
		a.draw()
	}
}

func (a *Application) handleKeyEvent(ev *tcell.EventKey) {
	if ev.Modifiers() == tcell.ModAlt {
		indexes := map[rune]int{'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9}
		tab := indexes[ev.Rune()]
		if tab == 0 {
			a.channelTab = ""
		} else {
			a.channelTab = a.server.ChannelNames()[tab-1]
		}
		return
		// @todo: more checks
	}

	switch ev.Key() {
	case tcell.KeyCtrlC:
		a.Stop()
		// @todo: end properly
	case tcell.KeyEnter:
		if len(a.inputText) == 0 {
			a.inputActive = !a.inputActive
		} else {
			a.server.SendMessage(utils.UnmarshalMessage(string(a.inputText)))
			a.inputText = make([]rune, 0)
			a.inputCursor = 0
		}
	default:
	}

	if a.inputActive {
		switch ev.Key() {
		case tcell.KeyBackspace:
			if a.inputCursor > 0 {
				a.inputText = append(a.inputText[:a.inputCursor-1], a.inputText[a.inputCursor:]...)
				a.inputCursor--
			}
		case tcell.KeyLeft:
			a.inputCursor--
			if a.inputCursor < 0 {
				a.inputCursor = 0
			}
		case tcell.KeyRight:
			a.inputCursor++
			if a.inputCursor > len(a.inputText) {
				a.inputCursor = len(a.inputText)
			}
		default:
			if unicode.IsPrint(ev.Rune()) {
				a.inputText = append(a.inputText, ' ')
				copy(a.inputText[a.inputCursor+1:], a.inputText[a.inputCursor:])
				a.inputText[a.inputCursor] = ev.Rune()
				a.inputCursor++
			}
		}
	}
}

func (a *Application) draw() {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	a.screen.SetStyle(defStyle)

	a.screen.Clear()

	a.drawTopBar()
	a.drawBottomBar()
	a.drawLogs()
	a.drawInput()

	a.screen.Show()
}

func (a *Application) drawTopBar() {
	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorReset)

	for col := range a.width {
		a.screen.SetContent(col, 0, ' ', nil, style)
		col++
	}

	text := fmt.Sprintf("RibbIRC v0.1.0 [%dx%d]", a.width, a.height)
	a.drawString(0, 0, text, style)
}

func (a *Application) drawBottomBar() {
	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorReset)

	for col := range a.width {
		a.screen.SetContent(col, a.height-2, ' ', nil, style)
		col++
	}

	text := "[0. Status]"
	for i, channel := range a.server.ChannelNames() {
		text += fmt.Sprintf(" [%d. %s]", i+1, channel)
	}

	a.drawString(0, a.height-2, text, style)
}

func (a *Application) drawLogs() {
	style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	var logs []utils.Log
	if a.channelTab == "" {
		logs = a.server.GetLogger().GetNLogs(a.height-3, 0)
	} else {
		logs = a.server.ChannelLogger(a.channelTab).GetNLogs(a.height-3, 0)
	}

	for row, log := range logs {
		text := fmt.Sprintf("%s > %s", log.Source, log.Text)
		a.drawString(0, row+1, text, style)
	}
}

func (a *Application) drawInput() {
	style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	a.drawString(0, a.height-1, string(a.inputText), style)

	if a.inputActive {
		a.screen.ShowCursor(a.inputCursor, a.height-1)
	} else {
		a.screen.HideCursor()
	}
}

func (a *Application) drawString(x int, y int, text string, style tcell.Style) {
	row := y
	col := x
	for _, r := range []rune(text) {
		a.screen.SetContent(col, row, r, nil, style)
		col++
	}
}
