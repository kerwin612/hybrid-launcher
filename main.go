package launcher

import (
    "os"
    "fmt"
    "net"
    "time"
    "syscall"
    "strconv"
    "runtime"
    "os/user"
    "os/exec"
    "net/http"
    "io/ioutil"
    "path/filepath"
    "github.com/getlantern/systray"
)

type OpenWithFunc func(string)

type Config struct {
    Port        int
    Open        bool
    Pid         string
    Icon        []byte
    Title       string
    Tooltip     string
    OpenWith    OpenWithFunc
    TrayOnReady func()
    RootHandler http.Handler
}

var pid string
var iconData []byte
var defaultIcon []byte = IconData
var defaultTitle string = "Hybrid Launcher"
var defaultTooltip string = "Hybrid Launcher Application"
var defaultOpenWith OpenWithFunc = func(url string) {
    var cmd string
    var args []string

    switch runtime.GOOS {
        case "windows":
            cmd = "cmd"
            args = []string{"/c", "start"}
        case "darwin":
            cmd = "open"
        default: // "linux", "freebsd", "openbsd", "netbsd"
            cmd = "xdg-open"
    }

    args = append(args, url)
    cmd_instance := exec.Command(cmd, args...)
    cmd_instance.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
    cmd_instance.Start()
}
var defaultRootHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello Hybrid Launcher")
})

func DefaultConfig() *Config {
    pid := _pid()
    return &Config{
        Port: 0,
        Pid: *pid,
        Open: true,
        TrayOnReady: nil,
        Icon: defaultIcon,
        Title: defaultTitle,
        Tooltip: defaultTooltip,
        OpenWith: defaultOpenWith,
        RootHandler: defaultRootHandler,
    }
}

func Exit() {
    systray.Quit()
    os.Remove(pid)
    go os.Exit(0)
}

func Addr(pid *string) *string {
    if pid == nil || *pid == "" {
        pid = _pid()
    }
    if _, err := os.Stat(*pid); err == nil {
        data, err := ioutil.ReadFile(*pid)
        if err != nil {
            panic(err)
        }
        addr := string(data)
        req, err := http.NewRequest("HEAD", addr, nil)
        if err != nil {
            panic(err)
        }
        client := &http.Client{ Timeout: time.Second * 1 }
        resp, err := client.Do(req)
        if err == nil && resp != nil {
            return &addr
        }
        os.Remove(*pid)
    }
    return nil
}

func Start() {
    StartWithConfig(DefaultConfig())
}

func StartWithConfig(c *Config) {

    if c == nil {
        Start()
        return
    }

    pid = c.Pid
    if pid == "" {
        pid = *_pid()
    }

    OpenWith := c.OpenWith
    if OpenWith == nil {
        OpenWith = defaultOpenWith
    }

    open := c.Open
    port := c.Port
    addr := Addr(&pid)

    if addr != nil && *addr != "" {
        if (open) {
            go OpenWith(*addr)
        }
        time.Sleep(time.Second * 1)
        os.Exit(0)
    }

    listener, err := net.Listen("tcp", ":" + strconv.Itoa(port))
    if err != nil {
        panic(err)
    }

    _addr := fmt.Sprintf("%s%d", "http://localhost:", listener.Addr().(*net.TCPAddr).Port)

    if err := os.MkdirAll(filepath.Dir(pid), 0775); err != nil {
        panic(err)
    }
    file, err := os.Create(pid)
    file.WriteString(_addr)

    go systray.Run(func() {

        SetIcon(c.Icon)
        SetTitle(c.Title)
        SetTooltip(c.Tooltip)

        rootHandler := c.RootHandler
        if rootHandler == nil {
            rootHandler = defaultRootHandler
        }

        http.Handle("/", rootHandler)

        http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "image/x-icon")
            w.WriteHeader(http.StatusOK)
            w.Write(iconData)
        })

        trayOnReady := c.TrayOnReady
        if trayOnReady == nil {
            trayOnReady = func() {
                mShow := systray.AddMenuItem("Show", "Show the app")
                mQuit := systray.AddMenuItem("Quit", "Quit the app")
                go func() {
                    for {
                        select {
                            case <-mShow.ClickedCh:
                                go OpenWith(_addr)
                            case <-mQuit.ClickedCh:
                                Exit()
                                return
                        }
                    }
                }()
            }
        }

        trayOnReady()

        if (open) {
            go OpenWith(_addr)
        }

    }, Exit)

    panic(http.Serve(listener, nil))

}

func SetIcon(icon []byte) {
    iconData = icon
    if iconData == nil {
        iconData = defaultIcon
    }
    systray.SetIcon(iconData)
}

func SetTitle(_title string) {
    title := _title
    if title == "" {
        title = defaultTitle
    }
    systray.SetTitle(title)
}

func SetTooltip(_tooltip string) {
    tooltip := _tooltip
    if tooltip == "" {
        tooltip = defaultTooltip
    }
    systray.SetTooltip(tooltip)
}

func _pid() *string {
    myself, error := user.Current()
    if error != nil {
        panic(error)
    }
    homedir := myself.HomeDir + "/.hl/"
    pid = homedir + ".pid"
    return &pid
}
