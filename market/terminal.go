package market

import (
    // "fmt"
    "time"
    "runtime"
    "strconv"
    "strings"
    "errors"
    "./core"
    "./exchanges"
    "./exchanges/connections"
    "./constants"
    "./strategies"
    "../utilities"
)

type Terminal struct {
    // Наследуем события
    utilities.Events
    // Свойства
    Exchanges map[string]*exchanges.Exchange
    Manifests map[string][]*connections.Manifest
    Strategies map[string]strategies.IStrategy
    indicators utilities.Collection
    // Канал связи с биржами
    chSignal chan *core.Signal
    // Канал связи со стратегиями
    chAction chan *core.StrategyAction
    // Каналы для передачи сообщений и ошибок
    chMsg, chErr chan interface{}
    // Канал на остановку терминала
    chStop chan bool
}

func NewTerminal(chMsg, chErr chan interface{}) (terminal *Terminal, err error) {
    // Выделяем память для терминала
    terminal = new(Terminal)
    // Инициализируем каналов связи
    terminal.chSignal = make(chan *core.Signal)
    terminal.chAction = make(chan *core.StrategyAction)
    terminal.chMsg = chMsg
    terminal.chErr = chErr
    terminal.chStop = make(chan bool)
    // Инициализируем необходимые свойства
    terminal.Exchanges = make(map[string]*exchanges.Exchange)
    terminal.Strategies = make(map[string]strategies.IStrategy)
    terminal.indicators = utilities.Collection { Name: "Indicators" }
    terminal.indicators.SetValue(constants.INDICATOR_EXCHANGES_COUNT, "0")
    terminal.indicators.SetValue(constants.INDICATOR_ACTIVE_EXCHANGES_COUNT, "0")
    terminal.indicators.AddAction(utilities.COLLECTION_EVENT_SET_VALUE, func(event string, data []interface{}, callback func(string)) {
        terminal.On(constants.EVENT_SET_INDICATOR, []interface{} { data[1], data[2] }, nil)
    })
    // Читаем манифесты
    var _err error
    terminal.Manifests, _err = connections.GetManifests()
    if _err != nil {
        err = errors.New(strings.Replace(constants.MSG_MANIFESTS_ERROR, constants.MSG_PLACE_ERROR, _err.Error(), 1))
        return
	}
    // Запускаем горутину чтения данных от бирж
    go func() {
        for {
            select {
            case <-terminal.chStop:
                return
            case action := <-terminal.chAction:
                var res bool
                switch action.Action {
                case core.STARTEGY_ACTION_BUY:
                    res = terminal.Exchanges[action.Order.Exchange].SendOrder(action.Order)
                case core.STARTEGY_ACTION_SELL:
                    res = terminal.Exchanges[action.Order.Exchange].SendOrder(action.Order)
                case core.STARTEGY_ACTION_CLOSE:
                    res = terminal.Exchanges[action.Order.Exchange].CloseOrder(action.Order)
                case core.STARTEGY_ACTION_UPDATE:
                    res = terminal.Exchanges[action.Order.Exchange].UpdateOrder(action.Order)
                }
                terminal.On(constants.EVENT_NEW_ACTION, []interface{} { action, res }, nil)
            case signal := <-terminal.chSignal:
                terminal.On(constants.EVENT_NEW_SIGNAL, []interface{} { signal }, nil)
                for _, strategy := range terminal.Strategies {
                    strategy.CalculateAction(signal)
                    terminal.On(constants.EVENT_CALCULATE_ACTION, []interface{} { strategy.GetName(), strategy.GetProperties() }, nil)
                }
            case <-time.After(time.Second):
                chErr<- core.NewError("Terminal", "Превышено время ожидания бирж", "")
            default:
                runtime.Gosched()
            }
        }
    }()
    // Инициализируем биржы
    for key, manifests := range terminal.Manifests {
        terminal.Exchanges[key] = exchanges.NewExchange(manifests, terminal.chSignal, chMsg, chErr)
        // Подписываемся на события
        terminal.Exchanges[key].AddAction(constants.EVENT_STARTED, func(event string, params []interface{}, callback func(string)) {
            terminal.calculateActiveExchanges()
        })
        //
        terminal.Exchanges[key].AddAction(constants.EVENT_STOPED, func(event string, params []interface{}, callback func(string)) {
            terminal.calculateActiveExchanges()
        })
        //
        terminal.Exchanges[key].AddAction(constants.EVENT_SET_INDICATOR, func(event string, data []interface{}, callback func(string)) {
            terminal.On(constants.EVENT_SET_INDICATOR, []interface{} { data[0], data[1] }, nil)
        })
    }
    terminal.SetIndicator(constants.INDICATOR_EXCHANGES_COUNT, strconv.FormatInt(int64(len(terminal.Exchanges)), 10))
    return
}

func (terminal *Terminal) calculateActiveExchanges() {
    var count int64
    count = 0
    for name, exchanges := range terminal.Exchanges {
        if exchanges.CountActiveConnection() > 0 {
            terminal.SetIndicator(name, "active")
            count++
        } else {
            terminal.SetIndicator(name, "inactive")
        }
    }
    terminal.SetIndicator(constants.INDICATOR_ACTIVE_EXCHANGES_COUNT, strconv.FormatInt(count, 10))
}

func (terminal *Terminal) AddStrategy(strategy strategies.IStrategy) {
    strategy.InitChans(terminal.chAction, terminal.chMsg, terminal.chErr)
    // Подписываемся на событие изменение индикатора и ретранслируем по событию терминала
    strategy.AddAction(constants.EVENT_SET_INDICATOR, func(event string, params []interface{}, callback func(string))  {
        terminal.On(constants.EVENT_SET_INDICATOR, []interface{} { params }, nil)
    })
    terminal.Strategies[strategy.GetName()] = strategy
}

func (terminal *Terminal) SetIndicator(indicator string, value string) {
    terminal.indicators.SetValue(indicator, value)
}

func (terminal *Terminal) GetIndicators() map[string]string {
    return terminal.indicators.Storage
}
