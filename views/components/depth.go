package components

import (
    // "sort"
    "strings"
    "strconv"

    "../borders"
    packUI "../../ui"
    "../../market/core"
    "github.com/nsf/termbox-go"
)

func ASKsDraw(x, y int, depth *core.Depth, symbol, exchange string, maxheight int) (height, width int) {
    return OrdersDraw(x, y, depth.GetAsks(), symbol, exchange, "ask", maxheight)
}

func BIDsDraw(x, y int, depth *core.Depth, symbol, exchange string, maxheight int) (height, width int) {
    return OrdersDraw(x, y, depth.GetBids(), symbol, exchange, "bid", maxheight)
}

func OrdersDraw(x, y int, orders []*core.Order, symbol, exchange, direction string, maxheight int) (height, width int) {
    column := 5
    height = len(orders) + 2
    if height > maxheight {
        height = maxheight + 2
    }
    // Стенки
    borders.GroupBoxDraw(strings.ToUpper(direction) + "/" + exchange + "/" + strings.ToUpper(symbol) + "/" + strconv.FormatInt(int64(len(orders)), 10), x, y, 32, height - 1, 1, termbox.ColorCyan)
    // Шапка данных
    packUI.Tprint(x + 2, y + 1, coldef, coldef, "# ")
    packUI.Tprint(x + 2 + column, y + 1, coldef, coldef, " Цена ")
    packUI.Tprint(x + 2 + column + 13, y + 1, coldef, coldef, " Объем ")
    // Данные
    for line, order := range orders {
        if line < height - 2 {
            packUI.Tprint(x + 2, y + 2 + line, coldef, coldef, (" " + strconv.FormatInt(int64(line), 10) + " "))
            packUI.Tprint(x + 2 + column, y + 2 + line, coldef, coldef, (" " + strconv.FormatFloat(order.Price,'f', -1, 64) + " "))
            packUI.Tprint(x + 2 + column + 13, y + 2 + line, coldef, coldef, (" " + strconv.FormatFloat(order.Amount,'f', -1, 64) + " "))
        }
    }
    width = 32
    return
}
// ask - цена покупки
// bid - цена продажи
