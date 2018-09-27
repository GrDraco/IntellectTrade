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

const coldef = termbox.ColorDefault

func QuotationsDraw(x, y int, quotations map[string]map[string]*strategies.Quotations) (height, width int) {
    line := 0
    column := 11
    // Шапка
    packUI.Tprint(x + 2, y + 2, coldef, coldef, "Инструмент ")
    var keysSymbol []string
    var keysExchange []string
    for key, _ := range quotations {
        keysSymbol = append(keysSymbol, key)
    }
    sort.Strings(keysSymbol)
    for i := len(keysSymbol) - 1; i >= 0; i-- {
        for key, _ := range quotations[keysSymbol[i]] {
            keysExchange = append(keysExchange, key)
        }
        sort.Strings(keysExchange)
        // Стенки
        borders.GroupBoxDraw("Стаканы", x, y, len(keysExchange) * 40, len(keysSymbol) + 2, 2, termbox.ColorCyan)
        column = 11
        // Шапка данных
        for j := len(keysExchange) - 1; j >= 0; j-- {
            packUI.Tprint(x + 2 + column, y + 1, coldef, coldef, (" " + strings.ToUpper(keysExchange[j]) + " "))
            packUI.Tprint(x + 2 + column, y + 2 , coldef, coldef, " Уровень ")
            packUI.Tprint(x + 2 + column + 8, y + 2 , coldef, coldef, " ASK ")
            packUI.Tprint(x + 2 + column + 21, y + 2 , coldef, coldef, " BID ")
            column = column + 34
        }
        // Данные
        column = 11
        for j := len(keysExchange) - 1; j >= 0; j-- {
            packUI.Tprint(x + 2, y + 3 + line, coldef, coldef, (strings.ToUpper(keysSymbol[i]) + " "))
            packUI.Tprint(x + 2 + column, y + 3 + line, coldef, coldef, (" " + strconv.FormatInt(int64(quotations[keysSymbol[i]][keysExchange[j]].IndexDepth), 10) + " "))
            if len(quotations[keysSymbol[i]][keysExchange[j]].Depth.Asks) > 0 {
                packUI.Tprint(x + 2 + column + 8, y + 3 + line, coldef, coldef, (" " + strconv.FormatFloat(quotations[keysSymbol[i]][keysExchange[j]].Depth.GetAsks()[quotations[keysSymbol[i]][keysExchange[j]].IndexDepth].Price, 'f', -1, 64) + " "))
            }
            if len(quotations[keysSymbol[i]][keysExchange[j]].Depth.Bids) > 0 {
                packUI.Tprint(x + 2 + column + 21, y + 3 + line, coldef, coldef, (" " + strconv.FormatFloat(quotations[keysSymbol[i]][keysExchange[j]].Depth.GetBids()[quotations[keysSymbol[i]][keysExchange[j]].IndexDepth].Price, 'f', -1, 64) + " "))
            }
            column = column + 34
        }
        line++
    }
    height = line + 3
    width = len(keysExchange) * 40
    return
}
