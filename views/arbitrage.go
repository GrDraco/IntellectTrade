package views

import (
    // "fmt"
    "sort"

    "./components"
    packUI "../ui"
    "../market/strategies"

    "github.com/nsf/termbox-go"
)

// Функция отрисовки стратегии arbitrage
func arbitrageDraw(console *packUI.Console, bodyX, bodyY int) {
    if console.Controls["Properties"] == nil {
        return
    }
    if !console.Controls["Properties"].(map[string]interface{})[strategies.PROPERTY_STARTED].(bool) {
        return
    }
    // Отображаем котиовки Полученные стратегий
    quotations := console.Controls["Properties"].(map[string]interface{})[strategies.PROPERTY_QUOTATIONS].(map[string]map[string]*strategies.Quotations)
    line := 3
    height, _ := components.QuotationsDraw(bodyX + 2, line, quotations)
    line = line + height + 1
    // Отобрабаем стакан указанный для отображения по комманде
    var keysSymbol []string
    for key, _ := range quotations {
        keysSymbol = append(keysSymbol, key)
    }
    sort.Strings(keysSymbol)
    var keysExchange []string
    for key, _ := range quotations[keysSymbol[0]] {
        keysExchange = append(keysExchange, key)
    }
    sort.Strings(keysExchange)
    heightDepths := 0
    column := bodyX + 2
    for _, exchange := range keysExchange {
        h_1, w_1 := components.ASKsDraw(column, line, quotations[keysSymbol[0]][exchange].Depth, keysSymbol[0], exchange, 10)
        h_2, _ := components.BIDsDraw(column, line + h_1 + 1, quotations[keysSymbol[0]][exchange].Depth, keysSymbol[0], exchange, 10)
        column = column + w_1 + 1
        if (h_1 + h_2) > heightDepths {
            heightDepths = h_1 + h_2
        }
    }
    line = line + heightDepths + 2
    // Отображаем найденные лучшие цены
    components.BestPricesDraw(bodyX + 2, line, console.Controls["Properties"].(map[string]interface{})[strategies.PROPERTY_BESTPRICES].(map[string]*strategies.BestPrices))
}
// Функция очистки стратегии arbitrage
func arbitrageClear(console *packUI.Console, bodyX, bodyY int) {
    w, _ := termbox.Size()
    packUI.Fill(bodyX + 2, 2, w - bodyX + 2, bodyY - 2, termbox.Cell{Ch: ' '})
}
// Функция обработки стратегии arbitrage
func arbitrageAction(console *packUI.Console, data []interface{}) {
    console.Controls["Properties"] = data[0]
}

func Arbitrage() *packUI.Console  {
    return packUI.NewConsole("arbitrage", `Торговая стратегия "Арбитраж"`, false, arbitrageDraw, arbitrageClear, arbitrageAction)
}
// console arbitrage
