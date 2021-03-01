package main

import (
	"time"

	"github.com/lxn/walk"
	declarative "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

const (
	mWidth  = 128
	mHeight = 128
)

func main() {
	// core.Init()
	// gui.Start()

	var mascot *walk.MainWindow

	declarative.MainWindow{
		Title:    "Mascot",
		AssignTo: &mascot,
		Size:     declarative.Size{Width: mWidth, Height: mHeight},
		Visible:  false,
		Layout:   declarative.VBox{},
		Children: []declarative.Widget{},
	}.Create()

	defaultStyle := win.GetWindowLong(mascot.Handle(), win.GWL_STYLE)
	mStyle := defaultStyle &^ win.WS_THICKFRAME &^ win.WS_SYSMENU &^ win.WS_CAPTION
	win.SetWindowLong(mascot.Handle(), win.GWL_STYLE, mStyle)

	go func() {
		time.Sleep(time.Duration(1) * time.Second)
		win.ShowWindow(mascot.Handle(), win.SW_SHOW)

		time.Sleep(time.Duration(2) * time.Second)
		xScreen := win.GetSystemMetrics(win.SM_CXSCREEN)
		yScreen := win.GetSystemMetrics(win.SM_CYSCREEN)
		win.SetWindowPos(
			mascot.Handle(),
			0,
			(xScreen-mHeight)/2,
			(yScreen-mHeight)/2,
			mHeight,
			mHeight,
			win.SWP_FRAMECHANGED,
		)

		time.Sleep(time.Duration(2) * time.Second)
		mascot.Close()
	}()

	mascot.Run()
}
