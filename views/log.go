package views

import (
    "strconv"
    "reflect"
    packUI "../ui"
    "../market/core"
    "../market/constants"

    "github.com/nsf/termbox-go"
)

const (
    CONSOLE_LOG = "log"
    CONSOLE_LOG_DISCRIPTION = "Онлайн поток работы программы"
)

// Функция отрисовки встроенной консоли логов
func logDraw(console *packUI.Console, bodyX, bodyY int) {
    // fmt.Println("redraw", console.Values["Action"])
    y := bodyY - 1
    packUI.Fill(bodyX + 2, 2, 1, y - 1, termbox.Cell{Ch: '>'})
    if console.Controls["Text"] != nil {
        newValue := console.Controls["Text"].(*packUI.Text)
        // Сдвигаем всю историю комманд вверх
        console.Controls[strconv.FormatInt(int64(y), 10)] = newValue
        for y = 1 ; y <= bodyY - 1; y++ {
            console.Controls[strconv.FormatInt(int64(y), 10)] = console.Controls[strconv.FormatInt(int64(y + 1), 10)]
        }
        // Отображаем историю комманд
        for y = 1 ; y <= bodyY - 1; y++ {
            control := console.Controls[strconv.FormatInt(int64(y), 10)]
            if control != nil {
                text := valueToString(control.(*packUI.Text))
                if text != nil {
                    packUI.Tprint(bodyX + 4, y + 1, text.FG, text.BG, text.Value.(string))
                }
            }
        }
    }
}
// Функция очистки встроенной консоли логов
func logClear(console *packUI.Console, bodyX, bodyY int) {
    w, _ := termbox.Size()
    packUI.Fill(bodyX + 4, 2, w - bodyX + 4, bodyY - 2, termbox.Cell{Ch: ' '})
}
// Функция обработки встроенной консоли логов
func logAction(console *packUI.Console, data []interface{}) {
    // fmt.Println(action)
    if data == nil {
        return
    }
    if len(data) >= 1 {
        console.Controls["Text"] = data[0]
    }
}

func Log() *packUI.Console  {
    return packUI.NewConsole(CONSOLE_LOG, CONSOLE_LOG_DISCRIPTION, false, logDraw, logClear, logAction)
}

func valueToString(text *packUI.Text) *packUI.Text {
    if reflect.TypeOf(text.Value).Kind() == reflect.String {
        return text
    }
    if reflect.TypeOf(text.Value).Elem().Name() == "Signal" {
        signal := text.Value.(*core.Signal)
        var res *packUI.Text
        switch signal.Entity {
            case constants.ENTITY_TICK:
                res = &packUI.Text { Value: text.Value.(*core.Signal).TickStr(true), BG: termbox.ColorDefault }
            case constants.ENTITY_DEPTH:
                res = &packUI.Text { Value: text.Value.(*core.Signal).DepthStr(true), BG: termbox.ColorDefault }
            case constants.ENTITY_CANDLE:
                res = &packUI.Text { Value: text.Value.(*core.Signal).CandlesStr(true), BG: termbox.ColorDefault }
        }
        if signal.TimeOut {
            res.FG = termbox.ColorRed
        } else {
            res.FG = termbox.ColorGreen
        }
        return res
    }
    if reflect.TypeOf(text.Value).Elem().Name() == "Message" {
        res := &packUI.Text { Value: text.Value.(*core.Message).FullString() }
        if text.Value.(*core.Message).Kind == core.MSG_WORNING {
            res.FG = termbox.ColorYellow
            res.BG = termbox.ColorYellow
        } else {
            res.FG = text.FG
            res.BG = text.BG
        }
        return res
    }
    if reflect.TypeOf(text.Value).Elem().Name() == "StrategyAction" {
        //TODO: Написать конвертацию действий стратегии для вывода в лог
    }
    return nil
}
