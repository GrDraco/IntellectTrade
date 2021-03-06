/*
    Exchange это общий класс для всех бирж, в котором есть все необходимое
    для работы с данными биржи
*/
package exchanges

import (
    // "errors"
    "strings"
    //"time"
    //"encoding/json"
    // "fmt"
    "../core"
    "../../utilities"
    "../constants"
    "./connections"
    //"./hitBTC"
)

type Exchange struct {
    // Наследуем события
    utilities.Events
    // Свойства
    Name string                                     // Название биржи
    indicators utilities.Collection
    // Коллекции
    connections map[string]map[string]connections.IConnection  // Коллекция соккетов для получения данных с бирж
    // Каналы
    chSignal chan *core.Signal
    chMsg, chErr chan interface{}
}

func NewExchange(manifests []*connections.Manifest, chSignal chan *core.Signal, chMsg , chErr chan interface{}) *Exchange {
    // Выделение памяти биржу
    exchange := new(Exchange)
    // Инициализация всех необходимых каналов связи и данных
    exchange.chSignal = chSignal
    exchange.chMsg = chMsg
    exchange.chErr = chErr
    exchange.connections = make(map[string]map[string]connections.IConnection)
    exchange.indicators = utilities.Collection { Name: "Indicators" }
    exchange.Name = manifests[0].Exchange
    // Инициализация каналов связи с сервером биржи
    for _, manifest := range manifests {
        // Передаем в манифест каналы связи
        manifest.ChSignal = chSignal
        manifest.ChMsg = chMsg
        manifest.ChErr = chErr
        // Создаем конекшен к бирже и запоминаем его в колекцию
        if exchange.connections[strings.ToLower(manifest.Entity)] == nil {
            exchange.connections[strings.ToLower(manifest.Entity)] = make(map[string]connections.IConnection)
        }
        exchange.connections[strings.ToLower(manifest.Entity)][strings.ToLower(manifest.Provider)] = connections.NewConnection(manifest)
    }
    // Заполняем стартовыми показаниями индикаторы
    for _, groupEntity := range exchange.connections {
        for _, connection := range groupEntity {
            exchange.SetIndicator(connection.GetName(), connection.GetStatus())
        }
    }
    // Реагируем на событие установки значениея у индикатора
    exchange.indicators.AddAction(utilities.COLLECTION_EVENT_SET_VALUE, func(event string, data []interface{}, callback func(string)) {
        exchange.On(constants.EVENT_SET_INDICATOR, []interface{} { data[1], data[2] }, nil)
    })
    return exchange
}

func (exchange *Exchange) msg(msg *core.Message) {
    if msg == nil {
        return
    }
    msg.Exchange = exchange.Name
    exchange.chMsg<-msg
}

func (exchange *Exchange) err(err *core.Message) {
    if err == nil {
        return
    }
    err.Exchange = exchange.Name
    exchange.chErr<-err
}

func (exchange *Exchange) CountActiveConnection() (count int64) {
    count = 0
    for _, groupEntity := range exchange.connections {
        for _, connection := range groupEntity {
            if connection.GetStatus() == connections.STATUS_STARTED {
                count++
            }
        }
    }
    return
}

func (exchange *Exchange) Turn(entity, provider string, start bool) (status, success bool) {
    status = false
    entity = strings.ToLower(entity)
    if entity == "" {
        exchange.err(core.NewError("", strings.Replace(constants.MSG_PARAMS_REQUIRED, constants.MSG_PLACE_PARAMS, "entity", 1), ""))
        success = false
        return
    }
    connection := exchange.connections[strings.ToLower(entity)][strings.ToLower(provider)]
    if connection == nil {
        exchange.err(core.NewError("", constants.MSG_CONNECTION_NOT_EXIST, ""))
        success = false
        return
    }
    // exchange.msg(core.NewMessage("connection.GetStatus()", connection.GetStatus(), ""))
    if start {
        if connection.Start() {
            exchange.On(constants.EVENT_STARTED, nil, nil)
            exchange.SetIndicator(connection.GetName(), connection.GetStatus())
            status = true
            success = true
            return
        }
    }
    if !start {
        if connection.Stop() {
            exchange.On(constants.EVENT_STOPED, nil, nil)
            exchange.SetIndicator(connection.GetName(), connection.GetStatus())
            status = false
            success = true
            return
        }
    }
    return
}

func (exchange *Exchange) SetValues(entity, provider string, values interface{}) bool {
    if entity == "" || values == nil {
        exchange.err(core.NewError("", strings.Replace(constants.MSG_PARAMS_REQUIRED, constants.MSG_PLACE_PARAMS, "entity, values", 1), ""))
        return false
    }
    connection := exchange.connections[strings.ToLower(entity)][strings.ToLower(provider)]
    if connection == nil {
        exchange.err(core.NewError("", constants.MSG_CONNECTION_NOT_EXIST, ""))
        return false
    }
    return connection.SetValues(values)
}

func CmdAmswer(cmd string, answer string) string {
    var str string
    str = strings.Replace(constants.CMD_TEMPLATE_ANSWER, constants.CMD_PLACE_CMD, cmd, 1)
    str = strings.Replace(str, constants.CMD_PLACE_ANSWER, answer, 1)
    return str
}
//TODO: Написать манифест по реализации отправки ордера на биржу
func (exchange *Exchange) SendOrder(order *core.Order) bool {
    // connection := exchange.connections[constants.ENTITY_NEW_ORDER]
    // if connection != nil {
    //     //!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
    //     if connection.SetValues(nil) {
    //     //!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
    //         return connection.Start()
    //     }
    // }
    return false
}
//TODO: Написать манифест по реализации реализацию закрытия ордера на биржу
func (exchange *Exchange) CloseOrder(order *core.Order) bool {
    // connection := exchange.connections[constants.ENTITY_CLOSE_ORDER]
    // if connection != nil {
    //     //!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
    //     if connection.SetValues(nil) {
    //     //!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
    //         return connection.Start()
    //     }
    // }
    return false
}
//TODO: Написать манифест по реализации реализацию обновление ордера на биржу
func (exchange *Exchange) UpdateOrder(order *core.Order) bool {
    // connection := exchange.connections[constants.ENTITY_UPDATE_ORDER]
    // if connection != nil {
    //     //!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
    //     if connection.SetValues(nil) {
    //     //!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
    //         return connection.Start()
    //     }
    // }
    return false
}

func (exchange *Exchange) SetIndicator(indicator string, value string) {
    exchange.indicators.SetValue(indicator, value)
}

func (exchange *Exchange) GetIndicators() map[string]string {
    return exchange.indicators.Storage
}

func (exchange *Exchange) GetManifest(entity, provider string) *connections.Manifest {
    return exchange.connections[entity][provider].GetManifest()
}
