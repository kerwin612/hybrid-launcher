package main

import (
    "os"
    "os/user"
    "net/http"
    "github.com/rakyll/statik/fs"
    "github.com/getlantern/systray"
    l "github.com/kerwin612/hybrid-launcher"
    _ "github.com/kerwin612/hybrid-launcher/examples/example4/statik"
)

var launcher *l.Launcher

func main() {

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
    c.Icon = IconData1
    c.Title = "Example"
    c.Tooltip = "Hybrid Launcher Example"
    c.RootHandler = http.StripPrefix("/", http.FileServer(statikFS))
    c.TrayOnReady = func() {

        mQuit := systray.AddMenuItem("Quit", "Quit the example app")
        go func() {
            <-mQuit.ClickedCh
            launcher.Exit()
        }()

        http.HandleFunc("/2cn", func (w http.ResponseWriter, r *http.Request) {
            mQuit.SetTitle("退出")
            mQuit.SetTooltip("退出程序")
            launcher.SetIcon(IconData2)
            launcher.SetTitle("混合应用")
            launcher.SetTooltip("混合应用示例")
            w.WriteHeader(200)
        })

    }

    launcher, err = l.NewWithConfig(c)
    if err != nil {
        panic(err)
    }

    launcher.StartAndOpen()

}
