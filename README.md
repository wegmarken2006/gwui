# gwui
## Example

``` go
package main

import (
	. "fmt"
	"bufio"
	"os"
	gw "github.com/wegmarken2006/gwui"
)


func main() {
	
	gc := gw.GuiCfg{}
	body := gc.GWInit("gwui Test")
	//mandatory: callback on body
	body.Callback(func(string) {})
	gc.Body = &body

	var bt2 gw.Elem
	bt1 := gc.GWButton("btn-primary", "bt1", "START1")

	lb1 := gc.GWLabel("lb1", "Label 1")

	ta1 := gc.GWTextArea("ta1", 12)
	//mandatory: callback on textarea
	ta1.Callback(func(string) {})

	bt2 = gc.GWButton("btn-primary", "bt2", "START2")
	bt2.Callback(func(string) { ta1.WriteTextArea("From Button 2\n") })

	i1 := gc.GWInputText("i1")
	i1.Callback(func(value string) { 
		text := Sprintf("From Input field: %s\n", value)
		ta1.WriteTextArea(text)
	})

	r1 := gc.GWRow("r1")
	r2 := gc.GWRow("r2")
	c1 := gc.GWCol("c1")
	c2 := gc.GWCol("c2")
	c3 := gc.GWCol("c3")
	c4 := gc.GWCol("c4")
	c5 := gc.GWCol("c5")
	c1.Add(lb1)
	c2.Add(bt1)
	c3.Add(i1)
	r1.Add(c1)
	r1.Add(c2)
	r1.Add(c3)

	c4.Add(bt2)
	c5.Add(ta1)
	r2.Add(c4)
	r2.Add(c5)

	bt1.Callback(func(string) {
		gc.GWChangeText(lb1, "Text changed")
		gc.GWChangeBackgroundColor(c1, "red")
	})

	body.Add(r1)
	body.Add(r2)

	gc.GWClose(body)

	gc.GWRun()

	reader := bufio.NewReader(os.Stdin)
	Println("Press:\n q<Enter> to exit")
	for {
			text, _ := reader.ReadString('\n')
			//cut final 0xd, 0xa
			text = text[:len(text)-2]
			switch text {
			case "q", "Q":
				os.Exit(0)
			}
		}

}
```
