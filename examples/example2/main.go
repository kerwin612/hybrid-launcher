package main

import (
    "net/http"
    "github.com/rakyll/statik/fs"
    "github.com/kerwin612/hybrid-launcher"
    _ "github.com/kerwin612/hybrid-launcher/examples/example2/statik"
)

func main() {

    statikFS, err := fs.New()
    if err != nil {
        panic(err)
    }

    c, err := launcher.DefaultConfig()
    if err != nil {
        panic(err)
    }

    c.Title = "Example"
    c.Tooltip = "Hybrid Launcher Example"
    c.RootHandler = http.StripPrefix("/", http.FileServer(statikFS))

    l, err := launcher.NewWithConfig(c)
    if err != nil {
        panic(err)
    }

    http.HandleFunc("/exit", func (w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
        l.Exit()
    })

    l.StartAndOpen()

}
