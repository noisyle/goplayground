package main

import (
	"time"

	"github.com/lxn/walk"
	declarative "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

const (
	SIZE_W = 128
	SIZE_H = 128
)

func main() {
	// core.Init()
	// gui.Start()

	var mascot *walk.MainWindow

	declarative.MainWindow{
		Title:    "Mascot",
		AssignTo: &mascot,
		Size:     declarative.Size{Width: SIZE_W, Height: SIZE_H},
		Visible:  false,
		Layout:   declarative.VBox{},
		Children: []declarative.Widget{},
	}.Create()

	defaultStyle := win.GetWindowLong(mascot.Handle(), win.GWL_STYLE)
	newStyle := defaultStyle &^ win.WS_THICKFRAME
	win.SetWindowLong(mascot.Handle(), win.GWL_STYLE, newStyle)

	go func() {
		time.Sleep(time.Duration(2) * time.Second)
		win.ShowWindow(mascot.Handle(), win.SW_SHOW)

		time.Sleep(time.Duration(2) * time.Second)
		xScreen := win.GetSystemMetrics(win.SM_CXSCREEN)
		yScreen := win.GetSystemMetrics(win.SM_CYSCREEN)
		win.SetWindowPos(
			mascot.Handle(),
			0,
			(xScreen-SIZE_W)/2,
			(yScreen-SIZE_H)/2,
			SIZE_W,
			SIZE_H,
			win.SWP_FRAMECHANGED,
		)

		time.Sleep(time.Duration(2) * time.Second)
		win.ShowWindow(mascot.Handle(), win.SW_HIDE)
	}()

	mascot.Run()

}
