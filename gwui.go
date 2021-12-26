package gwui

import (
	. "fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	fp "github.com/wegmarken2006/filepanic"
)

const WEBSOCKET_BUFFER_SIZE = 4096
const BUFFER_SIZE = 4096

var upgrader = websocket.Upgrader{
	ReadBufferSize:  WEBSOCKET_BUFFER_SIZE,
	WriteBufferSize: WEBSOCKET_BUFFER_SIZE,
}

const STARTING_PORT = 9000

const (
	ButtonT    = 0
	TextAreaT  = 1
	LabelT     = 2
	RowT       = 3
	ColT       = 4
	BodyT      = 5
	ITextT     = 6
	TabsT      = 7
	TPaneT     = 8
	ParagraphT = 9
)

type Elem struct {
	elType   int
	hStart   string
	hEnd     string
	html     string
	js       string
	id       string
	gs       *websocket.Conn
	SubElems []Elem
	subStart string
	subEnd   string
}
type GuiCfg struct {
	fh   fp.File
	fjs  fp.File
	fcss fp.File
	Body *Elem
}

func (e *Elem) WriteTextArea(text string) {
	if e.elType == TextAreaT && e.gs != nil {
		err := e.gs.WriteMessage(websocket.TextMessage, []byte(text))
		if err != nil {
			Println("Write Error", err)
		}
	}
	if e.gs == nil {
		Println("No WriteTextArea, Set", e.id, "Callback!")
	}
}

func (e *Elem) Add(n Elem) {
	e.html = e.html + n.html + n.hEnd
	e.js = e.js + n.js
	e.html = e.html + n.subStart
	for _, se := range n.SubElems {
		e.html = e.html + se.html + se.hEnd
		e.js = e.js + se.js
	}
	e.html = e.html + n.subEnd
}

func (e *Elem) Callback(fn func(string)) {
	addr := Sprintf("/%s", e.id)
	if e.elType == ButtonT {
		http.HandleFunc(addr, func(w http.ResponseWriter, r *http.Request) {
			fn("")
		})
	} else if e.elType == ITextT {
		http.HandleFunc(addr, func(w http.ResponseWriter, r *http.Request) {
			buf := make([]byte, BUFFER_SIZE)
			inp := r.Body
			inp.Read(buf)
			buf = buf[:r.ContentLength]
			fn(string(buf))
		})

	} else if e.elType == TextAreaT || e.elType == BodyT {
		http.HandleFunc(addr, func(w http.ResponseWriter, r *http.Request) {
			upgrader.CheckOrigin = func(r *http.Request) bool { return true }
			var err error
			e.gs, err = upgrader.Upgrade(w, r, nil)
			if err != nil {
				text := Sprintf("wsEndPoint %s, %s", e.id, err)
				Println(text)
			}
			fn("")
		})
	}
}

func (gc *GuiCfg) GWRun() {

	port := STARTING_PORT
	go func() {
		for {
			portStr := Sprintf(":%d", port)
			text := Sprintf("Serving on http://localhost%s", portStr)
			Println(text)
			err := http.ListenAndServe(portStr, nil)
			if err != nil {
				Println(err)
				port += 1
			} else {

				break
			}
		}
	}()
}

func (gc *GuiCfg) GWClose(body Elem) {

	gc.fh.Write([]byte(body.html))
	gc.fh.Write([]byte(body.hEnd))
	gc.fjs.Write([]byte(body.js))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/", http.FileServer(http.Dir("./static")))
	//Println("Close")

}

func (gc *GuiCfg) GWInit(title string) Elem {

	if _, err := os.Stat("./static"); os.IsNotExist(err) {
		Println("Folder ./static missing.")
		Println("Make ./static folder and copy bootstrap inside.")
		os.Exit(0)
	}
	gc.fh = fp.Create("static/index.html")
	gc.fjs = fp.Create("static/web2.js")
	gc.fcss = fp.Create("static/web2.css")

	hStart := Sprintf(`
	<!DOCTYPE html>
	<html lang="en">

	<head>
    <title>%s</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">

    <link href="/static/bootstrap/css/bootstrap.min.css" rel="stylesheet" media="screen">
    <link href="/static/web2.css" rel="stylesheet">
    <script type="text/javascript" src="https://code.jquery.com/jquery.js"></script>
    <script type="text/javascript" src="/static/bootstrap/js/bootstrap.min.js"></script>

	</head>

	<body>
	<div class="container">
	`, title)
	hEnd := `
	</div>
	</body>
	<script type="text/javascript" src="/static/web2.js"></script>
	
	</html>`
	e := Elem{hStart: hStart, hEnd: hEnd, html: hStart, id: "body", elType: BodyT}

	addr := Sprintf("/%s", e.id)
	e.js = Sprintf(`

	var addr1 = "ws://" + document.location.host + "%s";
	conn_%s = new WebSocket(addr1);
	conn_%s.onmessage = function (evt) {
		var messages = evt.data.split('@');
		var type = messages[0];
		var id = messages[1];
		var item = document.getElementById(id);
		if (type === "TEXT") {
			item.innerHTML = messages[2];
		}
		if (type === "COLOR") {
			var color = messages[2];
			item.style.color = color;
		}
		if (type === "BCOLOR") {
			var color = messages[2];
			item.style.backgroundColor  = color;
		}

		
	};
	`, addr, e.id, e.id)

	return e
	//gc.fh.Write([]byte(body))
}

func (gc *GuiCfg) GWChangeText(el Elem, text string) {
	if gc.Body.gs != nil {
		toSend := Sprintf("TEXT@%s@%s", el.id, text)
		gc.Body.gs.WriteMessage(websocket.TextMessage, []byte(toSend))
	} else {
		Println("No Change Text, Set", gc.Body.id, "Callback!")
	}
}

func (gc *GuiCfg) GWChangeColor(el Elem, text string) {
	if gc.Body.gs != nil {
		toSend := Sprintf("COLOR@%s@%s", el.id, text)
		gc.Body.gs.WriteMessage(websocket.TextMessage, []byte(toSend))
	} else {
		Println("No Change Color, Set", gc.Body.id, "Callback!")
	}
}

func (gc *GuiCfg) GWSetBackgroundColor(el Elem, text string) {
	js := Sprintf(`
	var item = document.getElementById(%s);
	item.style.color = %s;		
	`, el.id, text)
	el.js = el.js + js
}

func (gc *GuiCfg) GWChangeBackgroundColor(el Elem, text string) {
	if gc.Body.gs != nil {
		toSend := Sprintf("BCOLOR@%s@%s", el.id, text)
		gc.Body.gs.WriteMessage(websocket.TextMessage, []byte(toSend))
	} else {
		Println("No Change Back Color, Set", gc.Body.id, "Callback!")
	}
}

func (gc *GuiCfg) GWTabs(ids []string, texts []string) Elem {
	var elems []Elem
	hText := `
	<ul class="nav nav-tabs">
	`
	for ind, id := range ids {
		var linkType string
		if ind == 0 {
			linkType = "active"
		} else {
			linkType = ""
		}

		hText = Sprintf(`
		%s
		<li class="nav-item" role="presentation">
            <a class="nav-link %s" data-bs-toggle="tab" data-bs-target="#%s" id="t_%s">%s</a>
        </li>
		`, hText, linkType, id, id, texts[ind])

	}
	hText = Sprintf(`
	%s
	</ul>
	`, hText)
	tabs := Elem{hStart: hText, hEnd: "", html: hText, id: "tabs", elType: TabsT, js: ""}

	for ind, id := range ids {
		var paneType string
		if ind == 0 {
			paneType = "show active"
		} else {
			paneType = "fade"
		}
		hStart := Sprintf(`
		<div class="tab-pane %s" id="%s">
		`, paneType, id)
		hEnd := `
		</div>`
		e := Elem{hStart: hStart, hEnd: hEnd, html: hStart, id: id, elType: TPaneT, js: ""}
		elems = append(elems, e)
	}
	tabs.subStart = `<div class="tab-content">`
	tabs.subEnd = `</div>`
	tabs.SubElems = elems
	return tabs
}

func (gc *GuiCfg) GWParagraph(id string) Elem {
	hStart := Sprintf(`
	<p id="%s">`, id)
	hEnd := `
	</p>`
	e := Elem{hStart: hStart, hEnd: hEnd, html: hStart, id: id, elType: ParagraphT, js: ""}
	return e
}
func (gc *GuiCfg) GWRow(id string) Elem {
	hStart := Sprintf(`
	<div class="row" id="%s">`, id)
	hEnd := `
	</div>`
	e := Elem{hStart: hStart, hEnd: hEnd, html: hStart, id: id, elType: RowT, js: ""}
	return e
}

func (gc *GuiCfg) GWCol(id string) Elem {
	hStart := Sprintf(`
	<div class="col" id="%s">`, id)
	hEnd := `
	</div>`
	e := Elem{hStart: hStart, hEnd: hEnd, html: hStart, id: id, elType: ColT, js: ""}
	return e
}

func (gc *GuiCfg) GWButton(bType string, id string, text string) Elem {
	hText := Sprintf(`
	<button type="button" class="btn %s m-2" id="%s" onclick="%s_func()">%s</button>`, bType, id, id, text)
	//gc.fh.Write([]byte(hText))
	e := Elem{hStart: hText, hEnd: "", html: hText, id: id, elType: ButtonT}

	addr := Sprintf("/%s", e.id)
	e.js = Sprintf(`
	function %s_func() {
		xhr = new XMLHttpRequest();
		xhr.open("POST", "%s", true);
		xhr.send();
	}
	`, e.id, addr)

	return e
}

func (gc *GuiCfg) GWInputText(id string) Elem {
	hStart := Sprintf(`
	<input type="text" class="m-2" id="%s" name="%s" onkeypress="%s_func(event)">
	`, id, id, id)
	e := Elem{hStart: hStart, hEnd: "", html: hStart, id: id, elType: ITextT}
	addr := Sprintf("/%s", e.id)
	e.js = Sprintf(`
	function %s_func(e) {
		if(e.keyCode == 13) {
			xhr = new XMLHttpRequest();
			xhr.open("POST", "%s", true);
			var val = document.getElementById("%s").value;
			xhr.send(val);
		}
	}
	`, e.id, addr, id)

	return e
}

func (gc *GuiCfg) GWLabel(id string, text string) Elem {
	hText := Sprintf(`
	<label class="m-2" id=%s>%s</label>`, id, text)
	//gc.fh.Write([]byte(hText))
	e := Elem{hStart: hText, hEnd: "", html: hText, id: id, elType: LabelT, js: ""}
	return e
}

func (gc *GuiCfg) GWTextArea(id string, rows int) Elem {
	hText := Sprintf(`
	<div class="form-group mx-2" style="min-width: 90%c">
	<p><textarea class="form-control" id=%s rows="%d"></textarea></p>
	</div>`, '%', id, rows)
	//gc.fh.Write([]byte(hText))
	e := Elem{hStart: hText, hEnd: "", html: hText, id: id, elType: TextAreaT}

	addr := Sprintf("/%s", e.id)
	e.js = Sprintf(`
	var text = document.getElementById("%s");
	var addr1 = "ws://" + document.location.host + "%s";
	conn_%s = new WebSocket(addr1);
	conn_%s.onmessage = function (evt) {
		var edata = evt.data;
		var messages = edata.split('\n');
		for (var i = 0; i < messages.length; i++) {
			if (messages[i] != "") {
				var str = messages[i];
				str = text.value + str;
				diff = str.length - 4096;
				if (diff > 0) {
					text.value = str.slice(diff) + '\n';
				} else {
					text.value = str + '\n';
				}
			}
		}
		text.scrollTop = text.scrollHeight;
	};
	`, e.id, addr, e.id, e.id)

	return e
}
