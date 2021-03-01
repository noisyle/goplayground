package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()

	if drv, ok := fyne.CurrentApp().Driver().(desktop.Driver); ok {
		w := drv.CreateSplashWindow()
		image := canvas.NewImageFromFile("mascot/hiiro/shime1.png")
		image.FillMode = canvas.ImageFillOriginal
		w.SetContent(image)
		w.Show()

		go func() {
			time.Sleep(time.Second * 3)
			w.Close()
		}()
	}

	w := a.NewWindow("Hello")

	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	))

	w.ShowAndRun()
}
