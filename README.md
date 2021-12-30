# gwui
## UI library based on [Bootstrap 5.1.3](https://getbootstrap.com/) and [Bootstrap Icons 1.7.2](https://icons.getbootstrap.com/).
### UI is served on http://localhost, then the default browser can be automatically launched.
### HTTP POST used for UI -> logic messages.
### Websockets used for logic -> UI messages. 
### Intended for personal utility tools.

### Still  under development.
###
### App folder structure:
```
|-myapp
|   myapp.go
|   go.mod
|----static
|    | web2.css
|    |----bootstrap
|    |----bootstrap-icons
```
### Copy the static folder from the [examples](./examples) here.
### Optionally, use web2.css for further customizations.
###
### To compile the examples:
```
cd example\ksink
go mod tidy
go build .
```

### Minimal example:
```go
package main

import (
    . "fmt"

	gw "github.com/wegmarken2006/gwui"
)

func main() {
	gc := gw.GuiCfg{Port: 9000, BrowserStart: true}
	body := gc.GWB5Init("mini")

	//mandatory: callback on body
	body.Callback(func(string, int) {})
	gc.Body = &body

	bt1 := gc.GWB5ButtonNew("bt1", "primary", "Count")
	lb1 := gc.GWB5LabelNew("lb1", "0")

	count := 0
	bt1.Callback(func(string, int) {
		count++
		text := Sprintf("%d", count)
		lb1.ChangeText(text)
	})

	body.Add(lb1)
	body.Add(bt1)

	gc.GWClose(body)
	gc.GWRun()

	gc.GWWaitKeyFromCOnsole()
}
```
