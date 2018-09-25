package core

import (
    // "fmt"
    "strconv"
    "time"
    "github.com/satori/go.uuid"
    "../../utilities"
    "../../market/constants"
)

const (
    MSG_MESSAGE = "MESSAGE"
    MSG_WORNING = "WORNING"
    MSG_ERROR = "ERROR"
)

// Структура для передачи сообщений
type Message struct {
    Parent string       //
    Exchange string     //
    Message string      // Сообщение
    Description string  // Описание с рекоммендациями
    Kind string         // Цвет отображения сообщение характеризует важность
}

func (message *Message) SchortString() string {
    return "(" + message.Kind + ") " + message.Exchange + "." + message.Parent + ": " + message.Message
}

func (message *Message) FullString() string {
    if message.Description == "" {
        return message.SchortString()
    }
    return message.SchortString() + ", " + message.Description
}

func NewMessage(parent string, message string, description string) *Message {
    return &Message {
        Parent: parent,
        Message: message,
        Description: description,
        Kind: MSG_MESSAGE }
}

func NewWorning(parent string, message string, description string) *Message {
    return &Message {
        Parent: parent,
        Message: message,
        Description: description,
        Kind: MSG_WORNING }
}

func NewError(parent string, message string, description string) *Message {
    return &Message {
        Parent: parent,
        Message: message,
        Description: description,
        Kind: MSG_ERROR }
}

const (
    STARTEGY_ACTION_BUY = "buy"
    STARTEGY_ACTION_SELL = "sell"
    STARTEGY_ACTION_CLOSE = "close"
    STARTEGY_ACTION_UPDATE = "update"
)

// действие по стратегии
type StrategyAction struct {
    Action string                       // действие
    Order *Order                        // Ордер на покупку, продажу или закрытие
    Properties map[string]interface{}   // Свойста переданные из обекта стратегии
}

func NewAction(action string, order *Order, properties map[string]interface{}) *StrategyAction {
    strategyAction := &StrategyAction {
        Action: action,
        Order: order,
        Properties: properties }
    return strategyAction
}

const (
    DIRECTION_BUY = "buy"
    DIRECTION_SELL = "sell"
)

// Ордер
type Order struct {
    Id string
    // Данные отображаемые в стакане
    Price float64       // Цена
    Amount float64      // Кол-во
    Volume float64      // Объем
    // Остальные данные
    Exchange string     //
    Symbol string       // Торговая пара
    Direction string    // Направление
}

func NewOrder(exchange, symbol, direction string, price, amount float64) *Order {
    order := &Order {
        Price: price,
        Amount: amount,
        Exchange: exchange,
        Symbol: symbol,
        Direction: direction }
    order.Id = uuid.Must(uuid.NewV4()).String()
    return order
}
//Свеча
type Candle struct {
    Open float64        // Цена открытия за период
    Close float64       // Цена закрытя за период
    Min float64         // Минимальная цена за период
    Max float64         // Максимальная цена за период
    Volume float64      // Объем за период
    VolumeQuote float64 // ????
    Timestamp string    // Время
}
// Тик
type Tick struct {
    Ask float64         // Цена ask - запрос
    Bid float64         // Цена bid - предложение
    Volume float64      // Объем
    Timestamp string    // Время
}
// Стакан
type Depth struct {
    Asks map[float64]*Order   // Ордер на покупку
    Bids map[float64]*Order   // Ордер на продажу
}
//
func GetOrderByIndex(prices map[float64]*Order, index int) *Order {
    i := 0
    for _, order := range prices {
        if i == index {
            return order
        }
        i++
    }
    return nil
}
// Данные для торговли
type Signal struct {
    Ping int64          // Время запроса, характерно для API в милисекундах
    Connection string   // Имя коннекшена через который пришли данные
    Exchange string     // От какой биржы пришли данные
    Entity string       // Передеаваемая сущьность в свойстве Data
    Symbol string       // Торговая пара
    Data interface {}   // Полученные данные от биржы
    TimeRecd time.Time  // Время формирования сигнала
    TimeOut bool        // Превышено ли время допустимого ожидания
    Speed int64         // Время за которое сигнал пришол от формирования до передачи конечному потребителю
}

func (signal *Signal) calcSpeed() {
    if signal.Speed == 0 {
        signal.Speed = time.Now().Sub(signal.TimeRecd).Nanoseconds()
    }
}

func (signal *Signal) Tick() *Tick {
    signal.calcSpeed()
    if signal.Entity != constants.ENTITY_TICK {
        return nil
    }
    return signal.Data.(*Tick)
}

func (signal *Signal) TickStr(showSpeeds bool) string {
    signal.calcSpeed()
    if signal.Tick() == nil {
        return ""
    }
    if showSpeeds {
        return "(SIGNAL) " + signal.Exchange + "." + signal.Entity + ": " +
               // "Speed = " + strconv.FormatInt(signal.Speed/1000, 10) +
               "Ping = " + strconv.FormatInt(signal.Ping, 10) +
               "ms TimeOut = " + strconv.FormatBool(signal.TimeOut) +
               " Symbol = " + signal.Symbol +
               " Bid = " + utilities.FloatToString(signal.Tick().Bid, 8) +
               " Ask = " + utilities.FloatToString(signal.Tick().Ask, 8)
    }
    return "(SIGNAL) " + signal.Exchange + "." + signal.Entity +
           ": " + signal.Symbol +
           "TimeOut = " + strconv.FormatBool(signal.TimeOut) +
           " Bid = " + utilities.FloatToString(signal.Tick().Bid, 8) +
           " Ask = " + utilities.FloatToString(signal.Tick().Ask, 8)
}

func (signal *Signal) Depth() *Depth {
    signal.calcSpeed()
    if signal.Entity != constants.ENTITY_DEPTH {
        return nil
    }
    return signal.Data.(*Depth)
}

func (signal *Signal) DepthStr(showSpeeds bool) string {
    signal.calcSpeed()
    if signal.Depth() == nil {
        return ""
    }
    ordersBid := ""
    ordersAsk := ""
    i := 0
    for _, order := range signal.Depth().Bids {
        ordersBid += utilities.FloatToString(order.Price, 8) + "(" + utilities.FloatToString(order.Amount, 8) + ") "
        i++
        if i == 2 {
            break
        }
    }
    i = 0
    for _, order := range signal.Depth().Asks {
        ordersAsk += utilities.FloatToString(order.Price, 8) + "(" + utilities.FloatToString(order.Amount, 8) + ") "
        i++
        if i == 2 {
            break
        }
    }
    if showSpeeds {
        return "(SIGNAL) " + signal.Exchange + "." + signal.Entity + ": " +
               "Speed = " + strconv.FormatInt(signal.Speed/1000, 10) +
               "µs Ping = " + strconv.FormatInt(signal.Ping, 10) +
               "ms TimeOut = " + strconv.FormatBool(signal.TimeOut) +
               " Symbol = " + signal.Symbol +
               " 2_Bids = " + ordersBid +
               "2_Asks = " + ordersAsk
    }
    return "(SIGNAL) " + signal.Exchange + "." + signal.Entity +
           ": " + signal.Symbol +
           "TimeOut = " + strconv.FormatBool(signal.TimeOut) +
           " 2_Bids = " + ordersBid +
           "2_Asks = " + ordersAsk
}

func (signal *Signal) Candles() []*Candle {
    signal.calcSpeed()
    if signal.Entity != constants.ENTITY_CANDLE {
        return nil
    }
    return signal.Data.([]*Candle)
}

func (signal *Signal) CandlesStr(showSpeeds bool) string {
    signal.calcSpeed()
    if signal.Candles() == nil {
        return ""
    }
    candleStr := ""
    for _, candle := range signal.Candles() {
        candleStr += utilities.FloatToString(candle.Open, 8) + "(" + utilities.FloatToString(candle.Volume, 8) + ") "
    }
    if showSpeeds {
        return "(SIGNAL) " + signal.Exchange + "." + signal.Entity + ": " +
               "Speed = " + strconv.FormatInt(signal.Speed/1000, 10) +
               "µs Ping = " + strconv.FormatInt(signal.Ping, 10) +
               "ms TimeOut = " + strconv.FormatBool(signal.TimeOut) +
               " Symbol = " + signal.Symbol +
               " Open(Volume) = " + candleStr
    }
    return "(SIGNAL) " + signal.Exchange + "." + signal.Entity +
           ": " + signal.Symbol +
           "TimeOut = " + strconv.FormatBool(signal.TimeOut) +
           " Open(Volume) = " + candleStr
}
