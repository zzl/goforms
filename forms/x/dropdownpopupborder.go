package forms

import (
	"github.com/zzl/go-win32api/v2/win32"
	. "github.com/zzl/goforms/forms"
	"github.com/zzl/goforms/framework/virtual"
)

type DropdownPopupBorder interface {
	Window
}

type DropdownPopupBorderObject struct {
	WindowObject
	super *WindowObject

	Color   win32.COLORREF
	Visible bool
}

func NewDropdownPopupBorderObject() *DropdownPopupBorderObject {
	return virtual.New[DropdownPopupBorderObject]()
}

func (this *DropdownPopupBorderObject) Init() {
	this.super.Init()

	this.Color = win32.GetSysColor(win32.COLOR_HIGHLIGHT)
}

func (this *DropdownPopupBorderObject) WinProc(win *WindowObject, m *Message) error {
	if m.UMsg == win32.WM_ERASEBKGND {
		return m.SetHandledWithResult(1)
	}
	if m.UMsg == win32.WM_PAINT {
		var ps win32.PAINTSTRUCT
		win32.BeginPaint(win.Handle, &ps)
		if this.Visible {
			hdc := ps.Hdc
			//FillSolidRect(ps.Hdc, &ps.RcPaint, win32.RGB(255, 0, 0))
			var rc win32.RECT
			win32.GetClientRect(win.Handle, &rc)
			clr := win32.GetSysColor(win32.COLOR_HIGHLIGHT)
			//clr := win32.RGB(255, 0, 0)
			hPen := win32.CreatePen(win32.PS_SOLID, 1, clr)
			hOriPen := win32.SelectObject(hdc, win32.HGDIOBJ(hPen))
			win32.MoveToEx(hdc, 0, 0, nil)
			win32.LineTo(hdc, rc.Right-1, 0)
			win32.LineTo(hdc, rc.Right-1, rc.Bottom-1)
			win32.LineTo(hdc, 0, rc.Bottom-1)
			win32.LineTo(hdc, 0, 0)
			win32.SelectObject(hdc, hOriPen)
		}
		win32.EndPaint(win.Handle, &ps)
		return m.SetHandledWithResult(0)
	}
	if m.UMsg == win32.WM_NCHITTEST {
		return m.SetHandledWithResult(NegativeOne)
	}
	return win.CallOriWndProc(m)
}

func (this *DropdownPopupBorderObject) GetWindowClass() string {
	EnsurePlainWindowClassRegistered()
	return "goforms.plainwindow"
}

func (this *DropdownPopupBorderObject) GetDefaultStyle() WINDOW_STYLE {
	return win32.WS_POPUP //|win32.WS_VISIBLE
}

func (this *DropdownPopupBorderObject) GetDefaultExStyle() WINDOW_EX_STYLE {
	return win32.WS_EX_TOOLWINDOW | win32.WS_EX_TOPMOST
}

func (this *DropdownPopupBorderObject) SetBounds(left, top, width, height int) {
	cx := width
	cy := height
	win32.MoveWindow(this.Handle, int32(left), int32(top), int32(cx), int32(cy), win32.TRUE)
	hRgnA := win32.CreateRectRgn(0, 0, int32(cx), int32(cy))
	hRgnX := win32.CreateRectRgn(1, 1, int32(cx-1), int32(cy-1))
	hRgn := win32.CreateRectRgn(0, 0, 0, 0)
	win32.CombineRgn(hRgn, hRgnA, hRgnX, win32.RGN_DIFF)
	win32.SetWindowRgn(this.Handle, hRgn, 1)
	win32.DeleteObject(hRgnX)
	win32.DeleteObject(hRgnA)
}
