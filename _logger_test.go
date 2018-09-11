package main

import (
    "testing"
)

func TestPrintLog(t *testing.T) {
    test := new(TestLocation)
    test.Init("test logger")
    test.ViewDebug = true
    test.TestMethod()
}
