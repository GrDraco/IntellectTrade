package utilities

import (
    "golang.org/x/text/search"
    "golang.org/x/text/language"
    "strconv"
)

func SearchIndex(str string, substr string) (int, bool) {
    m := search.New(language.English, search.IgnoreCase)
    start, _ := m.IndexString(str, substr)
    if start == -1 {
        return start, false
    }
    return start, true
}

func FloatToString(input_num float64, modulo int) string {
    return strconv.FormatFloat(input_num, 'f', modulo, 64)
}
