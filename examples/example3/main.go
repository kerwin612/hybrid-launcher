package main

import (
    "io"
    "os"
    "log"
    "os/user"
    "net/http"
    "encoding/json"
    "github.com/rakyll/statik/fs"
    "github.com/getlantern/systray"
    "github.com/kerwin612/hybrid-launcher"
    _ "github.com/kerwin612/hybrid-launcher/examples/example3/statik"
)

var logger *log.Logger
var configFile string
var config struct {
    Port int
}

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

    lf, err := os.OpenFile(homedir + "log.txt", os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
    if err != nil {
        panic(err)
    }
    //defer lf.Close()

    logger = log.New(io.MultiWriter(lf, os.Stdout), "", log.LstdFlags)
    //logger.SetOutput(os.Stdout)

    configFile, err := os.Open(homedir + ".cfg")
    defer configFile.Close()
    port := 0
    if err == nil {
        jsonParser := json.NewDecoder(configFile)
        if err = jsonParser.Decode(&config); err == nil {
            logger.Println("config: ", config)
            port = config.Port
        } else {
            logger.Println("decode: ", err)
        }
    } else {
        logger.Println("open: ", err)
    }

    c := launcher.DefaultConfig()
    c.Pid = homedir + ".pid"
    c.Port = port
    c.Icon = IconData
    c.Title = "Example"
    c.Tooltip = "Hybrid Launcher Example"
    c.RootHandler = http.StripPrefix("/", http.FileServer(statikFS))
    c.TrayOnReady = func() {
        mQuitOrig := systray.AddMenuItem("Quit", "Quit the example app")
        go func() {
            <-mQuitOrig.ClickedCh
            launcher.Exit()
        }()
    }
    launcher.StartWithConfig(c)

}
