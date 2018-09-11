package ui

import (
    "github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
    "strings"
    "runtime"
    "time"
)

func Tprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func Fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

const (
    CONSOLE_HELP = "help"
    CONSOLE_HELP_DISCRIPTION = "Описание комманд"
    CONSOLE_MESSAGE = "console_message"
    CONSOLE_MESSAGE_DISCRIPTION = ""
    CONSOLE_LOG = "log"
    CONSOLE_LOG_DISCRIPTION = "Онлайн поток работы программы"
    CONSOLE_CURRENT_CMD = "console_current_cmd"
    CONSOLE_CURRENT_CMD_DISCRIPTION = ""

    MSG_PLACE_NAME = "{name}"
    MSG_CONSOLE_NOT_ACTIVETED = `Нет активной консоли"`
    MSG_CONSOLE_NOT_EXIST = `Консоль "` + MSG_PLACE_NAME + `" не существует`
)

type Console struct {
    BodyX, BodyY *int
    Name string
    Title string
    IsActive bool
    Values map[string]string
    Controls map[string]interface{}
    Action func(*Console, []interface{})
    Draw func(*Console, int, int)
    Clear func(*Console, int, int)
}

func NewConsole(name, title string, isActive bool, draw func(*Console, int, int), clear func(*Console, int, int), action func(*Console, []interface{})) *Console {
    console := Console {Name: name, Title: title, IsActive: isActive, Draw: draw, Clear: clear, Action: action }
    console.Values = make(map[string]string)
    console.Controls = make(map[string]interface{})
    return &console
}

func (console *Console) Redraw() {
    if console.IsActive {
        console.ClearDraw()
        // const coldef = termbox.ColorDefault
        // termbox.Clear(coldef, coldef)
        if console.Draw != nil {
            console.Draw(console, *console.BodyX, *console.BodyY)
        }
        time.Sleep(time.Duration(5)*time.Millisecond)
        termbox.Flush()
        time.Sleep(time.Duration(5)*time.Millisecond)
    }
}

func (console *Console) Execute(data []interface{}) {
    if console.IsActive {
        if console.Action != nil {
            console.Action(console, data)
        }
    }
}

func (console *Console) ClearDraw() {
    if console.Clear != nil {
        console.Clear(console, *console.BodyX, *console.BodyY)
    }
}

type Text struct {
    Value interface{}
    FG termbox.Attribute
    BG termbox.Attribute
}

type UI struct {
    chClose chan bool
    ChLog, ChErr, ChMsg chan interface{}
    MainConsole *Console
    BodyX, BodyY int
    ChClosed chan bool
    Values map[string]string
    Controls map[string]interface{}
    Consoles map[string]*Console
    Events map[string]func(termbox.Event)
    drawBaseConsole func(*UI)
}

func NewUI(console func(*UI)) *UI {
    ui := UI{}
    ui.drawBaseConsole = console
    ui.chClose = make(chan bool)
    ui.ChLog = make(chan interface{})
    ui.ChErr = make(chan interface{})
    ui.ChMsg = make(chan interface{})
    ui.ChClosed = make(chan bool)
    ui.Values = make(map[string]string)
    ui.Controls = make(map[string]interface{})
    ui.Events = make(map[string]func(termbox.Event))
    ui.Consoles = make(map[string]*Console)
    ui.AddConsole(NewConsole(CONSOLE_HELP, CONSOLE_HELP_DISCRIPTION, false, nil, nil, nil))
    ui.AddConsole(NewConsole(CONSOLE_MESSAGE, "", true, nil, nil, nil))
    ui.AddConsole(NewConsole(CONSOLE_LOG, CONSOLE_LOG_DISCRIPTION, false, nil, nil, nil))
    ui.AddConsole(NewConsole(CONSOLE_CURRENT_CMD, CONSOLE_CURRENT_CMD_DISCRIPTION, true, nil, nil, nil))
    // Инициализируем пакет по отрисовки консольного интерфейса
    err := termbox.Init()
    if err != nil {
        panic(err)
    }
    termbox.SetInputMode(termbox.InputEsc)
    // Проводим первую отрисовку интерфейса
    ui.RedrawAll()
    // Запускаем горутину на чтение и обработку клавиш
    go func() {
        for {
            select {
            case <-ui.chClose:
                termbox.Close()
                close(ui.ChClosed)
                return
            default:
                ui.readingKeybord()
                runtime.Gosched()
            }
        }
    }()
    // Запускаем горутину на чтение каналов ошибок и сообщений
    go func() {
		for {
			select {
			case <-ui.chClose:
				return
			case err := <-ui.ChErr:
				ui.Error(err)
            case log := <-ui.ChLog:
                ui.Log(log)
            case log := <-ui.ChMsg:
				ui.Message(log)
			default:
				runtime.Gosched()
			}
		}
	}()
    return &ui
}

func (ui *UI) RedrawAll() {
    // Рисуем основную консоль
    ui.Redraw()
    for _, console := range ui.Consoles {
        if console.IsActive {
            console.Redraw()
        }
    }
}

func (ui *UI) readingKeybord() {
	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventKey:
		for _, event := range ui.Events {
            event(ev)
        }
	case termbox.EventError:
		panic(ev.Err)
	}
}

func (ui *UI) Redraw() {
    const coldef = termbox.ColorDefault
    termbox.Clear(coldef, coldef)
    ui.drawBaseConsole(ui)
    termbox.Flush()
}

func (ui *UI) AddConsole(console *Console) {
    if console == nil {
        return
    }
    console.BodyX = &ui.BodyX
    console.BodyY = &ui.BodyY
    ui.Consoles[console.Name] = console
}

func (ui *UI) SetMainConsole(name string) {
    console := ui.Consoles[name]
    if console == nil {
        ui.Error(strings.Replace(MSG_CONSOLE_NOT_EXIST, MSG_PLACE_NAME, name, 1))
        return
    }
    // Обнуляем все значения контрола
    // for control, _ := range console.Values {
    //     console.Values[control] = ""
    // }
    if ui.MainConsole != nil {
        ui.MainConsole.IsActive = false
    }
    console.IsActive = true
    ui.MainConsole = console
}

func (ui *UI) SetControlValue(control, value string) {
    if ui.MainConsole == nil {
        ui.Error(MSG_CONSOLE_NOT_ACTIVETED)
    } else {
        ui.MainConsole.Values[control] = value
    }
}

func (ui *UI) Close() {
    close(ui.chClose)
}

func (ui *UI) Error(msg interface{}) {
    if msg != nil {
        ui.DrawConsole(CONSOLE_MESSAGE, []interface{} { &Text { Value: msg, FG: termbox.ColorRed } })
        ui.DrawConsole(CONSOLE_LOG, []interface{} { &Text { Value: msg, FG: termbox.ColorRed } })
    }
}

func (ui *UI) Message(msg interface{}) {
    if msg != nil {
        ui.DrawConsole(CONSOLE_MESSAGE, []interface{} { &Text { Value: msg, FG: termbox.ColorBlue } })
        ui.DrawConsole(CONSOLE_LOG, []interface{} { &Text { Value: msg, FG: termbox.ColorBlue } })
    }
}

func (ui *UI) Log(msg interface{}) {
    if msg != nil {
        ui.DrawConsole(CONSOLE_LOG, []interface{} { &Text { Value: msg, FG: termbox.ColorBlue } })
    }
}

func (ui *UI) DrawConsole(name string, data []interface{}) {
    console := ui.Consoles[name]
    if console.IsActive {
        console.Execute(data)
        console.Redraw()
    }
}
