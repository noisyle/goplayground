package gui

import (
	"image"
	"math"
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
	// var mImageView *walk.ImageView

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

	var rect win.RECT
	win.GetWindowRect(mWindow.Handle(), &rect)

	mImage, _ := walk.NewBitmapFromFile("mascot/hiiro/shime1.png")
	rgba, _ := mImage.ToImage()
	hbitmap, _ := hBitmapFromImage(image.Image(rgba), 96)

	hdcWnd := win.GetDC(mWindow.Handle())
	hdc := win.CreateCompatibleDC(hdcWnd)
	win.SelectObject(hdc, win.HGDIOBJ(hbitmap))
	defer win.DeleteDC(hdc)
	defer win.DeleteDC(hdcWnd)
	defer win.DeleteObject(win.HGDIOBJ(hbitmap))

	ptDst := win.POINT{X: rect.Left, Y: rect.Top}
	ptSrc := win.POINT{X: 0, Y: 0}
	sizeWnd := win.SIZE{CX: int32(mImage.Size().Width), CY: int32(mImage.Size().Height)}
	bf := win.BLENDFUNCTION{
		BlendOp:             0x1, // AC_SRC_OVER
		BlendFlags:          0,
		SourceConstantAlpha: 255,
		AlphaFormat:         win.AC_SRC_ALPHA,
	}

	updateLayeredWindow(mWindow.Handle(), hdcWnd, &ptDst, &sizeWnd, hdc, &ptSrc, 0, &bf, 0x2 /* ULW_ALPHA */)

	go func() {
		time.Sleep(time.Duration(1) * time.Second)
		win.ShowWindow(mWindow.Handle(), win.SW_SHOW)

		for i := 1; i < 10; i++ {
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
		}

		time.Sleep(time.Duration(1) * time.Second)
		mWindow.Close()
	}()

	mWindow.Run()
}

func hBitmapFromImage(im image.Image, dpi int) (win.HBITMAP, error) {
	var bi win.BITMAPV5HEADER
	bi.BiSize = uint32(unsafe.Sizeof(bi))
	bi.BiWidth = int32(im.Bounds().Dx())
	bi.BiHeight = -int32(im.Bounds().Dy())
	bi.BiPlanes = 1
	bi.BiBitCount = 32
	bi.BiCompression = win.BI_BITFIELDS
	const inchesPerMeter float64 = 39.37007874
	dpm := int32(math.Round(float64(dpi) * inchesPerMeter))
	bi.BiXPelsPerMeter = dpm
	bi.BiYPelsPerMeter = dpm
	// The following mask specification specifies a supported 32 BPP
	// alpha format for Windows XP.
	bi.BV4RedMask = 0x00FF0000
	bi.BV4GreenMask = 0x0000FF00
	bi.BV4BlueMask = 0x000000FF
	bi.BV4AlphaMask = 0xFF000000

	hdc := win.GetDC(0)
	defer win.ReleaseDC(0, hdc)

	var lpBits unsafe.Pointer

	// Create the DIB section with an alpha channel.
	hBitmap := win.CreateDIBSection(hdc, &bi.BITMAPINFOHEADER, win.DIB_RGB_COLORS, &lpBits, 0, 0)

	// Fill the image
	bitmapArray := (*[1 << 30]byte)(unsafe.Pointer(lpBits))
	i := 0
	for y := im.Bounds().Min.Y; y != im.Bounds().Max.Y; y++ {
		for x := im.Bounds().Min.X; x != im.Bounds().Max.X; x++ {
			r, g, b, a := im.At(x, y).RGBA()
			bitmapArray[i+3] = byte(a >> 8)
			bitmapArray[i+2] = byte(r >> 8)
			bitmapArray[i+1] = byte(g >> 8)
			bitmapArray[i+0] = byte(b >> 8)
			i += 4
		}
	}

	return hBitmap, nil
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
