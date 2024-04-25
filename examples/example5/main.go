package main

import (
    "os"
    "os/user"
    "os/exec"
    "syscall"
    "net/http"
    "github.com/rakyll/statik/fs"
    l "github.com/kerwin612/hybrid-launcher"
    _ "github.com/kerwin612/hybrid-launcher/examples/example5/statik"
)

var launcher *l.Launcher

func main() {

    http.HandleFunc("/exit", func (w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
        launcher.Exit()
    })

    statikFS, err := fs.New()
    if err != nil {
        panic(err)
    }

    myself, err := user.Current()
    if err != nil {
        panic(err)
    }
    homedir := myself.HomeDir + "/.hle/"
    if err := os.MkdirAll(homedir, 0775); err != nil {
        panic(err)
    }

    c, err := l.DefaultConfig()
    if err != nil {
        panic(err)
    }

    c.Pid = homedir + ".pid"
    c.Title = "Example"
    c.Tooltip = "Hybrid Launcher Example"
    c.RootHandler = http.StripPrefix("/", http.FileServer(statikFS))
    c.OpenWith = func(url string) {
        cmd_instance := exec.Command("cmd", []string{"/c", "start", url}...)
        cmd_instance.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
        cmd_instance.Start()
    }

    launcher, err = l.NewWithConfig(c)
    if err != nil {
        panic(err)
    }

    launcher.StartAndOpen()

}
