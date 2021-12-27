package gwui

import (
	"bufio"
	. "fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/browser"
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
	DDownT     = 10
	CardT      = 11
)

type Elem struct {
	elType   int
	hStart   string
	hEnd     string
	html     string
	js       string
	css      string
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
	Port int
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
	e.css = e.css + n.css

	e.html = e.html + n.subStart
	for _, se := range n.SubElems {
		e.html = e.html + se.html + se.hEnd
		e.js = e.js + se.js
		e.css = e.css + se.css
	}
	e.html = e.html + n.subEnd
}

func (e *Elem) Callback(fn func(string)) {
	addr := Sprintf("/%s", e.id)
	if e.elType == ButtonT {
		http.HandleFunc(addr, func(w http.ResponseWriter, r *http.Request) {
			fn("")
		})
	} else if e.elType == ITextT || e.elType == DDownT {
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
	const TRIES int = 5
	var serveURL string
	go func() {
		for ind := 0; ind < TRIES; ind++ {
			portStr := Sprintf(":%d", gc.Port)
			serveURL = Sprintf("http://localhost:%d", gc.Port)
			text := Sprintf("Serving on %s", serveURL)
			Println(text)
			err := http.ListenAndServe(portStr, nil)
			if err != nil {
				Println("Port busy, try another")
				gc.Port += 1
			} else {
				return
			}
		}
		Println("Something went wrong with the server")
	}()

	go func() {
		ind := 0
		for {
			if serveURL == "" {
				timeD := time.Duration(500) * time.Millisecond
				time.Sleep(timeD)
				continue
			}
			if ind >= TRIES {
				break
			}
			errb := browser.OpenURL(serveURL)
			if errb != nil {
				Println("Wrong URL")
			} else {
				text := Sprintf("Launching browser with %s ...", serveURL)
				Println(text)
				Println(" You may need to disable adblocker on localhost")
				return
			}
			ind++
		}
		Println("Something went wrong with the browser")
	}()

}

func (gc *GuiCfg) GWClose(body Elem) {

	gc.fh.Write([]byte(body.html))
	gc.fh.Write([]byte(body.hEnd))
	gc.fjs.Write([]byte(body.js))
	gc.fcss.Write([]byte(body.css))

	gc.fh.Close()
	gc.fjs.Close()
	gc.fcss.Close()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/", http.FileServer(http.Dir("./static")))
	//Println("Close")

}

func (gc *GuiCfg) GWWaitKeyFromCOnsole() {
	//Wait for a key press
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

func (gc *GuiCfg) GWB5Init(title string) Elem {

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

    <link href="/static/bootstrap/css/bootstrap.css" rel="stylesheet" media="screen">
    <link href="/static/web2.css" rel="stylesheet">
    <script type="text/javascript" src="https://code.jquery.com/jquery.js"></script>
    <script type="text/javascript" src="/static/bootstrap/js/bootstrap.bundle.js"></script>

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
		if (type === "FONTSIZE") {
			var fsize = messages[2];
			item.style.fontSize  = fsize;
		}
		if (type === "FONTFAMILY") {
			var font = messages[2];
			item.style.fontFamily  = font;
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
		Println("Failed Text change, Set", gc.Body.id, "Callback!")
	}
}

func (gc *GuiCfg) GWChangeFontFamily(el Elem, text string) {
	if gc.Body.gs != nil {
		toSend := Sprintf("FONTFAMILY@%s@%s", el.id, text)
		gc.Body.gs.WriteMessage(websocket.TextMessage, []byte(toSend))
	} else {
		Println("Failed Font change, Set", gc.Body.id, "Callback!")
	}
}

func (gc *GuiCfg) GWChangeColor(el Elem, text string) {
	if gc.Body.gs != nil {
		toSend := Sprintf("COLOR@%s@%s", el.id, text)
		gc.Body.gs.WriteMessage(websocket.TextMessage, []byte(toSend))
	} else {
		Println("Failed Color change, Set", gc.Body.id, "Callback!")
	}
}

func (gc *GuiCfg) GWSetBackgroundColor(el *Elem, text string) {
	var js string
	if el.elType == BodyT {
		js = Sprintf(`
		document.body.style.backgroundColor = "%s";		
		`, text)
	} else {
		js = Sprintf(`
		var item = document.getElementById("%s");
		item.style.backgroundColor = "%s";		
		`, el.id, text)
	}

	el.js = el.js + js
}

func (gc *GuiCfg) GWSetColor(el *Elem, text string) {
	js := Sprintf(`
	var item = document.getElementById("%s");
	item.style.color = "%s";		
	`, el.id, text)
	el.js = el.js + js
}

func (gc *GuiCfg) GWSetFontSize(el *Elem, text string) {
	js := Sprintf(`
	var item = document.getElementById("%s");
	item.style.fontSize = "%s";		
	`, el.id, text)
	el.js = el.js + js
}

func (gc *GuiCfg) GWSetFontFamily(el *Elem, text string) {
	js := Sprintf(`
	var item = document.getElementById("%s");
	item.style.fontFamily = "%s";		
	`, el.id, text)
	el.js = el.js + js
}

func (gc *GuiCfg) GWChangeBackgroundColor(el Elem, text string) {
	if gc.Body.gs != nil {
		toSend := Sprintf("BCOLOR@%s@%s", el.id, text)
		gc.Body.gs.WriteMessage(websocket.TextMessage, []byte(toSend))
	} else {
		Println("Failed Background Color change, Set", gc.Body.id, "Callback!")
	}
}
func (gc *GuiCfg) GWChangeFontSize(el Elem, text string) {
	if gc.Body.gs != nil {
		toSend := Sprintf("FONTSIZE@%s@%s", el.id, text)
		gc.Body.gs.WriteMessage(websocket.TextMessage, []byte(toSend))
	} else {
		Println("Failed Font Size change, Set", gc.Body.id, "Callback!")
	}
}

func (gc *GuiCfg) GWB5Tabs(ids []string, texts []string) Elem {
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

func (gc *GuiCfg) GWB5Card(id string, header string, title string) Elem {
	hStart := Sprintf(`
	<div class="card">
	<h5 class="card-header">%s</h5>
	<div class="card-body" id="%s">
	<h5 class="card-title">%s</h5>`, header, id, title)
	hEnd := `
	</div>
	</div>`
	e := Elem{hStart: hStart, hEnd: hEnd, html: hStart, id: id, elType: CardT, js: ""}
	return e
}

func (gc *GuiCfg) GWB5Row(id string) Elem {
	hStart := Sprintf(`
	<div class="row" id="%s">`, id)
	hEnd := `
	</div>`
	e := Elem{hStart: hStart, hEnd: hEnd, html: hStart, id: id, elType: RowT, js: ""}
	return e
}

func (gc *GuiCfg) GWB5Col(id string) Elem {
	hStart := Sprintf(`
	<div class="col" id="%s">`, id)
	hEnd := `
	</div>`
	e := Elem{hStart: hStart, hEnd: hEnd, html: hStart, id: id, elType: ColT, js: ""}
	return e
}

func (gc *GuiCfg) GWB5DropDown(bType string, id string, text string, list []string) Elem {
	hText := Sprintf(`
	<div class="dropdown">
  	<button class="btn %s m-2 dropdown-toggle" type="button" id="%s" data-bs-toggle="dropdown" aria-expanded="false" >
    %s
  	</button>
  	<ul class="dropdown-menu" aria-labelledby="%s" id="%s%s" onclick="%s_func(event)">`, bType, id, text, id, id, id, id)

	for _, elem := range list {
		hText = Sprintf(`%s
		<li><a class="dropdown-item" href="#">%s</a></li>`, hText, elem)
	}
	hText = hText + `
  	</ul>
  	</div>`

	e := Elem{hStart: hText, hEnd: "", html: hText, id: id, elType: DDownT}

	addr := Sprintf("/%s", e.id)
	e.js = Sprintf(`
	function %s_func(e) {
		xhr = new XMLHttpRequest();
		xhr.open("POST", "%s", true);
		var val = e.target.innerHTML;
		xhr.send(val);
	}
	`, e.id, addr)

	return e
}

func (gc *GuiCfg) GWB5Button(bType string, id string, text string) Elem {
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

func (gc *GuiCfg) GWB5InputText(id string) Elem {
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

func (gc *GuiCfg) GWB5Label(id string, text string) Elem {
	hText := Sprintf(`
	<label class="m-2" id=%s>%s</label>`, id, text)
	//gc.fh.Write([]byte(hText))
	e := Elem{hStart: hText, hEnd: "", html: hText, id: id, elType: LabelT, js: ""}
	return e
}

func (gc *GuiCfg) GWB5TextArea(id string, rows int) Elem {
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
