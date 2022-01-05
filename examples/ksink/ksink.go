package main

import (
	. "fmt"
	"time"

	gw "github.com/wegmarken2006/gwui"
)

func main() {

	gc := gw.GuiCfg{BrowserStart: true, PlotIncluded: true}
	body := gc.GWB5Init("gwui test")
	//mandatory: callback on body
	body.Callback(func(string, int) {})
	gc.Body = &body
	body.SetBackgroundImage("abstract.jpg", 95)

	tabs := gc.GWB5TabsNew([]string{"tb1", "tb2", "tab3"}, []string{"tab1", "tab2", "tab3"})

	//create first tab elements

	bt1 := gc.GWB5ButtonWithIconNew("bt1", "primary", "brush", "Change")
	lb1 := gc.GWB5LabelNew("lb1", "Change text color")

	bt2 := gc.GWB5ButtonWithIconNew("bt2", "secondary", "brush-fill", "Change")
	lb2 := gc.GWB5LabelNew("lb2", "Change text background color")

	bt3 := gc.GWB5ButtonWithIconNew("bt3", "success", "box-arrow-up-right", "Change")
	lb3 := gc.GWB5LabelNew("lb3", "Change text size")

	lb4 := gc.GWB5LabelNew("lb4", "Input text with Modal")
	it4 := gc.GWB5InputTextNew("i4")

	lb5 := gc.GWB5LabelNew("lb5", "Select font family")

	dd5 := gc.GWB5DropDownNew("dd5", "warning",
		"Font Family", []string{"arial", "verdana", "monospace"})

	lb10 := gc.GWB5LabelNew("lb10", "Another tab")

	md1 := gc.GWB5ModalNew("m1bt1", "m2bt2", "TEXT INPUT", "Are you sure", "yes", "no")

	ta1 := gc.GWB5TextAreaNew("ta1", 12)
	//mandatory: callback on textarea
	ta1.Callback(func(string, int) {})

	ta1.SetBackgroundColor("#ffe6e6")
	ta1.SetColor("blue")
	ta1.SetFontSize("small")
	ta1.SetFontFamily("monospace")

	cd1 := gc.GWB5CardNew("cd1", "Kitchen Sink", "Elements")
	cd1.SetBackgroundColor("#eeffee")

	// callbacks

	//first modal button
	md1.SubElems[0].Callback(func(string, int) {
		md1.ChanBool1 <- true
	})

	//second modal button
	md1.SubElems[1].Callback(func(string, int) {
		md1.ChanBool1 <- false
	})

	it4.Callback(func(strValue string, intValue int) {
		md1.B5ModalShow()
		yes := <-md1.ChanBool1
		if yes {
			text := Sprintf("From Input field: %s\n", strValue)
			ta1.WriteTextArea(text)
		}
	})

	dd5.Callback(func(strValue string, intValue int) {
		ta1.ChangeFontFamily(strValue)
	})

	bt1.Callback(func(string, int) {
		ta1.ChangeColor("red")
		bt1.ChangeText("Changed")
		bt1.ChangeToDisable()
	})

	bt2.Callback(func(string, int) {
		ta1.ChangeBackgroundColor("#66ffff")
		bt2.ChangeText("Changed")
	})

	bt3.Callback(func(string, int) {
		ta1.ChangeFontSize("large")
		bt3.ChangeText("Changed")
	})

	// fisrt tab
	// place elements in a grid
	// most external must be: container, row

	ct1 := gc.GWB5ContainerNew("ct1")
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

	ct1.Add(r1) //container
	cd1.Add(ct1)

	//Fisrt tab content
	tabs.SubElems[0].Add(cd1)

	//Second tab

	img1 := gc.GWImageNew("img1", "abstract.jpg", 50, 50)
	img2 := gc.GWImageNew("img2", "abstract.jpg", 100, 100)
	img3 := gc.GWImageNew("img3", "abstract.jpg", 200, 200)
	img4 := gc.GWImageNew("img3", "abstract.jpg", 300, 100)

	rs1 := gc.GWB5RangeSliderNew("rs1", 50.0, 0.0, 100.0, 1.0)

	pb1 := gc.GWB5PillBadgeNew("pb1", "danger", "50")

	rd1Choices := []string{"do nothing", "abstract.jpg", "sunset.jpg"}
	rd1 := gc.GWB5RadioNew([]string{"rd1", "rd2", "rd3"},
		rd1Choices, 1)

	fi1 := gc.GWB5FileInputNew("fi1", "Open")
	lb11 := gc.GWB5LabelNew("lb11", "No file")

	rs1.Callback(func(strValue string, intValue int) {
		pb1.ChangeText(strValue)
	})

	rd1.Callback(func(strValue string, intValue int) {
		if intValue > 0 {
			img4.ChangeImage(rd1Choices[intValue])
		}
	})

	fi1.Callback(func(strValue string, intValue int) {
		lb11.ChangeText(strValue)
	})

	//Grid: container, row(s), ...
	ct2 := gc.GWB5ContainerNew("ct2")
	r1t2 := gc.GWB5RowNew("r1t2")
	c1t2 := gc.GWB5ColSpanNew("c1t2", 2)
	c2t2 := gc.GWB5ColNew("c1t2")

	c1t2.Add(rd1)
	c2t2.Add(img1)
	c2t2.Add(img2)
	c2t2.Add(img3)
	c2t2.Add(img4)
	r1t2.Add(c1t2)
	r1t2.Add(c2t2)
	ct2.Add(r1t2)

	r2t2 := gc.GWB5RowNew("r2t2")
	c3t2 := gc.GWB5ColSpanNew("c3t2", 2)
	c4t2 := gc.GWB5ColNew("c4t2")
	c3t2.Add(pb1)
	c4t2.Add(rs1)
	r2t2.Add(c3t2)
	r2t2.Add(c4t2)
	ct2.Add(r2t2)

	r3t2 := gc.GWB5RowNew("r3t2")
	c5t2 := gc.GWB5ColNew("c5t2")
	c6t2 := gc.GWB5ColSpanNew("c6t2", 2)

	c5t2.Add(fi1)
	c6t2.Add(lb11)
	r3t2.Add(c5t2)
	r3t2.Add(c6t2)
	ct2.Add(r3t2)
	tabs.SubElems[1].Add(ct2)
	tabs.SubElems[1].SetBackgroundColor("white")

	//Third tab

	x2 := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	xs2 := []string{"aa", "bb", "cc", "dd", "ee"}
	ys2 := xs2
	y2Original := []float64{1.0, 2.0, 4.0, 8.0, 16.0}
	y2 := []float64{1.0, 2.0, 4.0, 8.0, 16.0}
	pl1 := gc.GWPlyPlotLine("pl1", x2, y2, "chart1", "xaxis", "yaxis", 300, 250)
	pl2 := gc.GWPlyPlotScatter("pl2", x2, y2, "chart2", "xaxis", "yaxis", 300, 250)
	pl3 := gc.GWPlyPlotVBar("pl3", xs2, y2, "chart3", "xaxis", "yaxis", 300, 250)
	pl4 := gc.GWPlyPlotHBar("pl4", y2Original, ys2, "chart4", "xaxis", "yaxis", 300, 250)
	rs2 := gc.GWB5RangeSliderNew("rs2", 1.0, 1.0, 100.0, 1.0)
	rs2.Callback(func(strValue string, intValue int) {
		for ind, yElem := range y2Original {
			y2[ind] = yElem + float64(intValue)
			gc.GWPlyPlotRedrawY(pl1, y2)
		}
	})

	ct3 := gc.GWB5ContainerNew("ct3")
	r1t3 := gc.GWB5RowNew("r1t3")
	//r2t3 := gc.GWB5RowNew("r2t3")
	c1t3 := gc.GWB5ColSpanNew("c1t3", 2)
	c2t3 := gc.GWB5ColNew("c1t3")
	c3t3 := gc.GWB5ColNew("c1t3")
	c1t3.Add(lb10)
	c2t3.Add(pl1)
	c2t3.Add(pl2)
	c2t3.Add(rs2)
	c3t3.Add(pl3)
	c3t3.Add(pl4)
	r1t3.Add(c1t3)
	r1t3.Add(c2t3)
	r1t3.Add(c3t3)
	ct3.Add(r1t3)

	tabs.SubElems[2].Add(ct3)
	tabs.SubElems[2].SetBackgroundColor("white")

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
