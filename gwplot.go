package gwui

import (
	. "fmt"

	"github.com/gorilla/websocket"
)

func (gc *GuiCfg) PlyDemoPlot() Elem {
	id := gc.idNew()
	hStart := Sprintf(`
	<div id="%s" style="width:600px;height:250px;"></div>`, id)
	e := Elem{gc: gc, hStart: hStart, hEnd: "", html: hStart, id: id, elType: ParagraphT, js: ""}
	e.js = Sprintf(`
	TESTER = document.getElementById('%s');
	var data = 
	Plotly.newPlot( TESTER, [{
	x: [1, 2, 3, 4, 5],
	y: [1, 2, 4, 8, 16] }], {
	margin: { t: 0 } } , {responsive: true});
	`, id)
	return e
}

func (gc *GuiCfg) plyPlotNumStr(id string, xVec []float64, yVec []string,
	pType string, mode string, title string, xTitle string, yTitle string,
	width int, height int) Elem {
	xStr := "["
	yStr := "["
	for ind, xElem := range xVec {
		xStr = Sprintf(`%s%7.2f, `, xStr, xElem)
		yStr = Sprintf(`%s"%s", `, yStr, yVec[ind])
	}
	xStr = xStr + "]"
	yStr = yStr + "]"
	hStart := Sprintf(`
	<div id="%s"></div>`, id)
	e := Elem{gc: gc, hStart: hStart, hEnd: "", html: hStart, id: id, elType: ParagraphT, js: ""}
	e.js = Sprintf(`
	PLOT%s = document.getElementById('%s');
	var data%s = [{
		x: %s,
		y: %s, 
		type: "%s",
		mode: "%s",
		orientation: "h"
	}];
	var layout%s = {
		title: "%s",
		xaxis: {title: "%s"},
		yaxis: {title: "%s"},
		width: %d,
		height: %d
	  };

	Plotly.newPlot( PLOT%s, data%s, layout%s);
	`, id, id, id, xStr, yStr, pType, mode, id, title, xTitle, yTitle, width, height, id, id, id)
	return e
}

func (gc *GuiCfg) plyPlotStrNum(id string, xVec []string, yVec []float64,
	pType string, mode string, title string, xTitle string, yTitle string,
	width int, height int) Elem {
	xStr := "["
	yStr := "["
	for ind, xElem := range xVec {
		xStr = Sprintf(`%s"%s", `, xStr, xElem)
		yStr = Sprintf(`%s%7.2f, `, yStr, yVec[ind])
	}
	xStr = xStr + "]"
	yStr = yStr + "]"
	hStart := Sprintf(`
	<div id="%s"></div>`, id)
	e := Elem{gc: gc, hStart: hStart, hEnd: "", html: hStart, id: id, elType: ParagraphT, js: ""}
	e.js = Sprintf(`
	PLOT%s = document.getElementById('%s');
	var data%s = [{
		x: %s,
		y: %s, 
		type: "%s",
		mode: "%s",
	}];
	var layout%s = {
		title: "%s",
		xaxis: {title: "%s"},
		yaxis: {title: "%s"},
		width: %d,
		height: %d
	  };

	Plotly.newPlot( PLOT%s, data%s, layout%s);
	`, id, id, id, xStr, yStr, pType, mode, id, title, xTitle, yTitle, width, height, id, id, id)
	return e
}

func (gc *GuiCfg) plyPlotNumMulti(id string, yVec [][]float64,
	traceNames []string, pType string, mode string, title string, yTitle string,
	width int, height int) Elem {

	lines := len(yVec)
	var data []string

	for line := 0; line < lines; line++ {
		yStr := "["
		for ind, _ := range yVec[line] {
			yStr = Sprintf("%s%7.2f, ", yStr, yVec[line][ind])
		}
		yStr = yStr + "]"
		datum := Sprintf(`{
			y: %s, 
			type: "%s",
			mode: "%s",
			name: "%s",
		},
		`, yStr, pType, mode, traceNames[line])
		data = append(data, datum)
	}

	hStart := Sprintf(`
	<div id="%s"></div>`, id)
	e := Elem{gc: gc, hStart: hStart, hEnd: "", html: hStart, id: id, elType: ParagraphT, js: ""}
	e.js = Sprintf(`
	PLOT%s = document.getElementById('%s');
	var data%s = 
		%s
	;

	var layout%s = {
		title: "%s",
		yaxis: {title: "%s"},
		width: %d,
		height: %d
	  };

	Plotly.newPlot( PLOT%s, data%s, layout%s);
	`, id, id, id, data, id, title, yTitle, width, height, id, id, id)
	return e
}

func (gc *GuiCfg) plyPlotNumNumMulti(id string, xVec []float64, yVec [][]float64,
	traceNames []string, pType string, mode string, title string, xTitle string, yTitle string,
	width int, height int) Elem {

	lines := len(yVec)
	var data []string

	for line := 0; line < lines; line++ {
		xStr := "["
		yStr := "["
		for ind, xElem := range xVec[line] {
			xStr = Sprintf("%s%7.2f, ", xStr, xElem)
			yStr = Sprintf("%s%7.2f, ", yStr, yVec[line][ind])
		}
		xStr = xStr + "]"
		yStr = yStr + "]"
		datum := Sprintf(`{
			x: %s,
			y: %s, 
			type: "%s",
			mode: "%s",
			name: "%s",
		},
		`, xStr, yStr, pType, mode, traceNames[line])
		data = append(data, datum)
	}

	hStart := Sprintf(`
	<div id="%s"></div>`, id)
	e := Elem{gc: gc, hStart: hStart, hEnd: "", html: hStart, id: id, elType: ParagraphT, js: ""}
	e.js = Sprintf(`
	PLOT%s = document.getElementById('%s');
	var data%s = 
		%s
	;

	var layout%s = {
		title: "%s",
		xaxis: {title: "%s"},
		yaxis: {title: "%s"},
		width: %d,
		height: %d
	  };

	Plotly.newPlot( PLOT%s, data%s, layout%s);
	`, id, id, id, data, id, title, xTitle, yTitle, width, height, id, id, id)
	return e
}

func (gc *GuiCfg) plyPlotNumNum(id string, xVec []float64, yVec []float64,
	pType string, mode string, title string, xTitle string, yTitle string,
	width int, height int) Elem {
	xStr := "["
	yStr := "["
	for ind, xElem := range xVec {
		xStr = Sprintf("%s%7.2f, ", xStr, xElem)
		yStr = Sprintf("%s%7.2f, ", yStr, yVec[ind])
	}
	xStr = xStr + "]"
	yStr = yStr + "]"
	hStart := Sprintf(`
	<div id="%s"></div>`, id)
	e := Elem{gc: gc, hStart: hStart, hEnd: "", html: hStart, id: id, elType: ParagraphT, js: ""}
	e.js = Sprintf(`
	PLOT%s = document.getElementById('%s');
	var data%s = [{
		x: %s,
		y: %s, 
		type: "%s",
		mode: "%s",
	}];
	var layout%s = {
		title: "%s",
		xaxis: {title: "%s"},
		yaxis: {title: "%s"},
		width: %d,
		height: %d
	  };

	Plotly.newPlot( PLOT%s, data%s, layout%s);
	`, id, id, id, xStr, yStr, pType, mode, id, title, xTitle, yTitle, width, height, id, id, id)
	return e
}

func (gc *GuiCfg) PlotVBar(x []string, y []float64, title string, xTitle string, yTitle string, width int, height int) Elem {
	id := gc.idNew()
	return gc.plyPlotStrNum(id, x, y, "bar", "s", title, xTitle, yTitle, width, height)
}

func (gc *GuiCfg) PlotHBar(x []float64, y []string, title string, xTitle string, yTitle string, width int, height int) Elem {
	id := gc.idNew()
	return gc.plyPlotNumStr(id, x, y, "bar", "s", title, xTitle, yTitle, width, height)
}

func (gc *GuiCfg) PlotLine(x []float64, y []float64, title string, xTitle string, yTitle string, width int, height int) Elem {
	id := gc.idNew()
	return gc.plyPlotNumNum(id, x, y, "", "lines", title, xTitle, yTitle, width, height)
}

func (gc *GuiCfg) PlotLineMultiY(x []float64, y [][]float64, title string, traceNames []string, xTitle string, yTitle string, width int, height int) Elem {
	id := gc.idNew()
	return gc.plyPlotNumNumMulti(id, x, y, traceNames, "", "lines", title, xTitle, yTitle, width, height)
}

func (gc *GuiCfg) PlotBoxplot(y [][]float64, title string, traceNames []string, yTitle string, width int, height int) Elem {
	id := gc.idNew()
	return gc.plyPlotNumMulti(id, y, traceNames, "box", "", title, yTitle, width, height)
}

func (gc *GuiCfg) PlotScatter(x []float64, y []float64, title string, xTitle string, yTitle string, width int, height int) Elem {
	id := gc.idNew()
	return gc.plyPlotNumNum(id, x, y, "scatter", "markers", title, xTitle, yTitle, width, height)
}

// PlotRedrawY updates a chart with new numeric y array
func (gc *GuiCfg) PlotRedrawY(el Elem, y []float64) {
	//Plotly.newPlot(element,charData,layout);
	if gc.Body.gs != nil {
		toSend := Sprintf("PREDRAWY@%s@%d", el.id, len(y))
		for _, yElem := range y {
			toSend = Sprintf("%s@%7.2f", toSend, yElem)
		}
		gc.mutex.Lock()
		defer gc.mutex.Unlock()
		gc.Body.gs.WriteMessage(websocket.TextMessage, []byte(toSend))
	} else {
		Println("Failed Redraw, Set", gc.Body.id, "Callback!")
	}
}

// PlotRedrawXsY updates a chart with new string x array
// and new numeric y array
func (gc *GuiCfg) PlotRedrawXsY(el Elem, x []string, y []float64) {
	//Plotly.newPlot(element,charData,layout);
	if gc.Body.gs != nil {
		toSend := Sprintf("PREDRAWXSY@%s@%d", el.id, len(y))
		for _, yElem := range y {
			toSend = Sprintf("%s@%7.2f", toSend, yElem)
		}
		for _, xElem := range x {
			toSend = Sprintf("%s@%s", toSend, xElem)
		}

		gc.mutex.Lock()
		defer gc.mutex.Unlock()
		gc.Body.gs.WriteMessage(websocket.TextMessage, []byte(toSend))
	} else {
		Println("Failed Redraw, Set", gc.Body.id, "Callback!")
	}
}
