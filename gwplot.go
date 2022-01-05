package gwui

import (
	. "fmt"
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

func (gc *GuiCfg) plyPlotIntStr(id string, xVec []int, yVec []string,
	pType string, mode string, title string, xTitle string, yTitle string,
	width int, height int) Elem {
	xStr := "["
	yStr := "["
	for ind, xElem := range xVec {
		xStr = Sprintf(`%s%d, `, xStr, xElem)
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

func (gc *GuiCfg) plyPlotStrInt(id string, xVec []string, yVec []int,
	pType string, mode string, title string, xTitle string, yTitle string,
	width int, height int) Elem {
	xStr := "["
	yStr := "["
	for ind, xElem := range xVec {
		xStr = Sprintf(`%s"%s", `, xStr, xElem)
		yStr = Sprintf(`%s%d, `, yStr, yVec[ind])
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

func (gc *GuiCfg) plyPlotIntInt(id string, xVec []int, yVec []int,
	pType string, mode string, title string, xTitle string, yTitle string,
	width int, height int) Elem {
	xStr := "["
	yStr := "["
	for ind, xElem := range xVec {
		xStr = Sprintf("%s%d, ", xStr, xElem)
		yStr = Sprintf("%s%d, ", yStr, yVec[ind])
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

func (gc *GuiCfg) GWPlyPlotVBar(id string, x []string, y []int, title string, xTitle string, yTitle string, width int, height int) Elem {
	return gc.plyPlotStrInt(id, x, y, "bar", "s", title, xTitle, yTitle, width, height)
}

func (gc *GuiCfg) GWPlyPlotHBar(id string, x []int, y []string, title string, xTitle string, yTitle string, width int, height int) Elem {
	return gc.plyPlotIntStr(id, x, y, "bar", "s", title, xTitle, yTitle, width, height)
}

func (gc *GuiCfg) GWPlyPlotLine(id string, x []int, y []int, title string, xTitle string, yTitle string, width int, height int) Elem {
	return gc.plyPlotIntInt(id, x, y, "", "lines", title, xTitle, yTitle, width, height)
}

func (gc *GuiCfg) GWPlyPlotScatter(id string, x []int, y []int, title string, xTitle string, yTitle string, width int, height int) Elem {
	return gc.plyPlotIntInt(id, x, y, "scatter", "markers", title, xTitle, yTitle, width, height)
}

func (gc *GuiCfg) GWPlyPlotRedraw(el Elem, x []int, y []int) {
	//Plotly.newPlot(element,charData,layout);
}
