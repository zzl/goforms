package gdi

import (
	"github.com/zzl/go-win32api/v2/win32"
	"math"
	"syscall"
)

//var xDbu, yDbu int32

var dbusCache map[win32.HFONT]win32.SIZE

func MeasureDbus(hFont win32.HFONT) (int32, int32) {
	if dbusCache == nil {
		dbusCache = make(map[win32.HFONT]win32.SIZE)
	}
	if size, ok := dbusCache[hFont]; ok {
		return size.Cx, size.Cy
	}
	hdc := win32.CreateCompatibleDC(0)
	hOriFont := win32.SelectObject(hdc, win32.HGDIOBJ(hFont))
	var tm win32.TEXTMETRIC
	win32.GetTextMetrics(hdc, &tm)
	height := tm.TmHeight
	var width int32
	if tm.TmPitchAndFamily&win32.TMPF_FIXED_PITCH != 0 {
		s := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		wsz, _ := syscall.UTF16FromString(s)
		var size win32.SIZE
		win32.GetTextExtentPoint32(hdc, &wsz[0], int32(len(wsz)-1), &size)
		width = int32(math.Round(float64(size.Cx) / float64(len(s))))
	} else {
		width = tm.TmAveCharWidth
	}
	win32.SelectObject(hdc, hOriFont)
	win32.DeleteDC(hdc)
	dbusCache[hFont] = win32.SIZE{width, height}
	return width, height
}

//func DluToPx(hWnd win32.HWND, x, y int) (int, int) {
//	if xDbu == 0 {
//		xDbu, yDbu = MeasureDbus(hWnd)
//	}
//	x32, _ := win32.MulDiv(int32(x), xDbu, 4)
//	y32, _ := win32.MulDiv(int32(y), yDbu, 8)
//	return int(x32), int(y32)
//}
