package components

import (
    "sort"
    "strings"
    "strconv"

    "../borders"
    packUI "../../ui"
    "../../market/strategies"
    "github.com/nsf/termbox-go"
)

func BestPricesDraw(x, y int, bestPrices map[string]*strategies.BestPrices) (height, width int) {
    line := 0
    column := 11
    var keysSymbol []string
    for key, _ := range bestPrices {
        keysSymbol = append(keysSymbol, key)
    }
    sort.Strings(keysSymbol)
    // Стенки
    borders.GroupBoxDraw("Лучшие цены из стаканов", x, y, 53, len(keysSymbol) + 1, 1, termbox.ColorCyan)
    // Шапка данных
    packUI.Tprint(x + 2, y + 1, coldef, coldef, "Инструмент")
    packUI.Tprint(x + 2 + column, y + 1, coldef, coldef, "ASK")
    packUI.Tprint(x + 2 + column + 11, y + 1, coldef, coldef, "Биржа")
    packUI.Tprint(x + 2 + column + 22, y + 1, coldef, coldef, "BID")
    packUI.Tprint(x + 2 + column + 33, y + 1, coldef, coldef, "Биржа")
    line++
    // Данные
    for i := len(keysSymbol) - 1; i >= 0; i-- {
        packUI.Tprint(x + 2, y + 2, coldef, coldef, strings.ToUpper(bestPrices[keysSymbol[i]].Symbol))
        packUI.Tprint(x + 2 + column, y + 2, coldef, coldef, strconv.FormatFloat(bestPrices[keysSymbol[i]].Ask.Price, 'f', -1, 64))
        packUI.Tprint(x + 2 + column + 11, y + 2, coldef, coldef, bestPrices[keysSymbol[i]].Ask.Exchange)
        packUI.Tprint(x + 2 + column + 22, y + 2, coldef, coldef, strconv.FormatFloat(bestPrices[keysSymbol[i]].Bid.Price, 'f', -1, 64))
        packUI.Tprint(x + 2 + column + 33, y + 2, coldef, coldef, bestPrices[keysSymbol[i]].Bid.Exchange)
        line++
    }
    height = line + 1
    width = 53
    return
}
