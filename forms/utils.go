package forms

import (
	"github.com/zzl/goforms/framework/consts"
	"log"
	"runtime"
	"strings"
	"syscall"
	"unsafe"

	"github.com/zzl/goforms/framework/scope"
	"github.com/zzl/goforms/framework/types"

	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/drawing"
	"github.com/zzl/goforms/drawing/colors"
)

var _defaultFont win32.HFONT

func GetDefaultFont() win32.HFONT {
	if _defaultFont == 0 {
		var ncm win32.NONCLIENTMETRICS
		ncm.CbSize = uint32(unsafe.Sizeof(ncm))
		win32.SystemParametersInfo(win32.SPI_GETNONCLIENTMETRICS,
			ncm.CbSize, unsafe.Pointer(&ncm), 0)
		_defaultFont = win32.CreateFontIndirect(&ncm.LfMessageFont)
	}
	return _defaultFont
}

type ClassOptions struct {
	Style           win32.WNDCLASS_STYLES
	BackgroundBrush win32.HBRUSH
	CursorResId     int
}

func MustRegisterClass(className string, wndProc WndProcFunc, options ClassOptions) {
	_, err := RegisterClass(className, wndProc, options)
	if err != nil {
		panic(err)
	}
}

func RegisterClass(className string, wndProc WndProcFunc, options ClassOptions) (win32.ATOM, error) {
	var pWndProc uintptr
	if wndProc == nil {
		pWndProc = wndProcCallback
	} else {
		pWndProc = syscall.NewCallback(wndProc)
	}

	cursorResId := options.CursorResId
	if cursorResId == 0 {
		cursorResId = int(uintptr(unsafe.Pointer(win32.IDC_ARROW)))
	}
	cursor, _ := win32.LoadCursor(0, win32.MAKEINTRESOURCE(uint16(cursorResId)))
	style := options.Style

	wc := win32.WNDCLASSEX{
		HInstance:     HInstance,
		LpszClassName: win32.StrToPwstr(className),
		LpfnWndProc:   pWndProc,
		Style:         win32.WNDCLASS_STYLES(style),
		HCursor:       cursor,
		HbrBackground: options.BackgroundBrush,
	}
	wc.CbSize = uint32(unsafe.Sizeof(wc))

	var err error
	retVal, errno := win32.RegisterClassEx(&wc)
	if retVal == 0 {
		err = errno
	}
	return retVal, err
}

//

type WindowOptions struct {
	//WinProcFunc WndProcFunc
	ClassName      string
	WindowName     string
	ParentHandle   HWND
	HMENU          win32.HMENU
	Style          WINDOW_STYLE
	StyleInclude   WINDOW_STYLE
	StyleExclude   WINDOW_STYLE
	ExStyle        WINDOW_EX_STYLE
	ExStyleExclude WINDOW_EX_STYLE
	ExStyleInclude WINDOW_EX_STYLE
	Width          int
	Height         int
	Left           int
	Top            int
	ControlId      uint16
}

func resolveWindowOptions(opts *WindowOptions) {
	left := int(win32.CW_USEDEFAULT)
	if opts.Left != 0 {
		left = opts.Left
		if left == consts.Zero {
			left = 0
		}
	}
	top := int(win32.CW_USEDEFAULT)
	if opts.Top != 0 {
		top = opts.Top
		if top == consts.Zero {
			top = 0
		}
	}
	width := int(win32.CW_USEDEFAULT)
	if opts.Width != 0 {
		width = opts.Width
		if width == consts.Zero {
			width = 0
		}
	}
	height := int(win32.CW_USEDEFAULT)
	if opts.Height != 0 {
		height = opts.Height
		if height == consts.Zero {
			height = 0
		}
	}
	if opts.HMENU == 0 && opts.ControlId != 0 {
		opts.HMENU = win32.HMENU(opts.ControlId)
	}
	//
	opts.Style = (opts.Style | opts.StyleInclude) &^ opts.StyleExclude
	opts.ExStyle = (opts.ExStyle | opts.ExStyleInclude) &^ opts.ExStyleExclude
}

func CreateWindow(className string,
	options WindowOptions, lpParam unsafe.Pointer) (HWND, error) {

	resolveWindowOptions(&options)

	hWnd, errno := win32.CreateWindowEx(
		options.ExStyle,
		win32.StrToPwstr(className),
		win32.StrToPwstr(options.WindowName),
		options.Style,
		int32(options.Left),
		int32(options.Top),
		int32(options.Width),
		int32(options.Height),
		options.ParentHandle,
		options.HMENU,
		HInstance,
		lpParam)

	var err error
	if hWnd == 0 {
		err = errno
	}
	return hWnd, err
}

func SendMessage[T_WPARAM, T_LPARAM types.PtrCompatible](
	hWnd HWND, Msg uint32, wParam T_WPARAM, lParam T_LPARAM) (uintptr, WIN32_ERROR) {
	return win32.SendMessage(hWnd, Msg, WPARAM(wParam), LPARAM(lParam))
}

func GetMessagePos() win32.POINT {
	dwPos := win32.GetMessagePos()
	x, y := int32(int16(win32.LOWORD(dwPos))), int32(int16(win32.HIWORD(dwPos)))
	return win32.POINT{X: x, Y: y}
}

func GetMessageClientPos(hWnd win32.HWND) win32.POINT {
	pos := GetMessagePos()
	win32.ScreenToClient(hWnd, &pos)
	return pos
}

func IsChildWindow(hWnd win32.HWND) bool {
	dwStyle, _ := win32.GetWindowLong(hWnd, win32.GWL_STYLE)
	if win32.WINDOW_STYLE(dwStyle)&win32.WS_CHILD != 0 {
		return true
	}
	return false
}

func ContainsWindow(hWndAncestor win32.HWND, hWndTest win32.HWND) bool {
	for hWndTest != 0 {
		if hWndTest == hWndAncestor {
			return true
		}
		hWndTest, _ = win32.GetParent(hWndTest)
	}
	return false
}

func drawSizeGrip(hWnd HWND, hdc win32.HDC) {
	var rc win32.RECT
	win32.GetClientRect(hWnd, &rc)
	pwsz, _ := syscall.UTF16PtrFromString("Status")
	hTheme := win32.OpenThemeData(hWnd, pwsz)
	//hTheme = 0
	if hTheme == 0 {
		size, _ := win32.GetSystemMetrics(win32.SM_CXVSCROLL)
		rc.Left = rc.Right - size
		rc.Top = rc.Bottom - size
		win32.DrawFrameControl(hdc, &rc, win32.DFC_SCROLL, win32.DFCS_SCROLLSIZEGRIP)
	} else {
		var size win32.SIZE
		win32.GetThemePartSize(hTheme, hdc,
			int32(win32.SP_GRIPPER), 0, nil, win32.TS_DRAW, &size)

		rc.Left = rc.Right - size.Cx
		rc.Top = rc.Bottom - size.Cy
		win32.DrawThemeBackground(hTheme, hdc, int32(win32.SP_GRIPPER),
			0, &rc, nil)

		win32.CloseThemeData(hTheme)
	}
}

func hitTestSizeGrip(hWnd HWND, msg *Message) {
	dwXy := win32.DWORD(msg.LParam)
	x, y := win32.LOWORD(dwXy), win32.HIWORD(dwXy)

	var rc win32.RECT
	win32.GetClientRect(hWnd, &rc)
	size, _ := win32.GetSystemMetrics(win32.SM_CXVSCROLL)
	rc.Left = rc.Right - size
	rc.Top = rc.Bottom - size

	pt := win32.POINT{X: int32(x), Y: int32(y)}
	win32.ScreenToClient(hWnd, &pt)
	if PtInRect(&rc, pt) {
		msg.Result = uintptr(win32.HTBOTTOMRIGHT)
		msg.Handled = true
	}
}

//var colorBrushMap map[win32.COLORREF]win32.HBRUSH

func handleCtlColor(win Window, hdc win32.HDC, backColor drawing.Color) win32.HBRUSH {
	hWnd := win.GetHandle()
	if backColor == colors.Null {
		return 0
	}
	if backColor == colors.Transparent {
		hWndParent, _ := win32.GetParent(hWnd)
		dwStyle, _ := win32.GetWindowLongPtr(hWndParent, win32.GWL_STYLE)
		if dwStyle&uintptr(win32.WS_CLIPCHILDREN) != 0 {

			var rc win32.RECT
			win32.GetClientRect(hWnd, &rc)
			hRgn := win32.CreateRectRgn(0, 0, rc.Right, rc.Bottom)
			win32.SelectClipRgn(hdc, hRgn)

			var ptOffset win32.POINT
			win32.MapWindowPoints(hWnd, hWndParent, &ptOffset, 1)
			var ptOrig win32.POINT
			win32.OffsetWindowOrgEx(hdc, ptOffset.X, ptOffset.Y, &ptOrig)
			SendMessage(hWndParent, win32.WM_ERASEBKGND, hdc, 0)

			win32.SelectClipRgn(hdc, 0)
			win32.DeleteObject(win32.HGDIOBJ(hRgn))

			win32.SetWindowOrgEx(hdc, ptOrig.X, ptOrig.Y, nil)
		}

		win32.SetBkMode(hdc, win32.TRANSPARENT)
		hbr := win32.GetStockObject(win32.NULL_BRUSH)
		return hbr
	} else {
		oldHbr := win.GetData(Data_BackColorBrush)
		if oldHbr != nil {
			win32.DeleteObject(oldHbr.(win32.HBRUSH))
		}
		win32.SetBkMode(hdc, win32.TRANSPARENT)
		hbr := win32.CreateSolidBrush(backColor.Win32Color())
		win.SetData(Data_BackColorBrush, hbr)
		return hbr
	}
}

func RectWithSize(left int, top int, width int, height int) Rect {
	return Rect{Left: left, Top: top, Right: left + width, Bottom: top + height}
}

func RectFromRECT(rc win32.RECT) Rect {
	return Rect{Left: int(rc.Left), Top: int(rc.Top),
		Right: int(rc.Right), Bottom: int(rc.Bottom)}
}

func RequestGC() {
	Dispatcher.Invoke(runtime.GC)
}

func ToSysColorBrush(color byte) win32.HBRUSH {
	return win32.HBRUSH(color + 1)
}

func StrToWsz(str string) []uint16 {
	var wcs []uint16
	var end = false
	for {
		pos := strings.IndexByte(str, '\000')
		if pos == -1 {
			pos = len(str)
			end = true
		}
		wsz, _ := syscall.UTF16FromString(str[:pos])
		wcs = append(wcs, wsz...)
		if end {
			break
		}
		str = str[pos+1:]
	}
	return wcs
}

func FillSolidRect(hdc win32.HDC, lpRect *win32.RECT, clr win32.COLORREF) {
	win32.SetBkColor(hdc, clr)
	win32.ExtTextOut(hdc, 0, 0, win32.ETO_OPAQUE, lpRect, nil, 0, nil)
}

var _enumedThreadWindows []Window
var _enumThreadWindowsCallback = syscall.NewCallback(
	func(hWnd HWND, lparam LPARAM) win32.LRESULT {
		win, ok := windowMap[hWnd]
		if ok {
			_enumedThreadWindows = append(_enumedThreadWindows, win)
		}
		return 1
	})

func GetTopWindows() []Window {
	threadId := win32.GetCurrentThreadId()
	_enumedThreadWindows = nil
	win32.EnumThreadWindows(win32.DWORD(threadId),
		_enumThreadWindowsCallback, 0)
	return _enumedThreadWindows
}

func IsFocusable(hWnd HWND) bool {
	if hWnd == 0 {
		return false
	}
	visible := win32.IsWindowVisible(hWnd)
	if visible == win32.FALSE {
		return false
	}
	enabled := win32.IsWindowEnabled(hWnd)
	if enabled == win32.FALSE {
		return false
	}
	dwStyle, _ := win32.GetWindowLong(hWnd, win32.GWL_STYLE)
	return dwStyle&int32(win32.WS_TABSTOP) != 0
}

func GetChildHandles(hWndParent HWND) []HWND {
	var childHwnds []HWND
	hWnd, _ := win32.GetWindow(hWndParent, win32.GW_CHILD)
	for hWnd != 0 {
		childHwnds = append(childHwnds, hWnd)
		hWnd, _ = win32.GetWindow(hWnd, win32.GW_HWNDNEXT)
	}
	return childHwnds
}

func GetDescendantHandles(hWndRoot HWND) []HWND {
	var descendants []HWND
	children := GetChildHandles(hWndRoot)
	for _, c := range children {
		descendants = append(descendants, c)
		if _, ok := GetWindow(c).(Container); ok {
			for _, cc := range GetDescendantHandles(c) {
				descendants = append(descendants, cc)
			}
		}
	}
	return descendants
}

func _findFirstFocusable(hWndParent HWND) HWND {
	childHwnds := GetChildHandles(hWndParent)
	for _, hWnd := range childHwnds {
		if IsFocusable(hWnd) {
			return hWnd
		}
		childFocusable := _findFirstFocusable(hWnd)
		if childFocusable != 0 {
			return childFocusable
		}
	}
	return 0
}

func FindFirstFocusable(hWndContainer HWND) HWND {
	return _findFirstFocusable(hWndContainer)
}

func HandleTabFocus(hWnd HWND) {
	hWndRoot := win32.GetAncestor(hWnd, win32.GA_ROOT)

	//dunno y..
	SendMessage(hWndRoot, win32.WM_UPDATEUISTATE,
		win32.MAKELONG(uint16(win32.UIS_CLEAR), uint16(win32.UISF_HIDEFOCUS)), 0)

	keyState := win32.GetKeyState(int32(win32.VK_SHIFT))
	shiftDown := keyState < 0
	hWndTabNext, _ := win32.GetNextDlgTabItem(hWndRoot, hWnd, win32.BoolToBOOL(shiftDown))
	win32.SetFocus(hWndTabNext)
}

func PtInRect(rect *win32.RECT, pt win32.POINT) bool {
	return pt.X >= rect.Left && pt.X <= rect.Right &&
		pt.Y >= rect.Top && pt.Y <= rect.Bottom
}

func ParseMouseMsgParams(wParam WPARAM, lParam LPARAM) (int32, int32, byte) {
	dw := win32.DWORD(lParam)
	x, y := win32.GET_X_LPARAM(dw), win32.GET_Y_LPARAM(dw)
	button := byte(0)
	if uint32(wParam)&uint32(win32.MK_LBUTTON) != 0 {
		button = 1
	} else if uint32(wParam)&uint32(win32.MK_RBUTTON) != 0 {
		button = 2
	}
	return x, y, button
}

func IsRectEmpty(rect *win32.RECT) bool {
	return rect.Right == rect.Left || rect.Bottom == rect.Top
}

func DeflateRect(rect *win32.RECT, x int32, y int32) {
	rect.Left += x
	rect.Right -= x
	rect.Top += y
	rect.Bottom -= y
}

func MessageBoxEx(hWnd HWND, text, caption string,
	uType win32.MESSAGEBOX_STYLE) (win32.MESSAGEBOX_RESULT, win32.WIN32_ERROR) {
	pwszText, _ := syscall.UTF16PtrFromString(text)
	pwszCaption, _ := syscall.UTF16PtrFromString(caption)
	return win32.MessageBox(hWnd, pwszText, pwszCaption, uType)
}

func GetWindowText(hwnd HWND) (string, win32.WIN32_ERROR) {
	textLen, _ := win32.GetWindowTextLength(hwnd)
	textLen += 1

	buf := make([]uint16, textLen)
	_, err := win32.GetWindowText(hwnd, &buf[0], textLen)
	return syscall.UTF16ToString(buf), err
}

func SetWindowText(hWnd HWND, text string) (bool, win32.WIN32_ERROR) {
	pwsz, _ := syscall.UTF16PtrFromString(text)
	ret, err := win32.SetWindowText(hWnd, pwsz)
	return ret != win32.FALSE, err
}

func MessageBox(text string, title string) {
	MessageBoxEx(HWndActive, text, title, 0)
}

func Alert(text string) {
	MessageBox(text, "提示")
}

func Info(text string) {
	MessageBoxEx(HWndActive, text, "提示", win32.MB_ICONINFORMATION)
}

func Warn(text string) {
	MessageBoxEx(HWndActive, text, "提示", win32.MB_ICONEXCLAMATION)
}

func Confirm(text string) bool {
	ret, _ := MessageBoxEx(HWndActive, text, "提示",
		win32.MB_ICONQUESTION|win32.MB_OKCANCEL)
	return ret == win32.IDOK
}

func ConfirmYesNoCancel(text string) win32.MESSAGEBOX_RESULT {
	ret, _ := MessageBoxEx(HWndActive, text, "提示",
		win32.MB_ICONQUESTION|win32.MB_YESNOCANCEL|win32.MB_DEFBUTTON3)
	return ret
}

func GetWindow(hWnd HWND) WindowInterface {
	win, ok := windowMap[hWnd]
	if !ok {
		return nil
	}
	return win
}

func GetActiveWin() Window {
	return GetWindow(HWndActive)
}

func LoadIconFromImageData(data []byte) win32.HICON {
	bmp, _ := drawing.NewBitmapFromBytes(nil, data)
	hIcon := bmp.GetHIcon()
	bmp.Dispose()
	return hIcon
}

func NewScope() *scope.Scope {
	return scope.NewScope()
}

func WithScope(scopedFunc scope.ScopedFunc) {
	scope.WithScope(scopedFunc)
}

func Pt(x, y int) Point {
	return Point{X: x, Y: y}
}

func Sz(cx, cy int) Size {
	return Size{Width: cx, Height: cy}
}

func configControlSize(control Control, size Size) {
	if size.Width == 0 || size.Height == 0 {
		pw, ph := control.GetPreferredSize(0, 0)
		if size.Width == 0 {
			size.Width = pw
		}
		if size.Height == 0 {
			size.Height = ph
		}
	}
	control.SetSize(size.Width, size.Height)
}

func resolveParentHandle(parentWin Window) win32.HWND {
	if parentWin == nil {
		parentWin = ContextContainer
	}
	if parentWin == nil {
		parentWin = GetActiveWin()
	}
	return parentWin.GetHandle()
}

func assertNoErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
