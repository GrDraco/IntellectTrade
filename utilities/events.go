package utilities

import (
    "reflect"
    "runtime"
    //"fmt"
    //"../market/core"
)

const (
    EVENT_ON = "event_on"
)

type Events struct {
    events map[string]map[string]func(string, []interface{}, func(string))
}

func GetFunctionName(i interface{}) string {
    return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func (e *Events) initEvent(event string) {
    if len(e.events) == 0 {
        e.events = make(map[string]map[string]func(string, []interface{}, func(string)))
    }
    e.events[event] = make(map[string]func(string, []interface{}, func(string)))
}

func (e *Events) AddAction(event string, action func(string, []interface{}, func(string))) {
    e.initEvent(event)
    actionName := GetFunctionName(action)
    e.events[event][actionName] = action
}

func (e *Events) DelAction(event string, action func(string, []interface{}, func(string))) {
    delete(e.events[event], GetFunctionName(action))
}

func (e *Events) On(event string, params []interface{}, callback func(string)) {
    for _, action := range e.events[EVENT_ON] {
        action(event, params, callback)
    }
    for _, action := range e.events[event] {
        action(event, params, callback)
    }
}
