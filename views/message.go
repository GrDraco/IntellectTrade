package views

import (
    packUI "../ui"

    "github.com/nsf/termbox-go"
)

const (
    CONSOLE_MESSAGE = "console_message"
    CONSOLE_MESSAGE_DISCRIPTION = ""
)

// Функция отрисовки сообщений
func messageDraw(console *packUI.Console, bodyX, bodyY int) {
    if console.Controls["Text"] != nil {
        text := valueToString(console.Controls["Text"].(*packUI.Text))
        if text != nil {
            packUI.Tprint(bodyX + 2, bodyY + 4, text.FG, text.BG, text.Value.(string))
        }
    }
}
// Функция очистки сообщений
func messageClear(console *packUI.Console, bodyX, bodyY int) {
    w, _ := termbox.Size()
    packUI.Fill(bodyX + 2, bodyY + 4, w - bodyX + 2, 1, termbox.Cell{Ch: ' '})
}
// Функция обработки сообщений
func messageAction(console *packUI.Console, data []interface{}) {
    if data == nil {
        return
    }
    if len(data) >= 1 {
        console.Controls["Text"] = data[0]
    }
}

func Message() *packUI.Console  {
    return packUI.NewConsole(CONSOLE_MESSAGE, CONSOLE_MESSAGE_DISCRIPTION, true, messageDraw, messageClear, messageAction)
}
