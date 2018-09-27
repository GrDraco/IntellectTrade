/*
    Manifest это класс с параметрами для биржы и работы их данными
*/
package connections

import (
    "errors"
    "strings"
    "io/ioutil"
    "os"
    "encoding/json"
    "github.com/satori/go.uuid"
    "time"
    // "fmt"
    "../../constants"
    "../../core"
    "../../../utilities"
)

type ArrayValues struct {
    Path []string             `json:"path"`
    ValuesIsArray bool        `json:"values_is_array"`
    // Пути к значениям внутри масивов
    ArrPrice []string         `json:"arr_price"`
    ArrAmount []string        `json:"arr_amount"`
    ArrVolume []string        `json:"arr_volume"`
    ArrTimestamp []string     `json:"arr_timestamp"`
    ArrOpen []string          `json:"arr_open"`
    ArrClose []string         `json:"arr_close"`
    ArrMin []string           `json:"arr_min"`
    ArrMax []string           `json:"arr_max"`
    ArrVolumeQuote []string   `json:"arr_volumeQuote"`

    // Позиции значений в масиве
    IndexPrice int64          `json:"index_price"`
    IndexAmount int64         `json:"index_amount"`
    IndexVolume int64         `json:"index_volume"`
    IndexTimestamp int64      `json:"index_timestamp"`
    IndexOpen int64           `json:"index_open"`
    IndexClose int64          `json:"index_close"`
    IndexMin int64            `json:"index_min"`
    IndexMax int64            `json:"index_max"`
    IndexVolumeQuote int64    `json:"index_volumeQuote"`
}

type Values struct {
    // Пути к значетиям
    Ask []string            `json:"ask"`
    Bid []string            `json:"bid"`
    Volume []string         `json:"volume"`
    Symbol []string         `json:"symbol"`
    Timestamp []string      `json:"timestamp"`
    Ping []string           `json:"ping"`
    // Пути к данным с масивами
    Asks ArrayValues        `json:"asks"`
    Bids ArrayValues        `json:"bids"`
    Candles ArrayValues     `json:"candles"`
}

type Failed struct {
    Message []string     `json:"message"`
    Description []string `json:"description"`
}

type Response struct {
    Success []string `json:"success"`
    Failed Failed    `json:"failed"`
    Values Values    `json:"values"`
    Ping int64
    JSON interface{}
}

// type Value struct {
//     Name string     `json:"name"`
//     Path []string   `json:"path"`
// }

const (
    PERIOD_SECONDS = "s"
    PERIOD_MINUTES = "m"
    PERIOD_HOURS = "h"
    PERIOD_DAY = "d"
    PERIOD_WEEKLY = "w"
)

type Manifest struct {
    // Задается пользователем
    Name string             `json:"name"`
    Exchange string         `json:"exchange"`
    Provider string         `json:"provider"`
    Entity string           `json:"entity"`
    URL string              `json:"url"`
    Origin string           `json:"origin"`
    RequestJSON interface{} `json:"request_json"`
    Response Response       `json:"response"`
    Regular bool            `json:"regular"`
    DataIsUpdates bool      `json:"data_is_updates"`
    Timing float64          `json:"timing"`
    TimingUnit string       `json:"timing_unit"` //s,m,h,d,w
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
    if manifest.TimingUnit == "" {
        return true
    }
    if manifest.Timing == 0 {
        return true
    }
    switch manifest.TimingUnit {
    case PERIOD_SECONDS:
        return time.Now().Sub(started).Seconds() >= manifest.Timing
    case PERIOD_MINUTES:
        return time.Now().Sub(started).Minutes() >= manifest.Timing
    case PERIOD_HOURS:
        return time.Now().Sub(started).Hours() >= manifest.Timing
    case PERIOD_DAY:
        return time.Now().Sub(started).Hours() >= manifest.Timing * 24
    case PERIOD_WEEKLY:
        return time.Now().Sub(started).Hours() >= manifest.Timing * 24 * 7
    default: return true
    }
}

// Инициализируем функцию приведения данных к общему виду
func (manifest *Manifest) Convertation() error {
    // manifest.TimeRecd = time.Now()
    // В зависимости от запрашиваемой сущности вызываем соответсвующую конвертацию
    switch manifest.Entity {
    case constants.ENTITY_TICK:
        return manifest.ConvertToTick()
    case constants.ENTITY_DEPTH:
        return manifest.ConvertToDepth()
    case constants.ENTITY_CANDLE:
        return manifest.ConvertToCandle()
    }
    return nil
}

func (manifest *Manifest) ConvertError() error {
    if manifest.Id == "" || manifest.Entity == "" || manifest.Response.JSON == nil || manifest.Convertation == nil {
        return errors.New("Manifest is not initialization")
    }
    if manifest.Response.JSON == nil {
        return errors.New("Incoming response.JSON is nil")
    }
    if manifest.Response.Success == nil {
        return errors.New("In JSON file `success` is nil")
    }
    if len(manifest.Response.Success) == 0 {
        return errors.New("In JSON file `success` not elements")
    }
    if manifest.Response.Failed.Message == nil {
        return errors.New("In JSON file `message` is nil")
    }
    if len(manifest.Response.Failed.Message) == 0 {
        return errors.New("In JSON file `message` not elements")
    }
    return nil
}

func (manifest *Manifest) ConvertToError() error {
    err := manifest.ConvertError()
    if err != nil {
        return err
    }
    if utilities.GetValue(manifest.Response.JSON, manifest.Response.Success) != nil {
        return nil
    }
    msgFailed := core.NewError(manifest.Entity,
                              utilities.ToString(utilities.GetValue(manifest.Response.JSON, manifest.Response.Failed.Message)),
                              utilities.ToString(utilities.GetValue(manifest.Response.JSON, manifest.Response.Failed.Description)))

    if msgFailed.Message != "" || msgFailed.Description != "" {
        return errors.New(msgFailed.Message + " " + msgFailed.Description)
    }
    return nil
}

func (manifest *Manifest) ConvertToTick() error {
    signal := new(core.Signal)
    err := manifest.ConvertToError()
    if err != nil {
        return err
    }
    signal = &core.Signal {
        Ping: manifest.Response.Ping,
        Connection: manifest.Name,
        Exchange: manifest.Exchange,
        Entity: constants.ENTITY_TICK,
        Symbol: clearSymbol(utilities.ToString(utilities.GetValue(manifest.Response.JSON, manifest.Response.Values.Symbol))),
        TimeRecd: time.Now(),
        DataIsUpdates: manifest.DataIsUpdates }
    signal.Data = &core.Tick {
        Ask: utilities.ToFloat(utilities.GetValue(manifest.Response.JSON, manifest.Response.Values.Ask)),
        Bid: utilities.ToFloat(utilities.GetValue(manifest.Response.JSON, manifest.Response.Values.Bid)),
        Volume: utilities.ToFloat(utilities.GetValue(manifest.Response.JSON, manifest.Response.Values.Volume)) }
    if manifest.Provider == constants.CONNECTION_API {
        if signal.Ping > manifest.TimeoutSignal {
            signal.TimeOut = true
        } else {
            signal.TimeOut = false
        }
    }
    manifest.ChSignal<-signal
    return nil
}

func (manifest *Manifest) ConvertToDepth() error {
    signal := new(core.Signal)
    err := manifest.ConvertToError()
    if err != nil {
        return err
    }
    signal = &core.Signal {
        Ping: manifest.Response.Ping,
        Connection: manifest.Name,
        Exchange: manifest.Exchange,
        Entity: constants.ENTITY_DEPTH,
        Symbol: clearSymbol(utilities.ToString(utilities.SearchValue(manifest.RequestJSON, "symbol"))),
        TimeRecd: time.Now(),
        DataIsUpdates: manifest.DataIsUpdates }
    depth := new(core.Depth)
    depth.Asks = make(map[float64]*core.Order)
    // Asks
    // fmt.Println("method", utilities.SearchValue(manifest.RequestJSON, "method"))
    arrAsk := utilities.GetValue(manifest.Response.JSON, manifest.Response.Values.Asks.Path)
    if arrAsk != nil {
        if manifest.Response.Values.Asks.ValuesIsArray {
            for i := 0; i < len(arrAsk.([]interface{})); i++ {
                depth.Asks[utilities.ToFloat(arrAsk.([]interface{})[i].([]interface{})[manifest.Response.Values.Asks.IndexPrice])] = &core.Order {
                    Price: utilities.ToFloat(arrAsk.([]interface{})[i].([]interface{})[manifest.Response.Values.Asks.IndexPrice]),
                    Amount: utilities.ToFloat(arrAsk.([]interface{})[i].([]interface{})[manifest.Response.Values.Asks.IndexAmount]),
                    Volume: utilities.ToFloat(arrAsk.([]interface{})[i].([]interface{})[manifest.Response.Values.Asks.IndexVolume])}
            }
        } else {
            for i := 0; i < len(arrAsk.([]interface{})); i++ {
                price := utilities.ToFloat(utilities.GetValue(arrAsk.([]interface{})[i], manifest.Response.Values.Asks.ArrPrice))
                depth.Asks[price] = &core.Order {
                    Price: price,
                    Amount: utilities.ToFloat(utilities.GetValue(arrAsk.([]interface{})[i], manifest.Response.Values.Asks.ArrAmount)),
                    Volume: utilities.ToFloat(utilities.GetValue(arrAsk.([]interface{})[i], manifest.Response.Values.Asks.ArrVolume))}
            }
        }
    }
    depth.Bids = make(map[float64]*core.Order)
    // Bids
    arrBid := utilities.GetValue(manifest.Response.JSON, manifest.Response.Values.Bids.Path)
    if arrBid != nil {
        if manifest.Response.Values.Bids.ValuesIsArray {
            for i := 0; i < len(arrBid.([]interface{})); i++ {
                depth.Bids[utilities.ToFloat(arrBid.([]interface{})[i].([]interface{})[manifest.Response.Values.Bids.IndexPrice])] = &core.Order {
                    Price: utilities.ToFloat(arrBid.([]interface{})[i].([]interface{})[manifest.Response.Values.Bids.IndexPrice]),
                    Amount: utilities.ToFloat(arrBid.([]interface{})[i].([]interface{})[manifest.Response.Values.Bids.IndexAmount]),
                    Volume: utilities.ToFloat(arrBid.([]interface{})[i].([]interface{})[manifest.Response.Values.Bids.IndexVolume])}
            }
        } else {
            for i := 0; i < len(arrBid.([]interface{})); i++ {
                price := utilities.ToFloat(utilities.GetValue(arrBid.([]interface{})[i], manifest.Response.Values.Bids.ArrPrice))
                depth.Bids[price] = &core.Order {
                    Price: price,
                    Amount: utilities.ToFloat(utilities.GetValue(arrBid.([]interface{})[i], manifest.Response.Values.Bids.ArrAmount)),
                    Volume: utilities.ToFloat(utilities.GetValue(arrBid.([]interface{})[i], manifest.Response.Values.Bids.ArrVolume))}
            }
        }
    }
    signal.Data = depth
    ///////
    if manifest.Provider == constants.CONNECTION_API {
        if signal.Ping > manifest.TimeoutSignal {
            signal.TimeOut = true
        } else {
            signal.TimeOut = false
        }
    }
    manifest.ChSignal<-signal
    return nil
}

func (manifest *Manifest) ConvertToCandle() error  {
    signal := new(core.Signal)
    err := manifest.ConvertToError()
    if err != nil {
        return err
    }
    signal = &core.Signal {
        Ping: manifest.Response.Ping,
        Connection: manifest.Name,
        Exchange: manifest.Exchange,
        Entity: constants.ENTITY_CANDLE,
        Symbol: clearSymbol(utilities.ToString(utilities.SearchValue(manifest.RequestJSON, "symbol"))),
        TimeRecd: time.Now(),
        DataIsUpdates: manifest.DataIsUpdates }
    candles := make([]*core.Candle, 0)
    var arr interface{}
    if len(manifest.Response.Values.Candles.Path) > 0 {
        arr = utilities.GetValue(manifest.Response.JSON, manifest.Response.Values.Candles.Path)
    } else {
        arr = manifest.Response.JSON
    }
    if arr != nil {
        if manifest.Response.Values.Candles.ValuesIsArray {
            for i := 0; i < len(arr.([]interface{})); i++ {
                candles = append(candles, &core.Candle {
                    Open: utilities.ToFloat(arr.([]interface{})[i].([]interface{})[manifest.Response.Values.Candles.IndexOpen]),
                    Close: utilities.ToFloat(arr.([]interface{})[i].([]interface{})[manifest.Response.Values.Candles.IndexClose]),
                    Min: utilities.ToFloat(arr.([]interface{})[i].([]interface{})[manifest.Response.Values.Candles.IndexMin]),
                    Max: utilities.ToFloat(arr.([]interface{})[i].([]interface{})[manifest.Response.Values.Candles.IndexMax]),
                    Volume: utilities.ToFloat(arr.([]interface{})[i].([]interface{})[manifest.Response.Values.Candles.IndexVolume]),
                    VolumeQuote: utilities.ToFloat(arr.([]interface{})[i].([]interface{})[manifest.Response.Values.Candles.IndexVolumeQuote]) })
            }
        } else {
            for i := 0; i < len(arr.([]interface{})); i++ {
                candles = append(candles, &core.Candle {
                    Open: utilities.ToFloat(utilities.GetValue(arr.([]interface{})[i], manifest.Response.Values.Candles.ArrOpen)),
                    Close: utilities.ToFloat(utilities.GetValue(arr.([]interface{})[i], manifest.Response.Values.Candles.ArrClose)),
                    Min: utilities.ToFloat(utilities.GetValue(arr.([]interface{})[i], manifest.Response.Values.Candles.ArrMin)),
                    Max: utilities.ToFloat(utilities.GetValue(arr.([]interface{})[i], manifest.Response.Values.Candles.ArrMax)),
                    Volume: utilities.ToFloat(utilities.GetValue(arr.([]interface{})[i], manifest.Response.Values.Candles.ArrVolume)),
                    VolumeQuote: utilities.ToFloat(utilities.GetValue(arr.([]interface{})[i], manifest.Response.Values.Candles.ArrVolumeQuote)) })
            }
        }
    }
    signal.Data = candles
    ///////
    if manifest.Provider == constants.CONNECTION_API {
        if signal.Ping > manifest.TimeoutSignal {
            signal.TimeOut = true
        } else {
            signal.TimeOut = false
        }
    }
    manifest.ChSignal<-signal
    return nil
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
