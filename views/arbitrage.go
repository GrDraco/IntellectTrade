package views

import (
    packUI "../ui"

    "github.com/nsf/termbox-go"
)

// Функция отрисовки стратегии arbitrage
func arbitrageDraw(console *packUI.Console, bodyX, bodyY int) {
    if console.Controls["Properties"] == nil {
        return
    }
    name := console.Controls["Properties"].(map[string]interface{})["test"].(string)
    packUI.Tprint(bodyX + 4, 2, coldef, coldef, name)
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
