package main

import (
    "github.com/kerwin612/hybrid-launcher"
)

func main() {
    l, err := launcher.New()
    if err != nil {
        panic(err)
    }
    panic(l.StartAndOpen())
}
