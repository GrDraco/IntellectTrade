package main

import (
    "testing"
    "time"
    "./market/exchanges"
    "./utilities"
)

// func TestWSocketGet(t *testing.T) {
//     test := new(TestLocation)
//     test.Init("test WSocket")
//     sub := exchanges.Subscription {
//         Name: "Exchange 'Test_Get'",
//         Url: "wss://api.hitbtc.com/api/2/ws",
//         Origin: "http://localhost/",
//         Data: `{ "method": "getCurrency", "params": { "currency": "ETH" }, "id": 123 }`,
//         Regular: false }
//     var ws = exchanges.NewWSocket(sub)
//     //go readingOut(ws.Out)
//     res := ws.Get()
//     _, found := utilities.SearchIndex(res, "ETH");
//     if !found {
//         t.Error("TestWSocketGet", "WSocketGet NOT WORK")
//     }
// }

func TestWSocketStart(t *testing.T) {
    test := new(TestLocation)
    test.Init("test WSocket")
    sub := exchanges.Subscription {
        Name: "Exchange 'Test_Start'",
        Url: "wss://api.hitbtc.com/api/2/ws",
        Origin: "http://localhost/",
        Data: `{ "method": "getCurrency", "params": { "currency": "ETH" }, "id": 123 }`,
        Regular: false }
    var ws = exchanges.NewWSocket(sub)
    ws.Start()
    i := 0
    for {        
        if i == 3 {
            return
        }
        test.PrintLog("TestWSocketStart", "reading...")
        select {
        case str := <-ws.Out:
            _, found := utilities.SearchIndex(str, "ETH");
            if !found {
                t.Error("TestWSocketStart", "WSocketGet NOT WORK")
            }
        case <-time.After(time.Second):
            test.PrintLog("TestWSocketStart", "Time Out")
            return
        }
        i++
    }
}
