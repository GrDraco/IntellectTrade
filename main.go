package main

import (
	// "fmt"
    "runtime"
    packUI "./ui"
    "./commands"
    "./market"
    "./market/strategies"
    marketConst "./market/constants"
    "./views"
    "github.com/nsf/termbox-go"
    // "github.com/golang-collections/collections/stack"
)

const (
    coldef = termbox.ColorDefault

    EVENT_MAIN_KEYS = "main_keys"

    PARAM_MAIN_CONSOLE = "main_console"

    MSG_SETTINGS_ERROR_APPLY = "Не удалось применить настройки, что-то пошло не так."
)

//////////////////////////////////////////////////////////////////////////////
var frame *packUI.UI
var settings *packUI.Settings
//////////////////////////////////////////////////////////////////////////////

// Функция обработки нажатия кнопок
func eventKeys(ev termbox.Event) {
    switch ev.Key {
    case termbox.KeyEsc:
        frame.Close()
    case termbox.KeyF2:
        frame.SetMainConsole(views.CONSOLE_LOG)
        frame.MainConsole.Redraw()
        return
    }
}
// Функция обратоки исполнения команды console
func cmdConsoleExecute(command *packUI.Command) bool {
    run := command.Params[commands.CMD_CONSOLE_PARAM_RUN].Value
    if run != "" {
        if run == commands.CMD_HELP {
            frame.SetMainConsole(views.CONSOLE_HELP)
            settings.ParamsMain[PARAM_MAIN_CONSOLE] = views.CONSOLE_HELP
        } else {
            frame.SetMainConsole(run)
            settings.ParamsMain[PARAM_MAIN_CONSOLE] = run
        }
        frame.RedrawAll()
        return true
    }
    return false
}
// Функция обратоки исполнения команды help
func cmdHelpExecute(command *packUI.Command) bool {
    cmd := command.Params[commands.CMD_HELP_PARAM_CMD].Value
    frame.SetMainConsole(views.CONSOLE_HELP)
    frame.MainConsole.Values[views.CONTROL_HELP_VIEW] = cmd
    frame.RedrawAll()
    return true
}
// Функция обратоки исполнения команды terminal
func cmdTerminalExecute(command *packUI.Command) bool {
    on := command.Params[commands.CMD_TERMINAL_PARAM_ON].Value
    off := command.Params[commands.CMD_TERMINAL_PARAM_OFF].Value
    name := command.Params[commands.CMD_TERMINAL_PARAM_NAME].Value
    entity := command.Params[commands.CMD_TERMINAL_PARAM_ENTITY].Value
    symbol := command.Params[commands.CMD_TERMINAL_PARAM_SYMBOL].Value
    terminal := command.Controls["terminal"].(*market.Terminal)
    // frame.Log("on " + on + " name " + name + " entity " + entity + " symbol " + symbol)
    if (on != "" || off != "") && name != "" && entity != "" {
        var start bool
        if on != "" {
            start = true
        }
        if off != "" {
            start = false
        }
        exchange := terminal.Exchanges[name]
        if exchange != nil {
            if status, success := exchange.Turn(entity, start); success {
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
func cmdSettingsExecute(command *packUI.Command) bool {
    load := command.Params[commands.CMD_SETTINGS_PARAM_LOAD].Value
    save := command.Params[commands.CMD_SETTINGS_PARAM_SAVE].Value
    def := command.Params[commands.CMD_SETTINGS_PARAM_DEFAULT].Value
    apply := command.Params[commands.CMD_SETTINGS_PARAM_APPLY].Value
    terminal := frame.Controls["terminal"].(*market.Terminal)
    if load != "" {
        err := settings.LoadFrom(load)
        if err != nil {
            frame.Error(err.Error())
            return false
        }
        frame.Consoles[views.CONSOLE_INDICATORS].Execute(nil)
        frame.Consoles[views.CONSOLE_INDICATORS].Redraw()
        if !settings.Apply(terminal) {
            return false
        }
    }
    if save == commands.CMD_VALUE_EMPTY {
        err := settings.Save()
        if err != nil {
            frame.Error(err.Error())
            return false
        }
    }
    if save != "" && save != commands.CMD_VALUE_EMPTY {
        err := settings.SaveAs(save)
        if err != nil {
            frame.Error(err.Error())
            return false
        }
    }
    if def != "" {
        allSettings, err := packUI.GetSettings()
        if err != nil {
            frame.Error(err.Error())
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
                frame.Error(err.Error())
                return false
            }
        }
    }
    if apply != "" {
        if !settings.Apply(terminal) {
            frame.Error(MSG_SETTINGS_ERROR_APPLY)
            return false
        }
    }
    frame.Consoles[views.CONSOLE_INDICATORS].Execute(nil)
    frame.Consoles[views.CONSOLE_INDICATORS].Redraw()
    return true
}
// Функция обработки исполнения команды strategy
func cmdStrategyExecute(command *packUI.Command) bool {
    name := command.Params[commands.CMD_STRATEGY_PARAM_NAME].Value
    // param := command.Params[commands.CMD_STRATEGY_PARAM_PARAM].Value
    // value := command.Params[commands.CMD_STRATEGY_PARAM_VALUE].Value
    on := command.Params[commands.CMD_STRATEGY_PARAM_ON].Value
    terminal := frame.Controls["terminal"].(*market.Terminal)
    if on != "" && name != "" {
        strategy := terminal.Strategies[name]
        if strategy != nil {
            strategy.Turn()
            // Сохраняем в настройках
            for property, value := range strategy.GetProperties() {
                if !strategy.GetKeysForSave()[property] {
                    continue
                }
                if settings.ParamsStartegies[name] == nil {
                    settings.ParamsStartegies[name] = make(map[string]interface{})
                }
                settings.ParamsStartegies[name][property] = value
            }
            return true
        }
    }
    return false
}
// Функция обработки события нового сигнала
func eventNewSignal(event string, params []interface{}, callback func(string)) {
    if params == nil {
        return
    }
    if len(params) >= 1 {
        frame.ChLog<- params[0]
    }
}
// Функция обработки события нового действия от стратегии
func eventNewAction(event string, params []interface{}, callback func(string)) {
    if params == nil {
        return
    }
    if len(params) >= 2 {
        if params[1].(bool) {
            frame.ChMsg<- params[0]
        } else {
            frame.ChErr<- params[0]
        }
    }
}
// Функция обработки события расчета действия от стратегии
func eventCalculateAction(event string, params []interface{}, callback func(string)) {
    if params == nil {
        return
    }
    if len(params) >= 2 {
        frame.DrawConsole(params[0].(string), []interface{} { params[1] })
    }
}
// Функция обработки события установки нового значения индикатора
func eventSetIndicator(event string, params []interface{}, callback func(string)) {
    if params == nil {
        return
    }
    if len(params) >= 2 {
        frame.Consoles[views.CONSOLE_INDICATORS].Execute(params)
        frame.Consoles[views.CONSOLE_INDICATORS].Redraw()
    }
}
// Функция загрузки сохранненых параметров
func loadCurrentSettings() error {
    settings = packUI.NewSettings("default")
    allSettings, err := packUI.GetSettings()
    if err != nil {
        return err
    }
    for _, set := range allSettings {
        if set.IsDefault {
            settings = set
            settings.Init()
            return nil
        }
    }
    return nil
}
//////////////////////////////////////////////////////////////////////////////
func main() {
    //
    err := loadCurrentSettings()
    // Создаем консольный интерфейс и описываем основйной интерфейс консоли
    frame = views.Frame(settings)
    //-------------------------------------------------------------------------
    // Инициализируем необходимые контролы
    frame.Controls[views.CONTROL_COMMAND_BOX] = new(packUI.EditBox)
    frame.Controls[views.CONTROL_COMMANDS] = make(map[string]*packUI.Command)
    //-------------------------------------------------------------------------
    // Настраиваем обработку основных кнопок приложения
    frame.Events[EVENT_MAIN_KEYS] = eventKeys
    // Настраиваем обработку событий клавиш
    frame.Events[views.CONTROL_COMMAND_BOX] = views.EventKeysCommandBox
    //-------------------------------------------------------------------------
    // Создаем консоль для ввода комманд
    frame.AddConsole(views.CommandBox())
    // Создаем консоль индикаторов
    frame.AddConsole(views.Indicators())
    // Создаем консоли стратегии arbitrage
    frame.AddConsole(views.Arbitrage())
    // Инициализируем уже встроенные консоли
    frame.AddConsole(views.CurrentCMD())
	// Настраиваем консоль сообщений
    frame.AddConsole(views.Message())
    // Настраиваем консоль помощи
    frame.AddConsole(views.Help())
	// Настраиваем консоль логов
    frame.AddConsole(views.Log())
    //-------------------------------------------------------------------------
	// Готовим команды
    // Комманда "console"
    frame.Controls[views.CONTROL_COMMANDS].(map[string]*packUI.Command)[commands.CMD_CONSOLE] = packUI.NewCommand(
        commands.CMD_CONSOLE, commands.CMD_CONSOLE_DISCRIPTION,
        []*packUI.Param { packUI.NewParam(commands.CMD_CONSOLE_PARAM_RUN, commands.CMD_CONSOLE_PARAM_RUN_DISCRIPTION,
                            commands.CMD_CONSOLE_PARAM_RUN_EXAMPLE, commands.CMD_CONSOLE_PARAM_RUN_ISFLAG,
                            commands.CMD_CONSOLE_PARAM_RUN_ALLOWED_EMPTY)},
        commands.CMD_CONSOLE_PARAMS_REQUIRED,
        commands.CMD_CONSOLE_DEFAULT_PARAM,
        cmdConsoleExecute)
    // Комманда "help"
    frame.Controls[views.CONTROL_COMMANDS].(map[string]*packUI.Command)[commands.CMD_HELP] = packUI.NewCommand(
        commands.CMD_HELP, commands.CMD_HELP_DISCRIPTION,
        []*packUI.Param {
            packUI.NewParam(commands.CMD_HELP_PARAM_CMD, commands.CMD_HELP_PARAM_CMD_DISCRIPTION,
                     commands.CMD_HELP_PARAM_CMD_EXAMPLE, commands.CMD_HELP_PARAM_CMD_ISFLAG,
                     commands.CMD_HELP_PARAM_CMD_ALLOWED_EMPTY)},
        commands.CMD_HELP_PARAMS_REQUIRED,
        commands.CMD_HELP_DEFAULT_PARAM,
        cmdHelpExecute)
    // Комманда "terminal"
    frame.Controls[views.CONTROL_COMMANDS].(map[string]*packUI.Command)[commands.CMD_TERMINAL] = packUI.NewCommand(
        commands.CMD_TERMINAL, commands.CMD_TERMINAL_DISCRIPTION,
        []*packUI.Param {
            packUI.NewParam(commands.CMD_TERMINAL_PARAM_ON, commands.CMD_TERMINAL_PARAM_ON_DISCRIPTION,
                     commands.CMD_TERMINAL_PARAM_ON_EXAMPLE, commands.CMD_TERMINAL_PARAM_ON_ISFLAG,
                     commands.CMD_TERMINAL_PARAM_ON_ALLOWED_EMPTY),
            packUI.NewParam(commands.CMD_TERMINAL_PARAM_OFF, commands.CMD_TERMINAL_PARAM_OFF_DISCRIPTION,
                     commands.CMD_TERMINAL_PARAM_OFF_EXAMPLE, commands.CMD_TERMINAL_PARAM_OFF_ISFLAG,
                     commands.CMD_TERMINAL_PARAM_OFF_ALLOWED_EMPTY),
            packUI.NewParam(commands.CMD_TERMINAL_PARAM_NAME, commands.CMD_TERMINAL_PARAM_NAME_DISCRIPTION,
                     commands.CMD_TERMINAL_PARAM_NAME_EXAMPLE, commands.CMD_TERMINAL_PARAM_NAME_ISFLAG,
                     commands.CMD_TERMINAL_PARAM_NAME_ALLOWED_EMPTY),
            packUI.NewParam(commands.CMD_TERMINAL_PARAM_ENTITY, commands.CMD_TERMINAL_PARAM_ENTITY_DISCRIPTION,
                     commands.CMD_TERMINAL_PARAM_ENTITY_EXAMPLE, commands.CMD_TERMINAL_PARAM_ENTITY_ISFLAG,
                     commands.CMD_TERMINAL_PARAM_ENTITY_ALLOWED_EMPTY),
            packUI.NewParam(commands.CMD_TERMINAL_PARAM_SYMBOL, commands.CMD_TERMINAL_PARAM_SYMBOL_DISCRIPTION,
                     commands.CMD_TERMINAL_PARAM_SYMBOL_EXAMPLE, commands.CMD_TERMINAL_PARAM_SYMBOL_ISFLAG,
                     commands.CMD_TERMINAL_PARAM_SYMBOL_ALLOWED_EMPTY)},
        commands.CMD_TERMINAL_PARAMS_REQUIRED,
        commands.CMD_TERMINAL_DEFAULT_PARAM,
        cmdTerminalExecute)
    // Комманда "settings"
    frame.Controls[views.CONTROL_COMMANDS].(map[string]*packUI.Command)[commands.CMD_SETTINGS] = packUI.NewCommand(
        commands.CMD_SETTINGS, commands.CMD_SETTINGS_DISCRIPTION,
        []*packUI.Param {
            packUI.NewParam(commands.CMD_SETTINGS_PARAM_LOAD, commands.CMD_SETTINGS_PARAM_LOAD_DISCRIPTION,
                     commands.CMD_SETTINGS_PARAM_LOAD_EXAMPLE, commands.CMD_SETTINGS_PARAM_LOAD_ISFLAG,
                     commands.CMD_SETTINGS_PARAM_LOAD_ALLOWED_EMPTY),
            packUI.NewParam(commands.CMD_SETTINGS_PARAM_SAVE, commands.CMD_SETTINGS_PARAM_SAVE_DISCRIPTION,
                     commands.CMD_SETTINGS_PARAM_SAVE_EXAMPLE, commands.CMD_SETTINGS_PARAM_SAVE_ISFLAG,
                     commands.CMD_SETTINGS_PARAM_SAVE_ALLOWED_EMPTY),
            packUI.NewParam(commands.CMD_SETTINGS_PARAM_DEFAULT, commands.CMD_SETTINGS_PARAM_DEFAULT_DISCRIPTION,
                     commands.CMD_SETTINGS_PARAM_DEFAULT_EXAMPLE, commands.CMD_SETTINGS_PARAM_DEFAULT_ISFLAG,
                     commands.CMD_SETTINGS_PARAM_DEFAULT_ALLOWED_EMPTY),
            packUI.NewParam(commands.CMD_SETTINGS_PARAM_APPLY, commands.CMD_SETTINGS_PARAM_APPLY_DISCRIPTION,
                     commands.CMD_SETTINGS_PARAM_APPLY_EXAMPLE, commands.CMD_SETTINGS_PARAM_APPLY_ISFLAG,
                     commands.CMD_SETTINGS_PARAM_APPLY_ALLOWED_EMPTY)},
        commands.CMD_SETTINGS_PARAMS_REQUIRED,
        commands.CMD_SETTINGS_DEFAULT_PARAM,
        cmdSettingsExecute)
    // Комманда "strategy"
    frame.Controls[views.CONTROL_COMMANDS].(map[string]*packUI.Command)[commands.CMD_STRATEGY] = packUI.NewCommand(
        commands.CMD_STRATEGY, commands.CMD_STRATEGY_DISCRIPTION,
        []*packUI.Param {
            packUI.NewParam(commands.CMD_STRATEGY_PARAM_NAME, commands.CMD_STRATEGY_PARAM_NAME_DISCRIPTION,
                     commands.CMD_STRATEGY_PARAM_NAME_EXAMPLE, commands.CMD_STRATEGY_PARAM_NAME_ISFLAG,
                     commands.CMD_STRATEGY_PARAM_NAME_ALLOWED_EMPTY),
            packUI.NewParam(commands.CMD_STRATEGY_PARAM_PARAM, commands.CMD_STRATEGY_PARAM_PARAM_DISCRIPTION,
                     commands.CMD_STRATEGY_PARAM_PARAM_EXAMPLE, commands.CMD_STRATEGY_PARAM_PARAM_ISFLAG,
                     commands.CMD_STRATEGY_PARAM_PARAM_ALLOWED_EMPTY),
            packUI.NewParam(commands.CMD_STRATEGY_PARAM_VALUE, commands.CMD_STRATEGY_PARAM_VALUE_DISCRIPTION,
                     commands.CMD_STRATEGY_PARAM_VALUE_EXAMPLE, commands.CMD_STRATEGY_PARAM_VALUE_ISFLAG,
                     commands.CMD_STRATEGY_PARAM_VALUE_ALLOWED_EMPTY),
            packUI.NewParam(commands.CMD_STRATEGY_PARAM_ON, commands.CMD_STRATEGY_PARAM_ON_DISCRIPTION,
                     commands.CMD_STRATEGY_PARAM_ON_EXAMPLE, commands.CMD_STRATEGY_PARAM_ON_ISFLAG,
                     commands.CMD_STRATEGY_PARAM_ON_ALLOWED_EMPTY)},
        commands.CMD_STRATEGY_PARAMS_REQUIRED,
        commands.CMD_STRATEGY_DEFAULT_PARAM,
        cmdStrategyExecute)
	//-------------------------------------------------------------------------
	// Устанавливаем текущую консоль
    if settings.ParamsMain[PARAM_MAIN_CONSOLE] == "" {
    	frame.SetMainConsole(views.CONSOLE_LOG)
    } else {
        frame.SetMainConsole(settings.ParamsMain[PARAM_MAIN_CONSOLE])
    }
    frame.RedrawAll()
    // Проверяем ошибки загрузки настроеек
    if err != nil {
        frame.Error(err.Error())
    }
	// Создаем терминал
    terminal, err := market.NewTerminal(frame.ChMsg, frame.ChErr)
    // запоминаем терминал для общего доступа
    frame.Controls["terminal"] = terminal
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
    // Применяем настройки к терминалу
    if !settings.Apply(terminal) {
        frame.Error(MSG_SETTINGS_ERROR_APPLY)
    }
    // Отображаем индикаторы
    frame.Consoles[views.CONSOLE_INDICATORS].Execute(nil)
    frame.Consoles[views.CONSOLE_INDICATORS].Redraw()
    frame.Controls[views.CONTROL_COMMANDS].(map[string]*packUI.Command)[commands.CMD_TERMINAL].Controls["terminal"] = terminal
    // Передаем данные об ошибке терминала в интерфейс
    if err != nil {
        frame.Error(err.Error())
    }
	// Ждем заверщения программы
	for {
		select {
		case <-frame.ChClosed:
			return
		default:
			runtime.Gosched()
		}
	}
}
//////////////////////////////////////////////////////////////////////////////

// ui.Indicators["Прибыль за день"] = "100"
// ui.Indicators["Активных бирж"] = "4"
// ui.Indicators["Активных стратегий"] = "1"
