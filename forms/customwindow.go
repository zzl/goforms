package forms

import (
	"github.com/zzl/goforms/drawing/colors"
	"unsafe"

	"github.com/zzl/go-win32api/v2/win32"
)

// CustomWindow is an interface for custom windows.
// Custom windows are user defined, not system provided windows.
type CustomWindow interface {
	Window // the parent interface

	IsMouseHovering() bool // checks if the mouse is inside the window's boundary
}

// CustomWindowSpi is an interface that provides additional methods
// specific to implementing a CustomWindow.
type CustomWindowSpi interface {
	WindowSpi       // the parent interface
	WinProcProvider // provides a custom WinProc

	OnCommand(msg *CommandMessage)                 // processes WM_COMMAND
	OnNotify(msg *NotifyMessage)                   // processes WM_NOTIFY
	OnContextMenu(sourceWin Window, pt Point) bool // processes WM_CONTEXT_MENU
	OnEraseBkgnd(hdc win32.HDC) bool               // processes WM_ERASEBKGND
	OnCtlColorStatic(msg *Message)                 // processes WM_CTLCOLORSTATIC
	OnDestroy()                                    // on WM_DESTROY

	OnSetFocus()  // on WM_SETFOCUS
	OnKillFocus() // on WM_KILLFOCUS

	OnMouseDown(x int32, y int32, button byte) // on WM_LBUTTONDOWN/WM_RBUTTONDOWN
	OnMouseUp(x int32, y int32, button byte)   // on WM_LBUTTONUP/WM_RBUTTONUP
	OnMouseMove(x int32, y int32, button byte) // on WM_MOUSEMOVE

	OnMouseEnter() // called when the mouse enters the window's boundary
	OnMouseLeave() // called when the mouse leaves the window's boundary

	OnKeyDown(args KeyEventArgs) // on WM_KEYDOWN
	OnKeyUp(args KeyEventArgs)   // on WM_KEYUP

	OnPaint(hdc win32.HDC, clipRect *win32.RECT) // on WM_PAINT

	// OnBubbleMessage processes reflection messages bubbled up from a child container.
	OnBubbleMessage(msg *Message)
}

// CustomWindowInterface is a composition of CustomWindow and CustomWindowSpi
type CustomWindowInterface interface {
	CustomWindow
	CustomWindowSpi
}

// CustomWindowObject implements the CustomWindowInterface
// It extends WindowObject.
type CustomWindowObject struct {

	// WindowObject is the parent struct.
	WindowObject

	// super is the special pointer to the parent struct.
	super *WindowObject

	mouseHovering bool // whether mouse is hovering over the window
}

// Init implements Window.Init
func (this *CustomWindowObject) Init() {
	this.super.Init()
}

// IsMouseHovering implements CustomWindow.IsMouseHovering
func (this *CustomWindowObject) IsMouseHovering() bool {
	return this.mouseHovering
}

// OnEraseBkgnd implements CustomWindowSpi.OnEraseBkgnd
func (this *CustomWindowObject) OnEraseBkgnd(hdc win32.HDC) bool {
	bgColor := this.GetBackColor()
	if bgColor == colors.Null {
		ret := win32.DefWindowProc(this.Handle,
			win32.WM_ERASEBKGND, hdc, 0)
		return ret != 0
	}
	if !bgColor.IsTransparent() {
		var rc win32.RECT
		win32.GetClipBox(hdc, &rc)
		FillSolidRect(hdc, &rc, bgColor.Win32Color())
	}
	return true
}

// OnCommand implements CustomWindowSpi.OnCommand
func (this *CustomWindowObject) OnCommand(info *CommandMessage) {
	//nop
}

// OnNotify implements CustomWindowSpi.OnNotify
func (this *CustomWindowObject) OnNotify(msg *NotifyMessage) {
	//?
}

// WinProc implements CustomWindowSpi.WinProc
func (this *CustomWindowObject) WinProc(winObj *WindowObject, m *Message) error {
	win := winObj.RealObject.(CustomWindowInterface)
	switch m.UMsg {
	case win32.WM_CONTEXTMENU: //to c
		sourceWin := GetWindow(m.WParam)
		pt := PointFromDWORD(win32.DWORD(m.LParam))
		if win.OnContextMenu(sourceWin, pt) {
			m.Handled = true
			return nil
		}
	case win32.WM_ERASEBKGND: //to cust control
		if win.OnEraseBkgnd(m.WParam) {
			m.Result = 1
		}
		m.Handled = true
		return nil
	case win32.WM_PAINT: //?
		var ps win32.PAINTSTRUCT
		win32.BeginPaint(this.Handle, &ps)
		win.OnPaint(ps.Hdc, &ps.RcPaint)
		win32.EndPaint(this.Handle, &ps)

	case win32.WM_CTLCOLORSTATIC: //to c
		checkReflectOrBubble(win, m)
		if m.Handled {
			return nil
		}
		win.OnCtlColorStatic(m)
		if m.Handled {
			return nil
		}
	case win32.WM_CTLCOLORBTN, win32.WM_CTLCOLOREDIT,
		win32.WM_CTLCOLORLISTBOX, win32.WM_CTLCOLORSCROLLBAR:
		checkReflectOrBubble(win, m)
		if m.Handled {
			return nil
		}
	case win32.WM_CTLCOLORDLG:
		println("??")
	case win32.WM_COMMAND: //to c
		win.OnCommand((*CommandMessage)(m))
		if m.Handled {
			return nil
		}
		checkReflectOrBubble(win, m)
		if m.Handled {
			return nil
		}
	case win32.WM_NOTIFY: //to c
		//win.OnNotify((*NotifyMessage)(m))
		win.OnNotify(&NotifyMessage{m})
		if m.Handled {
			return nil
		}
		checkReflectOrBubble(win, m)
		if m.Handled {
			return nil
		}
	case win32.WM_MEASUREITEM: //to c
		checkReflectOrBubble(win, m)
		if m.Handled {
			return nil
		}
	case win32.WM_DRAWITEM: //to c
		checkReflectOrBubble(win, m)
		if m.Handled {
			return nil
		}
	case win32.WM_DESTROY: //to c
		win.OnDestroy()
		winObj.callPreDispose()
	case win32.WM_NCDESTROY: //to c
		winObj.callDispose()
	//
	case win32.WM_SETFOCUS:
		win.OnSetFocus()
	case win32.WM_KILLFOCUS:
		win.OnKillFocus()
	case win32.WM_MOUSEMOVE:
		x, y, button := ParseMouseMsgParams(m.WParam, m.LParam)
		win.OnMouseMove(x, y, button)
		if !this.mouseHovering {
			this.mouseHovering = true
			win.OnMouseEnter()
			win32.SetTimer(this.Handle, CONTROL_TIMER_TRACK_LEAVE, 100, 0)
		}
	case win32.WM_LBUTTONDOWN, win32.WM_RBUTTONDOWN:
		x, y, button := ParseMouseMsgParams(m.WParam, m.LParam)
		win.OnMouseDown(x, y, button)
	case win32.WM_LBUTTONUP, win32.WM_RBUTTONUP:
		x, y, button := ParseMouseMsgParams(m.WParam, m.LParam)
		win.OnMouseUp(x, y, button)
	case win32.WM_KEYDOWN:
		win.OnKeyDown(NewKeyEventArgs(m.WParam, m.LParam))
	case win32.WM_KEYUP:
		win.OnKeyUp(NewKeyEventArgs(m.WParam, m.LParam))

	case win32.WM_TIMER:
		if m.WParam == CONTROL_TIMER_TRACK_LEAVE {
			hovered := true
			var rc win32.RECT
			win32.GetWindowRect(this.Handle, &rc)
			var pt win32.POINT
			win32.GetCursorPos(&pt)
			if pt.X < rc.Left || pt.X > rc.Right ||
				pt.Y < rc.Top || pt.Y > rc.Bottom {
				hovered = false
			}
			if !hovered {
				//?win32.KillTimer(this.Handle, RICH_EDIT_TIMER_TRACK_LEAVE)
				if this.mouseHovering {
					this.mouseHovering = false
					win.OnMouseLeave()
				}
			}
		}
	}
	return nil
}

// checkReflectOrBubble is a helper function for dispatching notification messages,
// first down to the message's source control, then to the parent of this window.
func checkReflectOrBubble(win CustomWindowInterface, m *Message) {
	var hWndCtrl HWND
	if m.UMsg == win32.WM_NOTIFY {
		pNmHdr := (*win32.NMHDR)(unsafe.Pointer(m.LParam))
		hWndCtrl = pNmHdr.HwndFrom
	} else if m.UMsg == win32.WM_MEASUREITEM {
		if m.WParam == 0 {
			return //??
		}
		controlId := int(m.WParam)
		ret, _ := win32.GetDlgItem(win.GetHandle(), int32(controlId))
		hWndCtrl = ret
		// measureitem special handling
		if _, ok := windowMap[hWndCtrl]; !ok {
			if ctrlWindow, ok := creatingControlMap[controlId]; ok {
				ctrlWindow.AsWindowObject().Handle = hWndCtrl
				windowMap[hWndCtrl] = ctrlWindow
			}
		}
	} else if m.UMsg == win32.WM_DRAWITEM {
		pdis := (*win32.DRAWITEMSTRUCT)(unsafe.Pointer(m.LParam))
		hWndCtrl = pdis.HwndItem
	} else { //?
		hWndCtrl = HWND(m.LParam)
	}
	ctrlWin, ok := windowMap[hWndCtrl]
	if ok {
		ctrlWin.OnReflectMessage(m)
		if m.Handled {
			return
		}
	}
	win.OnBubbleMessage(m)
}

// OnContextMenu implements CustomWindowSpi.OnContextMenu
func (this *CustomWindowObject) OnContextMenu(sourceWin Window, pt Point) bool {
	if sourceWin.GetContextMenu() != nil {
		sourceWin.GetContextMenu().Show(pt.X, pt.Y, sourceWin.GetHandle())
		return true
	}
	return false
}

// OnCtlColorStatic implements CustomWindowSpi.OnCtlColorStatic
func (this *CustomWindowObject) OnCtlColorStatic(msg *Message) {
	//
}

// OnSetFocus implements CustomWindowSpi.OnSetFocus
func (this *CustomWindowObject) OnSetFocus() {
	//
}

// OnKillFocus implements CustomWindowSpi.OnKillFocus
func (this *CustomWindowObject) OnKillFocus() {
	//
}

// OnMouseDown implements CustomWindowSpi.OnMouseDown
func (this *CustomWindowObject) OnMouseDown(x int32, y int32, button byte) {
	//
}

// OnMouseMove implements CustomWindowSpi.OnMouseMove
func (this *CustomWindowObject) OnMouseMove(x int32, y int32, button byte) {
	//
}

// OnMouseUp implements CustomWindowSpi.OnMouseUp
func (this *CustomWindowObject) OnMouseUp(x int32, y int32, button byte) {
	//
}

// OnMouseEnter implements CustomWindowSpi.OnMouseEnter
func (this *CustomWindowObject) OnMouseEnter() {
	//
}

// OnMouseLeave implements CustomWindowSpi.OnMouseLeave
func (this *CustomWindowObject) OnMouseLeave() {
	//
}

// OnDestroy implements CustomWindowSpi.OnDestroy
func (this *CustomWindowObject) OnDestroy() {
	//
}

// OnKeyDown implements CustomWindowSpi.OnKeyDown
func (this *CustomWindowObject) OnKeyDown(args KeyEventArgs) {
	//
}

// OnKeyUp implements CustomWindowSpi.OnKeyUp
func (this *CustomWindowObject) OnKeyUp(args KeyEventArgs) {
	//
}

// OnPaint implements CustomWindowSpi.OnPaint
func (this *CustomWindowObject) OnPaint(hdc win32.HDC, clipRect *win32.RECT) {
	//
}

// OnBubbleMessage implements CustomWindowSpi.OnBubbleMessage
func (this *CustomWindowObject) OnBubbleMessage(msg *Message) {
	if parent, ok := windowMap[this.GetParentHandle()].(CustomWindowSpi); ok {
		parent.OnBubbleMessage(msg)
	}
}
