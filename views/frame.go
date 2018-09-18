package views

import (
    packUI "../ui"
    "github.com/nsf/termbox-go"
)
const coldef = termbox.ColorDefault

var frame *packUI.UI
var settings *packUI.Settings
//
func Frame(_settings *packUI.Settings) *packUI.UI {
    settings = _settings
    frame = packUI.NewUI(CONSOLE_MESSAGE, CONSOLE_LOG, func (ui *packUI.UI) {
        w, h := termbox.Size()
        ui.BodyY = h - 8
        ui.BodyX = 30
        //┌ ┐ └ ┘ ─ │ ╎ ╌ ├ ┤
        lableProduct := "Intellect Trade"
        lableVer := "v0.1"
        lableIndicators := "Индикаторы"
        lableTitle := ""
        lableCommand := "Команда"
        lableQuit := "ESC выход"
        lableHelp := `"help" команды`
        lableConsole := `"console" вызовы окон`
        lableSettings := `"settings" настройки`
        lableF2 := `"F2 логи"`
        if ui.MainConsole != nil {
            lableTitle = ui.MainConsole.Title
        }
        // Свисающая шапка (head)
        termbox.SetCell(ui.BodyX + 3, 0, '│', coldef, coldef)
        termbox.SetCell(ui.BodyX + 3, 1, '└', coldef, coldef)
        packUI.Fill(ui.BodyX + 4, 1, w - ui.BodyX - 6, 1, termbox.Cell{Ch: '─'})
        termbox.SetCell(w - 3, 0, '│', coldef, coldef)
        termbox.SetCell(w - 3, 1, '┘', coldef, coldef)
        packUI.Tprint(((w + ui.BodyX) / 2) - (len([]rune(lableTitle)) / 2), 0, coldef, coldef, lableTitle)
        // Левая панель
        termbox.SetCell(2, 0, '│', coldef, coldef)
        termbox.SetCell(2, 1, '└', coldef, coldef)
        packUI.Fill(3, 1, ui.BodyX - 3, 1, termbox.Cell{Ch: '─'})
        termbox.SetCell(ui.BodyX, 0, '│', coldef, coldef)
        termbox.SetCell(ui.BodyX, 1, '┘', coldef, coldef)
        packUI.Tprint(((ui.BodyX + 3) / 2) - (len([]rune(lableIndicators)) / 2), 0, coldef, coldef, lableIndicators)
        packUI.Fill(2, 2, 1, ui.BodyY, termbox.Cell{Ch: '╎'})
        packUI.Fill(ui.BodyX, 2, 1, ui.BodyY, termbox.Cell{Ch: '╎'})
        packUI.Fill(2, ui.BodyY + 4, 1, h - ui.BodyY - 4, termbox.Cell{Ch: '╎'})
        packUI.Fill(ui.BodyX, ui.BodyY + 4, 1, h - ui.BodyY - 4, termbox.Cell{Ch: '╎'})

        // Footer левой панели
        packUI.Tprint(((ui.BodyX + 3) / 2) - (len(lableProduct) / 2), h - 3, termbox.ColorCyan, coldef, lableProduct)
        packUI.Tprint(((ui.BodyX + 3) / 2) - (len(lableVer) / 2), h - 2, termbox.ColorCyan, coldef, lableVer)
        // Коммандная панель
        packUI.Fill(0, ui.BodyY + 1, w, 1, termbox.Cell{Ch: '─'})
        packUI.Tprint(ui.BodyX + 2, ui.BodyY + 2, coldef, coldef, ">")
        packUI.Tprint(ui.BodyX - len([]rune(lableCommand)) - 1, ui.BodyY + 2, coldef, coldef, lableCommand)
        packUI.Fill(0, ui.BodyY + 3, w, 1, termbox.Cell{Ch: '─'})
        // Нижняя панель
        step := ui.BodyX + 4
        // help
        packUI.Tprint(step, h - 2, termbox.ColorCyan, coldef, lableHelp)
        packUI.Fill(step, h - 1, len([]rune(lableHelp)), 1, termbox.Cell{Ch: '─'})
        termbox.SetCell(step - 1, h - 1, '└', coldef, coldef)
        termbox.SetCell(step + len([]rune(lableHelp)), h - 1, '┘', coldef, coldef)
        // console
        step += len([]rune(lableHelp)) + 4
        packUI.Tprint(step, h - 2, termbox.ColorCyan, coldef, lableConsole)
        packUI.Fill(step, h - 1, len([]rune(lableConsole)), 1, termbox.Cell{Ch: '─'})
        termbox.SetCell(step - 1, h - 1, '└', coldef, coldef)
        termbox.SetCell(step + len([]rune(lableConsole)), h - 1, '┘', coldef, coldef)
        // settings
        step += len([]rune(lableConsole)) + 4
        packUI.Tprint(step, h - 2, termbox.ColorCyan, coldef, lableSettings)
        packUI.Fill(step, h - 1, len([]rune(lableSettings)), 1, termbox.Cell{Ch: '─'})
        termbox.SetCell(step - 1, h - 1, '└', coldef, coldef)
        termbox.SetCell(step + len([]rune(lableSettings)), h - 1, '┘', coldef, coldef)
        // F2
        step += len([]rune(lableConsole)) + 4
        packUI.Tprint(step, h - 2, termbox.ColorCyan, coldef, lableF2)
        packUI.Fill(step, h - 1, len([]rune(lableF2)), 1, termbox.Cell{Ch: '─'})
        termbox.SetCell(step - 1, h - 1, '└', coldef, coldef)
        termbox.SetCell(step + len([]rune(lableF2)), h - 1, '┘', coldef, coldef)
        // Quit
        step += len([]rune(lableF2)) + 4
        packUI.Tprint(step, h - 2, termbox.ColorCyan, coldef, lableQuit)
        packUI.Fill(step, h - 1, len([]rune(lableQuit)), 1, termbox.Cell{Ch: '─'})
        termbox.SetCell(step - 1, h - 1, '└', coldef, coldef)
        termbox.SetCell(step + len([]rune(lableQuit)), h - 1, '┘', coldef, coldef)
        // Пока оставил нижняя панель может пригодиться
        // termbox.SetCell(ui.BodyX + 3, ui.BodyY - 1, '┌', coldef, coldef)
        // termbox.SetCell(ui.BodyX + 3, ui.BodyY, '│', coldef, coldef)
        // Fill(ui.BodyX + 4, ui.BodyY - 1, w - ui.BodyX - 6, 1, termbox.Cell{Ch: '─'})
        // termbox.SetCell(w - 3, ui.BodyY - 1, '┐', coldef, coldef)
        // termbox.SetCell(w - 3, ui.BodyY, '│', coldef, coldef)
        // Fill(ui.BodyX + 3, ui.BodyY + 4, 1, h - ui.BodyY - 4, termbox.Cell{Ch: '│'})
        // Fill(w - 3, ui.BodyY + 4, 1, h - ui.BodyY - 4, termbox.Cell{Ch: '│'})
    })
    return frame
}
