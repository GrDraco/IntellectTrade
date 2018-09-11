package main

import (
    "testing"
    "./utilities"
)

func TestSearchIndex(t *testing.T) {
    test := new(TestLocation)
    test.Init("Test SearchIndex")
    test.Location.MethodName = "TestSearchIndex"
    _, found := utilities.SearchIndex("Aarhus", "A");
    if !found {
        t.Error("SearchIndex NOT WORK")
    }
}
