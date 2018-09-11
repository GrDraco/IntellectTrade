package main

import (
    "testing"
    // "strconv"
    "time"
    "fmt"
    "./market/core"
    "./market/exchanges"
    "./market/constants"
    "./market/exchanges/connections"
    // "./utilities"
)

func TestExchange(t *testing.T) {
    chSignal := make(chan *core.Signal)
    chMsg := make(chan interface{})
    chErr := make(chan interface{})
    chStop := make(chan bool)
    manifests, err := connections.GetManifests()
    if err != nil {
		fmt.Errorf("Ошибка чтения манифестов: " + err.Error())
	}
    go func() {
        for {
            select {
            case <-chStop:
                return
            case err := <-chErr:
                fmt.Println(err.(*core.Message).FullString())
            case msg := <-chMsg:
                fmt.Println(msg.(*core.Message).FullString())
            case signal := <-chSignal:
                // fmt.Println(signal.TickStr(true))
                fmt.Println(signal.DepthStr(true))
            }
        }
    }()
    var params map[string]interface{}
    // kucoin := exchanges.NewExchange(manifests["kucoin"], chSignal, chMsg, chErr)
    hitBTC := exchanges.NewExchange(manifests["hitbtc"], chSignal, chMsg, chErr)
    params = make(map[string]interface{})
    // -1-
    fmt.Println("-1-")
    // params["symbol"] = "BTC-USDT"
    // kucoin.SetValues(constants.ENTITY_TICK, params)
    // kucoin.SetValues(constants.ENTITY_DEPTH, params)
    params["symbol"] = "ETHBTC"
    // hitBTC.SetValues(constants.ENTITY_TICK, params)
    hitBTC.SetValues(constants.ENTITY_DEPTH, params)
    // kucoin.Turn(constants.ENTITY_TICK)
    // kucoin.Turn(constants.ENTITY_DEPTH)
    // hitBTC.Turn(constants.ENTITY_TICK)
    hitBTC.Turn(constants.ENTITY_DEPTH)
    time.Sleep(time.Duration(5)*time.Second)
    // -2-
    fmt.Println("-2-")
    // params["symbol"] = "ETH-USDT"
    // kucoin.SetValues(constants.ENTITY_TICK, params)
    // kucoin.SetValues(constants.ENTITY_DEPTH, params)
    params["symbol"] = "BTCUSD"
    // hitBTC.SetValues(constants.ENTITY_TICK, params)
    hitBTC.SetValues(constants.ENTITY_DEPTH, params)
    time.Sleep(time.Duration(5)*time.Second)
    // -3-
    fmt.Println("-3-")
    // kucoin.Turn(constants.ENTITY_TICK)
    // kucoin.Turn(constants.ENTITY_DEPTH)
    // hitBTC.Turn(constants.ENTITY_TICK)
    hitBTC.Turn(constants.ENTITY_DEPTH)
    time.Sleep(time.Duration(5)*time.Second)
    // kucoin.Turn(constants.ENTITY_TICK)
    // kucoin.Turn(constants.ENTITY_DEPTH)
    // hitBTC.Turn(constants.ENTITY_TICK)
    hitBTC.Turn(constants.ENTITY_DEPTH)
    time.Sleep(time.Duration(5)*time.Second)
    // kucoin.Turn(constants.ENTITY_TICK)
    // kucoin.Turn(constants.ENTITY_DEPTH)
    // fmt.Println("Kucoin is test SUCCESS")
    // hitBTC.Turn(constants.ENTITY_TICK)
    hitBTC.Turn(constants.ENTITY_DEPTH)
    fmt.Println("HitBTC is test SUCCESS")

    close(chStop)
}
