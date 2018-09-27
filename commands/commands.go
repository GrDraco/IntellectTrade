package commands

const (
    // Разделители
    TEXT_SEPARATOR_CMD = " "
    TEXT_SEPARATOR_CMD_PARAM = "--"
    TEXT_SEPARATOR_CMD_PARAM_VALUE = ":"
    // Значение флага
    CMD_VALUE_FLAG = "isFlag"
    CMD_VALUE_EMPTY = "isFlag"
    // Общий параметр для всех комманд
    CMD_ALL_PARAM_HELP = "help"
    CMD_ALL_PARAM_HELP_DISCRIPTION = ""
    CMD_ALL_PARAM_HELP_EXAMPLE = "Команда " + TEXT_SEPARATOR_CMD_PARAM + CMD_ALL_PARAM_HELP
    CMD_ALL_PARAM_HELP_ISFLAG = true
    CMD_ALL_PARAM_HELP_ALLOWED_EMPTY = false

    // CMD console
    CMD_CONSOLE = "console"
    CONSOLE_TEXT_START_CMD = CMD_CONSOLE + TEXT_SEPARATOR_CMD
    CMD_CONSOLE_DISCRIPTION = "Команда работы с консолями (окнами)"
    CMD_CONSOLE_PARAMS_REQUIRED = true
    CMD_CONSOLE_DEFAULT_PARAM = "run"
    // PARAMS console
    // EXAMPLE console --run:название_консоли
    CMD_CONSOLE_PARAM_RUN = "run"
    CMD_CONSOLE_PARAM_RUN_DISCRIPTION = "параметр указывающий какую консоль запустить"
    CMD_CONSOLE_PARAM_RUN_EXAMPLE = CONSOLE_TEXT_START_CMD + TEXT_SEPARATOR_CMD_PARAM + CMD_CONSOLE_PARAM_RUN + TEXT_SEPARATOR_CMD_PARAM_VALUE + "название_консоли"
    CMD_CONSOLE_PARAM_RUN_ISFLAG = false
    CMD_CONSOLE_PARAM_RUN_ALLOWED_EMPTY = false
    //---------------------------------------
    // CMD help
    CMD_HELP = "help"
    HELP_TEXT_START_CMD = CMD_HELP + TEXT_SEPARATOR_CMD
    CMD_HELP_DISCRIPTION = "Список комманд"
    CMD_HELP_PARAMS_REQUIRED = false
    CMD_HELP_DEFAULT_PARAM = ""
    // PARAMS help
    // EXAMPLE help cmd:название_команды
    CMD_HELP_PARAM_CMD = "cmd"
    CMD_HELP_PARAM_CMD_DISCRIPTION = ""
    CMD_HELP_PARAM_CMD_EXAMPLE = HELP_TEXT_START_CMD + TEXT_SEPARATOR_CMD_PARAM + CMD_HELP_PARAM_CMD + TEXT_SEPARATOR_CMD_PARAM_VALUE + "название_команды"
    CMD_HELP_PARAM_CMD_ISFLAG = false
    CMD_HELP_PARAM_CMD_ALLOWED_EMPTY = false


// Свои команды
    // CMD terminal
    CMD_TERMINAL = "terminal"
    TERMINAL_TEXT_START_CMD = CMD_TERMINAL + TEXT_SEPARATOR_CMD
    CMD_TERMINAL_DISCRIPTION = "Команда управления терминалом и биржами с которыми он умеет работать"
    CMD_TERMINAL_PARAMS_REQUIRED = true
    CMD_TERMINAL_DEFAULT_PARAM = ""
    // PARAMS terminal
    // EXAMPLE terminal --on --name:название_биржи --entity:тип_сигнала
    CMD_TERMINAL_PARAM_ON = "on"
    CMD_TERMINAL_PARAM_ON_DISCRIPTION = "Параметр старта биржи (после комманды с этим параметром будут поступать сигналы от указанной биржы и типа сигнала). "
    CMD_TERMINAL_PARAM_ON_EXAMPLE = TERMINAL_TEXT_START_CMD + TEXT_SEPARATOR_CMD_PARAM + CMD_TERMINAL_PARAM_ON +
                                    TEXT_SEPARATOR_CMD_PARAM + CMD_TERMINAL_PARAM_NAME + TEXT_SEPARATOR_CMD_PARAM_VALUE + "название_биржи " +
                                    TEXT_SEPARATOR_CMD_PARAM + CMD_TERMINAL_PARAM_ENTITY + TEXT_SEPARATOR_CMD_PARAM_VALUE + "тип_сигнала"
    CMD_TERMINAL_PARAM_ON_ISFLAG = true
    CMD_TERMINAL_PARAM_ON_ALLOWED_EMPTY = false
    // EXAMPLE terminal --off --name:название_биржи --entity:тип_сигнала
    CMD_TERMINAL_PARAM_OFF = "off"
    CMD_TERMINAL_PARAM_OFF_DISCRIPTION = "Параметр остановки биржи (после комманды с этим параметром будут поступать сигналы от указанной биржы и типа сигнала). "
    CMD_TERMINAL_PARAM_OFF_EXAMPLE = TERMINAL_TEXT_START_CMD + TEXT_SEPARATOR_CMD_PARAM + CMD_TERMINAL_PARAM_OFF +
                                    TEXT_SEPARATOR_CMD_PARAM + CMD_TERMINAL_PARAM_NAME + TEXT_SEPARATOR_CMD_PARAM_VALUE + "название_биржи " +
                                    TEXT_SEPARATOR_CMD_PARAM + CMD_TERMINAL_PARAM_ENTITY + TEXT_SEPARATOR_CMD_PARAM_VALUE + "тип_сигнала"
    CMD_TERMINAL_PARAM_OFF_ISFLAG = true
    CMD_TERMINAL_PARAM_OFF_ALLOWED_EMPTY = false
    //
    CMD_TERMINAL_PARAM_ENTITY = "entity"
    CMD_TERMINAL_PARAM_ENTITY_DISCRIPTION = "Сипользуется в связке с другими параметрами, передает название сигнала"
    CMD_TERMINAL_PARAM_ENTITY_EXAMPLE = "... " + TEXT_SEPARATOR_CMD_PARAM + CMD_TERMINAL_PARAM_ENTITY + TEXT_SEPARATOR_CMD_PARAM_VALUE + "тип_сигнала ..."
    CMD_TERMINAL_PARAM_ENTITY_ISFLAG = false
    CMD_TERMINAL_PARAM_ENTITY_ALLOWED_EMPTY = false
    //
    CMD_TERMINAL_PARAM_SYMBOL = "symbol"
    CMD_TERMINAL_PARAM_SYMBOL_DISCRIPTION = "Параметр активирования зевчения торговой пары"
    CMD_TERMINAL_PARAM_SYMBOL_EXAMPLE = TERMINAL_TEXT_START_CMD + TEXT_SEPARATOR_CMD_PARAM + CMD_TERMINAL_PARAM_NAME + TEXT_SEPARATOR_CMD_PARAM_VALUE + "название_биржи " +
                                      TEXT_SEPARATOR_CMD_PARAM + CMD_TERMINAL_PARAM_ENTITY + TEXT_SEPARATOR_CMD_PARAM_VALUE + "тип_сигнала " +
                                      TEXT_SEPARATOR_CMD_PARAM + CMD_TERMINAL_PARAM_SYMBOL + TEXT_SEPARATOR_CMD_PARAM_VALUE + "орговая_пара"
    CMD_TERMINAL_PARAM_SYMBOL_ISFLAG = false
    CMD_TERMINAL_PARAM_SYMBOL_ALLOWED_EMPTY = false
    // EXAMPLE terminal --name:название_биржи --entity:тип_сигнала --symbol:торговая_пара
    CMD_TERMINAL_PARAM_NAME = "name"
    CMD_TERMINAL_PARAM_NAME_DISCRIPTION = "Сипользуется в связке с другими параметрами, передает название биржи"
    CMD_TERMINAL_PARAM_NAME_EXAMPLE = "... " + TEXT_SEPARATOR_CMD_PARAM + CMD_TERMINAL_PARAM_NAME + TEXT_SEPARATOR_CMD_PARAM_VALUE + "название_биржи ..."
    CMD_TERMINAL_PARAM_NAME_ISFLAG = false
    CMD_TERMINAL_PARAM_NAME_ALLOWED_EMPTY = false

    // CMD settings
    CMD_SETTINGS = "settings"
    SETTINGS_TEXT_START_CMD = CMD_SETTINGS + TEXT_SEPARATOR_CMD
    CMD_SETTINGS_DISCRIPTION = "Команда управления настройками приложения"
    CMD_SETTINGS_PARAMS_REQUIRED = true
    CMD_SETTINGS_DEFAULT_PARAM = "load"
    // PARAMS settings
    // EXAMPLE settings --load:название_файла_настроек
    CMD_SETTINGS_PARAM_LOAD = "load"
    CMD_SETTINGS_PARAM_LOAD_DISCRIPTION = "Параметр загргузка все параметров приложения из указанного имени настроек"
    CMD_SETTINGS_PARAM_LOAD_EXAMPLE = SETTINGS_TEXT_START_CMD + TEXT_SEPARATOR_CMD_PARAM + CMD_SETTINGS_PARAM_LOAD +
                                      TEXT_SEPARATOR_CMD_PARAM_VALUE + "название_файла_настроек"
    CMD_SETTINGS_PARAM_LOAD_ISFLAG = false
    CMD_SETTINGS_PARAM_LOAD_ALLOWED_EMPTY = false
    //
    // EXAMPLE settings --save:название_файла_настроек
    CMD_SETTINGS_PARAM_SAVE = "save"
    CMD_SETTINGS_PARAM_SAVE_DISCRIPTION = "Параметр сохранения всех параметров приложения в указанный файл настроек"
    CMD_SETTINGS_PARAM_SAVE_EXAMPLE = SETTINGS_TEXT_START_CMD + TEXT_SEPARATOR_CMD_PARAM + CMD_SETTINGS_PARAM_SAVE +
                                      TEXT_SEPARATOR_CMD_PARAM_VALUE + "название_файла_настроек"
    CMD_SETTINGS_PARAM_SAVE_ISFLAG = false
    CMD_SETTINGS_PARAM_SAVE_ALLOWED_EMPTY = true
    //
    // EXAMPLE settings --default
    CMD_SETTINGS_PARAM_DEFAULT = "default"
    CMD_SETTINGS_PARAM_DEFAULT_DISCRIPTION = "Параметр установки текущего набора параметров как параметры по умолчанию при загрузке программы."
    CMD_SETTINGS_PARAM_DEFAULT_EXAMPLE = SETTINGS_TEXT_START_CMD + TEXT_SEPARATOR_CMD_PARAM + CMD_SETTINGS_PARAM_DEFAULT
    CMD_SETTINGS_PARAM_DEFAULT_ISFLAG = true
    CMD_SETTINGS_PARAM_DEFAULT_ALLOWED_EMPTY = false
    //
    // EXAMPLE settings --apply
    CMD_SETTINGS_PARAM_APPLY = "apply"
    CMD_SETTINGS_PARAM_APPLY_DISCRIPTION = "Параметр применения параметров записанных в текущей настройке."
    CMD_SETTINGS_PARAM_APPLY_EXAMPLE = SETTINGS_TEXT_START_CMD + TEXT_SEPARATOR_CMD_PARAM + CMD_SETTINGS_PARAM_APPLY
    CMD_SETTINGS_PARAM_APPLY_ISFLAG = true
    CMD_SETTINGS_PARAM_APPLY_ALLOWED_EMPTY = false

    // CMD strategy
    CMD_STRATEGY = "strategy"
    STRATEGY_TEXT_START_CMD = CMD_STRATEGY + TEXT_SEPARATOR_CMD
    CMD_STRATEGY_DISCRIPTION = "Команда управления настройками стратегии"
    CMD_STRATEGY_PARAMS_REQUIRED = true
    CMD_STRATEGY_DEFAULT_PARAM = ""
    // PARAMS strategy
    // EXAMPLE strategy --name:название_стратегии --param:название_свойства --value:значение_свойства
    CMD_STRATEGY_PARAM_NAME = "name"
    CMD_STRATEGY_PARAM_NAME_DISCRIPTION = "Параметр для указания имени стратегии по отношению к которой будут применены параметры"
    CMD_STRATEGY_PARAM_NAME_EXAMPLE = STRATEGY_TEXT_START_CMD + TEXT_SEPARATOR_CMD_PARAM + CMD_STRATEGY_PARAM_NAME + TEXT_SEPARATOR_CMD_PARAM_VALUE + "название_стратегии " +
                                      TEXT_SEPARATOR_CMD_PARAM + CMD_STRATEGY_PARAM_PARAM + TEXT_SEPARATOR_CMD_PARAM_VALUE + "название_свойства " +
                                      TEXT_SEPARATOR_CMD_PARAM + CMD_STRATEGY_PARAM_VALUE + TEXT_SEPARATOR_CMD_PARAM_VALUE + "значение_свойства"
    CMD_STRATEGY_PARAM_NAME_ISFLAG = false
    CMD_STRATEGY_PARAM_NAME_ALLOWED_EMPTY = false
    // EXAMPLE strategy --name:название_стратегии --symbol:торговая_пара --param:название_свойства --value:значение_свойства
    CMD_STRATEGY_PARAM_SYMBOL = "symbol"
    CMD_STRATEGY_PARAM_SYMBOL_DISCRIPTION = "Параметр указания для какой торговой пары будут указыватся параметры"
    CMD_STRATEGY_PARAM_SYMBOL_EXAMPLE = STRATEGY_TEXT_START_CMD + TEXT_SEPARATOR_CMD_PARAM + CMD_STRATEGY_PARAM_NAME + TEXT_SEPARATOR_CMD_PARAM_VALUE + "название_стратегии " +
                                      TEXT_SEPARATOR_CMD_PARAM + CMD_STRATEGY_PARAM_SYMBOL + TEXT_SEPARATOR_CMD_PARAM_VALUE + "торговая_пара " +
                                      TEXT_SEPARATOR_CMD_PARAM + CMD_STRATEGY_PARAM_PARAM + TEXT_SEPARATOR_CMD_PARAM_VALUE + "название_свойства " +
                                      TEXT_SEPARATOR_CMD_PARAM + CMD_STRATEGY_PARAM_VALUE + TEXT_SEPARATOR_CMD_PARAM_VALUE + "значение_свойства"
    CMD_STRATEGY_PARAM_SYMBOL_ISFLAG = false
    CMD_STRATEGY_PARAM_SYMBOL_ALLOWED_EMPTY = false
    // EXAMPLE strategy --name:название_стратегии --param:название_свойства --value:значение_свойства
    CMD_STRATEGY_PARAM_PARAM = "param"
    CMD_STRATEGY_PARAM_PARAM_DISCRIPTION = "Параметр указания названия свойства стратегии к которой необходимо установить значение"
    CMD_STRATEGY_PARAM_PARAM_EXAMPLE = STRATEGY_TEXT_START_CMD + TEXT_SEPARATOR_CMD_PARAM + CMD_STRATEGY_PARAM_NAME + TEXT_SEPARATOR_CMD_PARAM_VALUE + "название_стратегии " +
                                       TEXT_SEPARATOR_CMD_PARAM + CMD_STRATEGY_PARAM_SYMBOL + TEXT_SEPARATOR_CMD_PARAM_VALUE + "торговая_пара " +
                                       TEXT_SEPARATOR_CMD_PARAM + CMD_STRATEGY_PARAM_PARAM + TEXT_SEPARATOR_CMD_PARAM_VALUE + "название_свойства " +
                                       TEXT_SEPARATOR_CMD_PARAM + CMD_STRATEGY_PARAM_VALUE + TEXT_SEPARATOR_CMD_PARAM_VALUE + "значение_свойства"
    CMD_STRATEGY_PARAM_PARAM_ISFLAG = false
    CMD_STRATEGY_PARAM_PARAM_ALLOWED_EMPTY = false
    //
    // EXAMPLE strategy --name:название_стратегии --param:название_свойства --value:значение_свойства
    CMD_STRATEGY_PARAM_VALUE = "value"
    CMD_STRATEGY_PARAM_VALUE_DISCRIPTION = "Параметр для указания значения свойства стратегии"
    CMD_STRATEGY_PARAM_VALUE_EXAMPLE = STRATEGY_TEXT_START_CMD + TEXT_SEPARATOR_CMD_PARAM + CMD_STRATEGY_PARAM_NAME + TEXT_SEPARATOR_CMD_PARAM_VALUE + "название_стратегии " +
                                       TEXT_SEPARATOR_CMD_PARAM + CMD_STRATEGY_PARAM_SYMBOL + TEXT_SEPARATOR_CMD_PARAM_VALUE + "торговая_пара " +
                                       TEXT_SEPARATOR_CMD_PARAM + CMD_STRATEGY_PARAM_PARAM + TEXT_SEPARATOR_CMD_PARAM_VALUE + "название_свойства " +
                                       TEXT_SEPARATOR_CMD_PARAM + CMD_STRATEGY_PARAM_VALUE + TEXT_SEPARATOR_CMD_PARAM_VALUE + "значение_свойства"
    CMD_STRATEGY_PARAM_VALUE_ISFLAG = false
    CMD_STRATEGY_PARAM_VALUE_ALLOWED_EMPTY = false
    //
    // EXAMPLE strategy --name:название_стратегии --on
    CMD_STRATEGY_PARAM_ON = "on"
    CMD_STRATEGY_PARAM_ON_DISCRIPTION = "Параметр запуска работы стратегии."
    CMD_STRATEGY_PARAM_ON_EXAMPLE = STRATEGY_TEXT_START_CMD + TEXT_SEPARATOR_CMD_PARAM + CMD_STRATEGY_PARAM_NAME + TEXT_SEPARATOR_CMD_PARAM_VALUE + "название_стратегии " +
                                    TEXT_SEPARATOR_CMD_PARAM + CMD_STRATEGY_PARAM_ON
    CMD_STRATEGY_PARAM_ON_ISFLAG = true
    CMD_STRATEGY_PARAM_ON_ALLOWED_EMPTY = false
    //
    // EXAMPLE strategy --name:название_стратегии --off
    CMD_STRATEGY_PARAM_OFF = "off"
    CMD_STRATEGY_PARAM_OFF_DISCRIPTION = "Параметр остановки работы стратегии."
    CMD_STRATEGY_PARAM_OFF_EXAMPLE = STRATEGY_TEXT_START_CMD + TEXT_SEPARATOR_CMD_PARAM + CMD_STRATEGY_PARAM_NAME + TEXT_SEPARATOR_CMD_PARAM_VALUE + "название_стратегии " +
                                    TEXT_SEPARATOR_CMD_PARAM + CMD_STRATEGY_PARAM_OFF
    CMD_STRATEGY_PARAM_OFF_ISFLAG = true
    CMD_STRATEGY_PARAM_OFF_ALLOWED_EMPTY = false
    // CMD app
    CMD_APP = "app"
    APP_TEXT_START_CMD = CMD_STRATEGY + TEXT_SEPARATOR_CMD
    CMD_APP_DISCRIPTION = "Команда управления окном приложения"
    CMD_APP_PARAMS_REQUIRED = true
    CMD_APP_DEFAULT_PARAM = "redraw"
    // PARAMS app
    // EXAMPLE app --redraw
    CMD_APP_PARAM_REDRAW = "redraw"
    CMD_APP_PARAM_REDRAW_DISCRIPTION = "Параметр перезаписи экрана, это команда пмагает извабится от артифактов на экране"
    CMD_APP_PARAM_REDRAW_EXAMPLE = APP_TEXT_START_CMD + TEXT_SEPARATOR_CMD_PARAM + CMD_APP_PARAM_REDRAW
    CMD_APP_PARAM_REDRAW_ISFLAG = true
    CMD_APP_PARAM_REDRAW_ALLOWED_EMPTY = false
)
/*
    app redraw
    terminal --name:kucoin --entity:tick --symbol:BTC-USDT
    terminal --on --name:kucoin --entity:tick

    terminal --name:kucoin --entity:depth --symbol:BTC-USDT
    terminal --on --name:kucoin --entity:depth

    terminal --name:hitbtc --entity:tick --symbol:ETHBTC
    terminal --on --name:hitbtc --entity:tick

    terminal --name:hitbtc --entity:depth --symbol:ETHBTC
    terminal --on --name:hitbtc --entity:depth

    terminal --name:hitbtc --entity:candle --symbol:ETHBTC
    terminal --on --name:hitbtc --entity:candle

    strategy --name:arbitrage --on
    strategy --name:arbitrage --off
    strategy --name:arbitrage --symbol:def --param:index_depth --value:1
    strategy --name:arbitrage --symbol:def --param:limit --value:0
*/
