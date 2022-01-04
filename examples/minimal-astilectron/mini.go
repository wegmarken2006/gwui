package main

import (
	. "fmt"
	"log"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	gw "github.com/wegmarken2006/gwui"
)

func main() {
	gc := gw.GuiCfg{BrowserStart: false}
	body := gc.GWB5Init("mini")

	//mandatory: callback on body
	body.Callback(func(string, int) {})
	gc.Body = &body

	bt1 := gc.GWB5ButtonNew("bt1", "primary", "Count")
	lb1 := gc.GWB5LabelNew("lb1", "0")

	count := 0
	bt1.Callback(func(string, int) {
		count++
		text := Sprintf("%d", count)
		lb1.ChangeText(text)
	})

	body.Add(lb1)
	body.Add(bt1)

	gc.GWClose(body)
	gc.GWRun()

	// See https://github.com/asticode/go-astilectron/tree/master/example

	// Set logger
	l := log.New(log.Writer(), log.Prefix(), log.Flags())

	// Create astilectron

	a, err := astilectron.New(l, astilectron.Options{
		AppName:           "",
		BaseDirectoryPath: "static",
	})
	if err != nil {
		l.Fatal(Errorf("main: creating astilectron failed: %w", err))
	}
	defer a.Close()

	// Handle signals
	a.HandleSignals()

	// Start
	if err = a.Start(); err != nil {
		l.Fatal(Errorf("main: starting astilectron failed: %w", err))
	}

	Println("SERVE", gc, gc.ServeURL)
	// New window
	var w *astilectron.Window
	if w, err = a.NewWindow(gc.ServeURL, &astilectron.WindowOptions{
		Center: astikit.BoolPtr(true),
		Height: astikit.IntPtr(700),
		Width:  astikit.IntPtr(700),
	}); err != nil {
		l.Fatal(Errorf("main: new window failed: %w", err))
	}

	// Create windows
	if err = w.Create(); err != nil {
		l.Fatal(Errorf("main: creating window failed: %w", err))
	}

	// Blocking pattern
	a.Wait()

}
