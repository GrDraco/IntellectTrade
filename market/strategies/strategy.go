package strategies

import (
    "../../market/core"
)

type IStrategy interface {
    InitChans(chAction chan *core.StrategyAction, chMsg, chErr chan interface{})
    CalculateAction(signal *core.Signal) bool
    GetName() string
    GetProperties() map[string]interface{}
    GetProperty(name string) interface{}
    SetProperty(name string, value interface{}) bool
}

const (
    PROPERTY_NAME = "name"
)

type BaseStrategy struct {
    // Свойства
    Properties map[string]interface{}
    // Каналы
    chAction chan *core.StrategyAction
    chMsg, chErr chan interface{}
}

func (strategy *BaseStrategy) updateOrder(order *core.Order) bool {
    if strategy.chAction != nil {
        strategy.chAction<- core.NewAction(core.STARTEGY_ACTION_UPDATE, order, strategy.Properties)
        return true
    }
    return false
}

func (strategy *BaseStrategy) closeOrder(order *core.Order) bool {
    if strategy.chAction != nil {
        strategy.chAction<- core.NewAction(core.STARTEGY_ACTION_CLOSE, order, strategy.Properties)
        return true
    }
    return false
}

func (strategy *BaseStrategy) createOrder(order *core.Order) bool {
    if strategy.chAction != nil {
        strategy.chAction<- core.NewAction(order.Direction, order, strategy.Properties)
        return true
    }
    return false
}

func (strategy *BaseStrategy) createError(msg, exchange string) {
    if strategy.chErr != nil {
        message := core.NewError(strategy.GetName(), msg, "")
        message.Exchange = exchange
        strategy.chErr<- message
    }
}

func (strategy *BaseStrategy) createLog(msg, exchange string) {
    if strategy.chMsg != nil {
        message := core.NewMessage(strategy.GetName(), msg, "")
        message.Exchange = exchange
        strategy.chMsg<- message
    }
}

func (strategy *BaseStrategy) init(name string) {
    strategy.Properties = make(map[string]interface{})
    strategy.SetProperty(PROPERTY_NAME, name)
}

func (strategy *BaseStrategy) InitChans(chAction chan *core.StrategyAction, chMsg, chErr chan interface{}) {
    strategy.chAction = chAction
    strategy.chMsg = chMsg
    strategy.chErr = chErr
}

func (strategy *BaseStrategy) CalculateAction(signal *core.Signal) bool {
    // ЗАГОТОВКА ДЛЯ РАСЧЕТА СТРАТЕГИИ
    return false
}

func (strategy *BaseStrategy) GetName() string {
    name := strategy.GetProperty(PROPERTY_NAME)
    if name != nil {
        return name.(string)
    }
    return ""
}

func (strategy *BaseStrategy) GetProperties() map[string]interface{} {
    if strategy.Properties == nil {
        return nil
    }
    return strategy.Properties
}

func (strategy *BaseStrategy) GetProperty(name string) interface{} {
    if strategy.Properties == nil {
        return nil
    }
    return strategy.Properties[name]
}

func (strategy *BaseStrategy) SetProperty(name string, value interface{}) bool {
    if strategy.Properties == nil {
        return false
    }
    strategy.Properties[name] = value
    return true
}
