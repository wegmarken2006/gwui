package gwui

import (
	. "fmt"

	"github.com/gorilla/websocket"
)

func (gc *GuiCfg) GWPlyDemoPlot(id string) Elem {
	hStart := Sprintf(`
	<div id="%s" style="width:600px;height:250px;"></div>`, id)
	e := Elem{gc: gc, hStart: hStart, hEnd: "", html: hStart, id: id, elType: ParagraphT, js: ""}
	e.js = Sprintf(`
	TESTER = document.getElementById('%s');
	var data = 
	Plotly.newPlot( TESTER, [{
	x: [1, 2, 3, 4, 5],
	y: [1, 2, 4, 8, 16] }], {
	margin: { t: 0 } } );
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

func (gc *GuiCfg) GWPlyPlotVBar(id string, x []string, y []float64, title string, xTitle string, yTitle string, width int, height int) Elem {
	return gc.plyPlotStrNum(id, x, y, "bar", "s", title, xTitle, yTitle, width, height)
}

func (gc *GuiCfg) GWPlyPlotHBar(id string, x []float64, y []string, title string, xTitle string, yTitle string, width int, height int) Elem {
	return gc.plyPlotNumStr(id, x, y, "bar", "s", title, xTitle, yTitle, width, height)
}

func (gc *GuiCfg) GWPlyPlotLine(id string, x []float64, y []float64, title string, xTitle string, yTitle string, width int, height int) Elem {
	return gc.plyPlotNumNum(id, x, y, "", "lines", title, xTitle, yTitle, width, height)
}

func (gc *GuiCfg) GWPlyPlotScatter(id string, x []float64, y []float64, title string, xTitle string, yTitle string, width int, height int) Elem {
	return gc.plyPlotNumNum(id, x, y, "scatter", "markers", title, xTitle, yTitle, width, height)
}

func (gc *GuiCfg) GWPlyPlotRedrawY(el Elem, y []float64) {
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
