package utilities

// import (
//     // "fmt"
// )

type Collection struct {
    // Наследуем события
    Events
    //
    Name string
    // Коллекия индикаторов
    Storage map[string]string
}

const (
    COLLECTION_EVENT_SET_VALUE = "set_value"
)

func (collection *Collection) initStorage() {
    if collection.Storage == nil {
        collection.Storage = make(map[string]string)
    }
    if len(collection.Storage) == 0 {
        collection.Storage = make(map[string]string)
    }
}

func (collection *Collection) SetValue(name string, value string) {
    collection.initStorage()
    collection.Storage[name] = value
    collection.On(COLLECTION_EVENT_SET_VALUE, []interface{} { collection.Name, name, value }, nil)
}
