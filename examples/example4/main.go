package main

import (
    "os"
    "os/user"
    "net/http"
    "github.com/rakyll/statik/fs"
    "github.com/getlantern/systray"
    "github.com/kerwin612/hybrid-launcher"
    _ "github.com/kerwin612/hybrid-launcher/examples/example4/statik"
)

func main() {

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

    launcher.StartWithConfig(c)

}
