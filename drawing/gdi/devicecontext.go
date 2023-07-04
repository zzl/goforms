package gdi

import (
	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/drawing"
	"github.com/zzl/goforms/framework/consts"
	"syscall"
)

type DeviceContext struct {
	Handle win32.HDC

	hWndOwner win32.HWND
	bitmap0   win32.HBITMAP
	font0     win32.HFONT
}

func NewDeviceContext(hdc win32.HDC) *DeviceContext {
	d := &DeviceContext{Handle: hdc, hWndOwner: consts.NegativeOne}
	d.Init()
	return d
}

func NewWindowDc(hWnd win32.HWND) *DeviceContext {
	hDc := win32.GetWindowDC(hWnd)
	return &DeviceContext{Handle: hDc, hWndOwner: hWnd}
}

func NewClientDc(hWnd win32.HWND) *DeviceContext {
	hDc := win32.GetDC(hWnd)
	return &DeviceContext{Handle: hDc, hWndOwner: hWnd}
}

func NewCompatibleDc(hWnd win32.HWND) *DeviceContext {
	hDcRef := win32.GetDC(hWnd)
	hDc := win32.CreateCompatibleDC(hDcRef)
	win32.ReleaseDC(hWnd, hDcRef)
	return &DeviceContext{Handle: hDc}
}

func NewCompatibleDcAndBitmap(hWnd win32.HWND) (*DeviceContext, win32.HBITMAP) {
	hDcRef := win32.GetDC(hWnd)
	hDc := win32.CreateCompatibleDC(hDcRef)
	dc := &DeviceContext{Handle: hDc}
	var rc win32.RECT
	win32.GetWindowRect(hWnd, &rc)
	hbm := win32.CreateCompatibleBitmap(hDcRef, rc.Right-rc.Left, rc.Bottom-rc.Top)
	dc.SetBitmap(hbm)
	win32.ReleaseDC(hWnd, hDcRef)
	return dc, hbm
}

func (this *DeviceContext) Init() {
	//this.font0, _ = win32.GetCurrentObject(this.Handle, win32.OBJ_FONT)
}

func (this *DeviceContext) Dispose() {
	if this.bitmap0 != 0 {
		win32.SelectObject(this.Handle, this.bitmap0)
	}
	if this.font0 != 0 {
		win32.SelectObject(this.Handle, this.font0)
	}
	if this.hWndOwner == consts.NegativeOne {
		//nop
	} else if this.hWndOwner != 0 {
		win32.ReleaseDC(this.hWndOwner, this.Handle)
	} else {
		win32.DeleteDC(this.Handle)
	}
}

func (this *DeviceContext) SetTextColor(color drawing.Color) {
	_ = win32.SetTextColor(this.Handle, color.Win32Color())
}

func (this *DeviceContext) SetFont(font *Font) {
	hOriFont := win32.SelectObject(this.Handle, font.Handle)
	if this.font0 == 0 {
		this.font0 = win32.HFONT(hOriFont)
	}
}

func (this *DeviceContext) UnsetFont() {
	if this.font0 != 0 {
		win32.SelectObject(this.Handle, this.font0)
		this.font0 = 0
	}
}

func (this *DeviceContext) SetBitmap(bitmap win32.HBITMAP) {
	hOriBitmap := win32.SelectObject(this.Handle, bitmap)
	if this.bitmap0 == 0 {
		this.bitmap0 = hOriBitmap
	}
}

func (this *DeviceContext) UnsetBitmap() {
	if this.bitmap0 != 0 {
		win32.SelectObject(this.Handle, this.bitmap0)
		this.bitmap0 = 0
	}
}

func (this *DeviceContext) DrawText(text string, rect drawing.Rect,
	format win32.DRAW_TEXT_FORMAT) {
	wsz, _ := syscall.UTF16FromString(text)
	win32Rect := rect.Win32Rect()
	win32.DrawText(this.Handle, &wsz[0], int32(len(wsz)-1),
		&win32Rect, format)
}

func (this *DeviceContext) MeasureText(text string, size drawing.Size,
	format win32.DRAW_TEXT_FORMAT) drawing.Size {
	wsz, _ := syscall.UTF16FromString(text)
	rect := win32.RECT{
		0, 0, size.Width, size.Height,
	}
	format |= win32.DT_CALCRECT
	win32.DrawText(this.Handle, &wsz[0], int32(len(wsz)-1), &rect, format)
	return drawing.Size{rect.Right, rect.Bottom}
}

func (this *DeviceContext) SetBkMode(mode win32.BACKGROUND_MODE) {
	win32.SetBkMode(this.Handle, mode)
}
