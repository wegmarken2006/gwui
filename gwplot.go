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

func (gc *GuiCfg) plyPlotIntInt(id string, xVec []int, yVec []int,
	pType string, mode string) Elem {
	xStr := "["
	yStr := "["
	for ind, xElem := range xVec {
		xStr = Sprintf("%s%d, ", xStr, xElem)
		yStr = Sprintf("%s%d, ", yStr, yVec[ind])
	}
	xStr = xStr + "]"
	yStr = yStr + "]"
	hStart := Sprintf(`
	<div id="%s" style="width:600px;height:250px;"></div>`, id)
	e := Elem{gc: gc, hStart: hStart, hEnd: "", html: hStart, id: id, elType: ParagraphT, js: ""}
	e.js = Sprintf(`
	PLOT%s = document.getElementById('%s');
	var data = [{
		x: %s,
		y: %s, 
		type: "%s",
		mode: "%s",
	}];

	Plotly.newPlot( PLOT%s, data);
	`, id, id, xStr, yStr, pType, mode, id)
	return e
}

func (gc *GuiCfg) GWPlyPlotLine(id string, x []int, y []int) Elem {
	return gc.plyPlotIntInt(id, x, y, "", "lines")
}

func (gc *GuiCfg) GWPlyPlotScatter(id string, x []int, y []int) Elem {
	return gc.plyPlotIntInt(id, x, y, "scatter", "markers")
}
