package views

import (
    "strings"
    "../commands"
    packUI "../ui"
    "github.com/nsf/termbox-go"
)

const (
    CONTROL_COMMAND_BOX = "command_box"
    CONTROL_COMMANDS = "commands"
)



// Функция обработки нажатия клавиш для поля ввода комманд
func EventKeysCommandBox(ev termbox.Event) {
    commandBox := frame.Controls[CONTROL_COMMAND_BOX].(*packUI.EditBox)
    switch ev.Key {
    case termbox.KeyArrowUp:
        commandBox.ReadHistoryUp()
    case termbox.KeyArrowDown:
        commandBox.ReadHistoryDown()
    case termbox.KeyEnter:
        CMD(commandBox.Value())
    case termbox.KeyArrowLeft, termbox.KeyCtrlB:
        commandBox.MoveCursorOneRuneBackward()
    case termbox.KeyArrowRight, termbox.KeyCtrlF:
        commandBox.MoveCursorOneRuneForward()
    case termbox.KeyBackspace, termbox.KeyBackspace2:
        commandBox.DeleteRuneBackward()
    case termbox.KeyDelete, termbox.KeyCtrlD:
        commandBox.DeleteRuneForward()
    case termbox.KeyTab:
        commandBox.InsertRune('\t')
    case termbox.KeySpace:
        commandBox.InsertRune(' ')
    case termbox.KeyCtrlK:
        commandBox.DeleteTheRestOfTheLine()
    case termbox.KeyHome, termbox.KeyCtrlA:
        commandBox.MoveCursorToBeginningOfTheLine()
    case termbox.KeyEnd, termbox.KeyCtrlE:
        commandBox.MoveCursorToEndOfTheLine()
    default:
        if ev.Ch != 0 {
            commandBox.InsertRune(ev.Ch)
        }
    }
    frame.Consoles[CONTROL_COMMAND_BOX].Redraw()
}
// Функция отрисовки комманой строки
func commandBoxDraw(console *packUI.Console, bodyX, bodyY int) {
    commandBox := frame.Controls[CONTROL_COMMAND_BOX].(*packUI.EditBox)
    w, _ := termbox.Size()
    commandBox.Draw(bodyX + 4, bodyY + 2, w - bodyX, 1)
    termbox.SetCursor(bodyX + 4 + commandBox.CursorX(), bodyY + 2)
}
// Функция очистки коммандной строки
func commandBoxClear(console *packUI.Console, bodyX, bodyY int) {
    //frame.Controls[CONTROL_COMMAND_BOX].(*packUI.EditBox).Clear()
    w, _ := termbox.Size()
    packUI.Fill(bodyX + 4, bodyY + 2, w - bodyX + 4, 1, termbox.Cell{Ch: ' '})
}
//
func CommandBox() *packUI.Console {
    return packUI.NewConsole(CONTROL_COMMAND_BOX, "", true, commandBoxDraw, commandBoxClear, nil)
}

func viewCurrentCommand(cmd string) {
    frame.Values["CurrentCommand"] = cmd
    console := frame.Consoles[CONSOLE_CURRENT_CMD]
    if console.IsActive {
        console.Execute([]interface{} { cmd })
        console.Redraw()
    }
    frame.DrawConsole(CONSOLE_LOG, []interface{} { &packUI.Text { Value: cmd } })
}

func CMD(cmd string) {
    if frame == nil {
        return
    }
    if cmd == "" {
        return
    }
    frame.Values["LastCommand"] = frame.Values["CurrentCommand"]
    viewCurrentCommand(cmd)
    command, err := ParseCommand(cmd)
    if err == "" {
        err = command.Execute(frame, CONTROL_HELP_VIEW)
        if err == "" {
            // frame.Consoles[packUI.ConsoleMsg()].ClearDraw()
            frame.Message(strings.Replace(packUI.MSG_CMD_SUCCESS, packUI.MSG_PLACE_CMD, cmd, 1))
        } else {
            frame.Error(err)
        }
    } else {
        frame.Error(err)
    }
    frame.Controls[CONTROL_COMMAND_BOX].(*packUI.EditBox).Clear()
}

func ParseCommand(cmd string) (command *packUI.Command, err string) {
    if frame == nil {
        return
    }
    if cmd == "" {
        return
    }
    // Убиваем все пробелы
    cmdArray := strings.Replace(cmd, " ", "", -1)
    parts := strings.Split(cmdArray, "--")
    commandList := frame.Controls[CONTROL_COMMANDS].(map[string]*packUI.Command)
    // Находим команду в коллекции
    command = commandList[parts[0]]
    if command == nil {
        parts = strings.Split(cmd, " ")
        if len(parts) == 2 {
            command = commandList[parts[0]]
            if command == nil {
                err = strings.Replace(packUI.MSG_CMD_NOT_EXIST, packUI.MSG_PLACE_CMD, parts[0], 1)
                command = &packUI.Command { Cmd: cmd }
                return
            }
            command.Cmd = cmd
            command.ResetParams()
            if command.DefaultParam != "" {
                command.Params[command.DefaultParam].Value = parts[1]//strings.ToLower(parts[1])
                return
            }
        }
        err = strings.Replace(packUI.MSG_CMD_NOT_EXIST, packUI.MSG_PLACE_CMD, parts[0], 1)
        command = &packUI.Command { Cmd: cmd }
        return
    }
    command.Cmd = cmd
    command.ResetParams()
    // Проверяем на наличие параметров
    if command.ParamsRequired && len(parts) == 1 {
        err = strings.Replace(packUI.MSG_CMD_NOT_PARAMS, packUI.MSG_PLACE_CMD, command.Name, 1)
    }
    // Считываем параметры комманды
    for i, part := range parts {
        if i > 0 {
            paramParts := strings.Split(part, ":")
            param := command.Params[paramParts[0]]
            if param == nil{
                err = strings.Replace(packUI.MSG_CMD_INCORRECT_PARAMS, packUI.MSG_PLACE_CMD, command.Name, 1)
                err = strings.Replace(err, packUI.MSG_PLACE_PARAM, paramParts[0], 1)
                return
            }
            if len(paramParts) > 1 {
                if paramParts[1] == "" {
                    if !param.AllowedEmpty {
                        err = strings.Replace(packUI.MSG_CMD_PARAM_NOT_VALUE, packUI.MSG_PLACE_PARAM, param.Name, 1)
                        return
                    } else {
                        command.Params[paramParts[0]].Value = commands.CMD_VALUE_EMPTY
                    }
                }
                command.Params[paramParts[0]].Value = paramParts[1]//strings.ToLower(paramParts[1])
            } else {
                if param.IsFlag {
                    command.Params[paramParts[0]].Value = commands.CMD_VALUE_FLAG
                } else {
                    if !param.AllowedEmpty {
                        err = strings.Replace(packUI.MSG_CMD_PARAM_NOT_VALUE, packUI.MSG_PLACE_PARAM, param.Name, 1)
                        return
                    } else {
                        command.Params[paramParts[0]].Value = commands.CMD_VALUE_EMPTY
                    }
                }
            }
        }
    }
    return
}
