# HybridLauncher

## API
```golang
package main

import (
    "net/http"
    "github.com/kerwin612/hybrid-launcher"
)

func main() {

    lch, err := launcher.New()
    if err != nil {
        panic(err)
    }

    http.HandleFunc("/exit", func (w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
        lch.Exit()
    })

    panic(lch.StartAndOpen())

}
```  
For detailed instructions, see: [examples](./examples/)

## Try the example app
```bash
git clone git@github.com:kerwin612/hybrid-launcher.git
cd hybrid-launcher/examples/example1
go run main.go
```  

## Credits
* [getlantern/systray](https://github.com/getlantern/systray)
