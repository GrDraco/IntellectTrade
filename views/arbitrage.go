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
    height_1 := 0
    height_2 := 0
    width_1 := 0
    column := bodyX + 2
    for _, exchange := range keysExchange {
        height_1, width_1 = components.ASKsDraw(column, line, quotations[keysSymbol[0]][exchange].Depth, keysSymbol[0], exchange)
        height_2, _ = components.BIDsDraw(column, line + height_1 + 1, quotations[keysSymbol[0]][exchange].Depth, keysSymbol[0], exchange)
        column = column + width_1 + 1
    }
    line = line + height_1 + height_2 + 1
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
