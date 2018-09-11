package main

import (
	// "fmt"
    "errors"
    "strings"
    "runtime"
    "strconv"
    "reflect"
    "sort"
    "io/ioutil"
    "os"
    "encoding/json"
    packUI "./ui"
    "./commands"
    "./market"
    "./market/core"
    "./market/strategies"
    marketConst "./market/constants"
    "./utilities"
    "github.com/nsf/termbox-go"
    // "github.com/golang-collections/collections/stack"
)

const (
    PATH_SETTINGS = "./settings/"

    CONTROL_COMMAND_BOX = "command_box"
    CONTROL_COMMANDS = "commands"
    CONTROL_INDICATORS = "indicators"
    CONTROL_HELP_VIEW = "help_view"
    CONSOLE_INDICATORS = "commands"
    EVENT_MAIN_KEYS = "main_keys"

    MSG_PLACE_NAME = "{name}"
    MSG_PLACE_CMD = "{cmd}"
    MSG_PLACE_CONTROL = "{control}"
    MSG_PLACE_PARAM = "{param}"

    MSG_SETTINGS_NO_NAME = "Отсутсвет имя настроеек"

    MSG_CMD_SUCCESS = `Команда "` + MSG_PLACE_CMD + `" УСПЕШНО выполнена`
    MSG_CMD_FAILD = `Команда "` + MSG_PLACE_CMD + `" ПРОВАЛЕНА`
    MSG_CMD_NOT_EXIST = `Команды "` + MSG_PLACE_CMD + `" не существует`
    MSG_CMD_NOT_PARAM = `У комманды "` + MSG_PLACE_CMD + `" нет параметра "` + MSG_PLACE_PARAM + `", для получения список комманд используйте параметр "` + commands.TEXT_SEPARATOR_CMD_PARAM + commands.CMD_ALL_PARAM_HELP + `"`
    MSG_CMD_INCORRECT_PARAMS = `У комманды "` + MSG_PLACE_CMD + `" указан не корректно параметр "` + MSG_PLACE_PARAM + `", для получения описания комманды используйте параметр "` + commands.TEXT_SEPARATOR_CMD_PARAM + commands.CMD_ALL_PARAM_HELP + `"`
    MSG_CMD_NOT_PARAMS = `Комманда "` + MSG_PLACE_CMD + `" должна иметь параметры, для получения описания комманды используйте параметр "` + commands.TEXT_SEPARATOR_CMD_PARAM + commands.CMD_ALL_PARAM_HELP + `"`
    MSG_CMD_PARAM_NOT_VALUE = `Параметр "` + MSG_PLACE_PARAM + `" должен иметь значение, для получения описания комманды используйте параметр "` + commands.TEXT_SEPARATOR_CMD_PARAM + commands.CMD_ALL_PARAM_HELP + `"`
    MSG_CMD_NOT_EXECUTE = `Комманда "` + MSG_PLACE_CMD + `" не имеет исполняюемую функцию, команду не возможно исполнить.`
)

//////////////////////////////////////////////////////////////////////////////
var form *packUI.UI
var settings *Settings
//
type Settings struct {
    Name string                                             `json:"name"`
    IsDefault bool                                          `json:"is_default"`
    ParamsExchanges map[string]map[string]map[string]string `json:"params_exchanges"`
    StartedExchanges map[string]map[string]bool             `json:"started_exchanges"`
    StartedStartegies map[string]bool                       `json:"started_strategies"`
}
// Функция создания нового объекта настроеек
func NewSettings(name string) *Settings {
    settings := new(Settings)
    settings.Name = name
    settings.IsDefault = false
    settings.ParamsExchanges = make(map[string]map[string]map[string]string)
    settings.StartedExchanges = make(map[string]map[string]bool)
    settings.StartedStartegies = make(map[string]bool)
    return settings
}
func (settings *Settings) Init() {
    if settings.ParamsExchanges == nil {
        settings.ParamsExchanges = make(map[string]map[string]map[string]string)
    }
    if settings.StartedExchanges == nil {
        settings.StartedExchanges = make(map[string]map[string]bool)
    }
    if settings.StartedStartegies == nil {
        settings.StartedStartegies = make(map[string]bool)
    }
}
// Запись насироек на диск
func (settings *Settings) Save() error {
    return settings.SaveAs(settings.Name)
}
func (settings *Settings) SaveAs(name string) error {
    if name == "" {
        return errors.New(MSG_SETTINGS_NO_NAME)
    }
    lastName := settings.Name
    settings.Name = name
    if err := WriteSettings(settings, PATH_SETTINGS + name + ".json"); err != nil {
        return err
    }
    if settings.IsDefault && lastName != name {
        var settingsTemp *Settings
        var err error
        if settingsTemp, err = ReadSettings(PATH_SETTINGS + lastName + ".json"); err != nil {
            return err
        }
        settingsTemp.IsDefault = false
        if err := WriteSettings(settings, PATH_SETTINGS + lastName + ".json"); err != nil {
            return err
        }

    }
    return nil
}
// Чтение настроек из файла на диске
func (settings *Settings) Apply(terminal *market.Terminal) bool {
    if settings.ParamsExchanges == nil {
        return false
    }
    // Выставляем параметры согласно настроек
    for setExchange, entities := range settings.ParamsExchanges {
        exchange := terminal.Exchanges[setExchange]
        if exchange != nil {
            for entity, properties := range entities {
                if properties != nil {
                    params := make(map[string]interface{})
                    for param, value := range properties {
                        params[param] = value
                    }
                    res := exchange.SetValues(entity, params)
                    if !res {
                        return false
                    }
                }
            }
        }
    }
    // Запускаем биржи согласно настроек
    for setExchange, entities := range settings.StartedExchanges {
        exchange := terminal.Exchanges[setExchange]
        if exchange != nil {
            for entity, value := range entities {
                if value {
                    _, success := exchange.Turn(entity)
                    if !success {
                        return false
                    }
                }
            }
        }
    }
    // Запускаем стратегии согласно настроек
    return true
}
//
func GetSettingsPaths() (paths map[string]string, err error) {
    var files []os.FileInfo
    files, err = ioutil.ReadDir(PATH_SETTINGS)
    if err != nil {
        return
    }
    paths = make(map[string]string)
    for _, f := range files {
        paths[strings.Split(f.Name(), ".")[0]] = PATH_SETTINGS + f.Name()
    }
    return
}
// Чтение настроек из файла на диске
func ReadSettings(path string) (settings *Settings, err error) {
    var jsonFile []byte
    jsonFile, err = ioutil.ReadFile(path)
    if err != nil {
		return
	}
    settings = new(Settings)
    json.Unmarshal(jsonFile, settings)
    return
}
// Запись настроек в файл
func WriteSettings(settings *Settings, path string) error {
    var settingsJson []byte
    var err error
    settingsJson, err = json.Marshal(settings)
    if err != nil {
        return err
    }
    err = ioutil.WriteFile(path, settingsJson, 0644)
    if err != nil {
        return err
    }
    return nil
}
// Получение всех имеющихся настроеек
func GetSettings() (allSettings map[string]*Settings, err error) {
    var paths map[string]string
    var settings *Settings
    paths, err = GetSettingsPaths()
    if err != nil {
		return
	}
    allSettings = make(map[string]*Settings)
    for _, path := range paths {
        settings, err = ReadSettings(path)
        if err != nil {
    		return
    	}
        allSettings[settings.Name] = settings
    }
    return
}
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
func (command *Command) Execute(ui *packUI.UI) string {
    if command.Params[commands.CMD_ALL_PARAM_HELP] != nil {
        help := command.Params[commands.CMD_ALL_PARAM_HELP].Value
        if help != "" {
            ui.SetMainConsole(commands.CMD_HELP)
            ui.MainConsole.Values[CONTROL_HELP_VIEW] = command.Name
            form.MainConsole.Redraw()
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
//////////////////////////////////////////////////////////////////////////////
const coldef = termbox.ColorDefault
// Функция отрисовки основной формы
func formDraw(ui *packUI.UI) {
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
    // ПОка осталив нижняя панель может пригодиться
    // termbox.SetCell(ui.BodyX + 3, ui.BodyY - 1, '┌', coldef, coldef)
    // termbox.SetCell(ui.BodyX + 3, ui.BodyY, '│', coldef, coldef)
    // Fill(ui.BodyX + 4, ui.BodyY - 1, w - ui.BodyX - 6, 1, termbox.Cell{Ch: '─'})
    // termbox.SetCell(w - 3, ui.BodyY - 1, '┐', coldef, coldef)
    // termbox.SetCell(w - 3, ui.BodyY, '│', coldef, coldef)
    // Fill(ui.BodyX + 3, ui.BodyY + 4, 1, h - ui.BodyY - 4, termbox.Cell{Ch: '│'})
    // Fill(w - 3, ui.BodyY + 4, 1, h - ui.BodyY - 4, termbox.Cell{Ch: '│'})
}
// Функция обработки нажатия кнопок
func eventKeys(ev termbox.Event) {
    switch ev.Key {
    case termbox.KeyEsc:
        form.Close()
    case termbox.KeyF2:
        form.SetMainConsole(packUI.CONSOLE_LOG)
        form.MainConsole.Redraw()
        return
    }
}
// Функция обработки нажатия клавиш для поля ввода комманд
func eventKeysCommandBox(ev termbox.Event) {
    commandBox := form.Controls[CONTROL_COMMAND_BOX].(*packUI.EditBox)
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
    form.Consoles[CONTROL_COMMAND_BOX].Redraw()
}
// Функция отрисовки комманой строки
func commandBoxDraw(console *packUI.Console, bodyX, bodyY int) {
    commandBox := form.Controls[CONTROL_COMMAND_BOX].(*packUI.EditBox)
    w, _ := termbox.Size()
    commandBox.Draw(bodyX + 4, bodyY + 2, w - bodyX, 1)
    termbox.SetCursor(bodyX + 4 + commandBox.CursorX(), bodyY + 2)
}
// Функция очистки коммандной строки
func commandBoxClear(console *packUI.Console, bodyX, bodyY int) {
    //form.Controls[CONTROL_COMMAND_BOX].(*packUI.EditBox).Clear()
    w, _ := termbox.Size()
    packUI.Fill(bodyX + 4, bodyY + 2, w - bodyX + 4, 1, termbox.Cell{Ch: ' '})
}
// Функция отрисовки индикаторов
func indicatorsDraw(console *packUI.Console, bodyX, bodyY int) {
    if console.Controls[CONTROL_INDICATORS] == nil {
        return
    }
    indicators := console.Controls[CONTROL_INDICATORS].(map[string]map[string]string)
    var keysGroup []string
    for key, _ := range indicators {
        keysGroup = append(keysGroup, key)
    }
    sort.Strings(keysGroup)
    k := 0
    for i:=len(keysGroup) - 1; i >= 0; i-- {
        var keysIndicator []string
        for key, _ := range indicators[keysGroup[i]] {
            keysIndicator = append(keysIndicator, key)
        }
        sort.Strings(keysIndicator)
        colForm := termbox.ColorCyan
        for j:=len(keysIndicator) - 1; j >= 0; j-- {
            termbox.SetCell(3, bodyY - k - 1, '│', colForm, coldef)
            termbox.SetCell(3, bodyY - k - 2, '│', colForm, coldef)
            termbox.SetCell(bodyX - 1, bodyY - k - 1, '│', colForm, coldef)
            termbox.SetCell(bodyX - 1, bodyY - k - 2, '│', colForm, coldef)
            if j > 0 {
                termbox.SetCell(3, bodyY - k, '└', colForm, coldef)
                termbox.SetCell(bodyX - 1, bodyY - k, '┘', colForm, coldef)
            } else {
                termbox.SetCell(3, bodyY - k, '├', colForm, coldef)
                termbox.SetCell(bodyX - 1, bodyY - k, '┤', colForm, coldef)
            }
            packUI.Fill(5, bodyY - k, bodyX - 7, 1, termbox.Cell{Ch: '╌', Fg: colForm})
            packUI.Tprint(5, bodyY - k, coldef, coldef, keysIndicator[j] + " ")
            packUI.Tprint(bodyX - len([]rune(indicators[keysGroup[i]][keysIndicator[j]])) - 3, bodyY - k, coldef, coldef, " " + indicators[keysGroup[i]][keysIndicator[j]])
            k++
        }
        k++
        packUI.Fill(3, bodyY - k, bodyX - 3, 1, termbox.Cell{Ch: '─', Fg: colForm})
        groupTitle := "┤" + keysGroup[i] + "├"
        packUI.Tprint(4, bodyY - k, colForm, coldef, groupTitle)
        termbox.SetCell(3, bodyY - k, '┌', colForm, coldef)
        termbox.SetCell(bodyX - 1, bodyY - k, '┐', colForm, coldef)
        k = k + 2
    }
}
// Функция очистки индикаторов
func indicatorsClear(console *packUI.Console, bodyX, bodyY int) {
    // _, h := termbox.Size()
    packUI.Fill(3, 2, bodyX - 3, bodyY - 1, termbox.Cell{Ch: ' '})
}
// Функция обработки индикаторов
func indicatorsAction(console *packUI.Console, data []interface{}) {
    indicators := make(map[string]map[string]string)
    settingsTitle := "Настройки"
    terminalTitle := "Терминал"
    indicators[settingsTitle] = make(map[string]string)
    indicators[terminalTitle] = make(map[string]string)
    // Записываем индикаторы по настройкм
    indicators[settingsTitle]["Название"] = settings.Name
    indicators[settingsTitle]["Автозагрузка"] = strconv.FormatBool(settings.IsDefault)
    // Записываем индикаторы по терминалу
    terminal := form.Controls["terminal"]
    if terminal != nil {
        indicators[terminalTitle] = utilities.AppendMap(indicators[terminalTitle], terminal.(*market.Terminal).GetIndicators())
    }
    // Записываем индикаторы по биржам
    for _, exchange := range terminal.(*market.Terminal).Exchanges {
        indicators[exchange.Name] = make(map[string]string)
        indicators[exchange.Name] = utilities.AppendMap(indicators[exchange.Name], exchange.GetIndicators())
    }
    console.Controls[CONTROL_INDICATORS] = indicators
}
// Функция отрисовки стратегии arbitrage
func arbitrageDraw(console *packUI.Console, bodyX, bodyY int) {
    if console.Controls["Properties"] == nil {
        return
    }
    name := console.Controls["Properties"].(map[string]interface{})["test"].(string)
    packUI.Tprint(bodyX + 4, 2, coldef, coldef, name)
}
// Функция очистки стратегии arbitrage
func arbitrageClear(console *packUI.Console, bodyX, bodyY int) {
    w, _ := termbox.Size()
    packUI.Fill(bodyX + 2, 2, w - bodyX + 2, bodyY - 2, termbox.Cell{Ch: ' '})
}
// Функция обработки стратегии arbitrage
func arbitrageAction(console *packUI.Console, data []interface{}) {
    console.Controls["Properties"] = data[0]
}
// Функция отрисовки текущей комманды
func currentCMDDraw(console *packUI.Console, bodyX, bodyY int) {
    // packUI.Tprint(bodyX + 2, bodyY, coldef, coldef, ">")
    packUI.Tprint(bodyX + 4, bodyY, coldef, coldef, form.Values["CurrentCommand"])
}
// Функция очистки текущей комманды
func currentCMDClear(console *packUI.Console, bodyX, bodyY int) {
    w, _ := termbox.Size()
    packUI.Fill(bodyX + 4, bodyY, w - bodyX + 2, 1, termbox.Cell{Ch: ' '})
}
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
// Функция отрисовки встроенной консоли помощи
func helpDraw(console *packUI.Console, bodyX, bodyY int) {
    view := console.Values[CONTROL_HELP_VIEW]
    commandList := form.Controls[CONTROL_COMMANDS].(map[string]*Command)
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
                for _, con := range form.Consoles  {
                    if con.IsActive {
                        continue
                    }
                    packUI.Tprint(indent + x, 9 + line, coldef, coldef, con.Name + ", ")
                    x = x + len([]rune(con.Name + ", "))
                }
            }
        } else {
            form.Error(strings.Replace(MSG_CMD_NOT_EXIST, MSG_PLACE_CMD, view, 1))
        }
    }
}
// Функция очитски встроенной консоли помощи
func helpClear(console *packUI.Console, bodyX, bodyY int) {
    w, _ := termbox.Size()
    packUI.Fill(bodyX + 2, 2, w - bodyX + 4, bodyY - 2, termbox.Cell{Ch: ' '})
}
// Функция отрисовки встроенной консоли логов
func logDraw(console *packUI.Console, bodyX, bodyY int) {
    // fmt.Println("redraw", console.Values["Action"])
    y := bodyY - 1
    packUI.Fill(bodyX + 2, 2, 1, y - 1, termbox.Cell{Ch: '>'})
    if console.Controls["Text"] != nil {
        newValue := console.Controls["Text"].(*packUI.Text)
        // Сдвигаем всю историю комманд вверх
        console.Controls[strconv.FormatInt(int64(y), 10)] = newValue
        for y = 1 ; y <= bodyY - 1; y++ {
            console.Controls[strconv.FormatInt(int64(y), 10)] = console.Controls[strconv.FormatInt(int64(y + 1), 10)]
        }
        // Отображаем историю комманд
        for y = 1 ; y <= bodyY - 1; y++ {
            control := console.Controls[strconv.FormatInt(int64(y), 10)]
            if control != nil {
                text := valueToString(control.(*packUI.Text))
                if text != nil {
                    packUI.Tprint(bodyX + 4, y + 1, text.FG, text.BG, text.Value.(string))
                }
            }
        }
    }
}
// Функция очистки встроенной консоли логов
func logClear(console *packUI.Console, bodyX, bodyY int) {
    w, _ := termbox.Size()
    packUI.Fill(bodyX + 4, 2, w - bodyX + 4, bodyY - 2, termbox.Cell{Ch: ' '})
}
// Функция обработки встроенной консоли логов
func logAction(console *packUI.Console, data []interface{}) {
    // fmt.Println(action)
    if data == nil {
        return
    }
    if len(data) >= 1 {
        console.Controls["Text"] = data[0]
    }
}
// Функция обратоки исполнения команды console
func cmdConsoleExecute(command *Command) bool {
    run := command.Params[commands.CMD_CONSOLE_PARAM_RUN].Value
    if run != "" {
        if run == commands.CMD_HELP {
            form.SetMainConsole(packUI.CONSOLE_HELP)
        } else {
            form.SetMainConsole(run)
        }
        form.RedrawAll()
        return true
    }
    return false
}
// Функция обратоки исполнения команды help
func cmdHelpExecute(command *Command) bool {
    cmd := command.Params[commands.CMD_HELP_PARAM_CMD].Value
    form.SetMainConsole(packUI.CONSOLE_HELP)
    form.MainConsole.Values[CONTROL_HELP_VIEW] = cmd
    form.RedrawAll()
    return true
}
// Функция обратоки исполнения команды terminal
func cmdTerminalExecute(command *Command) bool {
    on := command.Params[commands.CMD_TERMINAL_PARAM_ON].Value
    name := command.Params[commands.CMD_TERMINAL_PARAM_NAME].Value
    entity := command.Params[commands.CMD_TERMINAL_PARAM_ENTITY].Value
    symbol := command.Params[commands.CMD_TERMINAL_PARAM_SYMBOL].Value
    terminal := command.Controls["terminal"].(*market.Terminal)
    // form.Log("on " + on + " name " + name + " entity " + entity + " symbol " + symbol)
    if on != "" && name != "" && entity != "" {
        exchange := terminal.Exchanges[name]
        if exchange != nil {
            if status, success := exchange.Turn(entity); success {
                // Сохраняем в настройках
                if settings.StartedExchanges[name] == nil {
                    settings.StartedExchanges[name] = make(map[string]bool)
                }
                settings.StartedExchanges[name][entity] = status
                return true
            }
        }
    }
    if name != "" && entity != "" && symbol != "" {
        exchange := terminal.Exchanges[name]
        if exchange != nil {
            params := make(map[string]interface{})
            params["symbol"] = symbol
            if exchange.SetValues(entity, params) {
                // Сохраняем в настройках
                if settings.ParamsExchanges[name] == nil {
                    settings.ParamsExchanges[name] = make(map[string]map[string]string)
                    settings.ParamsExchanges[name][entity] = make(map[string]string)
                }
                settings.ParamsExchanges[name][entity]["symbol"] = symbol
                return true
            }
        }
    }
    return false
}
// Функция обработки исполнения команды settings
func cmdSettingsExecute(command *Command) bool {
    load := command.Params[commands.CMD_SETTINGS_PARAM_LOAD].Value
    save := command.Params[commands.CMD_SETTINGS_PARAM_SAVE].Value
    def := command.Params[commands.CMD_SETTINGS_PARAM_DEFAULT].Value
    terminal := form.Controls["terminal"].(*market.Terminal)
    if load != "" {
        var err error
        settings, err = ReadSettings(PATH_SETTINGS + load + ".json")
        if err != nil {
            form.Error(err.Error())
            return false
        }
        if !settings.Apply(terminal) {
            return false
        }
    }
    if save == commands.CMD_VALUE_EMPTY {
        err := settings.Save()
        if err != nil {
            form.Error(err.Error())
            return false
        }
    }
    if save != "" && save != commands.CMD_VALUE_EMPTY {
        err := settings.SaveAs(save)
        if err != nil {
            form.Error(err.Error())
            return false
        }
    }
    if def != "" {
        allSettings, err := GetSettings()
        if err != nil {
            form.Error(err.Error())
            return false
        }
        settings.IsDefault = true
        for _, set := range allSettings {
            if set.Name == settings.Name {
                set.IsDefault = true
            } else {
                set.IsDefault = false
            }
            err := set.Save()
            if err != nil {
                form.Error(err.Error())
                return false
            }
        }
    }
    form.Consoles[CONSOLE_INDICATORS].Execute(nil)
    form.Consoles[CONSOLE_INDICATORS].Redraw()
    return true
}
// Функция обработки события нового сигнала
func eventNewSignal(event string, params []interface{}, callback func(string)) {
    if params == nil {
        return
    }
    if len(params) >= 1 {
        form.ChLog<- params[0]
    }
}
// Функция обработки события нового действия от стратегии
func eventNewAction(event string, params []interface{}, callback func(string)) {
    if params == nil {
        return
    }
    if len(params) >= 2 {
        if params[1].(bool) {
            form.ChMsg<- params[0]
        } else {
            form.ChErr<- params[0]
        }
    }
}
// Функция обработки события расчета действия от стратегии
func eventCalculateAction(event string, params []interface{}, callback func(string)) {
    if params == nil {
        return
    }
    if len(params) >= 2 {
        form.DrawConsole(params[0].(string), []interface{} { params[1] })
    }
}
// Функция обработки события установки нового значения индикатора
func eventSetIndicator(event string, params []interface{}, callback func(string)) {
    if params == nil {
        return
    }
    if len(params) >= 2 {
        form.Consoles[CONSOLE_INDICATORS].Execute(params)
        form.Consoles[CONSOLE_INDICATORS].Redraw()
    }
}
// Функция загрузки сохранненых параметров
func loadCurrentSettings(terminal *market.Terminal) {
    settings = NewSettings("default")
    allSettings, err := GetSettings()
    if err != nil {
        form.Error(err.Error())
        return
    }
    for _, set := range allSettings {
        if set.IsDefault {
            settings = set
            settings.Init()
            settings.Apply(terminal)
            return
        }
    }
}
//////////////////////////////////////////////////////////////////////////////
func main() {
    // Создаем консольный интерфейс и описываем основйной интерфейс консоли
    form = packUI.NewUI(formDraw)
    //-------------------------------------------------------------------------
    // Инициализируем необходимые контролы
    form.Controls[CONTROL_COMMAND_BOX] = new(packUI.EditBox)
    form.Controls[CONTROL_COMMANDS] = make(map[string]*Command)
    //-------------------------------------------------------------------------
    // Настраиваем обработку основных кнопок приложения
    form.Events[EVENT_MAIN_KEYS] = eventKeys
    // Настраиваем обработку событий клавиш
    form.Events[CONTROL_COMMAND_BOX] = eventKeysCommandBox
    //-------------------------------------------------------------------------
    // Создаем консоль для ввода комманд
    form.AddConsole(packUI.NewConsole(CONTROL_COMMAND_BOX, "", true,
                                      commandBoxDraw, commandBoxClear, nil))
    // Создаем консоль индикаторов
    form.AddConsole(packUI.NewConsole(CONSOLE_INDICATORS, "", true,
                                      indicatorsDraw, indicatorsClear, indicatorsAction))
    // Создаем консоли стратегии arbitrage
    form.AddConsole(packUI.NewConsole("arbitrage", `Торговая стратегия "Арбитраж"`, false,
                                      arbitrageDraw, arbitrageClear, arbitrageAction))
    // Инициализируем уже встроенные консоли
	form.Consoles[packUI.CONSOLE_CURRENT_CMD].Draw = currentCMDDraw
    form.Consoles[packUI.CONSOLE_CURRENT_CMD].Clear = currentCMDClear
	// Настраиваем консоль сообщений
	form.Consoles[packUI.CONSOLE_MESSAGE].Draw = messageDraw
	form.Consoles[packUI.CONSOLE_MESSAGE].Clear = messageClear
    form.Consoles[packUI.CONSOLE_MESSAGE].Action = messageAction
    // Настраиваем консоль помощи
	form.Consoles[packUI.CONSOLE_HELP].Draw = helpDraw
    form.Consoles[packUI.CONSOLE_HELP].Clear = helpClear
	// Настраиваем консоль логов
	form.Consoles[packUI.CONSOLE_LOG].Draw = logDraw
    form.Consoles[packUI.CONSOLE_LOG].Clear = logClear
    form.Consoles[packUI.CONSOLE_LOG].Action = logAction
    //-------------------------------------------------------------------------
	// Готовим команды
    // Комманда "console"
    form.Controls[CONTROL_COMMANDS].(map[string]*Command)[commands.CMD_CONSOLE] = NewCommand(
        commands.CMD_CONSOLE, commands.CMD_CONSOLE_DISCRIPTION,
        []*Param { NewParam(commands.CMD_CONSOLE_PARAM_RUN, commands.CMD_CONSOLE_PARAM_RUN_DISCRIPTION,
                            commands.CMD_CONSOLE_PARAM_RUN_EXAMPLE, commands.CMD_CONSOLE_PARAM_RUN_ISFLAG,
                            commands.CMD_CONSOLE_PARAM_RUN_ALLOWED_EMPTY)},
        commands.CMD_CONSOLE_PARAMS_REQUIRED,
        commands.CMD_CONSOLE_DEFAULT_PARAM,
        cmdConsoleExecute)
    // Комманда "help"
    form.Controls[CONTROL_COMMANDS].(map[string]*Command)[commands.CMD_HELP] = NewCommand(
        commands.CMD_HELP, commands.CMD_HELP_DISCRIPTION,
        []*Param {
            NewParam(commands.CMD_HELP_PARAM_CMD, commands.CMD_HELP_PARAM_CMD_DISCRIPTION,
                     commands.CMD_HELP_PARAM_CMD_EXAMPLE, commands.CMD_HELP_PARAM_CMD_ISFLAG,
                     commands.CMD_HELP_PARAM_CMD_ALLOWED_EMPTY)},
        commands.CMD_HELP_PARAMS_REQUIRED,
        commands.CMD_HELP_DEFAULT_PARAM,
        cmdHelpExecute)
    // Комманда "terminal"
    form.Controls[CONTROL_COMMANDS].(map[string]*Command)[commands.CMD_TERMINAL] = NewCommand(
        commands.CMD_TERMINAL, commands.CMD_TERMINAL_DISCRIPTION,
        []*Param {
            NewParam(commands.CMD_TERMINAL_PARAM_ON, commands.CMD_TERMINAL_PARAM_ON_DISCRIPTION,
                     commands.CMD_TERMINAL_PARAM_ON_EXAMPLE, commands.CMD_TERMINAL_PARAM_ON_ISFLAG,
                     commands.CMD_TERMINAL_PARAM_ON_ALLOWED_EMPTY),
            NewParam(commands.CMD_TERMINAL_PARAM_NAME, commands.CMD_TERMINAL_PARAM_NAME_DISCRIPTION,
                     commands.CMD_TERMINAL_PARAM_NAME_EXAMPLE, commands.CMD_TERMINAL_PARAM_NAME_ISFLAG,
                     commands.CMD_TERMINAL_PARAM_NAME_ALLOWED_EMPTY),
            NewParam(commands.CMD_TERMINAL_PARAM_ENTITY, commands.CMD_TERMINAL_PARAM_ENTITY_DISCRIPTION,
                     commands.CMD_TERMINAL_PARAM_ENTITY_EXAMPLE, commands.CMD_TERMINAL_PARAM_ENTITY_ISFLAG,
                     commands.CMD_TERMINAL_PARAM_ENTITY_ALLOWED_EMPTY),
            NewParam(commands.CMD_TERMINAL_PARAM_SYMBOL, commands.CMD_TERMINAL_PARAM_SYMBOL_DISCRIPTION,
                     commands.CMD_TERMINAL_PARAM_SYMBOL_EXAMPLE, commands.CMD_TERMINAL_PARAM_SYMBOL_ISFLAG,
                     commands.CMD_TERMINAL_PARAM_SYMBOL_ALLOWED_EMPTY)},
        commands.CMD_TERMINAL_PARAMS_REQUIRED,
        commands.CMD_TERMINAL_DEFAULT_PARAM,
        cmdTerminalExecute)
    // Комманда "settings"
    form.Controls[CONTROL_COMMANDS].(map[string]*Command)[commands.CMD_SETTINGS] = NewCommand(
        commands.CMD_SETTINGS, commands.CMD_SETTINGS_DISCRIPTION,
        []*Param {
            NewParam(commands.CMD_SETTINGS_PARAM_LOAD, commands.CMD_SETTINGS_PARAM_LOAD_DISCRIPTION,
                     commands.CMD_SETTINGS_PARAM_LOAD_EXAMPLE, commands.CMD_SETTINGS_PARAM_LOAD_ISFLAG,
                     commands.CMD_SETTINGS_PARAM_LOAD_ALLOWED_EMPTY),
            NewParam(commands.CMD_SETTINGS_PARAM_SAVE, commands.CMD_SETTINGS_PARAM_SAVE_DISCRIPTION,
                     commands.CMD_SETTINGS_PARAM_SAVE_EXAMPLE, commands.CMD_SETTINGS_PARAM_SAVE_ISFLAG,
                     commands.CMD_SETTINGS_PARAM_SAVE_ALLOWED_EMPTY),
            NewParam(commands.CMD_SETTINGS_PARAM_DEFAULT, commands.CMD_SETTINGS_PARAM_DEFAULT_DISCRIPTION,
                     commands.CMD_SETTINGS_PARAM_DEFAULT_EXAMPLE, commands.CMD_SETTINGS_PARAM_DEFAULT_ISFLAG,
                     commands.CMD_SETTINGS_PARAM_DEFAULT_ALLOWED_EMPTY)},
        commands.CMD_SETTINGS_PARAMS_REQUIRED,
        commands.CMD_SETTINGS_DEFAULT_PARAM,
        cmdSettingsExecute)
	//-------------------------------------------------------------------------
	// Устанавливаем текущую консоль
	form.SetMainConsole(packUI.CONSOLE_LOG)
    form.RedrawAll()
	// Создаем терминал
    terminal, err := market.NewTerminal(form.ChMsg, form.ChErr)
    loadCurrentSettings(terminal)
    // запоминаем терминал для общего доступа
    form.Controls["terminal"] = terminal
    // Создаем стратегии
    terminal.AddStrategy(strategies.NewArbitrage("arbitrage"))
    // Обрабатываем событие на появление нового сигнала
    terminal.AddAction(marketConst.EVENT_NEW_SIGNAL, eventNewSignal)
    // Обрабатываем событие на сформирование действие тратегией
    terminal.AddAction(marketConst.EVENT_NEW_ACTION, eventNewAction)
    // Обрабатываем событие на сформирование действие тратегией
    terminal.AddAction(marketConst.EVENT_CALCULATE_ACTION, eventCalculateAction)
    // Обрабатываем событие на установку нового значения у индикатора
    terminal.AddAction(marketConst.EVENT_SET_INDICATOR, eventSetIndicator)
    // Отображаем индикаторы
    form.Consoles[CONSOLE_INDICATORS].Execute(nil)
    form.Consoles[CONSOLE_INDICATORS].Redraw()
    form.Controls[CONTROL_COMMANDS].(map[string]*Command)[commands.CMD_TERMINAL].Controls["terminal"] = terminal
    // Передаем данные об ошибке терминала в интерфейс
    if err != nil {
        form.Error(err.Error())
    }
	// Ждем заверщения программы
	for {
		select {
		case <-form.ChClosed:
			return
		default:
			runtime.Gosched()
		}
	}
}
//////////////////////////////////////////////////////////////////////////////

func valueToString(text *packUI.Text) *packUI.Text {
    if reflect.TypeOf(text.Value).Kind() == reflect.String {
        return text
    }
    if reflect.TypeOf(text.Value).Elem().Name() == "Signal" {
        signal := text.Value.(*core.Signal)
        var res *packUI.Text
        switch signal.Entity {
            case marketConst.ENTITY_TICK:
                res = &packUI.Text { Value: text.Value.(*core.Signal).TickStr(true), BG: termbox.ColorDefault }
            case marketConst.ENTITY_DEPTH:
                res = &packUI.Text { Value: text.Value.(*core.Signal).DepthStr(true), BG: termbox.ColorDefault }
            case marketConst.ENTITY_CANDLE:
                res = &packUI.Text { Value: text.Value.(*core.Signal).CandlesStr(true), BG: termbox.ColorDefault }
        }
        if signal.TimeOut {
            res.FG = termbox.ColorRed
        } else {
            res.FG = termbox.ColorGreen
        }
        return res
    }
    if reflect.TypeOf(text.Value).Elem().Name() == "Message" {
        res := &packUI.Text { Value: text.Value.(*core.Message).FullString() }
        if text.Value.(*core.Message).Kind == core.MSG_WORNING {
            res.FG = termbox.ColorYellow
            res.BG = termbox.ColorYellow
        } else {
            res.FG = text.FG
            res.BG = text.BG
        }
        return res
    }
    return nil
}

func viewCurrentCommand(cmd string) {
    form.Values["CurrentCommand"] = cmd
    console := form.Consoles[packUI.CONSOLE_CURRENT_CMD]
    if console.IsActive {
        console.Execute([]interface{} { cmd })
        console.Redraw()
    }
    form.DrawConsole(packUI.CONSOLE_LOG, []interface{} { &packUI.Text { Value: cmd } })
}

func CMD(cmd string) {
    if form == nil {
        return
    }
    if cmd == "" {
        return
    }
    form.Values["LastCommand"] = form.Values["CurrentCommand"]
    viewCurrentCommand(cmd)
    command, err := ParseCommand(cmd)
    if err == "" {
        err = command.Execute(form)
        if err == "" {
            // form.Consoles[packUI.CONSOLE_MESSAGE].ClearDraw()
            form.Message(strings.Replace(MSG_CMD_SUCCESS, MSG_PLACE_CMD, cmd, 1))
        } else {
            form.Error(err)
        }
    } else {
        form.Error(err)
    }
    form.Controls[CONTROL_COMMAND_BOX].(*packUI.EditBox).Clear()
}

func ParseCommand(cmd string) (command *Command, err string) {
    if form == nil {
        return
    }
    if cmd == "" {
        return
    }
    // Убиваем все пробелы
    cmdArray := strings.Replace(cmd, " ", "", -1)
    parts := strings.Split(cmdArray, "--")
    commandList := form.Controls[CONTROL_COMMANDS].(map[string]*Command)
    // Находим команду в коллекции
    command = commandList[parts[0]]
    if command == nil {
        parts = strings.Split(cmd, " ")
        if len(parts) == 2 {
            command = commandList[parts[0]]
            if command == nil {
                err = strings.Replace(MSG_CMD_NOT_EXIST, MSG_PLACE_CMD, parts[0], 1)
                command = &Command { Cmd: cmd }
                return
            }
            command.Cmd = cmd
            command.ResetParams()
            if command.DefaultParam != "" {
                command.Params[command.DefaultParam].Value = strings.ToLower(parts[1])
                return
            }
        }
        err = strings.Replace(MSG_CMD_NOT_EXIST, MSG_PLACE_CMD, parts[0], 1)
        command = &Command { Cmd: cmd }
        return
    }
    command.Cmd = cmd
    command.ResetParams()
    // Проверяем на наличие параметров
    if command.ParamsRequired && len(parts) == 1 {
        err = strings.Replace(MSG_CMD_NOT_PARAMS, MSG_PLACE_CMD, command.Name, 1)
    }
    // Считываем параметры комманды
    for i, part := range parts {
        if i > 0 {
            paramParts := strings.Split(part, ":")
            param := command.Params[paramParts[0]]
            if param == nil{
                err = strings.Replace(MSG_CMD_INCORRECT_PARAMS, MSG_PLACE_CMD, command.Name, 1)
                err = strings.Replace(err, MSG_PLACE_PARAM, paramParts[0], 1)
                return
            }
            if len(paramParts) > 1 {
                if paramParts[1] == "" {
                    if !param.AllowedEmpty {
                        err = strings.Replace(MSG_CMD_PARAM_NOT_VALUE, MSG_PLACE_PARAM, param.Name, 1)
                        return
                    } else {
                        command.Params[paramParts[0]].Value = commands.CMD_VALUE_EMPTY
                    }
                }
                command.Params[paramParts[0]].Value = strings.ToLower(paramParts[1])
            } else {
                if param.IsFlag {
                    command.Params[paramParts[0]].Value = commands.CMD_VALUE_FLAG
                } else {
                    if !param.AllowedEmpty {
                        err = strings.Replace(MSG_CMD_PARAM_NOT_VALUE, MSG_PLACE_PARAM, param.Name, 1)
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

// ui.Indicators["Прибыль за день"] = "100"
// ui.Indicators["Активных бирж"] = "4"
// ui.Indicators["Активных стратегий"] = "1"
