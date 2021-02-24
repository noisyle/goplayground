package gui

import (
	"log"
	"shimeji/daemon"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

func Start() {
	mw, err := walk.NewMainWindow()
	if err != nil {
		log.Fatal(err)
	}

	icon, err := walk.Resources.Image("icon.png")
	if err != nil {
		log.Fatal(err)
	}

	ni, err := walk.NewNotifyIcon(mw)
	if err != nil {
		log.Fatal(err)
	}
	defer ni.Dispose()

	if err := ni.SetIcon(icon); err != nil {
		log.Fatal(err)
	}
	if err := ni.SetToolTip("左键查看信息，右键打开菜单。"); err != nil {
		log.Fatal(err)
	}

	mascot := declarative.MainWindow{
		AssignTo: &mw,
		Size:     declarative.Size{Width: 128, Height: 128},
		Layout:   declarative.VBox{},
		Children: []declarative.Widget{},
	}

	bm, err := walk.NewBitmapFromFile("mascot/hiiro/shime1.png")
	log.Print(bm.Size().Width)

	ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}

		daemon.SendCmd(&daemon.Cmd{Op: "left_click"})
		go mascot.Run()

		if err := ni.ShowCustom(
			"Walk NotifyIcon Example",
			"There are multiple ShowX methods sporting different icons.",
			icon); err != nil {

			log.Fatal(err)
		}
	})

	exitAction := walk.NewAction()
	if err := exitAction.SetText("E&xit"); err != nil {
		log.Fatal(err)
	}
	exitAction.Triggered().Attach(func() {
		walk.App().Exit(0)
	})
	if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
		log.Fatal(err)
	}

	if err := ni.SetVisible(true); err != nil {
		log.Fatal(err)
	}

	mw.Run()

}
