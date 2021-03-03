package gui

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"

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

	declarative.MainWindow{
		AssignTo: &mWindow,
		Size:     declarative.Size{Width: mWidth, Height: mHeight},
		Visible:  false,
		Layout:   declarative.VBox{MarginsZero: true},
	}.Create()

	// 去除标题菜单和边框
	defaultStyle := win.GetWindowLong(mWindow.Handle(), win.GWL_STYLE)
	mStyle := defaultStyle &^ win.WS_THICKFRAME &^ win.WS_SYSMENU &^ win.WS_CAPTION
	win.SetWindowLong(mWindow.Handle(), win.GWL_STYLE, mStyle)
	win.SetWindowLong(mWindow.Handle(), win.GWL_EXSTYLE, win.WS_EX_TOPMOST|win.WS_EX_LAYERED)

	// 移动到屏幕中央
	win.SetWindowPos(
		mWindow.Handle(),
		0,
		(win.GetSystemMetrics(win.SM_CXSCREEN)-mWidth)/2,
		(win.GetSystemMetrics(win.SM_CYSCREEN)-mHeight)/2,
		mWidth,
		mHeight,
		0,
	)

	updateWindow(mWindow.Handle(), "mascot/hiiro/shime1.png")
	win.ShowWindow(mWindow.Handle(), win.SW_SHOW)

	go func() {
		time.Sleep(time.Duration(1) * time.Second)

		for i := 0; i < 10; i++ {
			time.Sleep(time.Duration(200) * time.Millisecond)
			var rect win.RECT
			win.GetWindowRect(mWindow.Handle(), &rect)
			win.SetWindowPos(
				mWindow.Handle(),
				0,
				rect.Left+10,
				rect.Top,
				mWidth,
				mHeight,
				0,
			)
			updateWindow(mWindow.Handle(), fmt.Sprintf("mascot/hiiro/shime%d.png", i+2))
		}

		time.Sleep(time.Duration(1) * time.Second)
		mWindow.Close()
	}()

	mWindow.Run()
}

func updateWindow(hWnd win.HWND, filePath string) {
	hbitmap, _ := newHBITMAPFromFile(filePath, 96)

	hdcWnd := win.GetDC(hWnd)
	hdc := win.CreateCompatibleDC(hdcWnd)
	win.SelectObject(hdc, win.HGDIOBJ(hbitmap))
	defer win.DeleteDC(hdc)
	defer win.DeleteDC(hdcWnd)
	defer win.DeleteObject(win.HGDIOBJ(hbitmap))

	var rect win.RECT
	win.GetWindowRect(hWnd, &rect)
	ptDst := win.POINT{X: rect.Left, Y: rect.Top}
	ptSrc := win.POINT{X: 0, Y: 0}
	sizeWnd := win.SIZE{CX: int32(mWidth), CY: int32(mHeight)}
	bf := win.BLENDFUNCTION{
		BlendOp:             0x1, // AC_SRC_OVER
		BlendFlags:          0,
		SourceConstantAlpha: 255,
		AlphaFormat:         win.AC_SRC_ALPHA,
	}

	updateLayeredWindow(hWnd, hdcWnd, &ptDst, &sizeWnd, hdc, &ptSrc, 0, &bf, 0x2 /* ULW_ALPHA */)
}

type Error struct {
	message string
}

func (err *Error) Error() string {
	return err.message
}

func newHBITMAPFromFile(filePath string, dpi int) (win.HBITMAP, error) {
	var si win.GdiplusStartupInput
	si.GdiplusVersion = 1
	if status := win.GdiplusStartup(&si, nil); status != win.Ok {
		return 0, &Error{message: fmt.Sprintf("GdiplusStartup failed with status '%s'", status)}
	}
	defer win.GdiplusShutdown()

	var gpBmp *win.GpBitmap
	if status := win.GdipCreateBitmapFromFile(syscall.StringToUTF16Ptr(filePath), &gpBmp); status != win.Ok {
		return 0, &Error{message: fmt.Sprintf("GdipCreateBitmapFromFile failed with status '%s' for file '%s'", status, filePath)}
	}
	defer win.GdipDisposeImage((*win.GpImage)(gpBmp))

	var hBmp win.HBITMAP
	if status := win.GdipCreateHBITMAPFromBitmap(gpBmp, &hBmp, 0); status != win.Ok {
		return 0, &Error{message: fmt.Sprintf("GdipCreateHBITMAPFromBitmap failed with status '%s' for file '%s'", status, filePath)}
	}

	return hBmp, nil
}

func updateLayeredWindow(hWnd win.HWND, hdcDst win.HDC, pptDst *win.POINT, psize *win.SIZE, hdcSrc win.HDC, pptSrc *win.POINT, crKey win.COLORREF, pblend *win.BLENDFUNCTION, dwFlags uint32) bool {
	libuser32, err := syscall.LoadLibrary("user32.dll")
	defer syscall.FreeLibrary(libuser32)
	if err != nil {
		panic("加载 user32.dll UpdateLayeredWindow 失败")
	}
	updateLayeredWindow, err := syscall.GetProcAddress(libuser32, "UpdateLayeredWindow")
	ret, _, _ := syscall.Syscall9(updateLayeredWindow, 9,
		uintptr(hWnd),
		uintptr(hdcDst),
		uintptr(unsafe.Pointer(pptDst)),
		uintptr(unsafe.Pointer(psize)),
		uintptr(hdcSrc),
		uintptr(unsafe.Pointer(pptSrc)),
		uintptr(crKey),
		uintptr(unsafe.Pointer(pblend)),
		uintptr(dwFlags))
	if err != nil {
		panic("调用 user32.dll UpdateLayeredWindow 失败")
	}
	return ret != 0
}
