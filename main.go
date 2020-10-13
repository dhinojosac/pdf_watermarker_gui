// Package main provides various examples of Fyne API capabilities
package main

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/widget"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	pdf "github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

var labelSetWatermark = widget.NewLabel("Watermark Text:")
var entrySetWatermark = widget.NewEntry()
var HBSetWatermark = widget.NewHBox(labelSetWatermark, entrySetWatermark)
var labelFileName = widget.NewLabel("No PDF file selected")

//var progressWM = widget.NewProgressBarInfinite()

func StampTime() string {
	loc, err := time.LoadLocation("America/Santiago")
	t := time.Now()
	if err == nil {
		t = t.In(loc)
	}
	//fmt.Println(t.Format("Mon, 02 Jan 2006 15:04:05 -0700"))
	return t.Format("02-Jan-2006 | 15:04:05")
}

func hasWatermarks(inFile string) bool {
	ok, err := api.HasWatermarksFile(inFile, nil)
	if err != nil {
		fmt.Printf("Checking for watermarks: %s: %v\n", inFile, err)
	}
	return ok
}

func stampWatermark(inFile, t string) {
	outFile := strings.Replace(inFile, ".pdf", "", -1) + "_wm.pdf"
	onTop := true // we are testing stamps
	msg := "StampWaterMark"

	// Check for existing stamps.
	if ok := hasWatermarks(inFile); ok {
		fmt.Printf("Watermarks found: %s\n", inFile)
	}
	// Stamp all pages.
	wm, err := pdf.ParseTextWatermarkDetails(t, "op:0.2, sc:1.0 rel, off: -30 -0", onTop)
	if err != nil {
		fmt.Printf("%s %s: %v\n", msg, outFile, err)
	}
	if err := api.AddWatermarksFile(inFile, outFile, nil, wm, nil); err != nil {
		fmt.Printf("%s %s: %v\n", msg, outFile, err)
	}

	timestamp := StampTime()
	// Stamp all pages.
	wm, err = pdf.ParseTextWatermarkDetails(timestamp, "op:0.2, sc:0.7 rel, off: 60 -30", onTop)
	if err != nil {
		fmt.Printf("%s %s: %v\n", msg, outFile, err)
	}
	if err = api.AddWatermarksFile(outFile, outFile, nil, wm, nil); err != nil {
		fmt.Printf("%s %s: %v\n", msg, outFile, err)
	}

	// Check for existing stamps.
	if ok := hasWatermarks(outFile); !ok {
		fmt.Printf("No watermarks found: %s\n", outFile)
	}

}

func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}

func watermarkPage(a fyne.App, win fyne.Window) fyne.CanvasObject {
	var inFile string

	entrySetWatermark.SetPlaceHolder("Enter text to add as watermark")

	return widget.NewVBox(
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Welcome to PDF WaterMarker", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewHBox(
			widget.NewButton("File Open PDF", func() {
				w := fyne.CurrentApp().NewWindow("Search your PDF file")
				//w.SetContent(widget.NewScrollContainer(img))
				w.Resize(fyne.NewSize(900, 900))
				w.Show()
				fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
					if err == nil && reader == nil {
						return
					}
					if err != nil {
						dialog.ShowError(err, win)
						return
					}

					inFile = strings.Replace(reader.URI().String(), "file://", "", -1)
					fmt.Println(inFile)
					labelFileName.SetText("PDF file name: " + reader.Name())

					w.Close()
				}, w)
				fd.SetFilter(storage.NewExtensionFileFilter([]string{".pdf"}))
				fd.Show()
			}),
			labelFileName,
		),
		layout.NewSpacer(),
		HBSetWatermark,
		layout.NewSpacer(),
		widget.NewButton("Stamp WaterMark", func() {
			fmt.Println("Stamp WaterMark!")
			stampWatermark(inFile, entrySetWatermark.Text)
		}),
	)

}

func main() {
	a := app.NewWithID("cl.dhinojosac.pdfwatermarker")
	//a.SetIcon(theme.FyneLogo())

	w := a.NewWindow("PDF WaterMarker")
	w.Resize(fyne.Size{600, 300})
	//w.SetFixedSize(true)

	/*
		newItem := fyne.NewMenuItem("New", nil)
		otherItem := fyne.NewMenuItem("Other", nil)
		otherItem.ChildMenu = fyne.NewMenu("",
			fyne.NewMenuItem("Project", func() { fmt.Println("Menu New->Other->Project") }),
			fyne.NewMenuItem("Mail", func() { fmt.Println("Menu New->Other->Mail") }),
		)
		newItem.ChildMenu = fyne.NewMenu("",
			fyne.NewMenuItem("File", func() { fmt.Println("Menu New->File") }),
			fyne.NewMenuItem("Directory", func() { fmt.Println("Menu New->Directory") }),
			otherItem,
		)
		settingsItem := fyne.NewMenuItem("Settings", func() { fmt.Println("Menu Settings") })

		cutItem := fyne.NewMenuItem("Cut", func() {
			shortcutFocused(&fyne.ShortcutCut{
				Clipboard: w.Clipboard(),
			}, w)
		})
		copyItem := fyne.NewMenuItem("Copy", func() {
			shortcutFocused(&fyne.ShortcutCopy{
				Clipboard: w.Clipboard(),
			}, w)
		})
		pasteItem := fyne.NewMenuItem("Paste", func() {
			shortcutFocused(&fyne.ShortcutPaste{
				Clipboard: w.Clipboard(),
			}, w)
		})
		findItem := fyne.NewMenuItem("Find", func() { fmt.Println("Menu Find") })

		helpMenu := fyne.NewMenu("Help", fyne.NewMenuItem("Help", func() { fmt.Println("Help Menu") }))
		mainMenu := fyne.NewMainMenu(
			// a quit item will be appended to our first menu
			fyne.NewMenu("File", newItem, fyne.NewMenuItemSeparator(), settingsItem),
			fyne.NewMenu("Edit", cutItem, copyItem, pasteItem, fyne.NewMenuItemSeparator(), findItem),
			helpMenu,
		)
		w.SetMainMenu(mainMenu)
	*/

	w.SetMaster()

	w.SetContent(watermarkPage(a, w))

	w.ShowAndRun()
	//a.Preferences().SetInt(preferenceCurrentTab, tabs.CurrentTabIndex())
}
