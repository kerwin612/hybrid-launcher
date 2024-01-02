# HybridLauncher

## API
```golang
package main

import (
    "net/http"
    "github.com/kerwin612/hybrid-launcher"
)

func main() {

    http.HandleFunc("/exit", func (w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
        launcher.Exit()
    })

    launcher.Start()

}
```  
For detailed instructions, see: [example](./example/)

## Try the example app
```bash
git clone git@github.com:kerwin612/hybrid-launcher.git
cd hybrid-launcher/example
build.cmd
example
```  

## Credits
* [getlantern/systray](https://github.com/getlantern/systray)
