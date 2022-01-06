package main

import (
	. "fmt"

	"github.com/wegmarken2006/gwui"
)

func main() {
	gc := gwui.GuiCfg{BrowserStart: true}
	body := gc.Init("mini")

	bt1 := gc.ButtonNew("primary", "Count")
	lb1 := gc.LabelNew("0")

	count := 0
	bt1.Callback(func(string, int) {
		count++
		text := Sprintf("%d", count)
		lb1.ChangeText(text)
	})

	body.Add(lb1)
	body.Add(bt1)

	gc.Close(body)
	gc.Run()

	gc.WaitKeyFromCOnsole()
}
