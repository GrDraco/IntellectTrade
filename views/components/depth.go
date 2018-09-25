package components

import (
    "strings"
    "strconv"

    "../borders"
    packUI "../../ui"
    "../../market/core"
    "github.com/nsf/termbox-go"
)

func ASKsDraw(x, y int, depth *core.Depth, symbol, exchange string) (height, width int) {
    return OrdersDraw(x, y, depth.Asks, symbol, exchange, "ask")
}

func BIDsDraw(x, y int, depth *core.Depth, symbol, exchange string) (height, width int) {
    return OrdersDraw(x, y, depth.Bids, symbol, exchange, "bid")
}

func OrdersDraw(x, y int, orders map[float64]*core.Order, symbol, exchange, direction string) (height, width int) {
    line := 0
    column := 5
    maxDepth := 10
    // Стенки
    borders.GroupBoxDraw(strings.ToUpper(direction) + "/" + exchange + "/" + strings.ToUpper(symbol) + "/" + strconv.FormatInt(int64(len(orders)), 10), x, y, 32, maxDepth + 1, 1, termbox.ColorCyan)
    // Шапка данных
    packUI.Tprint(x + 2, y + 1, coldef, coldef, "# ")
    packUI.Tprint(x + 2 + column, y + 1, coldef, coldef, " Цена ")
    packUI.Tprint(x + 2 + column + 13, y + 1, coldef, coldef, " Объем ")
    // Данные
    iPrice := 0
    for _, order := range orders {
        if iPrice <= maxDepth - 1 {
            packUI.Tprint(x + 2, y + 2 + iPrice, coldef, coldef, (" " + strconv.FormatInt(int64(iPrice), 10) + " "))
            packUI.Tprint(x + 2 + column, y + 2 + iPrice, coldef, coldef, (" " + strconv.FormatFloat(order.Price,'f', -1, 64) + " "))
            packUI.Tprint(x + 2 + column + 13, y + 2 + iPrice, coldef, coldef, (" " + strconv.FormatFloat(order.Amount,'f', -1, 64) + " "))
        }
        iPrice++
        line++
    }
    height = maxDepth + 2
    width = 32
    return
}
