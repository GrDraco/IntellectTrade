package borders

import (
    packUI "../../ui"
    "github.com/nsf/termbox-go"
)

const coldef = termbox.ColorDefault

func GroupBoxDraw(title string, x, y, width, height, head int, colorBorder termbox.Attribute) {
    packUI.Fill(x, y, width, 1, termbox.Cell{Ch: '─', Fg: colorBorder})
    groupTitle := "┤" + title + "├"
    packUI.Tprint(x + 1, y, colorBorder, coldef, groupTitle)
    termbox.SetCell(x, y, '┌', colorBorder, coldef)
    termbox.SetCell(x + width, y, '┐', colorBorder, coldef)
    for j := 1; j <= height; j++ {
        if head > j - 1 && head >= 1 {
            termbox.SetCell(x, y + j, '│', colorBorder, coldef)
            termbox.SetCell(x + width, y + j, '│', colorBorder, coldef)
        } else {
            termbox.SetCell(x, y + j, '├', colorBorder, coldef)
            termbox.SetCell(x + width, y + j, '┤', colorBorder, coldef)
            packUI.Fill(x + 1, y + j, width - 1, 1, termbox.Cell{Ch: '╌', Fg: colorBorder})
        }
    }
    termbox.SetCell(x, y + height, '└', colorBorder, coldef)
    termbox.SetCell(x + width, y + height, '┘', colorBorder, coldef)
}
