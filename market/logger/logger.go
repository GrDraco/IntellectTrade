package logger

import (
    "log"
    "strings"
    "reflect"
)

type Location struct {
    // Свойства
    EntityName string   // Сущьность от которой это сообщение
    ClassName string    // Класс(структура) в которой вызывается это сообщение
    ViewDebug bool      //
    // Каналы
    Status chan string
}

func (location *Location) SetLocation(entityName string, variable interface{}) {
    log.SetFlags(0)
    location.Status = make(chan string)
    location.EntityName = entityName
    location.ClassName = reflect.TypeOf(variable).String()
}

func (location *Location) print(methodName string, message string, isErr bool, isDebug bool) {
    startMsg := "//LOG: "
    finishMsg := "/////\n\n"
    if isDebug {
        startMsg = "//DEBUG: "
        finishMsg = "///////\n\n"
    }
    msg := startMsg + strings.ToUpper(location.EntityName) + "(" + location.ClassName + "." + methodName + ")//\n" +
        "   " + message + "\n" + finishMsg
    if isErr {
        log.Fatal("\x1b[31;1m" + msg + "\x1b[0m")
    } else {
        log.Printf(msg)
    }
    go func(){
        location.Status<- msg
    }()
}

func (location *Location) IsInit() bool {
    if location.EntityName == "" || location.ClassName == "" {
        return false
    } else {
        return true
    }
}

func (location *Location) PrintLog(methodName string, message string, _isDebug ...bool) {
    isDebug := false
    if len(_isDebug) > 0 {
        isDebug = _isDebug[0]
    }
    if location.ViewDebug || !isDebug {
        if location.IsInit() {
            location.print(methodName, message, false, isDebug)
        } else {
            location.print(methodName, "NOT INIT LOCATION", true, isDebug)
        }
    }
}

func (location *Location) PrintError(methodName string, message string, _isDebug ...bool) {
    isDebug := false
    if len(_isDebug) > 0 {
        isDebug = _isDebug[0]
    }
    if location.ViewDebug || !isDebug {
        if location.IsInit() {
            location.print(methodName, message, true, isDebug)
        } else {
            location.print(methodName, "NOT INIT LOCATION", true, isDebug)
        }
    }
}
