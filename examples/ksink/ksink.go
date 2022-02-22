package main

import (
	. "fmt"
	"strings"
	"time"

	"github.com/wegmarken2006/gwui"
)

func main() {

	gc := gwui.GuiCfg{BrowserStart: true, PlotIncluded: true, ExitOnWindowClose: false}
	body := gc.Init("gwui test")
	body.SetBackgroundImage("abstract.jpg", 95)

	tabs := gc.TabsNew([]string{"tab1", "tab2", "tab3"})

	//create first tab elements

	bt1 := gc.ButtonWithIconNew("primary", "brush", "Start")
	lb1 := gc.LabelNew("Start TextArea write loop")

	bt2 := gc.ButtonWithIconNew("secondary", "brush-fill", "Change")
	lb2 := gc.LabelNew("Change text background color")

	bt3 := gc.ButtonWithIconNew("success", "box-arrow-up-right", "Change")
	lb3 := gc.LabelNew("Change text size")

	lb4 := gc.LabelNew("Input text with Modal")
	it4 := gc.InputTextNew("text")

	lb5 := gc.LabelNew("Select font family")

	dd5 := gc.DropDownNew("warning",
		"Font Family", []string{"arial", "verdana", "monospace"})

	md1 := gc.ModalNew("TEXT INPUT", "Are you sure", "yes", "no")
	md2 := gc.ModaPasswordlNew("Type anything", "Submit")

	ta1 := gc.TextAreaNew(12)
	//mandatory: callback on textarea to handle incoming messages
	ta1.Callback(func(string, int) {})

	ta1.SetBackgroundColor("#ffe6e6")
	ta1.SetColor("blue")
	ta1.SetFontSize("small")
	ta1.SetFontFamily("monospace")

	cd1 := gc.CardNew("Kitchen Sink", "Elements")
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

	md2.SubElems[0].Callback(func(strValue string, intValue int) {
		text := Sprintf("Password: %s\n", strValue)
		ta1.WriteTextArea(text)
	})

	it4.Callback(func(strValue string, intValue int) {
		md1.ModalShow()
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
		md2.ModalShow()
		go textLoop(&ta1)
		bt1.ChangeText("Started")
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

	ct1 := gc.ContainerNew()
	r1 := gc.RowNew()
	r11 := gc.RowNew()
	r12 := gc.RowNew()
	r13 := gc.RowNew()
	r14 := gc.RowNew()
	r15 := gc.RowNew()
	c1 := gc.ColNew()
	c2 := gc.ColNew()
	c11 := gc.ColNew()
	c12 := gc.ColNew()
	c21 := gc.ColNew()
	c22 := gc.ColNew()
	c31 := gc.ColNew()
	c32 := gc.ColNew()
	c41 := gc.ColNew()
	c42 := gc.ColNew()
	c51 := gc.ColNew()
	c52 := gc.ColNew()
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

	img1 := gc.ImageNew("abstract.jpg", 50, 50)
	img2 := gc.ImageNew("abstract.jpg", 100, 100)
	img3 := gc.ImageNew("abstract.jpg", 200, 200)
	img4 := gc.ImageNew("abstract.jpg", 300, 100)

	rs1 := gc.RangeSliderNew(50.0, 0.0, 100.0, 1.0)

	pb1 := gc.PillBadgeNew("danger", "50")

	rd1Choices := []string{"do nothing", "abstract.jpg", "sunset.jpg"}
	rd1 := gc.RadioNew(rd1Choices, 1)

	fi1 := gc.FileInputNew("Open")
	lb11 := gc.LabelNew("No file")

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

	//Use Grid function instead of Col and Row directly

	par1 := gc.ParagraphNew()
	par1.Add(img1)
	par1.Add(img2)
	par1.Add(img3)
	par1.Add(img4)
	colst21 := []*gwui.Elem{&rd1, &par1}
	spanst21 := []int{2, 0}
	colst22 := []*gwui.Elem{&pb1, &rs1}
	spanst22 := []int{2, 0}
	colst23 := []*gwui.Elem{&fi1, &lb11}
	spanst23 := []int{0, 2}

	rowst2 := []gwui.ColSpans{{Elems: colst21, Spans: spanst21},
		{Elems: colst22, Spans: spanst22},
		{Elems: colst23, Spans: spanst23}}
	ct2 := gc.GridNew(rowst2)

	tabs.SubElems[1].Add(ct2)
	tabs.SubElems[1].SetBackgroundColor("white")

	//Third tab

	x2 := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	xs2 := []string{"aa", "bb", "cc", "dd", "ee"}
	ys2 := xs2
	y2Original := []float64{1.0, 2.0, 4.0, 8.0, 16.0}
	y2 := []float64{1.0, 2.0, 4.0, 8.0, 16.0}
	ym2 := [][]float64{{1.0, 2.0, 4.0, 8.0, 16.0},
		{2.0, 4.0, 8.0, 16.0, 18.0},
		{5.0, 8.0, 16.0, 18.0, 19.0},
	}

	pl1 := gc.PlotLine(x2, y2, "Line", "xaxis", "yaxis", 250, 250)
	pl2 := gc.PlotScatter(x2, y2, "Scatter", "xaxis", "yaxis", 250, 250)
	pl3 := gc.PlotVBar(xs2, y2, "VBar", "xaxis", "yaxis", 250, 250)
	pl4 := gc.PlotHBar(y2Original, ys2, "HBar", "xaxis", "yaxis", 250, 250)
	pl5 := gc.PlotLineMultiY(x2, ym2, "Multi",
		[]string{"traceA", "tB", "tC"}, "xaxis", "yaxis", 250, 250)
	pl6 := gc.PlotBoxplot(ym2, "Boxplot",
		[]string{"setA", "sB", "sC"}, "yaxis", 250, 250)

	rs2 := gc.RangeSliderNew(1.0, 1.0, 100.0, 1.0)

	upper := true
	rs2.Callback(func(strValue string, intValue int) {
		for ind, yElem := range y2Original {
			y2[ind] = yElem + float64(intValue)
			gc.PlotRedrawY(pl1, y2)
			gc.PlotRedrawY(pl2, y2)
			if upper {
				for ind, elem := range xs2 {
					xs2[ind] = strings.ToUpper(elem)
				}
				upper = false
			} else {
				for ind, elem := range xs2 {
					xs2[ind] = strings.ToLower(elem)
				}
				upper = true
			}
			gc.PlotRedrawXsY(pl3, xs2, y2)
		}
	})

	ct3 := gc.ContainerNew()
	r1t3 := gc.RowNew()
	c1t3 := gc.ColNew()
	c2t3 := gc.ColNew()
	c3t3 := gc.ColNew()
	c1t3.Add(pl1)
	c1t3.Add(pl2)
	c1t3.Add(rs2)
	c2t3.Add(pl3)
	c2t3.Add(pl4)
	c3t3.Add(pl5)
	c3t3.Add(pl6)
	r1t3.Add(c1t3)
	r1t3.Add(c2t3)
	r1t3.Add(c3t3)
	ct3.Add(r1t3)

	tabs.SubElems[2].Add(ct3)
	tabs.SubElems[2].SetBackgroundColor("white")

	//final body additions; add modal to body directly
	body.Add(tabs)
	body.Add(md1)
	body.Add(md2)

	gc.Close(body)
	gc.Run()

	gc.WaitKeyFromCOnsole()

}

func textLoop(textArea *gwui.Elem) {
	ind := 0
	for {
		ind++
		text := Sprintf("%d: All work and no play ...\n", ind)
		textArea.WriteTextArea(text)
		timeD := time.Duration(3000) * time.Millisecond
		time.Sleep(timeD)
	}
}
