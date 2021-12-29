package main

import (
	. "fmt"
	"time"

	gw "github.com/wegmarken2006/gwui"
)

func main() {

	//create elements

	gc := gw.GuiCfg{Port: 9000, BrowserStart: true}
	body := gc.GWB5Init("gwui test")
	//mandatory: callback on body
	body.Callback(func(string) {})
	gc.Body = &body
	body.SetBackgroundColor("#ccffcc")

	bt1 := gc.GWB5ButtonNew("btn-primary", "bt1", "Change")
	lb1 := gc.GWB5LabelNew("lb1", "Change text color")

	bt2 := gc.GWB5ButtonNew("btn-secondary", "bt2", "Change")
	lb2 := gc.GWB5LabelNew("lb2", "Change text background color")

	bt3 := gc.GWB5ButtonNew("btn-success", "bt3", "Change")
	lb3 := gc.GWB5LabelNew("lb3", "Change text size")

	lb4 := gc.GWB5LabelNew("lb4", "Input text with Modal")
	it4 := gc.GWB5InputTextNew("i4")

	lb5 := gc.GWB5LabelNew("lb5", "Select font family")
	dd5 := gc.GWB5DropDownNew("btn-warning", "dd5",
		"Font Family", []string{"arial", "verdana", "monospace"})

	lb10 := gc.GWB5LabelNew("lb5", "Another tab")

	md1 := gc.GWB5ModalNew("m1bt1", "m2bt2", "TEXT INPUT", "Are you sure", "yes", "no")

	ta1 := gc.GWB5TextAreaNew("ta1", 12)
	//mandatory: callback on textarea
	ta1.Callback(func(string) {})

	ta1.SetBackgroundColor("#ffe6e6")
	ta1.SetColor("blue")
	ta1.SetFontSize("small")
	ta1.SetFontFamily("monospace")

	cd1 := gc.GWB5CardNew("cd1", "Kitchen Sink", "Elements")
	cd1.SetBackgroundColor("#eeffee")
	cd1.SetColor("green")

	tabs := gc.GWB5TabsNew([]string{"tb1", "tb2"}, []string{"tab1", "tab2"})

	// callbacks

	//first modal button
	md1.SubElems[0].Callback(func(string) {
		md1.ChanBool1 <- true
	})

	//second modal button
	md1.SubElems[1].Callback(func(string) {
		md1.ChanBool1 <- false
	})

	it4.Callback(func(value string) {
		md1.B5ModalShow()
		yes := <-md1.ChanBool1
		if yes {
			text := Sprintf("From Input field: %s\n", value)
			ta1.WriteTextArea(text)
		}

	})

	dd5.Callback(func(value string) {
		ta1.ChangeFontFamily(value)
	})

	bt1.Callback(func(string) {
		ta1.ChangeColor("red")
		bt1.ChangeText("Changed")
		bt1.ChangeToDisable()
	})

	bt2.Callback(func(string) {
		ta1.ChangeBackgroundColor("#66ffff")
		bt2.ChangeText("Changed")
	})

	bt3.Callback(func(string) {
		ta1.ChangeFontSize("large")
		bt3.ChangeText("Changed")
	})

	// place elements in a grid

	r1 := gc.GWB5RowNew("r1")
	r11 := gc.GWB5RowNew("r11")
	r12 := gc.GWB5RowNew("r12")
	r13 := gc.GWB5RowNew("r13")
	r14 := gc.GWB5RowNew("r14")
	r15 := gc.GWB5RowNew("r15")
	c1 := gc.GWB5ColNew("c1")
	c2 := gc.GWB5ColNew("c2")
	c11 := gc.GWB5ColNew("c11")
	c12 := gc.GWB5ColNew("c12")
	c21 := gc.GWB5ColNew("c21")
	c22 := gc.GWB5ColNew("c22")
	c31 := gc.GWB5ColNew("c31")
	c32 := gc.GWB5ColNew("c32")
	c41 := gc.GWB5ColNew("c41")
	c42 := gc.GWB5ColNew("c42")
	c51 := gc.GWB5ColNew("c51")
	c52 := gc.GWB5ColNew("c52")
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

	cd1.Add(r1)

	//Fisrt tab content
	tabs.SubElems[0].Add(cd1)

	//Second tab content
	tabs.SubElems[1].Add(lb10)

	//final body additions; add modal to body directly
	body.Add(tabs)
	body.Add(md1)

	gc.GWClose(body)
	gc.GWRun()

	//Processing simulation: continuous write in textarea
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
