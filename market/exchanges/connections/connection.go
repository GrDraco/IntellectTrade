/*
    Connection это класс для связи с сервером биржи
    он умеет связываться по разным каналам, пока это
    1. api
    2. websocket
*/
package connections

import (
    // "fmt"
    // "strconv"
    "strings"
    "runtime"
    "time"
    "../../core"
    "../../constants"
    "../../../utilities"
)

type IConnection interface {
    Start() bool
    Stop() bool
    GetName() string
    GetStatus() string
    GetManifest() *Manifest
    SetValues(values interface{}) bool
}

const (
    STATUS_STARTED = "started"
    STATUS_STOPED = "stoped"
    STATUS_NOT_ACTIVATED = "not_activated"
    STATUS_DEACTIVATED = "deactivated"
    STATUS_ACTIVATED = "activated"
)

type BaseConnection struct {
    // Свойства
    manifest *Manifest
    timeout float64
    // Флаги
    fWorking bool
    fNewValues bool
    fSuccessStart bool
    fSuccessStop bool
    status string
    // Каналы
    chKill chan bool
    // Базовы функции, которые необходимо реализовать в каждом конекшене
    init func() error
    send func() error
    close func() error
}

func (connection *BaseConnection) createError(msg string) {
    if connection.manifest.ChErr != nil {
        message := core.NewError(connection.manifest.Entity, msg, "")
        message.Exchange = connection.manifest.Exchange
        connection.manifest.ChErr<- message
    }
}

func (connection *BaseConnection) createLog(msg string) {
    if connection.manifest.ChMsg != nil {
        message := core.NewMessage(connection.manifest.Entity, msg, "")
        message.Exchange = connection.manifest.Exchange
        connection.manifest.ChMsg<- message
    }
}

// Без этого метода базовый функционал не будет работать, он запускает горутину и работает пока не будет вызван метод деактивации
func (connection *BaseConnection) activate(manifest *Manifest) bool {
    connection.manifest = manifest
    connection.chKill = make(chan bool)
    connection.fWorking = false
    connection.fNewValues = false
    connection.fSuccessStart = false
    connection.fSuccessStop = false
    // Запускаем горутину
    finishActivate := false
    go func() {
        var started time.Time
        first := true
        for {
            // fmt.Println("connection.status", connection.status)
            select {
            case <-connection.chKill:
                if err := connection.close(); err != nil {
                    connection.createError(err.Error())
                    return
                }
                connection.createLog(connection.manifest.Messages["CONNECTION_DEACTIVATED"])
                return
            default:
                if connection.fNewValues {
                    // Предварительная инициализация запроса если успешно
                    if err := connection.init(); err != nil {
                        connection.fNewValues = false
                        connection.createError(err.Error())
                    } else {
                        params := utilities.ValuesToString(utilities.GetValues(connection.manifest.RequestJSON))
                        connection.createLog(strings.Replace(constants.MSG_SET_PARAMS, constants.MSG_PLACE_PARAMS, params, 1))
                        connection.fNewValues = false
                    }
                }
                // Если флаг в состоянии работает то выполняем метод отправки
                if connection.fWorking {
                    if first || connection.manifest.IsTiming(started) {
                        // Когда производится инициализация send надо в нем происать вызов метода api.manifest.Convertation()
                        started = time.Now()
                        if err := connection.send(); err != nil {
                            // В случае ошибки все останавливаем
                            connection.fWorking = false
                            connection.fSuccessStart = false
                            connection.createError(err.Error())
                            } else {
                                connection.manifest.Response.Ping = time.Now().Sub(started).Nanoseconds()/1000000
                                if connection.fSuccessStart {
                                    // Выполняется единоразово после старта
                                    connection.status = STATUS_STARTED
                                    first = false
                                    connection.fSuccessStart = false
                                    connection.fSuccessStop = true
                                    connection.createLog(connection.manifest.Messages["CONNECTION_STARTED"])
                                }
                            }
                    }
                } else {
                    if connection.fSuccessStop {
                        // Выполняется единоразово после стопа
                        connection.status = STATUS_STOPED
                        connection.fSuccessStart = true
                        connection.fSuccessStop = false
                        connection.createLog(connection.manifest.Messages["CONNECTION_STOPED"])
                    }
                }
                // Передача ресурсов другим горутинам
                runtime.Gosched()
            }
            if connection.status == "" {
                connection.status = STATUS_ACTIVATED
                finishActivate = true
                connection.createLog(connection.manifest.Messages["CONNECTION_ACTIVATED"])
            }
        }
    }()
    if connection.stopUntilSuccess(&finishActivate, false, connection.manifest.Messages["CONNECTION_NOT_ACTIVATED"]) {
        return true
    }
    return false
}

// После деактивации не рекоммендуется активировать т.к. возникают коллизии, необходимо пересоздать конекшен заново
func (connection *BaseConnection) deactivate() {
    if connection.chKill != nil {
        close(connection.chKill)
    }
}

func (connection *BaseConnection) stopUntilSuccess(flag *bool, resultInverse bool, msgFailed string) bool {
    timeStart := time.Now()
    // fmt.Println(msgFailed, *flag)
    for {
        if resultInverse {
            if *flag == false {
                // fmt.Println(msgFailed, *flag)
                return true
            }
        } else {
            if *flag == true {
                // fmt.Println(msgFailed, *flag)
                return true
            }
        }
        timeCurrent := time.Now();
        if timeCurrent.Sub(timeStart).Seconds() >= connection.manifest.TimeoutEntity {
            connection.createError(constants.MSG_TIMEOUT + " " + msgFailed)
            return false
        }
        runtime.Gosched()
    }
}

func (connection *BaseConnection) SetValues(values interface{}) bool {
    // connection.createLog("connection.status " + connection.status)
    // Проверяем разрешенность установки значений
    if connection.GetStatus() == STATUS_NOT_ACTIVATED {
        connection.createError(connection.manifest.Messages["CONNECTION_NOT_ACTIVATED"])
        return false
    }
    connection.manifest.RequestJSON = utilities.ReplaceValues(connection.manifest.RequestJSON, values)
    connection.fNewValues = true
    params := utilities.ValuesToString(utilities.GetValues(connection.manifest.RequestJSON))
    if connection.stopUntilSuccess(&connection.fNewValues, true, strings.Replace(constants.MSG_PARAMS_NOT_SET, constants.MSG_PLACE_PARAMS, params, 1)) {
        return true
    }
    return false
}

func (connection *BaseConnection) Start() bool {
    // connection.createLog(connection.status)
    time.Sleep(time.Duration(1)*time.Second)
    if connection.status == STATUS_STARTED {
        return true
    }
    // Проверяем разрешенность старта
    if connection.status != STATUS_ACTIVATED &&
       connection.status != STATUS_STOPED {
           connection.createError(connection.manifest.Messages["CONNECTION_NOT_ACTIVATED"])
           return false
    }
    // Если требуется регулярное получение данных то зацикоиваем чтение данных
    if connection.manifest.Regular {
        connection.fWorking = true
        connection.fSuccessStart = true
        if connection.stopUntilSuccess(&connection.fSuccessStop, false, connection.manifest.Messages["CONNECTION_NOT_STARTED"]) {
            return true
        }
        return false
    } else {
        if err := connection.send(); err != nil {
            connection.createError(err.Error())
            return false
        }
    }
    return true
}

func (connection *BaseConnection) Stop() bool {
    // Проверяем разрешенность стопа
    if connection.status == STATUS_STOPED ||
       connection.status == STATUS_ACTIVATED {
        return true
    }
    if connection.status != STATUS_STARTED {
        return false
    }
    connection.fWorking = false
    if connection.stopUntilSuccess(&connection.fSuccessStart, false, connection.manifest.Messages["CONNECTION_NOT_STOPED"]) {
        return true
    }
    return false
}

func (connection *BaseConnection) GetName() string {
    return connection.manifest.Name
}

func (connection *BaseConnection) GetStatus() string {
    if connection.status == "" {
        return STATUS_NOT_ACTIVATED
    }
    return connection.status
}

func (connection *BaseConnection) GetManifest() *Manifest {
    return connection.manifest
}

func NewConnection(manifest *Manifest) IConnection {
    // Создаем провайдер данных
    switch strings.ToLower(manifest.Provider) {
    case constants.CONNECTION_WEBSOCKET:
        return NewWSocket(manifest)
    case constants.CONNECTION_API:
        return NewApi(manifest)
    default:
        return nil
    }
}
