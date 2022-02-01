package gwui

import (
	"bufio"
	. "fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
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
	ContainerT = 17
	FileInputT = 18
	PlotT      = 19
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
	callback  bool
	gc        *GuiCfg
}
type GuiCfg struct {
	fh                fp.File
	fjs               fp.File
	fcss              fp.File
	mutex             sync.Mutex
	idCnt             int
	Body              *Elem
	ServeURL          string
	BrowserStart      bool `default:"true"`
	PlotIncluded      bool `default:"false"`
	ExitOnWindowClose bool `default:"false"`
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
	if e.callback {
		Println("Callback already attached to elem", e.id, "type", e.elType)
		return
	}
	e.callback = true
	go func() {
		addr := Sprintf("/%s", e.id)
		if e.elType == ButtonT {
			http.HandleFunc(addr, func(w http.ResponseWriter, r *http.Request) {
				fn("", 0)
			})
		} else if e.elType == ITextT || e.elType == DDownT || e.elType == RSliderT || e.elType == RadioT || e.elType == FileInputT {
			http.HandleFunc(addr, func(w http.ResponseWriter, r *http.Request) {
				buf := make([]byte, BUFFER_SIZE)
				inp := r.Body
				inp.Read(buf)
				strValue := string(buf[:r.ContentLength])
				if e.elType == DDownT {
					//change button text with selected item
					e.ChangeText(strValue)

				}
				if e.elType == ITextT || e.elType == DDownT || e.elType == FileInputT {
					if e.elType == FileInputT {
						strValue = filepath.Base(strValue)
					}
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
				gc := e.gc
				gc.reader(e)
				fn("", 0)
			})
		}
	}()
}

func (gc *GuiCfg) reader(el *Elem) {
	if el.elType != BodyT || !gc.ExitOnWindowClose {
		return
	}
	var conn *websocket.Conn = el.gs
	for {
		// read in a message
		_, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				//text := Sprintf("reader %s, %d, %s", el.id, el.elType, err)
				//Println(text)
				return
			}
		}
		pStr := string(p)
		if pStr == "CLOSE" {
			Println("Browser window closed, exiting.")
			os.Exit(0)
		}
	}
}

// idNew generates a unique id
func (gc *GuiCfg) idNew() string {
	gc.idCnt++
	idStr := Sprintf("ID%d", gc.idCnt)
	return idStr
}

// Run starts the server and launches the browser.
func (gc *GuiCfg) Run() {
	const TRIES int = 5
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		Println(err)
		os.Exit(0)
	}

	go func() {
		err := http.Serve(listener, nil)
		if err != nil {
			Println("Error Serving:", err)
		}
	}()
	gc.ServeURL = Sprintf("http://%s", listener.Addr())
	text := Sprintf("Serving on %s", gc.ServeURL)
	Println(text)

	if gc.BrowserStart {
		go func() {
			ind := 0
			for {
				if gc.ServeURL == "" {
					timeD := time.Duration(500) * time.Millisecond
					time.Sleep(timeD)
					continue
				}
				if ind >= TRIES {
					break
				}
				errb := browser.OpenURL(gc.ServeURL)
				if errb != nil {
					Println("Wrong URL")
				} else {
					text := Sprintf("Launching browser with %s ...", gc.ServeURL)
					Println(text)
					Println(" You may need to disable adblocker on localhost")
					return
				}
				ind++
			}
			Println("Something went wrong with the browser")
		}()
	}

	//wait until body callback completed
	count := 0
	for {
		if gc.Body.gs != nil {
			return
		}
		timeD := time.Duration(100) * time.Millisecond
		time.Sleep(timeD)
		count++
		if count > 200 {
			Println("20 seconds passed and no body callback")
			return
		}
	}
}

// Close writes and closes web2.html and web2.js, where the
// UI is implemented.
func (gc *GuiCfg) Close(body Elem) {

	gc.fh.Write([]byte(body.html))
	gc.fh.Write([]byte(body.hEnd))
	gc.fjs.Write([]byte(body.js))

	gc.fh.Close()
	gc.fjs.Close()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/", http.FileServer(http.Dir("./static")))
}

func (gc *GuiCfg) WaitKeyFromCOnsole() {
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

// Init creates web2.html and web2.js files and returns the Body element.
func (gc *GuiCfg) Init(title string) Elem {

	gc.idCnt = 0 //unique id counter
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
	`, title)

	plotScript := ""
	if gc.PlotIncluded {
		plotScript = `<script type="text/javascript" src="/static/plotly/plotly-2.8.3.min.js"></script> `
	}
	hEnd := Sprintf(`
	</body>

	%s
	<script type="text/javascript" src="/static/web2.js"></script>
	</html>`, plotScript)

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
			item.src = src;
		}
		else if (type === "PREDRAWY") {
			var len = parseInt(messages[2], 10);
			var yVec = new Float32Array(len);
			for (ind = 0; ind < len; ind++) {
				yVec[ind] = parseFloat(messages[3+ind]);
			}
			window["data" + id][0].y = yVec;
			Plotly.newPlot(window["PLOT"+id], window["data"+id], window["layout"+id]);
		}
		else if (type === "PREDRAWXSY") {
			var len = parseInt(messages[2], 10);
			var yVec = new Float32Array(len);
			var xVec = new Array(len);;
			for (ind = 0; ind < len; ind++) {
				yVec[ind] = parseFloat(messages[3+ind]);
			}
			for (ind = 0; ind < len; ind++) {
				xVec[ind] = messages[len+3+ind];
			}
			window["data" + id][0].y = yVec;
			window["data" + id][0].x = xVec;
			Plotly.newPlot(window["PLOT"+id], window["data"+id], window["layout"+id]);
		}
		else if (type === "CLICK") {
			var fun = window[id+"_func"];
			fun();
			//item.click();
		}
	};
	window.onbeforeunload = function(e) {
		conn_%s.send("CLOSE");
	};

	`, addr, e.id, e.id, e.id)

	//attach a callback to body to handle logic -> gui messages
	e.Callback(func(string, int) {})
	gc.Body = &e

	return e
}

func (el *Elem) ModalShow() {
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

// Click forces a click.
func (el *Elem) Click() {
	gc := el.gc

	if gc.Body.gs != nil {
		toSend := Sprintf("CLICK@%s", el.id)
		gc.mutex.Lock()
		defer gc.mutex.Unlock()
		gc.Body.gs.WriteMessage(websocket.TextMessage, []byte(toSend))
	} else {
		Println("Failed Click, Set", gc.Body.id, "Callback!")
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
func (el *Elem) SetBackgroundImage(fileName string, opacityPerc int) {
	var js string
	if el.elType == BodyT {
		js = Sprintf(`
		document.body.style.backgroundImage = "url('static/%s')";	
		document.body.style.backgroundSize = "cover";	
		document.body.style.opacity = "%d%c";	
		`, fileName, opacityPerc, '%')
	} else {
		js = Sprintf(`
		var item = document.getElementById("%s");
		item.style.backgroundImage = "url('%s')";		
		item.style.backgroundSize = "cover";	
		item.style.opacity = "%d%c";		
		`, el.id, fileName, opacityPerc, '%')
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

// TabsNew creates a nav-tabs;  pass a vector of tab texts;
// contained tabs are returned as SubElems.
func (gc *GuiCfg) TabsNew(texts []string) Elem {
	var ids []string
	for i := 0; i < len(texts); i++ {
		ids = append(ids, gc.idNew())
	}
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

func (gc *GuiCfg) RadioNew(text []string, checkedInd int) Elem {
	var ids []string
	for i := 0; i < len(text); i++ {
		ids = append(ids, gc.idNew())
	}
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

// FileInputNew allows to select a file; pass the label text
func (gc *GuiCfg) FileInputNew(text string) Elem {
	id := gc.idNew()
	hStart := Sprintf(`
	<div class="input-group m-2">
  	<button class="btn btn-outline-secondary" type="button" id="%s" onclick="%s_func()">%s</button>
  	<input type="file" class="form-control" id="%sfile" aria-describedby="%s" aria-label="Upload">
    </div>`, id, id, text, id, id)

	e := Elem{gc: gc, hStart: hStart, hEnd: "", html: hStart, id: id, elType: FileInputT, js: ""}
	addr := Sprintf("/%s", e.id)
	e.js = Sprintf(`
	function %s_func(e) {
		xhr = new XMLHttpRequest();
		xhr.open("POST", "%s", true);
		var val = document.getElementById("%sfile").value;
		xhr.send(val);
	}
	`, e.id, addr, e.id)
	return e
}

// ParagraphNew creates a paragraph.
func (gc *GuiCfg) ParagraphNew() Elem {
	id := gc.idNew()
	hStart := Sprintf(`
	<p id="%s">`, id)
	hEnd := `
	</p>`
	e := Elem{gc: gc, hStart: hStart, hEnd: hEnd, html: hStart, id: id, elType: ParagraphT, js: ""}
	return e
}

// ParagraphNew creates an image tag; pass  the name of the image file,
// the size. The image file must reside in /static
func (gc *GuiCfg) ImageNew(fileName string, width int, height int) Elem {
	id := gc.idNew()
	hStart := Sprintf(`
	<img id="%s" src="static/%s" alt="missing img" width="%d" height="%d">
	`, id, fileName, width, height)
	e := Elem{gc: gc, hStart: hStart, hEnd: "", html: hStart, id: id, elType: ImageT, js: ""}
	return e
}

// ModalNew creates a Modal Dialog; pass the dialog
// title and general text, button texts.
func (gc *GuiCfg) ModalNew(title string, text string, bt1Text string, bt2Text string) Elem {
	id1 := gc.idNew()
	id2 := gc.idNew()
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

// CardNew creates a card; pass header and title text.
func (gc *GuiCfg) CardNew(header string, title string) Elem {
	id := gc.idNew()
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

// RangeSliderNew creates a slider, pass initial, min, max, step values.
func (gc *GuiCfg) RangeSliderNew(initial float32, min float32, max float32, step float32) Elem {
	id := gc.idNew()
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

// PillBadgeNew creates a pill badge, pass the
// type (same as Button), the text.
func (gc *GuiCfg) PillBadgeNew(bType string, text string) Elem {
	id := gc.idNew()
	hStart := Sprintf(`
	<span id="%s" class="badge rounded-pill bg-%s">%s</span>
	`, id, bType, text)

	e := Elem{gc: gc, hStart: hStart, hEnd: "", html: hStart, id: id, elType: PBadgeT, js: ""}
	return e
}

// ContainerNew creates a container.
func (gc *GuiCfg) ContainerNew() Elem {
	id := gc.idNew()
	hStart := Sprintf(`
	<div class="container" id="%s">`, id)
	hEnd := `
	</div>`
	e := Elem{gc: gc, hStart: hStart, hEnd: hEnd, html: hStart, id: id, elType: ContainerT, js: ""}
	return e
}

// Colspans is a struct useful to build grids
type ColSpans struct {
	Elems []*Elem
	Spans []int
}

// GridNew creates a container with a row col grid; content is input
// through an array of ColsSpans (one array element for every row).
// Use nil as *Elem for empty column.
// Use 0 for no span.
func (gc *GuiCfg) GridNew(colSpans []ColSpans) Elem {
	ct := gc.ContainerNew()
	for _, colSpan := range colSpans {
		row := gc.RowNew()
		for ind, elem := range colSpan.Elems {
			var col Elem
			span := colSpan.Spans[ind]
			if span != 0 {
				col = gc.ColSpanNew(span)
			} else {
				col = gc.ColNew()
			}
			if elem != nil {
				col.Add(*elem)
			}
			row.Add(col)
		}
		ct.Add(row)
	}
	return ct
}

// RowNew creates a row.
func (gc *GuiCfg) RowNew() Elem {
	id := gc.idNew()
	hStart := Sprintf(`
	<div class="row" id="%s">`, id)
	hEnd := `
	</div>`
	e := Elem{gc: gc, hStart: hStart, hEnd: hEnd, html: hStart, id: id, elType: RowT, js: ""}
	return e
}

// ColSpanNew creates a col with fixed width, pass the span (1, 2, 4, 6, 12).
func (gc *GuiCfg) ColSpanNew(span int) Elem {
	id := gc.idNew()
	hStart := Sprintf(`
	<div class="col-%d align-self-center" id="%s">`, span, id)
	hEnd := `
	</div>`
	e := Elem{gc: gc, hStart: hStart, hEnd: hEnd, html: hStart, id: id, elType: ColT, js: ""}
	return e
}

// ColNew creates a col.
func (gc *GuiCfg) ColNew() Elem {
	id := gc.idNew()
	hStart := Sprintf(`
	<div class="col align-self-center" id="%s">`, id)
	hEnd := `
	</div>`
	e := Elem{gc: gc, hStart: hStart, hEnd: hEnd, html: hStart, id: id, elType: ColT, js: ""}
	return e
}

// DropDownNew creates a button dropdown; pass the type (same as button),
// the button text and the list of options.
func (gc *GuiCfg) DropDownNew(bType string, text string, list []string) Elem {
	id := gc.idNew()
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

// ButtonNew creates a button, pass the B5 button type,
// the button text.
// B5 button types are: "primary", "secondary", "success", "danger",
// "warning", "info", "light", "dark".
func (gc *GuiCfg) ButtonNew(bType string, text string) Elem {
	id := gc.idNew()
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

// ButtonNew creates a button, pass the B5 type, the button text.
func (gc *GuiCfg) ButtonWithIconNew(bType string, iconName string, text string) Elem {
	id := gc.idNew()
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

// InputTextNew creates a input text field.
// Pass in input the placeholder string..
func (gc *GuiCfg) InputTextNew(placeholder string) Elem {
	id := gc.idNew()
	hStart := Sprintf(`
	<input type="text" class="m-2" id="%s" name="%s" placeholder="%s" style="text-align: right;" onkeypress="%s_func(event)">
	`, id, id, placeholder, id)
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

// LabelNew creates a label; pass the label text.
func (gc *GuiCfg) LabelNew(text string) Elem {
	id := gc.idNew()
	hText := Sprintf(`
	<label class="m-2" id=%s>%s</label>`, id, text)
	//gc.fh.Write([]byte(hText))
	e := Elem{gc: gc, hStart: hText, hEnd: "", html: hText, id: id, elType: LabelT, js: ""}
	return e
}

// TextAreaNew creates a textarea; the number of rows.
// Remember to attach a callback to handle output om the area.
func (gc *GuiCfg) TextAreaNew(rows int) Elem {
	id := gc.idNew()
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
		if ((messages.length == 1) && (messages[0] == '')){
			text.value = "";
		}
		text.scrollTop = text.scrollHeight;
	};
	`, e.id, addr, e.id, e.id)

	return e
}
