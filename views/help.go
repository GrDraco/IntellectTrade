package views

import (
    "strconv"
    "strings"
    "../commands"
    packUI "../ui"
    "github.com/nsf/termbox-go"
)

const (
    CONSOLE_HELP = "help"
    CONSOLE_HELP_DISCRIPTION = "Описание комманд"

    CONTROL_HELP_VIEW = "help_view"
)

// Функция отрисовки встроенной консоли помощи
func helpDraw(console *packUI.Console, bodyX, bodyY int) {
    view := console.Values[CONTROL_HELP_VIEW]
    commandList := frame.Controls[CONTROL_COMMANDS].(map[string]*packUI.Command)
    if view == "" {
        i := 0
        for _, command := range commandList{
            packUI.Tprint(bodyX + 4, i + 3, coldef, coldef, command.Name + " - " + command.Description)
            i = i + 2
        }
    } else {
        command := commandList[view]
        if command != nil {
            indent := bodyX + 4
            packUI.Tprint(indent, 3, coldef, coldef, command.Name + " - " + command.Description)
            packUI.Tprint(indent, 5, coldef, coldef, "параметр по умолчанию: " + command.DefaultParam)
            packUI.Tprint(indent, 6, coldef, coldef, "обязательное наличие параметров: " + strconv.FormatBool(command.ParamsRequired))
            line := 0
            for _, param := range command.Params {
                packUI.Tprint(indent, 8 + line, coldef, coldef, commands.TEXT_SEPARATOR_CMD_PARAM + param.Name + ": " + param.Description)
                packUI.Tprint(indent + len([]rune(param.Name)) + 4, 9 + line, coldef, coldef, "флаг: " + strconv.FormatBool(param.IsFlag))
                packUI.Tprint(indent + len([]rune(param.Name)) + 4, 10 + line, coldef, coldef, "может быть пустым: " + strconv.FormatBool(param.AllowedEmpty))
                packUI.Tprint(indent + len([]rune(param.Name)) + 4, 11 + line, coldef, coldef, "пример: " + param.Example)
                line = line + 4
            }
            if command.Name == commands.CMD_CONSOLE {
                titleCon := "Перечень консолей: "
                packUI.Tprint(indent, 9 + line, coldef, coldef, titleCon)
                x := len([]rune(titleCon))
                for _, con := range frame.Consoles  {
                    if con.IsActive {
                        continue
                    }
                    packUI.Tprint(indent + x, 9 + line, coldef, coldef, con.Name + ", ")
                    x = x + len([]rune(con.Name + ", "))
                }
            }
        } else {
            frame.Error(strings.Replace(packUI.MSG_CMD_NOT_EXIST, packUI.MSG_PLACE_CMD, view, 1))
        }
    }
}
// Функция очитски встроенной консоли помощи
func helpClear(console *packUI.Console, bodyX, bodyY int) {
    w, _ := termbox.Size()
    packUI.Fill(bodyX + 2, 2, w - bodyX + 4, bodyY - 2, termbox.Cell{Ch: ' '})
}

func Help() *packUI.Console  {
    return packUI.NewConsole(CONSOLE_HELP, CONSOLE_HELP_DISCRIPTION, false, helpDraw, helpClear, nil)
}
