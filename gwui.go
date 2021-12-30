package gwui

import (
	"bufio"
	. "fmt"
	"net/http"
	"os"
	"sync"
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
	ModalT     = 12
	ImageT     = 13
	RSliderT   = 14
	PBadgeT    = 15
	RadioT     = 16
)

type Elem struct {
	elType    int
	hStart    string
	hEnd      string
	html      string
	js        string
	id        string
	gs        *websocket.Conn
	SubElems  []Elem
	subStart  string
	subEnd    string
	ChanBool1 chan bool
	gc        *GuiCfg
}
type GuiCfg struct {
	fh           fp.File
	fjs          fp.File
	fcss         fp.File
	mutex        sync.Mutex
	Body         *Elem
	Port         int
	BrowserStart bool
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

// Event handler for elements; it runs in its own thread.
// It takes a user function in input.
// It passes a string or a int value to the user function, depending on what
// is received from a POST
func (e *Elem) Callback(fn func(string, int)) {
	go func() {
		addr := Sprintf("/%s", e.id)
		if e.elType == ButtonT {
			http.HandleFunc(addr, func(w http.ResponseWriter, r *http.Request) {
				fn("", 0)
			})
		} else if e.elType == ITextT || e.elType == DDownT || e.elType == RSliderT || e.elType == RadioT {
			http.HandleFunc(addr, func(w http.ResponseWriter, r *http.Request) {
				buf := make([]byte, BUFFER_SIZE)
				inp := r.Body
				inp.Read(buf)
				strValue := string(buf[:r.ContentLength])
				if e.elType == DDownT {
					//change button text with selected item
					e.ChangeText(strValue)

				}
				if e.elType == ITextT || e.elType == DDownT {
					fn(strValue, 0)
				} else if e.elType == RSliderT || e.elType == RadioT {
					var intValue int
					Sscanf(strValue, "%d", &intValue)
					fn(strValue, intValue)
				}
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
				fn("", 0)
			})
		}
	}()
}

// GWRun starts the server and launches the browser.
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

	if gc.BrowserStart {
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
}

// GWClose writes and closes web2.html and web2.js, where the
// UI is implemented.
func (gc *GuiCfg) GWClose(body Elem) {

	gc.fh.Write([]byte(body.html))
	gc.fh.Write([]byte(body.hEnd))
	gc.fjs.Write([]byte(body.js))

	gc.fh.Close()
	gc.fjs.Close()

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

// GWB5Init creates web2.html and wen2.js files and returns the Body element.
func (gc *GuiCfg) GWB5Init(title string) Elem {

	if _, err := os.Stat("./static"); os.IsNotExist(err) {
		Println("Folder ./static missing.")
		Println("Make ./static folder and copy bootstrap inside.")
		os.Exit(0)
	}
	gc.fh = fp.Create("static/index.html")
	gc.fjs = fp.Create("static/web2.js")
	_, errCss := os.Stat("static/web2.css")
	if errCss != nil {
		if os.IsNotExist(errCss) {
			gc.fcss = fp.Create("static/web2.css")
			gc.fcss.Close()
		}
	}

	hStart := Sprintf(`
	<!DOCTYPE html>
	<html lang="en">

	<head>
    <title>%s</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">

    <link href="/static/bootstrap/css/bootstrap.css" rel="stylesheet" media="screen">
    <link href="/static/web2.css" rel="stylesheet">
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
	e := Elem{gc: gc, hStart: hStart, hEnd: hEnd, html: hStart, id: "body", elType: BodyT}

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
		else if (type === "COLOR") {
			var color = messages[2];
			item.style.color = color;
		}
		else if (type === "BCOLOR") {
			var color = messages[2];
			item.style.backgroundColor  = color;
		}
		else if (type === "FONTSIZE") {
			var fsize = messages[2];
			item.style.fontSize  = fsize;
		}
		else if (type === "FONTFAMILY") {
			var font = messages[2];
			item.style.fontFamily  = font;
		}
		else if (type === "MODALSHOW") {
			var modal = new bootstrap.Modal(item); 
			modal.show();
		}
		else if (type === "ENABLE") {
			var enable = messages[2];
			if (enable === "ENABLE") {
				item.disabled = false;
			}
			else  {
				item.disabled = true;
			}
		}
		else if (type === "IMAGE") {
			src = 'static/' + messages[2]; 
			item.src = src
		}
		
	};
	`, addr, e.id, e.id)

	return e
}

func (el *Elem) B5ModalShow() {
	gc := el.gc
	if el.elType == ModalT && gc.Body.gs != nil {
		toSend := Sprintf("MODALSHOW@%s@%s", el.id, "dummy")
		gc.mutex.Lock()
		defer gc.mutex.Unlock()
		gc.Body.gs.WriteMessage(websocket.TextMessage, []byte(toSend))
	}
	if gc.Body.gs == nil {
		Println("Failed Modal Show, Set", gc.Body.id, "Callback!")
	}
}

// ChangeToDisable changes on the run an element status to disable.
func (el *Elem) ChangeToDisable() {
	gc := el.gc
	if gc.Body.gs != nil {
		toSend := Sprintf("ENABLE@%s@%s", el.id, "DISABLE")
		gc.mutex.Lock()
		defer gc.mutex.Unlock()
		gc.Body.gs.WriteMessage(websocket.TextMessage, []byte(toSend))
	} else {
		Println("Failed Disable, Set", gc.Body.id, "Callback!")
	}
}

// ChangeToEnable changes on the run an element status to enable.
func (el *Elem) ChangeToEnable() {
	gc := el.gc

	if gc.Body.gs != nil {
		toSend := Sprintf("ENABLE@%s@%s", el.id, "ENABLE")
		gc.mutex.Lock()
		defer gc.mutex.Unlock()
		gc.Body.gs.WriteMessage(websocket.TextMessage, []byte(toSend))
	} else {
		Println("Failed Disable, Set", gc.Body.id, "Callback!")
	}
}

func (el *Elem) WriteTextArea(text string) {
	gc := el.gc
	if el.elType == TextAreaT && el.gs != nil {
		gc.mutex.Lock()
		defer gc.mutex.Unlock()
		err := el.gs.WriteMessage(websocket.TextMessage, []byte(text))
		if err != nil {
			Println("Write Error", err)
		}
	}
	if el.gs == nil {
		Println("No WriteTextArea, Set", el.id, "Callback!")
	}
}

// ChangeImage changes the image element src with the fileName passed.
func (el *Elem) ChangeImage(fileName string) {
	gc := el.gc
	if gc.Body.gs != nil {
		var toSend string
		toSend = Sprintf("IMAGE@%s@%s", el.id, fileName)
		gc.mutex.Lock()
		defer gc.mutex.Unlock()
		gc.Body.gs.WriteMessage(websocket.TextMessage, []byte(toSend))
	} else {
		Println("Failed Image change, Set", gc.Body.id, "Callback!")
	}
}

// ChangeText changes on the run an element text.
func (el *Elem) ChangeText(text string) {
	gc := el.gc
	if gc.Body.gs != nil {
		var toSend string
		if el.elType == ButtonT {
			toSend = Sprintf("TEXT@%stext@%s", el.id, text)
		} else {
			toSend = Sprintf("TEXT@%s@%s", el.id, text)
		}
		gc.mutex.Lock()
		defer gc.mutex.Unlock()
		gc.Body.gs.WriteMessage(websocket.TextMessage, []byte(toSend))
	} else {
		Println("Failed Text change, Set", gc.Body.id, "Callback!")
	}
}

// ChangeFontFamily changes on the run an element font family.
func (el *Elem) ChangeFontFamily(text string) {
	gc := el.gc
	if gc.Body.gs != nil {
		toSend := Sprintf("FONTFAMILY@%s@%s", el.id, text)
		gc.mutex.Lock()
		defer gc.mutex.Unlock()
		gc.Body.gs.WriteMessage(websocket.TextMessage, []byte(toSend))
	} else {
		Println("Failed Font change, Set", gc.Body.id, "Callback!")
	}
}

// ChangeColor changes on the run an element color.
func (el *Elem) ChangeColor(text string) {
	gc := el.gc
	if gc.Body.gs != nil {
		toSend := Sprintf("COLOR@%s@%s", el.id, text)
		gc.mutex.Lock()
		defer gc.mutex.Unlock()
		gc.Body.gs.WriteMessage(websocket.TextMessage, []byte(toSend))
	} else {
		Println("Failed Color change, Set", gc.Body.id, "Callback!")
	}
}

// SetBackgroundColor sets an element background color.
func (el *Elem) SetBackgroundColor(text string) {
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

// SetBackgroundImage sets an element background image.
// The image must reside in the static/ folder.
// Opacity is expressed as %; pass 100 for no transparency.
func (el *Elem) SetBackgroundImage(fileName string, opacity int) {
	var js string
	if el.elType == BodyT {
		js = Sprintf(`
		document.body.style.backgroundImage = "url('static/%s')";	
		document.body.style.backgroundSize = "cover";	
		document.body.style.opacity = "%d%c";	
		`, fileName, opacity, '%')
	} else {
		js = Sprintf(`
		var item = document.getElementById("%s");
		item.style.backgroundImage = "url('%s')";		
		item.style.backgroundSize = "cover";	
		item.style.opacity = "%d%c";		
		`, el.id, fileName, opacity, '%')
	}
	el.js = el.js + js
}

// SetToDisable sets an element status to disabled.
func (el *Elem) SetToDisable() {
	js := Sprintf(`
	var item = document.getElementById("%s");
	item.disabled = true;		
	`, el.id)
	el.js = el.js + js
}

// SetToEnable sets an element status to enabled.
func (el *Elem) SetToEnable() {
	js := Sprintf(`
	var item = document.getElementById("%s");
	item.disabled = false;		
	`, el.id)
	el.js = el.js + js
}

// SetColor sets an element foreground color.
func (el *Elem) SetColor(text string) {
	js := Sprintf(`
	var item = document.getElementById("%s");
	item.style.color = "%s";		
	`, el.id, text)
	el.js = el.js + js
}

// SetFontSize sets an element font size.
func (el *Elem) SetFontSize(text string) {
	js := Sprintf(`
	var item = document.getElementById("%s");
	item.style.fontSize = "%s";		
	`, el.id, text)
	el.js = el.js + js
}

// GWSetFontFamily sets an element font family.
func (el *Elem) SetFontFamily(text string) {
	js := Sprintf(`
	var item = document.getElementById("%s");
	item.style.fontFamily = "%s";		
	`, el.id, text)
	el.js = el.js + js
}

// ChangeBackgroundColor changes on the run an element background color.
func (el *Elem) ChangeBackgroundColor(text string) {
	gc := el.gc
	if gc.Body.gs != nil {
		toSend := Sprintf("BCOLOR@%s@%s", el.id, text)
		gc.mutex.Lock()
		defer gc.mutex.Unlock()
		gc.Body.gs.WriteMessage(websocket.TextMessage, []byte(toSend))
	} else {
		Println("Failed Background Color change, Set", gc.Body.id, "Callback!")
	}
}

// ChangeFontSize changes on the run an element font size.
func (el *Elem) ChangeFontSize(text string) {
	gc := el.gc
	if gc.Body.gs != nil {
		toSend := Sprintf("FONTSIZE@%s@%s", el.id, text)
		gc.mutex.Lock()
		defer gc.mutex.Unlock()
		gc.Body.gs.WriteMessage(websocket.TextMessage, []byte(toSend))
	} else {
		Println("Failed Font Size change, Set", gc.Body.id, "Callback!")
	}
}

// GWB5TabsNew creates a nav-tabs; pass a vector of unique ids;
// pass a vector of tab texts; contained tabs are returned as SubElems.
func (gc *GuiCfg) GWB5TabsNew(ids []string, texts []string) Elem {
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
	tabs := Elem{gc: gc, hStart: hText, hEnd: "", html: hText, id: "tabs", elType: TabsT, js: ""}

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
		e := Elem{gc: gc, hStart: hStart, hEnd: hEnd, html: hStart, id: id, elType: TPaneT, js: ""}
		elems = append(elems, e)
	}
	tabs.subStart = `<div class="tab-content">`
	tabs.subEnd = `</div>`
	tabs.SubElems = elems
	return tabs
}

func (gc *GuiCfg) GWB5RadioNew(ids []string, text []string, checkedInd int) Elem {
	hStart := ""
	checked := ""
	singleId := ids[0] + "radio"
	for ind, id := range ids {
		if ind == checkedInd {
			checked = "checked"
		} else {
			checked = ""
		}
		hStart = Sprintf(`
		%s
		<div class="form-check">
		<input class="form-check-input" type="radio" value="%d" name="%s" id="%s" %s onclick="%s_func(event)">
		<label class="form-check-label" for="%s">
		%s
		</label>
	  	</div>
		`, hStart, ind, singleId, id, checked, singleId, id, text[ind])
	}

	e := Elem{gc: gc, hStart: hStart, hEnd: "", html: hStart, id: singleId, elType: RadioT, js: ""}
	addr := Sprintf("/%s", singleId)
	e.js = Sprintf(`
	function %s_func(e) {
		xhr = new XMLHttpRequest();
		xhr.open("POST", "%s", true);
		var val = e.target.value;
		xhr.send(val);
	}`, singleId, addr)
	return e
}

// GWParagraphNew creates a paragraph; pass a unique identifier.
func (gc *GuiCfg) GWParagraphNew(id string) Elem {
	hStart := Sprintf(`
	<p id="%s">`, id)
	hEnd := `
	</p>`
	e := Elem{gc: gc, hStart: hStart, hEnd: hEnd, html: hStart, id: id, elType: ParagraphT, js: ""}
	return e
}

// GWParagraphNew creates an image tag; pass a unique identifier,
// the name of the image file, the size. The image file must reside in /static
func (gc *GuiCfg) GWImageNew(id string, fileName string, width int, height int) Elem {
	hStart := Sprintf(`
	<img id="%s" src="static/%s" alt="missing img" width="%d" height="%d">
	`, id, fileName, width, height)
	e := Elem{gc: gc, hStart: hStart, hEnd: "", html: hStart, id: id, elType: ImageT, js: ""}
	return e
}

// GWB5ModalNew creates a Modal Dialog; pass 2 unique identifiers for the buttons
// title and general text, button texts.
func (gc *GuiCfg) GWB5ModalNew(id1 string, id2 string,
	title string, text string, bt1Text string, bt2Text string) Elem {
	hStart := Sprintf(`
	<div class="modal" tabindex="-1" id="%s%s">
	<div class="modal-dialog">
	  <div class="modal-content">
		<div class="modal-header">
		  <h5 class="modal-title">%s</h5>
		</div>
		<div class="modal-body">
		  <p>%s</p>
		</div>
		<div class="modal-footer">
		  <button type="button" class="btn btn-primary" data-bs-dismiss="modal" onclick="%s_func()">%s</button>
		  <button type="button" class="btn btn-secondary" data-bs-dismiss="modal" onclick="%s_func()">%s</button>
		</div>
	  </div>
	</div>
	</div>`, id1, id2, title, text, id1, bt1Text, id2, bt2Text)

	ch := make(chan bool)
	e := Elem{gc: gc, hStart: hStart, hEnd: "", html: hStart, id: id1 + id2,
		elType: ModalT, js: "", ChanBool1: ch}
	e1 := Elem{gc: gc, id: id1, elType: ButtonT}
	e2 := Elem{gc: gc, id: id2, elType: ButtonT}
	e.SubElems = []Elem{e1, e2}
	addr1 := Sprintf("/%s", id1)
	addr2 := Sprintf("/%s", id2)
	e.js = Sprintf(`
	function %s_func() {
		xhr = new XMLHttpRequest();
		xhr.open("POST", "%s", true);
		xhr.send();
	}
	function %s_func() {
		xhr = new XMLHttpRequest();
		xhr.open("POST", "%s", true);
		xhr.send();
	}
	`, id1, addr1, id2, addr2)

	return e
}

// GWB5CardNew creates a card; pass a unique identifier, header and title text.
func (gc *GuiCfg) GWB5CardNew(id string, header string, title string) Elem {
	hStart := Sprintf(`
	<div class="card">
	<h5 class="card-header">%s</h5>
	<div class="card-body" id="%s">
	<h5 class="card-title">%s</h5>`, header, id, title)
	hEnd := `
	</div>
	</div>`
	e := Elem{gc: gc, hStart: hStart, hEnd: hEnd, html: hStart, id: id, elType: CardT, js: ""}
	return e
}

// GWB5RangeSliderNew creates a slider, pass a unique identifier, initial, min, max, step values.
func (gc *GuiCfg) GWB5RangeSliderNew(id string, initial float32, min float32, max float32, step float32) Elem {
	hStart := Sprintf(`
	<input id="%s" type="range" class="form-range" min="%f" max="%f" step="%f" onchange="%s_func()">
	`, id, min, max, step, id)
	e := Elem{gc: gc, hStart: hStart, hEnd: "", html: hStart, id: id, elType: RSliderT, js: ""}

	addr := Sprintf("/%s", e.id)
	e.js = Sprintf(`
	function %s_func(e) {
		xhr = new XMLHttpRequest();
		xhr.open("POST", "%s", true);
		var val = document.getElementById("%s").value; 
		xhr.send(val);
	}
	document.getElementById("%s").value = "%f"; 
	`, e.id, addr, e.id, e.id, initial)
	return e
}

// GWB5PillBadgeNew creates a pill badge, pass a unique identifier, type, text.
func (gc *GuiCfg) GWB5PillBadgeNew(id string, bType string, text string) Elem {
	hStart := Sprintf(`
	<span id="%s" class="badge rounded-pill bg-%s">%s</span>
	`, id, bType, text)

	e := Elem{gc: gc, hStart: hStart, hEnd: "", html: hStart, id: id, elType: PBadgeT, js: ""}
	return e
}

// GWB5RowNew creates a row, pass a unique identifier.
func (gc *GuiCfg) GWB5RowNew(id string) Elem {
	hStart := Sprintf(`
	<div class="row" id="%s">`, id)
	hEnd := `
	</div>`
	e := Elem{gc: gc, hStart: hStart, hEnd: hEnd, html: hStart, id: id, elType: RowT, js: ""}
	return e
}

// GWB5ColNew creates a col, pass a unique identifier.
func (gc *GuiCfg) GWB5ColNew(id string) Elem {
	hStart := Sprintf(`
	<div class="col" id="%s">`, id)
	hEnd := `
	</div>`
	e := Elem{gc: gc, hStart: hStart, hEnd: hEnd, html: hStart, id: id, elType: ColT, js: ""}
	return e
}

// GWB5DropDownNew creates a button dropdown; pass the B5 button type,
// a unique identifier, the button text and the list of options.
func (gc *GuiCfg) GWB5DropDownNew(id string, bType string, text string, list []string) Elem {
	hText := Sprintf(`
	<div class="dropdown">
  	<button class="btn-%s m-2 dropdown-toggle" type="button" id="%s" data-bs-toggle="dropdown" aria-expanded="false" >
    %s
  	</button>
  	<ul class="dropdown-menu" aria-labelledby="%s" id="%s%s" onclick="%s_func(event)">
	`, bType, id, text, id, id, id, id)

	for _, elem := range list {
		hText = Sprintf(`%s
		<li><a class="dropdown-item" href="#">%s</a></li>`, hText, elem)
	}
	hText = hText + `
  	</ul>
  	</div>`

	e := Elem{gc: gc, hStart: hText, hEnd: "", html: hText, id: id, elType: DDownT}

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

// GWB5ButtonNew creates a button, pass the B5 type, a unique identifier, the button text.
func (gc *GuiCfg) GWB5ButtonNew(id string, bType string, text string) Elem {
	hText := Sprintf(`
	<button type="button" class="btn btn-%s m-2" id="%s" onclick="%s_func()">
	<span id="%stext">%s</span></button>`, bType, id, id, id, text)
	//gc.fh.Write([]byte(hText))
	e := Elem{gc: gc, hStart: hText, hEnd: "", html: hText, id: id, elType: ButtonT}

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

// GWB5ButtonNew creates a button, pass the B5 type, a unique identifier, the button text.
func (gc *GuiCfg) GWB5ButtonWithIconNew(id string, bType string, iconName string, text string) Elem {
	hText := Sprintf(`
	<button type="button" class="btn btn-%s m-2" id="%s" onclick="%s_func()">
	<span class="btn-label">
	<img src="static/bootstrap-icons/%s.svg" alt="" width="16" height="16"></i></span>
	<span id="%stext">%s</span></button>`, bType, id, id, iconName, id, text)
	//gc.fh.Write([]byte(hText))
	e := Elem{gc: gc, hStart: hText, hEnd: "", html: hText, id: id, elType: ButtonT}

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

// GWB5InputTextNew creates a input text field; pass a unique identifier.
func (gc *GuiCfg) GWB5InputTextNew(id string) Elem {
	hStart := Sprintf(`
	<input type="text" class="m-2" id="%s" name="%s" onkeypress="%s_func(event)">
	`, id, id, id)
	e := Elem{gc: gc, hStart: hStart, hEnd: "", html: hStart, id: id, elType: ITextT}
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

// GWB5LabelNew creates a label; pass a unique identifier and the label text.
func (gc *GuiCfg) GWB5LabelNew(id string, text string) Elem {
	hText := Sprintf(`
	<label class="m-2" id=%s>%s</label>`, id, text)
	//gc.fh.Write([]byte(hText))
	e := Elem{gc: gc, hStart: hText, hEnd: "", html: hText, id: id, elType: LabelT, js: ""}
	return e
}

// GWB5TextAreaNew creates a textarea; pass a unique identifier and the number of rows.
// Remember to attach a callback to handle output om the area.
func (gc *GuiCfg) GWB5TextAreaNew(id string, rows int) Elem {
	hText := Sprintf(`
	<div class="form-group mx-2" style="min-width: 90%c">
	<p><textarea class="form-control" id=%s rows="%d"></textarea></p>
	</div>`, '%', id, rows)
	//gc.fh.Write([]byte(hText))
	e := Elem{gc: gc, hStart: hText, hEnd: "", html: hText, id: id, elType: TextAreaT}

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
