package ui

import (
    "strings"
    "../commands"
)

const (
    MSG_PLACE_CMD = "{cmd}"
    MSG_PLACE_CONTROL = "{control}"
    MSG_PLACE_PARAM = "{param}"

    MSG_CMD_SUCCESS = `Команда "` + MSG_PLACE_CMD + `" УСПЕШНО выполнена`
    MSG_CMD_FAILD = `Команда "` + MSG_PLACE_CMD + `" ПРОВАЛЕНА`
    MSG_CMD_NOT_EXIST = `Команды "` + MSG_PLACE_CMD + `" не существует`
    MSG_CMD_NOT_PARAM = `У комманды "` + MSG_PLACE_CMD + `" нет параметра "` + MSG_PLACE_PARAM + `", для получения список комманд используйте параметр "` + commands.TEXT_SEPARATOR_CMD_PARAM + commands.CMD_ALL_PARAM_HELP + `"`
    MSG_CMD_INCORRECT_PARAMS = `У комманды "` + MSG_PLACE_CMD + `" указан не корректно параметр "` + MSG_PLACE_PARAM + `", для получения описания комманды используйте параметр "` + commands.TEXT_SEPARATOR_CMD_PARAM + commands.CMD_ALL_PARAM_HELP + `"`
    MSG_CMD_NOT_PARAMS = `Комманда "` + MSG_PLACE_CMD + `" должна иметь параметры, для получения описания комманды используйте параметр "` + commands.TEXT_SEPARATOR_CMD_PARAM + commands.CMD_ALL_PARAM_HELP + `"`
    MSG_CMD_PARAM_NOT_VALUE = `Параметр "` + MSG_PLACE_PARAM + `" должен иметь значение, для получения описания комманды используйте параметр "` + commands.TEXT_SEPARATOR_CMD_PARAM + commands.CMD_ALL_PARAM_HELP + `"`
    MSG_CMD_NOT_EXECUTE = `Комманда "` + MSG_PLACE_CMD + `" не имеет исполняюемую функцию, команду не возможно исполнить.`
)

// Структура параметра комманды
type Param struct {
    Name string
    Description string
    Example string
    IsFlag bool
    AllowedEmpty bool
    Value string
}
// Функция создания нового параметра для комманды
func NewParam(name, description, example string, isFlag, allowedEmpty bool) *Param {
    return &Param { Name: name, Description: description, Example: example, IsFlag: isFlag, AllowedEmpty: allowedEmpty }
}
// Структура комманды
type Command struct {
    Name string
    Description string
    Cmd string
    ParamsRequired bool
    Params map[string]*Param
    DefaultParam string
    Controls map[string]interface{}
    exeFunc func(*Command) bool
}
// Функция создания новой комманды
func NewCommand(name, description string, params []*Param, paramsRequired bool, defaultParam string, exeFunc func(*Command)bool) *Command{
    command := Command { Name: name, Description: description, ParamsRequired: paramsRequired, DefaultParam: defaultParam, exeFunc: exeFunc }
    command.Params = make(map[string]*Param)
    command.Controls = make(map[string]interface{})
    for _, param := range params {
        command.Params[param.Name] = param
    }
    command.Params[commands.CMD_ALL_PARAM_HELP] = NewParam(commands.CMD_ALL_PARAM_HELP, commands.CMD_ALL_PARAM_HELP_DISCRIPTION,
                                                           commands.CMD_ALL_PARAM_HELP_EXAMPLE, commands.CMD_ALL_PARAM_HELP_ISFLAG,
                                                           commands.CMD_ALL_PARAM_HELP_ALLOWED_EMPTY)
    return &command
}
// Функция сброса параметров комманды
func (command *Command) ResetParams() {
    command.Cmd = ""
    for _, param := range command.Params {
        param.Value = ""
    }
}
// Функция исполнения комманды
func (command *Command) Execute(ui *UI, CONTROL_HELP_VIEW string) string {
    if command.Params[commands.CMD_ALL_PARAM_HELP] != nil {
        help := command.Params[commands.CMD_ALL_PARAM_HELP].Value
        if help != "" {
            ui.SetMainConsole(commands.CMD_HELP)
            ui.MainConsole.Values[CONTROL_HELP_VIEW] = command.Name
            ui.MainConsole.Redraw()
            return ""
        }
    }
    if command.exeFunc == nil {
        return strings.Replace(MSG_CMD_NOT_EXECUTE, MSG_PLACE_CMD, command.Name, 1)
    }
    if command.exeFunc(command) {
        return ""
    }
    return strings.Replace(MSG_CMD_FAILD, MSG_PLACE_CMD, command.Name, 1)
}
