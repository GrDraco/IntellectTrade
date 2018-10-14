/*
    WSocket это класс по работе с сокетами
*/
package connections

import (
    // "fmt"
    "errors"
    // "../../../utilities"
    "golang.org/x/net/websocket"
)

type WSocket struct {
    // Наследования
    BaseConnection
    // Свойства
    socket *websocket.Conn
    err error
}
//SET wSocket.manifest.Request.JSON = utilities.ReplaceValues(wSocket.manifest.Request.JSON, data[0])
func NewWSocket(manifest *Manifest) *WSocket {
    // Выделение памяти под сокет
    wSocket := &WSocket {}
    // Инициализация функци работы запроса данных
    wSocket.init = func() error {
        // Создание канала сокета
        wSocket.socket, wSocket.err = websocket.Dial(wSocket.manifest.URL, "", wSocket.manifest.Origin)
        if wSocket.err != nil {
            return wSocket.err
        }
        // Отправляем по сокету данные
        wSocket.err = websocket.JSON.Send(wSocket.socket, wSocket.manifest.Request.JSON)
        if wSocket.err != nil {
            return wSocket.err
        }
        // params := make(map[string]interface{})
        // params["method"] = "snapshotOrderbook"
        // wSocket.manifest.Request.JSON = utilities.ReplaceValues(wSocket.manifest.Request.JSON, params)
        // wSocket.err = websocket.JSON.Send(wSocket.socket, wSocket.manifest.Request.JSON)
        // if wSocket.err != nil {
        //     return wSocket.err
        // }
        return nil
    }
    wSocket.send = func() error {
        if wSocket.socket == nil {
            return errors.New(wSocket.manifest.Messages["CONNECTION_NOT_PARAMS"])
        }
        wSocket.err = websocket.JSON.Receive(wSocket.socket, &wSocket.manifest.Response.JSON)
    	if wSocket.err != nil {
            return wSocket.err
        }
        return nil
    }
    wSocket.close = func() error {
        wSocket.err = wSocket.socket.Close()
        if wSocket.err != nil {
            return wSocket.err
        }
        return nil
    }
    // Активируем базовый функционал
    if wSocket.activate(manifest) {
        return wSocket
    }
    return nil
}

// func (wS *WSocket) send(data string) {
//     if _, wS.Err = wS.Socket.Write([]byte(data)); wS.Err != nil {
//         wS.PrintError("send", wS.Err.Error())
//     } else {
//         wS.PrintLog("send", "Send message -> " + data, true)
//     }
// }
//
// func (wS *WSocket) get() string {
//     var msg = make([]byte, 512)
//     var n int
//     if n, wS.Err = wS.Socket.Read(msg); wS.Err != nil {
//         wS.PrintError("Get", wS.Err.Error())
//         return ""
//     } else {
//         str := fmt.Sprintf("%s", msg[:n])
//         wS.PrintLog("Get", "Received -> " + str, true)
//         return str
//     }
// }




// func (wS *WSocket) Stop() {
//     wS.IsStop = true
// }
/*func main() {
    log.SetFlags(0)
    origin := "http://localhost/"
    //url := "wss://stream.binance.com:9443/ws/bnbbtc@depth"
    url := "wss://api.hitbtc.com/api/2/ws"
    ws, err := websocket.Dial(url, "", origin)
    if err != nil {
        log.Fatal(err)
    } else {
        log.Printf("Connected to: %s", url)
    }
    sendData := `{ "method": "subscribeTicker", "params": { "symbol": "ETHBTC" }, "id": 123 }`
    //sendData := `{ "method": "getCurrency", "params": { "currency": "ETH" }, "id": 123 }`
    if _, err := ws.Write([]byte(sendData)); err != nil {
        log.Fatal(err)
    } else {
        log.Printf("Send message: %s", sendData)
    }
    for {
        var msg = make([]byte, 512)
        var n int
        if n, err = ws.Read(msg); err != nil {
            log.Fatal(err)
        }
        log.Printf("Received: %s", msg[:n])
    }
}*/
