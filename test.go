package main

import "fmt"

func _main() {
    for i:=0; i<10; i++ {
        go func (i int) {
            fmt.Println("Got", i)
        }(i)
    }
    fmt.Scanln()
}
