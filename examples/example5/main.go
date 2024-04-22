package main

import (
    "os"
    "os/user"
    "os/exec"
    "net/http"
    "github.com/rakyll/statik/fs"
    "github.com/kerwin612/hybrid-launcher"
    _ "github.com/kerwin612/hybrid-launcher/examples/example5/statik"
)

func main() {

    http.HandleFunc("/exit", func (w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
        launcher.Exit()
    })

    statikFS, err := fs.New()
    if err != nil {
        panic(err)
    }

    myself, error := user.Current()
    if error != nil {
        panic(error)
    }
    homedir := myself.HomeDir + "/.hle/"
    if err := os.MkdirAll(homedir, 0775); err != nil {
        panic(err)
    }

    c := launcher.DefaultConfig()
    c.Pid = homedir + ".pid"
    c.Title = "Example"
    c.Tooltip = "Hybrid Launcher Example"
    c.RootHandler = http.StripPrefix("/", http.FileServer(statikFS))
    c.OpenWith = func(url string) {
        exec.Command("explorer", []string{url}...).Start()
    }
    launcher.StartWithConfig(c)

}
