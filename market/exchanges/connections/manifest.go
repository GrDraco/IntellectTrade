/*
    Manifest это класс с параметрами для биржы и работы их данными
*/
package connections

import (
    "os"
    // "fmt"
    "time"
    "errors"
    "strings"
    "reflect"
    "io/ioutil"
    "encoding/json"
    "github.com/satori/go.uuid"

    "../../constants"
    "../../core"
    "../../../utilities"
)

type Path struct {
    // Используется один из следующих
    //или
    Str []string        `json:"str"`            // Путь по именам
    //или
    Int []int64         `json:"int"`            // Путь по позициям в массиве
    //или
    SubValues []Value   `json:"sub_values"`     // Если значение имеет под уровни
}

type Design struct {
    // Тип определяющий тип данных значения
    Kind reflect.Kind   `json:"kind"`           // 17-Array 6-int64 14-float64 24-string
    Path []Path         `json:"path"`           // Путь к данным
    Value Path          `json:"value"`          // Путь к значению если он пустой значит значение берется напрямую из пути Path
    Check *Value        `json:"check"`
    CheckInt int64      `json:"check_int"`      // Значение для сверки типа данных с полученными для цифры
    CheckStr string     `json:"check_str"`      // Значение для сверки типа данных с полученными для строки
}

type Value struct {
    Name string         `json:"name"`           //
    Design Design       `json:"design"`         //
}

type Failed struct {
    Message []string     `json:"message"`       //
    Description []string `json:"description"`   //
}

type Response struct {
    Parent *Manifest
    // Хранилище данных от биржи
    JSON interface{}
    // Считываются из манифеста json
    SkipResponse int64  `json:"skip_response"`  // Сколько надо пропустить запросов
    Success []string    `json:"success"`        // Путь объекту наличее которого говорит о успешном получении данных
    Failed Failed       `json:"failed"`         // Структура определения ошибок от биржи
    Empty Value         `json:"empty"`          // Структура определения пустого ответа, указывает какое значение и где будет в ответе указывающие на пустой ответ
    IsUpdates bool      `json:"is_updates"`     // Являются ли пришедшие данные обновлениями или они целостностные
    // Values Values       `json:"values"`      //
    Values []Value       `json:"values"`        // Это карта данных лежащих в JSON
    // Рассчитываются в процессе обработки
    Ping int64                                  // Пинг
    Index int64                                 // Подсчет пришедших ответов
}

func (response *Response) skipResponse() (bool, error) {
    var res bool
    // Подсчитываем ответы
    response.Index++
    // Пропускаем если пришел пустой ответ
    if response.Empty.Design.Kind == reflect.Int64 {
        value, err := toInt64(response.JSON, response.Empty)
        if err != nil {
            return true, err
        }
        if value == response.Empty.Design.CheckInt {
            res = true
        }
    }
    // Пропускаем если указзаное количество по пропуску ен достигнуто
    if response.Index <= response.SkipResponse {
        res = true
    }
    if res {
        // if response.Parent.ChMsg != nil {
        //     message := core.NewMessage(response.Parent.Entity, strings.Replace(constants.MSG_CONNECTION_SKIP_RESPONSE, constants.MSG_PLACE_N, utilities.IntToString(response.Index), 1), "")
        //     message.Exchange = response.Parent.Exchange
        //     response.Parent.ChMsg<- message
        // }
        return true, nil
    } else {
        return false, nil
    }
}

func (response *Response) ToError() error {
    err := response.Parent.CheckError()
    if err != nil {
        return err
    }
    var value interface{}
    saccess := true
    if len(response.Success) > 0 {
        value, err = utilities.GetValueByStr(response.JSON, response.Success)
        if err != nil {
            return err
        }
        if value == nil {
            saccess = false
        }
    } else {
        value, err = utilities.GetValueByStr(response.JSON, response.Failed.Message)
        if err != nil {
            return err
        }
        if value != nil {
            saccess = false
        }
    }
    if !saccess {
        var description, message interface{}
        message, err = utilities.GetValueByStr(response.JSON, response.Failed.Message)
        if err != nil {
            return err
        }
        if len(response.Failed.Description) > 0 {
            description, err = utilities.GetValueByStr(response.JSON, response.Failed.Description)
            if err != nil {
                return err
            }
        }
        msgFailed := core.NewError(response.Parent.Entity,
            utilities.ToString(message),
            utilities.ToString(description))

        if msgFailed.Message != "" || msgFailed.Description != "" {
            return errors.New("//-SERVER- " + msgFailed.Message + " " + msgFailed.Description + " -SERVER-//")
        }
    }
    return nil
}

func (response *Response) ToSymbol(obj Value) (string, error) {
    // Определяем тип данных
    switch obj.Design.Kind {
    case reflect.Int64:
        // Получаем значение
        value, err := toInt64(response.JSON, obj)
        if err != nil {
            return "", err
        }
        // Приводим к нуному типу
        return clearSymbol(utilities.IntToString(value)), nil
    case reflect.Float64:
        // Получаем значение
        value, err := toFloat64(response.JSON, obj)
        if err != nil {
            return "", err
        }
        // Приводим к нуному типу
        return clearSymbol(utilities.FloatToString(value, 0)), nil
    case reflect.String:
        // Получаем значение
        value, err := toString(response.JSON, obj)
        if err != nil {
            return "", err
        }
        // Приводим к нуному типу
        return clearSymbol(value), nil
    }
    return "", nil
}

func (response *Response) ToTick() error {
    signal := new(core.Signal)
    err := response.ToError()
    if err != nil {
        return err
    }
    signal = &core.Signal {
        Ping: response.Ping,
        Connection: response.Parent.Name,
        Exchange: response.Parent.Exchange,
        Entity: constants.ENTITY_TICK,
        Symbol: response.Parent.Request.ToSymbol(),
        TimeRecd: time.Now(),
        DataIsUpdates: response.IsUpdates }
    tick := &core.Tick{}
    for _, obj := range response.Values {
        // Определяем тип данных
        switch strings.ToLower(obj.Name) {
        case OBJ_ASK:
            tick.Ask, err = toFloat64(response.JSON, obj)
        case OBJ_BID:
            tick.Bid, err = toFloat64(response.JSON, obj)
        case OBJ_VOLUME:
            tick.Volume, err = toFloat64(response.JSON, obj)
        case OBJ_SYMBOL:
            tick.Symbol, err = toString(response.JSON, obj)
        case OBJ_TIMESTAMP:
            tick.Timestamp, err = toString(response.JSON, obj)
        }
        if err != nil {
            return err
        }
    }
    signal.Data = tick
    if response.Parent.Provider == constants.CONNECTION_API {
        if signal.Ping > response.Parent.TimeoutSignal {
            signal.TimeOut = true
        } else {
            signal.TimeOut = false
        }
    }
    response.Parent.ChSignal<-signal
    return nil
}

func (response *Response) ToOrders(obj Value, symbol string) (map[float64]*core.Order, error) {
    // Проверяем на тип данных
    if obj.Design.Kind != reflect.Array {
        return nil, errors.New("Mismatched type in ToOrders")
    }
    // Получаем массив значений
    values, err := getValueByPath(response.JSON, obj.Design.Path)
    if err != nil {
        return nil, err
    }
    orders := make(map[float64]*core.Order)
    for _, value := range values.([]interface{}) {
        // Если данные не проходят проверку то пропускаем
        var err error
        var checked bool
        if obj.Design.Check != nil {
            checked, err = isChecked(value, *obj.Design.Check)
            if err != nil {
                return nil, err
            }
        } else {
            checked = true
        }
        if !checked {
            continue
        }
        // Получаем значение
        order := &core.Order {}
        for _, valPath := range obj.Design.Value.SubValues {
            switch strings.ToLower(valPath.Name) {
            case OBJ_PRICE:
                order.Price, err = toFloat64(value, valPath)
            case OBJ_AMOUNT:
                order.Amount, err = toFloat64(value, valPath)
            }
            if err != nil {
                return nil, err
            }
        }
        order.Symbol = symbol
        orders[order.Price] = order
    }
    return orders, nil
}

func (response *Response) ToDepth() error {
    signal := new(core.Signal)
    var err error
    err = response.ToError()
    if err != nil {
        return err
    }
    signal = &core.Signal {
        Ping: response.Ping,
        Connection: response.Parent.Name,
        Exchange: response.Parent.Exchange,
        Entity: constants.ENTITY_DEPTH,
        Symbol: response.Parent.Request.ToSymbol(),
        TimeRecd: time.Now(),
        DataIsUpdates: response.IsUpdates }
    signal.Data = new(core.Depth)
    for _, obj := range response.Values {
        // Определяем тип данных
        switch strings.ToLower(obj.Name) {
        case OBJ_ASKS:
            signal.Data.(*core.Depth).Asks, err = response.ToOrders(obj, signal.Symbol)
        case OBJ_BIDS:
            signal.Data.(*core.Depth).Bids, err = response.ToOrders(obj, signal.Symbol)
        }
        // fmt.Println(signal.Data)
        if err != nil {
            return err
        }
    }
    ///////
    if response.Parent.Provider == constants.CONNECTION_API {
        if signal.Ping > response.Parent.TimeoutSignal {
            signal.TimeOut = true
        } else {
            signal.TimeOut = false
        }
    }
    response.Parent.ChSignal<-signal
    return nil
}

func (response *Response) ToCandles() error  {
    signal := new(core.Signal)
    err := response.ToError()
    if err != nil {
        return err
    }
    signal = &core.Signal {
        Ping: response.Ping,
        Connection: response.Parent.Name,
        Exchange: response.Parent.Exchange,
        Entity: constants.ENTITY_CANDLE,
        Symbol: response.Parent.Request.ToSymbol(),
        TimeRecd: time.Now(),
        DataIsUpdates: response.IsUpdates }
    candles := make([]*core.Candle, 0)
    for _, obj := range response.Values {
        // Определяем тип данных
        if strings.ToLower(obj.Name) == OBJ_CANDLES {
            // Проверяем на тип данных
            if obj.Design.Kind != reflect.Array {
                return errors.New("Mismatched type in ToCandles")
            }
            // Получаем массив значений
            values, err := getValueByPath(response.JSON, obj.Design.Path)
            if err != nil {
                return err
            }
            for _, value := range values.([]interface{}) {
                // Если данные не проходят проверку то пропускаем
                var err error
                var checked bool
                if obj.Design.Check != nil {
                    checked, err = isChecked(value, *obj.Design.Check)
                    if err != nil {
                        return err
                    }
                } else {
                    checked = true
                }
                if !checked {
                    continue
                }
                // Получаем значение
                candle := &core.Candle {}
                for _, valPath := range obj.Design.Value.SubValues {
                    switch strings.ToLower(valPath.Name) {
                    case OBJ_OPEN:
                        candle.Open, err = toFloat64(value, valPath)
                    case OBJ_CLOSE:
                        candle.Close, err = toFloat64(value, valPath)
                    case OBJ_MIN:
                        candle.Min, err = toFloat64(value, valPath)
                    case OBJ_MAX:
                        candle.Max, err = toFloat64(value, valPath)
                    case OBJ_VOLUME:
                        candle.Volume, err = toFloat64(value, valPath)
                    case OBJ_VOLUMEQUOTE:
                        candle.VolumeQuote, err = toFloat64(value, valPath)
                    case OBJ_TIMESTAMP:
                        candle.Timestamp, err = toString(value, valPath)
                    }
                    if err != nil {
                        return err
                    }
                }
                candle.Symbol = signal.Symbol
                candles = append(candles, candle)
            }
        }
    }
    signal.Data = candles
    ///////
    if response.Parent.Provider == constants.CONNECTION_API {
        if signal.Ping > response.Parent.TimeoutSignal {
            signal.TimeOut = true
        } else {
            signal.TimeOut = false
        }
    }
    response.Parent.ChSignal<-signal
    return nil
}

const (
    PERIOD_SECONDS = "s"
    PERIOD_MINUTES = "m"
    PERIOD_HOURS = "h"
    PERIOD_DAY = "d"
    PERIOD_WEEKLY = "w"

    OBJ_ASKS = "asks"
    OBJ_BIDS = "bids"
    OBJ_CANDLES = "candles"
    OBJ_TICK = "tick"

    OBJ_BID = "bid"
    OBJ_ASK = "ask"
    OBJ_PRICE = "price"
    OBJ_AMOUNT = "amount"
    OBJ_SYMBOL = "symbol"
    OBJ_TIMESTAMP = "timestamp"
    OBJ_OPEN = "open"
    OBJ_CLOSE = "close"
    OBJ_MIN = "min"
    OBJ_MAX = "max"
    OBJ_VOLUME = "volume"
    OBJ_VOLUMEQUOTE = "volumequote"
)

type Request struct {
    SymbolField string   `json:"symbol_field"`
    Regular bool         `json:"regular"`
    Timing float64       `json:"timing"`
    TimingUnit string    `json:"timing_unit"` //s,m,h,d,w
    JSON interface{}     `json:"json"`
}

func (request *Request) ToSymbol() string {
    return clearSymbol(utilities.ToString(utilities.SearchValue(request.JSON, (request.SymbolField))))
}

type Manifest struct {
    // Задается пользователем
    Name string                 `json:"name"`
    Exchange string             `json:"exchange"`
    Provider string             `json:"provider"`
    Entity string               `json:"entity"`
    URL string                  `json:"url"`
    Origin string               `json:"origin"`
    Request Request             `json:"request"`
    Response Response           `json:"response"`
    // Инициализируется дополнительно
    Id string
    // Каналы для передачи данных
    ChSignal chan *core.Signal
    ChMsg, ChErr chan interface{}
    // Коллекция соббщений адаптированных под данный манифест
    Messages map[string]string
    // Порог времени получение данных для сигнала, превышение которго не допустимо, в милисекундах
    TimeoutSignal int64
    // Порог времени обработки данных, превышение которых говорит о том что что-то пошло не так, в секундах
    TimeoutEntity float64
    // Время формирования сигнала
    // TimeRecd time.Time
}

func (manifest *Manifest) Init() {
    // Проводим дополнительную инициализацию
    manifest.Id = uuid.Must(uuid.NewV4()).String()
    manifest.Name = strings.ToLower(manifest.Name)
    manifest.Exchange = strings.ToLower(manifest.Exchange)
    manifest.Entity = strings.ToLower(manifest.Entity)
    manifest.Provider = strings.ToLower(manifest.Provider)
    manifest.TimeoutSignal = 200
    manifest.TimeoutEntity = 5
    manifest.Response.Parent = manifest
    manifest.Response.Index = 0
    var json interface{}
    // Инициализируем переменую куда будет записыватся ответ от биржи
    manifest.Response.JSON = json
    //
    manifest.ChSignal = make(chan *core.Signal)
    // Инициализируем сообщения
    manifest.Messages = make(map[string]string)
    //--
    manifest.Messages["CONNECTION_NOT_ACTIVATED"] = strings.Replace(constants.MSG_CONNECTION_NOT_ACTIVATED, constants.MSG_PLACE_NAME, manifest.Entity, 1)
    //--
    manifest.Messages["CONNECTION_ACTIVATED"] = strings.Replace(constants.MSG_CONNECTION_ACTIVATED, constants.MSG_PLACE_NAME, manifest.Entity, 1)
    //--
    manifest.Messages["CONNECTION_DEACTIVATED"] = strings.Replace(constants.MSG_CONNECTION_DEACTIVATED, constants.MSG_PLACE_NAME, manifest.Entity, 1)
    //--
    manifest.Messages["CONNECTION_NOT_PARAMS"] = strings.Replace(constants.MSG_CONNECTION_NOT_PARAMS, constants.MSG_PLACE_NAME, manifest.Entity, 1)
    //--
    CONNECTION_STOPED := strings.Replace(constants.MSG_CONNECTION_STOPED, constants.MSG_PLACE_NAME, manifest.Entity, 1)
    CONNECTION_STOPED = strings.Replace(CONNECTION_STOPED, constants.MSG_PLACE_PROVIDER, manifest.Provider, 1)
    manifest.Messages["CONNECTION_STOPED"] = CONNECTION_STOPED
    //--
    CONNECTION_NOT_STOPED := strings.Replace(constants.MSG_CONNECTION_NOT_STOPED, constants.MSG_PLACE_NAME, manifest.Entity, 1)
    CONNECTION_NOT_STOPED = strings.Replace(CONNECTION_NOT_STOPED, constants.MSG_PLACE_PROVIDER, manifest.Provider, 1)
    manifest.Messages["CONNECTION_NOT_STOPED"] = CONNECTION_NOT_STOPED
    //--
    CONNECTION_STARTED := strings.Replace(constants.MSG_CONNECTION_STARTED, constants.MSG_PLACE_NAME, manifest.Entity, 1)
    CONNECTION_STARTED = strings.Replace(CONNECTION_STARTED, constants.MSG_PLACE_PROVIDER, manifest.Provider, 1)
    manifest.Messages["CONNECTION_STARTED"] = CONNECTION_STARTED
    //--
    CONNECTION_NOT_STARTED := strings.Replace(constants.MSG_CONNECTION_NOT_STARTED, constants.MSG_PLACE_NAME, manifest.Entity, 1)
    CONNECTION_NOT_STARTED = strings.Replace(CONNECTION_NOT_STARTED, constants.MSG_PLACE_PROVIDER, manifest.Provider, 1)
    manifest.Messages["CONNECTION_NOT_STARTED"] = CONNECTION_NOT_STARTED
    //--
    CONNECTION_DISCONNECTED_TO := strings.Replace(constants.MSG_CONNECTION_DISCONNECTED_TO, constants.MSG_PLACE_NAME, manifest.Entity, 1)
    CONNECTION_DISCONNECTED_TO = strings.Replace(CONNECTION_DISCONNECTED_TO, constants.MSG_PLACE_URL, manifest.URL, 1)
    manifest.Messages["CONNECTION_DISCONNECTED_TO"] = CONNECTION_DISCONNECTED_TO
    //--
    CONNECTION_CONNECTED_TO := strings.Replace(constants.MSG_CONNECTION_CONNECTED_TO, constants.MSG_PLACE_NAME, manifest.Entity, 1)
    CONNECTION_CONNECTED_TO = strings.Replace(CONNECTION_CONNECTED_TO, constants.MSG_PLACE_URL, manifest.URL, 1)
    manifest.Messages["CONNECTION_CONNECTED_TO"] = CONNECTION_CONNECTED_TO
}

func (manifest *Manifest) IsTiming(started time.Time) bool {
    if manifest.Request.TimingUnit == "" {
        return true
    }
    if manifest.Request.Timing == 0 {
        return true
    }
    switch manifest.Request.TimingUnit {
    case PERIOD_SECONDS:
        return time.Now().Sub(started).Seconds() >= manifest.Request.Timing
    case PERIOD_MINUTES:
        return time.Now().Sub(started).Minutes() >= manifest.Request.Timing
    case PERIOD_HOURS:
        return time.Now().Sub(started).Hours() >= manifest.Request.Timing
    case PERIOD_DAY:
        return time.Now().Sub(started).Hours() >= manifest.Request.Timing * 24
    case PERIOD_WEEKLY:
        return time.Now().Sub(started).Hours() >= manifest.Request.Timing * 24 * 7
    default: return true
    }
}

// Инициализируем функцию приведения данных к общему виду
func (manifest *Manifest) Convertation() error {
    // manifest.TimeRecd = time.Now()
    // Пропускаем ответы в кол-ве указанном в манифесте
    skip, err := manifest.Response.skipResponse()
    if err != nil {
        return err
    }
    if skip {
        return nil
    }
    // В зависимости от запрашиваемой сущности вызываем соответсвующую конвертацию
    switch manifest.Entity {
    case constants.ENTITY_TICK:
        return manifest.Response.ToTick()
    case constants.ENTITY_DEPTH:
        return manifest.Response.ToDepth()
    case constants.ENTITY_CANDLE:
        return manifest.Response.ToCandles()
    }
    return nil
}

func (manifest *Manifest) CheckError() error {
    if manifest.Id == "" || manifest.Entity == "" || manifest.Response.JSON == nil || manifest.Convertation == nil {
        return errors.New("Manifest is not initialization")
    }
    if manifest.Response.JSON == nil {
        return errors.New("Incoming response.JSON is nil")
    }
    if manifest.Response.Values == nil {
        return errors.New("In manifest Values is nil")
    }
    if len(manifest.Response.Values) == 0 {
        return errors.New("In manifest count of Values is 0")
    }
    for _, value := range manifest.Response.Values {
        if value.Name == "" {
            return errors.New("In manifest Name of Values is empty")
        }
        // if value.Design == nil {
        //     return errors.New("In manifest [" + value.Name + "]: Design is nil")
        // }
        if value.Design.Kind == 0 {
            return errors.New("In manifest [" + value.Name + "]: Design.Kind is 0")
        }
    }
    return nil
}

func getValueByPath(data interface{}, paths []Path) (res interface{}, err error) {
    res = data
    for _, path := range paths {
        if len(path.Int) > 0 {
            res, err = utilities.GetValueByInt(res, path.Int)
            if err != nil {
                return
            }
        } else {
            if len(path.Str) > 0 {
                res, err = utilities.GetValueByStr(res, path.Str)
                if err != nil {
                    return
                }
            }
        }
    }
    return
}

func isChecked(data interface{}, check Value) (bool, error) {
    if data != nil {
        switch check.Design.Kind {
        case reflect.String:
            value, err := toString(data, check)
            if err != nil {
                return false, err
            }
            if value != check.Design.CheckStr {
                return false, nil
            }
        case reflect.Int64:
            value, err := toInt64(data, check)
            if err != nil {
                return false, err
            }
            if value != check.Design.CheckInt {
                return false, nil
            }
        }
    }
    return true, nil
}

func toString(data interface{}, obj Value) (string, error) {
    if obj.Design.Kind != reflect.String {
        return "", errors.New("Mismatched type in ToString")
    }
    if len(obj.Design.Path) > 0 {
        var err error
        for _, path := range obj.Design.Path {
            if len(path.Int) > 0 {
                data, err = utilities.GetValueByInt(data, path.Int)
                if err != nil {
                    return "", err
                }
            } else {
                if len(path.Str) > 0 {
                    data, err = utilities.GetValueByStr(data, path.Str)
                    if err != nil {
                        return "", err
                    }
                }
            }
        }
        return utilities.ToString(data), nil
    }
    return "", nil
}

func toFloat64(data interface{}, obj Value) (float64, error) {
    // Проверяем на тип данных
    if obj.Design.Kind != reflect.Float64 {
        return 0, errors.New("Mismatched type in ToFloat64")
    }
    // Получаем значение
    value, err := getValueByPath(data, obj.Design.Path)
    if err != nil {
        return 0, err
    }
    // Приводим к нуному типу
    return utilities.ToFloat(value), nil
}

func toInt64(data interface{}, obj Value) (int64, error) {
    // Проверяем на тип данных
    if obj.Design.Kind != reflect.Int64 {
        return 0, errors.New("Mismatched type in ToInt64")
    }
    // Получаем значение
    value, err := getValueByPath(data, obj.Design.Path)
    if err != nil {
        return 0, err
    }
    // Приводим к нуному типу
    return utilities.ToInt(value), nil
}

func clearSymbol(symbol string) string {
    return strings.Replace(symbol, "-", "", -1)
}

func GetManifestPaths() (paths map[string]string, err error) {
    corePath := "./manifests/"
    var files []os.FileInfo
    files, err = ioutil.ReadDir(corePath)
    if err != nil {
        return
    }
    paths = make(map[string]string)
    for _, f := range files {
        paths[strings.Split(f.Name(), ".")[0]] = corePath + f.Name()
    }
    return
}

func ReadManifest(path string) (manifest *Manifest, err error) {
    var jsonFile []byte
    jsonFile, err = ioutil.ReadFile(path)
    if err != nil {
		return
	}
    manifest = new(Manifest)
    json.Unmarshal(jsonFile, manifest)
    manifest.Init()
    return
}

func GetManifests() (manifests map[string][]*Manifest, err error) {
    var paths map[string]string
    var manifest *Manifest
    paths, err = GetManifestPaths()
    if err != nil {
		return
	}
    manifests = make(map[string][]*Manifest)
    for _, path := range paths {
        manifest, err = ReadManifest(path)
        if err != nil {
    		return
    	}
        manifests[manifest.Exchange] = append(manifests[manifest.Exchange], manifest)
    }
    return
}

// func Equals(firstData interface{}, secondData interface{}) bool {
//     return true
// }
