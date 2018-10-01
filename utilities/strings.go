package utilities

import (
    "golang.org/x/text/search"
    "golang.org/x/text/language"
)

func SearchIndex(str string, substr string) (int, bool) {
    m := search.New(language.English, search.IgnoreCase)
    start, _ := m.IndexString(str, substr)
    if start == -1 {
        return start, false
    }
    return start, true
}
