package main

import (
	. "fmt"
	"time"

	gw "github.com/wegmarken2006/gwui"
)

func main() {

	gc := gw.GuiCfg{Port: 9000}
	body := gc.GWB5Init("gwui test")
	//mandatory: callback on body
	body.Callback(func(string) {})
	gc.Body = &body

	bt1 := gc.GWB5Button("btn-primary", "bt1", "Change")
	lb1 := gc.GWB5Label("lb1", "Change text color")

	bt2 := gc.GWB5Button("btn-secondary", "bt2", "Change")
	lb2 := gc.GWB5Label("lb2", "Change text background color")

	bt3 := gc.GWB5Button("btn-success", "bt3", "Change")
	lb3 := gc.GWB5Label("lb3", "Change text size")

	lb4 := gc.GWB5Label("lb4", "Input text")
	it4 := gc.GWB5InputText("i4")

	lb5 := gc.GWB5Label("lb5", "Select font family")
	dd5 := gc.GWB5DropDown("btn-warning", "dd5",
		"Font Family", []string{"arial", "verdana", "monospace"})

	lb10 := gc.GWB5Label("lb5", "Another tab")

	ta1 := gc.GWB5TextArea("ta1", 12)
	//mandatory: callback on textarea
	ta1.Callback(func(string) {})

	it4.Callback(func(value string) {
		text := Sprintf("From Input field: %s\n", value)
		ta1.WriteTextArea(text)
	})

	gc.GWSetBackgroundColor(&ta1, "#ffe6e6")
	gc.GWSetColor(&ta1, "blue")
	gc.GWSetFontSize(&ta1, "small")
	gc.GWSetFontFamily(&ta1, "monospace")

	bt1.Callback(func(string) {
		gc.GWChangeColor(ta1, "red")
		gc.GWChangeText(bt1, "Changed")
	})
	bt2.Callback(func(string) {
		gc.GWChangeBackgroundColor(ta1, "#66ffff")
		gc.GWChangeText(bt2, "Changed")
	})
	bt3.Callback(func(string) {
		gc.GWChangeFontSize(ta1, "large")
		gc.GWChangeText(bt3, "Changed")
	})

	r1 := gc.GWB5Row("r1")

	r11 := gc.GWB5Row("r11")
	r12 := gc.GWB5Row("r12")
	r13 := gc.GWB5Row("r13")
	r14 := gc.GWB5Row("r14")
	r15 := gc.GWB5Row("r15")
	c1 := gc.GWB5Col("c1")
	c2 := gc.GWB5Col("c2")
	c11 := gc.GWB5Col("c11")
	c12 := gc.GWB5Col("c12")
	c21 := gc.GWB5Col("c21")
	c22 := gc.GWB5Col("c22")
	c31 := gc.GWB5Col("c31")
	c32 := gc.GWB5Col("c32")
	c41 := gc.GWB5Col("c41")
	c42 := gc.GWB5Col("c42")
	c51 := gc.GWB5Col("c51")
	c52 := gc.GWB5Col("c52")
	c11.Add(lb1)
	c12.Add(bt1)
	c21.Add(lb2)
	c22.Add(bt2)
	c31.Add(lb3)
	c32.Add(bt3)
	c41.Add(lb4)
	c42.Add(it4)
	c51.Add(lb5)
	c52.Add(dd5)
	r11.Add(c11)
	r11.Add(c12)
	r12.Add(c21)
	r12.Add(c22)
	r13.Add(c31)
	r13.Add(c32)
	r14.Add(c41)
	r14.Add(c42)
	r15.Add(c51)
	r15.Add(c52)
	c1.Add(r11)
	c1.Add(r12)
	c1.Add(r13)
	c1.Add(r14)
	c1.Add(r15)

	c2.Add(ta1)

	r1.Add(c1)
	r1.Add(c2)

	tabs := gc.GWB5Tabs([]string{"tb1", "tb2"}, []string{"tab1", "tab2"})

	//Fisrt tab content
	tabs.SubElems[0].Add(r1)

	//Second tab content
	tabs.SubElems[1].Add(lb10)
	body.Add(tabs)

	gc.GWClose(body)
	gc.GWRun()

	//Continuous write in textarea
	timeD := time.Duration(5000) * time.Millisecond
	time.Sleep(timeD)
	go func() {
		ind := 0
		for {
			ind++
			text := Sprintf("%d: All work and no play ...\n", ind)
			ta1.WriteTextArea(text)
			timeD := time.Duration(3000) * time.Millisecond
			time.Sleep(timeD)
		}

	}()

	gc.GWWaitKeyFromCOnsole()

}