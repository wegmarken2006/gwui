package main

import (
	"bufio"
	. "fmt"
	"os"

	gw "github.com/wegmarken2006/gwui"
)

func main() {

	gc := gw.GuiCfg{Port: 9000}
	body := gc.GWB5Init("gwui Test")
	//mandatory: callback on body
	body.Callback(func(string) {})
	gc.Body = &body

	var bt2 gw.Elem
	bt1 := gc.GWB5Button("btn-primary", "bt1", "START1")

	lb1 := gc.GWB5Label("lb1", "Label 1")

	ta1 := gc.GWB5TextArea("ta1", 12)
	//mandatory: callback on textarea
	ta1.Callback(func(string) {})

	bt2 = gc.GWB5Button("btn-primary", "bt2", "START2")
	bt2.Callback(func(string) { ta1.WriteTextArea("From Button 2\n") })

	i1 := gc.GWB5InputText("i1")
	i1.Callback(func(value string) {
		text := Sprintf("From Input field: %s\n", value)
		ta1.WriteTextArea(text)
	})
	//gc.GWSetBackgroundColor(&i1, "BLUE")
	gc.GWSetBackgroundColor(&lb1, "#ffe6e6")
	gc.GWSetColor(&lb1, "red")
	gc.GWSetFontSize(&lb1, "xx-large")

	r1 := gc.GWB5Row("r1")
	r2 := gc.GWB5Row("r2")
	c1 := gc.GWB5Col("c1")
	c2 := gc.GWB5Col("c2")
	c3 := gc.GWB5Col("c3")
	c4 := gc.GWB5Col("c4")
	c5 := gc.GWB5Col("c5")
	p1 := gc.GWParagraph("p1")
	c1.Add(lb1)
	c2.Add(bt1)

	c3.Add(i1)
	p1.Add(lb1)
	p1.Add(bt1)
	p1.Add(i1)
	r1.Add(p1)

	c4.Add(bt2)
	c5.Add(ta1)
	r2.Add(c4)
	r2.Add(c5)

	bt1.Callback(func(string) {
		gc.GWChangeText(lb1, "Text changed")
		gc.GWChangeBackgroundColor(i1, "red")
		gc.GWChangeFontSize(i1, "xx-large")
	})

	//body.Add(r1)
	//body.Add(r2)

	tabs := gc.GWB5Tabs([]string{"tb1", "tb2"}, []string{"tab1", "tab2"})

	tabs.SubElems[0].Add(r1)
	tabs.SubElems[1].Add(r2)
	body.Add(tabs)
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
