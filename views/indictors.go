package views

import (
    "sort"
    "strconv"

    "./borders"
    "../market"
    "../utilities"
    packUI "../ui"
    "github.com/nsf/termbox-go"
)

const (
    CONTROL_INDICATORS = "indicators"

    CONSOLE_INDICATORS = "commands"
)

// Функция отрисовки индикаторов
func indicatorsDraw(console *packUI.Console, bodyX, bodyY int) {
    if console.Controls[CONTROL_INDICATORS] == nil {
        return
    }
    indicators := console.Controls[CONTROL_INDICATORS].(map[string]map[string]string)
    var keysGroup []string
    for key, _ := range indicators {
        keysGroup = append(keysGroup, key)
    }
    sort.Strings(keysGroup)
    k := 0
    for i:=len(keysGroup) - 1; i >= 0; i-- {
        var keysIndicator []string
        for key, _ := range indicators[keysGroup[i]] {
            keysIndicator = append(keysIndicator, key)
        }
        sort.Strings(keysIndicator)
        // Рисуем стенки формы
        borders.GroupBoxDraw(keysGroup[i], 3, bodyY - len(keysIndicator) - k, bodyX - 4, len(keysIndicator), -1 , termbox.ColorCyan)
        // colForm := termbox.ColorCyan
        for j:=len(keysIndicator) - 1; j >= 0; j-- {
            packUI.Tprint(5, bodyY - k, coldef, coldef, keysIndicator[j] + " ")
            packUI.Tprint(bodyX - len([]rune(indicators[keysGroup[i]][keysIndicator[j]])) - 3, bodyY - k, coldef, coldef, " " + indicators[keysGroup[i]][keysIndicator[j]])
            k++
        }
        k++
        k = k + 2
    }
}
// Функция очистки индикаторов
func indicatorsClear(console *packUI.Console, bodyX, bodyY int) {
    // _, h := termbox.Size()
    packUI.Fill(3, 2, bodyX - 3, bodyY - 1, termbox.Cell{Ch: ' '})
}
// Функция обработки индикаторов
func indicatorsAction(console *packUI.Console, data []interface{}) {
    indicators := make(map[string]map[string]string)
    settingsTitle := "Настройки"
    terminalTitle := "Терминал"
    // Записываем индикаторы по настройкм
    if settings != nil {
        indicators[settingsTitle] = make(map[string]string)
        indicators[settingsTitle]["Название"] = settings.Name
        indicators[settingsTitle]["Автозагрузка"] = strconv.FormatBool(settings.IsDefault)
    }
    // Записываем индикаторы по терминалу
    terminal := frame.Controls["terminal"]
    if terminal != nil {
        indicators[terminalTitle] = make(map[string]string)
        indicators[terminalTitle] = utilities.AppendMap(indicators[terminalTitle], terminal.(*market.Terminal).GetIndicators())
    }
    // Записываем индикаторы по биржам
    for _, exchange := range terminal.(*market.Terminal).Exchanges {
        indicators[exchange.Name] = make(map[string]string)
        indicators[exchange.Name] = utilities.AppendMap(indicators[exchange.Name], exchange.GetIndicators())
    }
    // Записываем индикаторы по стратегиям
    for _, strategy := range terminal.(*market.Terminal).Strategies {
        indicators[strategy.GetName()] = make(map[string]string)
        indicators[strategy.GetName()] = utilities.AppendMap(indicators[strategy.GetName()], strategy.GetIndicators())
    }
    console.Controls[CONTROL_INDICATORS] = indicators
}

func Indicators() *packUI.Console {
    return packUI.NewConsole(CONSOLE_INDICATORS, "", true, indicatorsDraw, indicatorsClear, indicatorsAction)
}
