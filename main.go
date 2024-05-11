package launcher

import (
    "os"
    "fmt"
    "net"
    "time"
    "errors"
    "strconv"
    "runtime"
    "os/user"
    "os/exec"
    "net/http"
    "io/ioutil"
    "path/filepath"
    "github.com/getlantern/systray"
)

type Config struct {
    Ip          string
    Port        int
    Pid         string
    Icon        []byte
    Title       string
    Tooltip     string
    TrayOnReady func()
    OpenWith    func(string)
    RootHandler http.Handler
}

type Launcher struct{
    config      Config
    listener    net.Listener
}

var defaultIcon []byte = IconData
var defaultTitle string = "Hybrid Launcher"
var defaultTooltip string = "Hybrid Launcher Application"
var defaultOpenWith func(string) = func(url string) {
    var cmd string

    switch runtime.GOOS {
        case "windows":
            cmd = "explorer"
        case "darwin":
            cmd = "open"
        default: // "linux", "freebsd", "openbsd", "netbsd"
            cmd = "xdg-open"
    }

    exec.Command(cmd, []string{url}...).Start()
}
var defaultPid func() (string, error) = func() (string, error) {
    cur, err := user.Current()
    if err != nil {
        return "", err
    }
    return filepath.Join(filepath.Join(cur.HomeDir, ".hl"), ".pid"), nil
}
var defaultRootHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello Hybrid Launcher")
})

func isStarted(pid string) (string, error) {
    _, err := os.Stat(pid)
    if err != nil {
        return "", err
    }

    data, err := ioutil.ReadFile(pid)
    if err != nil {
        return "", err
    }
    addr := string(data)
    req, err := http.NewRequest("HEAD", addr, nil)
    if err != nil {
        return "", err
    }
    client := &http.Client{ Timeout: time.Second * 1 }
    resp, err := client.Do(req)
    if err == nil && resp != nil {
        return addr, nil
    }
    os.Remove(pid)
    return "", err
}

func DefaultConfig() (*Config, error) {
    pid, err := defaultPid()
    if err != nil {
        return nil, err
    }
    return &Config{
        Ip: "",
        Port: 0,
        Pid: pid,
        TrayOnReady: nil,
        Icon: defaultIcon,
        Title: defaultTitle,
        Tooltip: defaultTooltip,
        OpenWith: defaultOpenWith,
        RootHandler: defaultRootHandler,
    }, nil
}

func New() (*Launcher, error) {
    cfg, err := DefaultConfig()
    if err != nil {
        return nil, err
    }
    return NewWithConfig(cfg)
}

func NewWithConfig(c *Config) (*Launcher, error) {

    if c == nil {
        return New()
    }

    cfg := *c

    if cfg.Pid == "" {
        pid, err := defaultPid()
        if err != nil {
            return nil, err
        }
        cfg.Pid = pid
    }

    if cfg.Icon == nil {
        cfg.Icon = defaultIcon
    }

    if cfg.Title == "" {
        cfg.Title = defaultTitle
    }

    if cfg.Tooltip == "" {
        cfg.Tooltip = defaultTooltip
    }

    if cfg.OpenWith == nil {
        cfg.OpenWith = defaultOpenWith
    }

    if cfg.RootHandler == nil {
        cfg.RootHandler = defaultRootHandler
    }

    addr, _ := isStarted(cfg.Pid)
    if addr != "" {
        return nil, errors.New("There are instances started, see address in pid file.")
    }

    ln, err := net.Listen("tcp", net.JoinHostPort(cfg.Ip, strconv.Itoa(cfg.Port)))
    if err != nil {
        return nil, err
    }

    if cfg.Ip == "" {
        cfg.Ip = "localhost"
    }

    cfg.Port = ln.Addr().(*net.TCPAddr).Port

    return &Launcher{
        config: cfg,
        listener: ln,
    }, nil

}

func (l *Launcher) start(isOpen bool) error {

    addr, _ := isStarted(l.config.Pid)
    if addr != "" {
        return errors.New("There are instances started, see address in pid file.")
    }

    if err := os.MkdirAll(filepath.Dir(l.config.Pid), 0775); err != nil {
        return err
    }

    file, err := os.Create(l.config.Pid)
    defer file.Close()
    if err != nil {
        return err
    }

    file.WriteString(l.Addr())

    go func(){
        panic(http.Serve(l.listener, nil))
    }()

    systray.Run(
        func() {

            l.SetIcon(l.config.Icon)
            l.SetTitle(l.config.Title)
            l.SetTooltip(l.config.Tooltip)

            http.Handle("/", l.config.RootHandler)

            http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("Content-Type", "image/x-icon")
                w.WriteHeader(http.StatusOK)
                w.Write(l.config.Icon)
            })

            if l.config.TrayOnReady == nil {
                l.config.TrayOnReady = func() {
                    mShow := systray.AddMenuItem("Show", "Show the app")
                    mQuit := systray.AddMenuItem("Quit", "Quit the app")
                    go func() {
                        for {
                            select {
                                case <-mShow.ClickedCh:
                                    go l.Open()
                                case <-mQuit.ClickedCh:
                                    l.Exit()
                                    return
                            }
                        }
                    }()
                }
            }

            l.config.TrayOnReady()

            if (isOpen) {
                go l.Open()
            }

        },
        func() {
            os.Remove(l.config.Pid)
            go os.Exit(0)
        },
    )

    return nil

}

func (l *Launcher) StartAndOpen() error {
    return l.start(true)
}

func (l *Launcher) Start() error {
    return l.start(false)
}

func (l *Launcher) Addr() string {
    return fmt.Sprintf("http://%s:%d", l.config.Ip, l.config.Port)
}

func (l *Launcher) Open() {
    go l.config.OpenWith(l.Addr())
}

func (l *Launcher) SetIcon(icon []byte) {
    iconData := icon
    if iconData == nil {
        iconData = defaultIcon
    }
    l.config.Icon = iconData
    systray.SetIcon(l.config.Icon)
}

func (l *Launcher) SetTitle(title string) {
    if title == "" {
        title = defaultTitle
    }
    l.config.Title = title
    systray.SetTitle(l.config.Title)
}

func (l *Launcher) SetTooltip(tooltip string) {
    if tooltip == "" {
        tooltip = defaultTooltip
    }
    l.config.Tooltip = tooltip
    systray.SetTooltip(l.config.Tooltip)
}

func (l *Launcher) Exit() {
    systray.Quit()
}
