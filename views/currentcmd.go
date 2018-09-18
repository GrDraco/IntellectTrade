package views

import (
    packUI "../ui"
    "github.com/nsf/termbox-go"
)

const (
    CONSOLE_CURRENT_CMD = "console_current_cmd"
    CONSOLE_CURRENT_CMD_DISCRIPTION = ""
)

// Функция отрисовки текущей комманды
func currentCMDDraw(console *packUI.Console, bodyX, bodyY int) {
    // packUI.Tprint(bodyX + 2, bodyY, coldef, coldef, ">")
    packUI.Tprint(bodyX + 4, bodyY, coldef, coldef, frame.Values["CurrentCommand"])
}
// Функция очистки текущей комманды
func currentCMDClear(console *packUI.Console, bodyX, bodyY int) {
    w, _ := termbox.Size()
    packUI.Fill(bodyX + 4, bodyY, w - bodyX + 2, 1, termbox.Cell{Ch: ' '})
}

func CurrentCMD() *packUI.Console  {
    return packUI.NewConsole(CONSOLE_CURRENT_CMD, CONSOLE_CURRENT_CMD_DISCRIPTION, false, currentCMDDraw, currentCMDClear, nil)
}
