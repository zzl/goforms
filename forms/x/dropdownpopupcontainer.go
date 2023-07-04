package forms

import (
	"github.com/zzl/go-win32api/v2/win32"
	. "github.com/zzl/goforms/forms"
	"github.com/zzl/goforms/framework/virtual"
	"unsafe"
)

// todo: esc hide? or relay msg?
type DropdownPopupContainer interface {
	Window

	GetOnDeactivate() *SimpleEvent

	UpdateSize()

	SetDropdownRect(rect Rect)
}

type DropdownPopupContainerSpi interface {
	onActivate(wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT
}

type DropdownPopupContainerObject struct {
	WindowObject
	super *WindowObject

	Popup     DropdownPopup
	HasBorder bool
	NoAnim    bool

	//GetPopupBoundsCallback func() (int, int, int, int)
	DropdownRect       Rect
	DropdownRightAlign bool

	hWndDropdownControl    win32.HWND
	hWndTopWin             win32.HWND
	hWndPopupControl       win32.HWND
	hWndPopupControlParent win32.HWND

	border       *DropdownPopupBorderObject
	OnDeactivate SimpleEvent
}

func NewDropdownPopupContainerObject() *DropdownPopupContainerObject {
	return virtual.New[DropdownPopupContainerObject]()
}

func (this *DropdownPopupContainerObject) Init() {
	this.super.Init()
}

func (this *DropdownPopupContainerObject) Dispose() {
	this.super.Dispose()
}

func (this *DropdownPopupContainerObject) SetDropdownRect(rect Rect) {
	this.DropdownRect = rect
}

func dropPopupClassProc(hWnd win32.HWND, uMsg uint32,
	wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {

	if uMsg == win32.WM_NCCALCSIZE {
		if wParam != 0 {
			style, _ := win32.GetWindowLong(hWnd, win32.GWL_STYLE)
			if style&int32(win32.WS_BORDER) != 0 {
				params := (*win32.NCCALCSIZE_PARAMS)(unsafe.Pointer(lParam))
				rc0 := &params.Rgrc[0]
				rc0.Right += 1
				rc0.Bottom += 1
				return 0
			}
		}
		return win32.DefWindowProc(hWnd, uMsg, wParam, lParam)
	}

	if uMsg == win32.WM_PRINT {
		ret := win32.DefWindowProc(hWnd, uMsg, wParam, lParam)
		this := GetWindow(hWnd).(*DropdownPopupContainerObject)
		if this.HasBorder {
			hdc := win32.HDC(wParam)
			clr := win32.GetSysColor(win32.COLOR_3DDKSHADOW)
			var rc win32.RECT
			win32.GetWindowRect(hWnd, &rc)
			hPen := win32.CreatePen(win32.PS_SOLID, 1, clr)
			hOriPen := win32.SelectObject(hdc, win32.HGDIOBJ(hPen))
			win32.MoveToEx(hdc, 0, 0, nil)
			win32.LineTo(hdc, rc.Right-1, 0)
			win32.LineTo(hdc, rc.Right-1, rc.Bottom-1)
			win32.LineTo(hdc, 0, rc.Bottom-1)
			win32.LineTo(hdc, 0, 0)
			win32.SelectObject(hdc, hOriPen)
			//win32.ReleaseDC(hWnd, hdc)
		}
		return ret
	}

	if uMsg == win32.WM_ACTIVATE {
		spi := GetWindow(hWnd).(DropdownPopupContainerSpi)
		return spi.onActivate(wParam, lParam)
	}

	return win32.DefWindowProc(hWnd, uMsg, wParam, lParam)
}

func (this *DropdownPopupContainerObject) GetOnDeactivate() *SimpleEvent {
	return &this.OnDeactivate
}

func (this *DropdownPopupContainerObject) onActivate(wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
	//return 0 //xx
	state := win32.LOWORD(win32.DWORD(wParam))
	if state == uint16(win32.WA_INACTIVE) {
		//win32.ShowWindow(this.hWndPopupControl, win32.SW_HIDE)
		//win32.SetParent(this.hWndPopupControl, this.hWndPopupControlParent)
		//win32.ShowWindow(this.Handle, win32.SW_HIDE)
		//if this.HasBorder {
		//	this.border.Destroy()
		//}
		this.OnDeactivate.Fire(this, &SimpleEventInfo{})
		//HActiveWin = this.hWndTopWin
		//win32.SetFocus(this.hWndDropdownControl)
	} else {
		win32.SendMessage(this.hWndTopWin, win32.WM_NCACTIVATE, 1, 0)
		win32.DefWindowProc(this.Handle, win32.WM_ACTIVATE, wParam, lParam)
		win32.SendMessage(this.hWndTopWin, win32.WM_NCACTIVATE, 1, 0)
		HWndActive = this.hWndTopWin
	}
	return 0
}

var dropPopupClass string

func (this *DropdownPopupContainerObject) EnsureClassRegistered() {
	if dropPopupClass == "" {
		dropPopupClass = "goforms.dropdownpopup"
		_, _ = RegisterClass(dropPopupClass, dropPopupClassProc, ClassOptions{
			BackgroundBrush: win32.HBRUSH(win32.COLOR_WINDOW + 1),
			Style:           win32.CS_DROPSHADOW, //?|win32.CS_HREDRAW|win32.CS_VREDRAW,
		})
	}
}

func (this *DropdownPopupContainerObject) GetWindowClass() string {
	return dropPopupClass
}

func (this *DropdownPopupContainerObject) GetDefaultStyle() WINDOW_STYLE {
	style := win32.WS_POPUP | win32.WS_CLIPCHILDREN
	if this.HasBorder {
		style |= win32.WS_BORDER
	}
	return style
}

func (this *DropdownPopupContainerObject) GetDefaultExStyle() WINDOW_EX_STYLE {
	return win32.WS_EX_TOOLWINDOW | win32.WS_EX_TOPMOST
}

func (this *DropdownPopupContainerObject) CreateFor(hWndDropdown win32.HWND) error {
	this.hWndDropdownControl = hWndDropdown
	this.hWndTopWin = win32.GetAncestor(hWndDropdown, win32.GA_ROOT)

	this.hWndPopupControl = this.Popup.GetControl().GetHandle()
	this.hWndPopupControlParent, _ = win32.GetParent(this.hWndPopupControl)

	err := this.Create(WindowOptions{
		ParentHandle: this.hWndTopWin,
		//StyleExclude: win32.WS_VISIBLE,
	})

	if this.HasBorder {
		this.border = NewDropdownPopupBorderObject()
		err = this.border.Create(WindowOptions{
			ParentHandle: this.Handle,
		})
		win32.ShowWindow(this.border.Handle, win32.SW_SHOWNOACTIVATE)
	}

	return err
}

func (this *DropdownPopupContainerObject) Close() {
	win32.ShowWindow(this.hWndPopupControl, win32.SW_HIDE)
	win32.SetParent(this.hWndPopupControl, this.hWndPopupControlParent)
	if this.HasBorder {
		this.border.Destroy()
	}
	this.Destroy()
}

func (this *DropdownPopupContainerObject) Show() {
	win32.SetParent(this.hWndPopupControl, this.Handle)

	this.Popup.NotifyBeforeShow()

	//
	x, y, cx, cy := this.getPopupBounds()
	xPopup, yPopup, cxPopup, cyPopup := 0, 0, cx, cy
	if this.HasBorder {
		xPopup, yPopup = 1, 1
		cxPopup -= 2
		cyPopup -= 2
	}

	win32.SetWindowPos(this.hWndPopupControl, 0,
		int32(xPopup), int32(yPopup), int32(cxPopup), int32(cyPopup),
		win32.SWP_SHOWWINDOW|win32.SWP_NOZORDER)

	//this.Popup.NotifyBeforeShow()

	var bAnim win32.BOOL
	if !this.NoAnim {
		win32.SystemParametersInfo(win32.SPI_GETCOMBOBOXANIMATION,
			0, unsafe.Pointer(&bAnim), 0)
	}

	//bAnim = 0
	if bAnim == 0 {
		win32.SetWindowPos(this.Handle, win32.HWND_TOP,
			int32(x), int32(y), int32(cx), int32(cy), win32.SWP_SHOWWINDOW)
	} else {
		win32.SetWindowPos(this.Handle, win32.HWND_TOP,
			int32(x), int32(y), int32(cx), int32(cy), win32.SWP_HIDEWINDOW)

		//win32.SendMessage(this.hWndTopWin, win32.WM_NCACTIVATE, 1, 0)
		win32.AnimateWindow(this.Handle, 180, win32.AW_SLIDE|win32.AW_VER_POSITIVE)
		//win32.SendMessage(this.hWndTopWin, win32.WM_NCACTIVATE, 1, 0)

		win32.RedrawWindow(this.Handle, nil, 0, win32.RDW_FRAME|win32.RDW_INVALIDATE)
	}

	win32.SetFocus(this.hWndPopupControl)
	//win32.SendMessage(this.hWndTopWin, win32.WM_NCACTIVATE, 1, 0)
	this.Popup.NotifyAfterShow()

	if this.HasBorder {
		this.border.Visible = true
		this.border.SetBounds(x, y, x+cx, y+cy)
		//win32.SetWindowPos(this.border.GetHandle(), win32.HWND_TOP,
		//	x, y, cx, cy, win32.SWP_SHOWWINDOW|win32.SWP_NOACTIVATE)
	}
	//win32.SendMessage(this.hWndTopWin, win32.WM_NCACTIVATE, 1, 0)
}

func (this *DropdownPopupContainerObject) getPopupBounds() (int, int, int, int) {
	//if this.GetPopupBoundsCallback != nil {
	//	return this.GetPopupBoundsCallback()
	//}
	var rc win32.RECT
	if this.DropdownRect.Right != this.DropdownRect.Left {
		rc = this.DropdownRect.ToRECT()
	} else {
		win32.GetWindowRect(this.hWndDropdownControl, &rc)
	}
	x, y := int(rc.Left), int(rc.Bottom)
	cxDdControl := int(rc.Right - rc.Left)

	cxScreen, _ := win32.GetSystemMetrics(win32.SM_CXSCREEN)
	cyScreen, _ := win32.GetSystemMetrics(win32.SM_CYSCREEN)

	width := cxDdControl
	maxWidth := int(cxScreen) - x
	maxHeight := int(cyScreen) - y

	popupWidth, maxPopupWidth, maxPopupHeight := width, maxWidth, maxHeight
	if this.HasBorder {
		popupWidth -= 2
		maxPopupWidth -= 2
		maxPopupHeight -= 2
	}

	cx, cy := this.Popup.GetPopupSize(popupWidth, maxPopupWidth, maxPopupHeight)
	if this.HasBorder {
		cx += 2
		cy += 2
	}

	if cx > maxWidth {
		cx = maxWidth
	}
	if cy > maxHeight {
		cy = maxHeight
	}
	if cx < width && !this.DropdownRightAlign {
		cx = width
	}
	//if this.DropdownRightAlign {
	x = int(rc.Right) - cx
	//}
	return x, y, cx, cy
}

func (this *DropdownPopupContainerObject) OnNotify(info *NotifyMessage) {
	pNmhdr := info.GetNMHDR()
	if pNmhdr.Code == win32.TVN_SELCHANGED {
		//println("???@@")
	}
}

func (this *DropdownPopupContainerObject) UpdateSize() {

	x, y, cx, cy := this.getPopupBounds()
	xPopup, yPopup, cxPopup, cyPopup := 0, 0, cx, cy
	if this.HasBorder {
		xPopup, yPopup = 1, 1
		cxPopup -= 2
		cyPopup -= 2
	}

	var swpFlags = win32.SWP_NOMOVE | win32.SWP_NOZORDER |
		win32.SWP_NOACTIVATE | win32.SWP_NOREDRAW

	if this.HasBorder {
		this.border.Visible = false
		this.border.Update()
		//win32.SetWindowPos(this.border.GetHandle(), 0, 0, 0, cx + 2, cy + 2, swpFlags)
	}

	//
	win32.SetWindowPos(this.hWndPopupControl, 0,
		int32(xPopup), int32(yPopup), int32(cxPopup), int32(cyPopup), swpFlags)
	win32.SetWindowPos(this.Handle, 0, 0, 0, int32(cx), int32(cy), swpFlags)

	if this.HasBorder {
		this.border.Visible = true
		this.border.SetBounds(x, y, -cx, -cy)
		//win32.SetWindowPos(this.border.GetHandle(), win32.HWND_TOP,
		//	x, y, cx, cy, win32.SWP_SHOWWINDOW|win32.SWP_NOACTIVATE)

	}

	win32.RedrawWindow(this.Handle, nil, 0,
		win32.RDW_ERASE|win32.RDW_FRAME|
			win32.RDW_INVALIDATE|win32.RDW_ALLCHILDREN|win32.RDW_UPDATENOW)

}
