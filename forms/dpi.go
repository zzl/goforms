package forms

import "github.com/zzl/go-win32api/v2/win32"

var DpiAware bool
var Dpi int32 = 96

func SetDpiAware() {
	win32.SetProcessDPIAware()
	DpiAware = true
	hDc := win32.GetDC(0)
	Dpi = win32.GetDeviceCaps(hDc, win32.LOGPIXELSX)
	win32.ReleaseDC(0, hDc)
}

func DpiScale(value int) int {
	if Dpi != 96 {
		value = int(win32.MulDiv(int32(value), Dpi, 96))
	}
	return value
}

func DpiSize(width, height int) Size {
	return Size{DpiScale(width), DpiScale(height)}
}

func DpiPoint(x, y int) Point {
	return Point{DpiScale(x), DpiScale(y)}
}

func DpiUnscale(value int) int {
	if Dpi != 96 {
		value = int(win32.MulDiv(int32(value), 96, Dpi))
	}
	return value
}
