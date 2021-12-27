# gwui
## UI library based on Bootstrap 5.1.3

## Minimal example:
```go
package main

import (
    . "fmt"

	gw "github.com/wegmarken2006/gwui"
)

func main() {
	gc := gw.GuiCfg{Port: 9000}
	body := gc.GWB5Init("mini")

	//mandatory: callback on body
	body.Callback(func(string) {})
	gc.Body = &body

	bt1 := gc.GWB5Button("btn-primary", "bt1", "Count")
	lb1 := gc.GWB5Label("lb1", "0")

	count := 0
	bt1.Callback(func(string) {
		count++
		text := Sprintf("%d", count)
		gc.GWChangeText(lb1, text)
	})

	body.Add(lb1)
	body.Add(bt1)

	gc.GWClose(body)
	gc.GWRun()

	gc.GWWaitKeyFromCOnsole()
}
```

## App folder structure:
```
|-myapp
|   myapp.go
|   go.mod
|----static
|    | web2.css
|    |----bootstrap
```
## Copy the static folder from the [examples](./examples) here.
## Optionally, use web2.css for further customizations.
##
## To compile the examples:
```
cd example\ksink
go mod tidy
go build .
```
