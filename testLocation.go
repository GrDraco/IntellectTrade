package main

import (
    "./market/logger"
)

type TestLocation struct {
    logger.Location
}

func (t *TestLocation) Init(entityName string) {
    t.Location.SetLocation(entityName, t)
}

func (t *TestLocation) TestMethod() {
    t.PrintLog("TestMethod", "Test Good Message DEBUG", true)
    t.PrintLog("TestMethod", "Test Good Message LOG")
    t.PrintError("TestMethod", "Test Error Message ERROR for DEBUG", true)
    t.PrintError("TestMethod", "Test Error Message ERROR for LOG")
}
