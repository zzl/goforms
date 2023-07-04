package forms

import (
	"log"
	"syscall"
	"unsafe"

	"github.com/zzl/go-win32api/v2/win32"
)

func MeasureText(hWnd HWND, text string) (int, int) {
	hdc := win32.GetDC(hWnd)
	hFont, _ := SendMessage(hWnd, win32.WM_GETFONT, 0, 0)
	hOriFont := win32.SelectObject(hdc, hFont)
	cx, cy := MeasureDcText(hdc, text)
	win32.SelectObject(hdc, hOriFont)
	win32.ReleaseDC(hWnd, hdc)
	return cx, cy
}

func MeasureDcText(hdc win32.HDC, text string) (int, int) {
	wsz, _ := syscall.UTF16FromString(text)
	var rc win32.RECT

	dtFlags := win32.DT_CALCRECT | win32.DT_LEFT | win32.DT_NOCLIP
	win32.DrawText(hdc, &wsz[0], int32(len(wsz)), &rc, dtFlags)

	return int(rc.Right), int(rc.Bottom)
}

func MeasureText2(hWnd HWND, text string, text2 string) (int, int) {
	hdc := win32.GetDC(hWnd)
	hFont, _ := SendMessage(hWnd, win32.WM_GETFONT, 0, 0)
	hOriFont := win32.SelectObject(hdc, hFont)
	cx, cy := MeasureDcText2(hdc, text, text2)
	win32.SelectObject(hdc, hOriFont)
	win32.ReleaseDC(hWnd, hdc)
	return cx, cy
}

func MeasureDcText2(hdc win32.HDC, text string, text2 string) (int, int) {
	wsz, _ := syscall.UTF16FromString(text)
	var rc win32.RECT

	dtFlags := win32.DT_CALCRECT | win32.DT_LEFT | win32.DT_NOCLIP
	win32.DrawText(hdc, &wsz[0], int32(len(wsz)), &rc, dtFlags)

	wsz2, _ := syscall.UTF16FromString(text2)
	var rc2 win32.RECT
	win32.DrawText(hdc, &wsz2[0], int32(len(wsz2)), &rc2, dtFlags)

	return int(max(rc.Right, rc2.Right)), int(max(rc.Bottom, rc2.Bottom))
}

type comboBoxMetrics struct {
	Height   int
	RcItem   win32.RECT
	RcButton win32.RECT
}

// todo: font aware?
var _comboBoxMetrics comboBoxMetrics

func GetComboBoxMetrics(hWndParent HWND) *comboBoxMetrics {
	if _comboBoxMetrics.Height != 0 {
		return &_comboBoxMetrics
	}
	hWnd, errno := win32.CreateWindowEx(
		0,
		win32.StrToPwstr("COMBOBOX"),
		nil,
		win32.WS_CHILD|WINDOW_STYLE(win32.CBS_DROPDOWN),
		0,
		0,
		64,
		32,
		hWndParent,
		0,
		HInstance,
		unsafe.Pointer(uintptr(0)))
	if hWnd == 0 {
		log.Fatal(errno)
	}

	var cbi win32.COMBOBOXINFO
	cbi.CbSize = uint32(unsafe.Sizeof(cbi))
	bOk, errno := win32.GetComboBoxInfo(hWnd, &cbi)
	if bOk == 0 {
		log.Fatal(errno)
	}
	var rc win32.RECT
	win32.GetWindowRect(hWnd, &rc)
	win32.DestroyWindow(hWnd)
	_comboBoxMetrics.RcItem = cbi.RcItem
	_comboBoxMetrics.RcButton = cbi.RcButton
	_comboBoxMetrics.Height = int(rc.Bottom-rc.Top) - 1

	return &_comboBoxMetrics
}
