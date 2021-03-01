package gui

import (
	"fmt"
	"time"

	"github.com/lxn/walk"
	declarative "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

const (
	mWidth  = 128
	mHeight = 128
)

func CreateMascot() {
	var mWindow *walk.MainWindow
	var mImageView *walk.ImageView

	declarative.MainWindow{
		AssignTo: &mWindow,
		Size:     declarative.Size{Width: mWidth, Height: mHeight},
		Visible:  false,
		Layout:   declarative.VBox{MarginsZero: true},
		Children: []declarative.Widget{
			declarative.ImageView{
				AssignTo: &mImageView,
				Image:    "mascot/hiiro/shime1.png",
				MinSize:  declarative.Size{Width: mWidth, Height: mHeight},
			},
		},
	}.Create()

	// 去除标题菜单和边框
	defaultStyle := win.GetWindowLong(mWindow.Handle(), win.GWL_STYLE)
	mStyle := defaultStyle &^ win.WS_THICKFRAME &^ win.WS_SYSMENU &^ win.WS_CAPTION
	win.SetWindowLong(mWindow.Handle(), win.GWL_STYLE, mStyle)

	// 移动到屏幕中央
	xScreen := win.GetSystemMetrics(win.SM_CXSCREEN)
	yScreen := win.GetSystemMetrics(win.SM_CYSCREEN)
	win.SetWindowPos(
		mWindow.Handle(),
		0,
		(xScreen-mWidth)/2,
		(yScreen-mHeight)/2,
		mWidth,
		mHeight,
		0,
	)

	go func() {
		time.Sleep(time.Duration(1) * time.Second)
		win.ShowWindow(mWindow.Handle(), win.SW_SHOW)

		var rect win.RECT
		win.GetWindowRect(mWindow.Handle(), &rect)
		for i := 1; i < 10; i++ {
			time.Sleep(time.Duration(200) * time.Millisecond)
			win.SetWindowPos(
				mWindow.Handle(),
				0,
				rect.Left+int32(i*10),
				rect.Top,
				mWidth,
				mHeight,
				0,
			)
			image, _ := walk.NewImageFromFileForDPI(fmt.Sprintf("mascot/hiiro/shime%d.png", i), 96)
			mImageView.SetImage(image)
		}

		time.Sleep(time.Duration(1) * time.Second)
		mWindow.Close()
	}()

	mWindow.Run()
}
