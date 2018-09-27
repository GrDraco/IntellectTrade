package strategies

import (
    // "fmt"
    "strconv"
    "reflect"

    "../../utilities"
    "../../market/core"
    "../../market/constants"
)

type IStrategy interface {
    Turn(start bool)
    InitChans(chAction chan *core.StrategyAction, chMsg, chErr chan interface{})
    CalculateAction(signal *core.Signal) bool
    GetName() string
    GetIndicators() map[string]string
    GetProperties() map[string]interface{}
    GetKeysForSave() map[string]bool
    GetProperty(name string) interface{}
    SetProperty(name string, value interface{}) bool
    GetPropertySymbol(symbol, property string) interface{}
    SetPropertySymbol(symbol, property string, value interface{}) bool
    On(event string, params []interface{}, callback func(string))
    AddAction(event string, action func(string, []interface{}, func(string)))
    DelAction(event string, action func(string, []interface{}, func(string)))
}
const (
    PROPERTY_NAME = "name"
    PROPERTY_STARTED = "started"
    PROPERTY_SYMBOLS = "conditions_symbol"
    PROPERTY_DEF_SYMBOL = "def"

    MSG_NOT_EXIST_CALCFUNC = "Не инициализирована функция расчета стратегии"
    MSG_SUCCESS_SETPROPERTY = "Успешно установлен параметр"
    MSG_TURNOFF = "Стратения успешно выключена"
    MSG_TURNON = "Стратения успешно включена"
    MSG_ACTION_UPDATEORDER = "Стратегия сформировала действие на обновление ордера"
    MSG_ACTION_CLOSEORDER = "Стратегия сформировала действие на закрытие ордера"
    MSG_ACTION_CREATEORDER = "Стратегия сформировала действие на создание ордера"
)

type BaseStrategy struct {
    // Наследуем события
    utilities.Events
    // Свойства
    Properties map[string]interface{}
    KeysForSave map[string]bool
    Indicators utilities.Collection
    // Каналы
    chAction chan *core.StrategyAction
    chMsg, chErr chan interface{}
    //
    calculate func(signal *core.Signal) bool
}

func (strategy *BaseStrategy) updateOrder(order *core.Order) bool {
    if strategy.chAction != nil {
        strategy.chAction<- core.NewAction(core.STARTEGY_ACTION_UPDATE, order, strategy.Properties)
        strategy.createLog(MSG_ACTION_UPDATEORDER)
        return true
    }
    return false
}

func (strategy *BaseStrategy) closeOrder(order *core.Order) bool {
    if strategy.chAction != nil {
        strategy.chAction<- core.NewAction(core.STARTEGY_ACTION_CLOSE, order, strategy.Properties)
        strategy.createLog(MSG_ACTION_CLOSEORDER)
        return true
    }
    return false
}

func (strategy *BaseStrategy) createOrder(order *core.Order) bool {
    if strategy.chAction != nil {
        strategy.chAction<- core.NewAction(order.Direction, order, strategy.Properties)
        strategy.createLog(MSG_ACTION_CREATEORDER)
        return true
    }
    return false
}

func (strategy *BaseStrategy) createError(msg string) {
    if strategy.chErr != nil {
        message := core.NewError(strategy.GetProperty(PROPERTY_NAME).(string), msg, "")
        message.Exchange = "Strategy"
        strategy.chErr<- message
    }
}

func (strategy *BaseStrategy) createLog(msg string) {
    if strategy.chMsg != nil {
        message := core.NewMessage(strategy.GetProperty(PROPERTY_NAME).(string), msg, "")
        message.Exchange = "Strategy"
        strategy.chMsg<- message
    }
}

func (strategy *BaseStrategy) init(name string) {
    // Инициализируем коллекции
    strategy.Properties = make(map[string]interface{})
    strategy.KeysForSave = make(map[string]bool)
    strategy.Indicators = utilities.Collection { Name: "Indicators" }
    strategy.Indicators.AddAction(utilities.COLLECTION_EVENT_SET_VALUE, func(event string, data []interface{}, callback func(string)) {
        strategy.On(constants.EVENT_SET_INDICATOR, []interface{} { data[1], data[2] }, nil)
    })
    // Запоминаем ключи по которым можно сохранять параметры
    // все те которые не входят в данную коллекцию не будут сохранятся в параметрах программы
    strategy.KeysForSave[PROPERTY_NAME] = true
    strategy.KeysForSave[PROPERTY_STARTED] = true
    strategy.KeysForSave[PROPERTY_SYMBOLS] = true
    // Установка названия стратегии
    strategy.SetProperty(PROPERTY_NAME, name)
    strategy.SetProperty(PROPERTY_STARTED, false)
    // Иницмализация коллекции торговых условий для каждой торговой пары
    // strategy.SetProperty(PROPERTY_SYMBOLS, make(map[string]map[string]interface{}))
}

func (strategy *BaseStrategy) Turn(start bool) {
    // value := strategy.GetProperty(PROPERTY_STARTED).(bool)
    // strategy.SetProperty(PROPERTY_STARTED, !value)
    // if !value {
    //     strategy.createLog(MSG_TURNON)
    // } else {
    //     strategy.createLog(MSG_TURNOFF)
    // }
    strategy.SetProperty(PROPERTY_STARTED, start)
    if start {
        strategy.createLog(MSG_TURNON)
    } else {
        strategy.createLog(MSG_TURNOFF)
    }
}

func (strategy *BaseStrategy) InitChans(chAction chan *core.StrategyAction, chMsg, chErr chan interface{}) {
    strategy.chAction = chAction
    strategy.chMsg = chMsg
    strategy.chErr = chErr
}

func (strategy *BaseStrategy) CalculateAction(signal *core.Signal) bool {
    if strategy.GetProperty(PROPERTY_STARTED).(bool) {
        if strategy.calculate != nil {
            return strategy.calculate(signal)
        } else {
            strategy.createError(MSG_NOT_EXIST_CALCFUNC)
        }
    }
    return false
}

func (strategy *BaseStrategy) GetName() string {
    name := strategy.GetProperty(PROPERTY_NAME)
    if name != nil {
        return name.(string)
    }
    return ""
}

func (strategy *BaseStrategy) GetIndicators() map[string]string {
    return strategy.Indicators.Storage
}

func (strategy *BaseStrategy) GetProperties() map[string]interface{} {
    return strategy.Properties
}

func (strategy *BaseStrategy) GetKeysForSave() map[string]bool {
    return strategy.KeysForSave
}

func (strategy *BaseStrategy) GetProperty(name string) interface{} {
    return strategy.Properties[name]
}

func (strategy *BaseStrategy) SetProperty(name string, value interface{}) bool {
    if strategy.Properties == nil {
        strategy.Properties = make(map[string]interface{})
    }
    strategy.Properties[name] = value
    // Выставляем индикаторы
    if name == PROPERTY_STARTED {
        strategy.Indicators.SetValue(PROPERTY_STARTED, strconv.FormatBool(value.(bool)))
    }
    return true
}

func (strategy *BaseStrategy) GetPropertySymbol(symbol, property string) interface{} {
    if strategy.GetProperty(PROPERTY_SYMBOLS).(map[string]interface{})[symbol] == nil {
        return nil
    }
    return strategy.GetProperty(PROPERTY_SYMBOLS).(map[string]interface{})[symbol].(map[string]interface{})[property]
}

func (strategy *BaseStrategy) SetPropertySymbol(symbol, property string, value interface{}) bool {
    if strategy.GetProperty(PROPERTY_SYMBOLS) == nil {
        strategy.SetProperty(PROPERTY_SYMBOLS, make(map[string]interface{}))
    }
    if strategy.GetProperty(PROPERTY_SYMBOLS).(map[string]interface{})[symbol] == nil {
        strategy.GetProperty(PROPERTY_SYMBOLS).(map[string]interface{})[symbol] = make(map[string]interface{})
    }
    strategy.GetProperty(PROPERTY_SYMBOLS).(map[string]interface{})[symbol].(map[string]interface{})[property] = value
    // Выставляем индикаторы
    indicator := ""
    switch reflect.TypeOf(value).Kind() {
    case reflect.Int:
        indicator = strconv.FormatInt(int64(value.(int)), 10)
    case reflect.Float64:
        indicator = strconv.FormatFloat(value.(float64), 'f', -1, 64)
    case reflect.String:
        indicator = value.(string)
    }
    if indicator != "" {
        strategy.Indicators.SetValue(symbol + "_" + property, indicator)
    }
    strategy.createLog(MSG_SUCCESS_SETPROPERTY + ": " + symbol + "." + property)
    return true
}
