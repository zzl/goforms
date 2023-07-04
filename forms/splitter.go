package forms

import (
	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/framework/virtual"
	"github.com/zzl/goforms/layouts"
	"log"
	"unsafe"

	"github.com/zzl/go-win32api/v2/win32"
)

type Splitter interface {
	Control
}

type SplitterObject struct {
	ControlObject
	super *ControlObject

	Width int

	dragInfo *splitterDragInfo
}

type NewSplitter struct {
	Parent Container
	Name   string
	Pos    Point
	Size   Size
	Width  int
}

func (me NewSplitter) Create(extraOpts ...*WindowOptions) Splitter {
	splitter := NewSplitterObject()
	splitter.name = me.Name
	splitter.Width = me.Width

	opts := utils.OptionalArg(extraOpts)
	opts.Left = me.Pos.X
	opts.Top = me.Pos.Y

	opts.ParentHandle = resolveParentHandle(me.Parent)
	err := splitter.Create(*opts)
	assertNoErr(err)
	configControlSize(splitter, me.Size)

	return splitter
}

func NewSplitterObject() *SplitterObject {
	return virtual.New[SplitterObject]()
}

type splitterDragInfo struct {
	rootHwnd HWND
	rootHdc  win32.HDC
	rootCx   int

	rcStart win32.RECT
	rcDraw  win32.RECT

	ptStart Point
	ptMove  Point
}

const DefaultSplitterWidth = 4

var _splitterClassRegistered = false

func ensureSplitterClassRegistered() {
	if _splitterClassRegistered {
		return
	}
	_, err := RegisterClass("goforms.splitter", win32.DefWindowProc, ClassOptions{
		BackgroundBrush: 0,
		CursorResId:     int(uintptr(unsafe.Pointer(win32.IDC_SIZEWE))),
	})
	if err != nil {
		log.Fatal(err)
	}

	_splitterClassRegistered = true
}

func (this *SplitterObject) Init() {
	this.super.Init()
	this.WinProcFunc = splitterWndProc
}

func (this *SplitterObject) EnsureClassRegistered() {
	ensureSplitterClassRegistered()
}

func (this *SplitterObject) GetWindowClass() string {
	return "goforms.splitter"
}

func splitterWndProc(win *WindowObject, m *Message) error {
	this := win.RealObject.(*SplitterObject)
	switch m.UMsg {
	case win32.WM_LBUTTONDOWN:
		this.onLButtonDown(m.WParam, m.LParam)
		return m.SetHandledWithResult(0)
	case win32.WM_MOUSEMOVE:
		this.onMouseMove(m.WParam, m.LParam)
		return m.SetHandledWithResult(0)
	case win32.WM_LBUTTONUP:
		this.onLButtonUp(m.WParam, m.LParam)
		return m.SetHandledWithResult(0)
	case win32.WM_GETDLGCODE:
		return m.SetHandledWithResult(win32.LRESULT(win32.DLGC_WANTMESSAGE))
	case win32.WM_KEYDOWN:
		this.onKeyDown(m.WParam, m.LParam)
		return m.SetHandledWithResult(0)
	}
	return nil
}

func (this *SplitterObject) onKeyDown(wParam WPARAM, lParam LPARAM) {
	if wParam == 27 {
		println("esc..")
		this.dragEnd(true)
	} else if win32.VIRTUAL_KEY(wParam) == win32.VK_TAB {
		gwCmd := uint32(win32.GW_HWNDNEXT)
		ret := win32.GetKeyState(int32(win32.VK_SHIFT))
		if ret < 0 {
			gwCmd = uint32(win32.GW_HWNDPREV)
		}
		hWndTarget, _ := win32.GetWindow(this.Handle, win32.GET_WINDOW_CMD(gwCmd))
		win32.SetFocus(hWndTarget)
	}
}

func (this *SplitterObject) dragStart(ptDown Point) {
	di := &splitterDragInfo{}
	di.ptStart = ptDown

	pattern := make([]int16, 8)
	for n := 0; n < 8; n++ {
		pattern[n] = (int16)(0x5555 << (n & 1))
	}
	hBmp := win32.CreateBitmap(8, 8, 1, 1, unsafe.Pointer(&pattern[0]))
	hWnd := this.GetRootWindow().GetHandle()

	var rc win32.RECT
	win32.GetClientRect(hWnd, &rc)
	di.rootCx = int(rc.Right)

	di.rootHwnd = hWnd
	di.rootHdc = win32.GetDC(hWnd)

	var lb win32.LOGBRUSH
	lb.LbStyle = win32.BS_PATTERN
	lb.LbHatch = uintptr(hBmp)
	hBrush := win32.CreateBrushIndirect(&lb)
	win32.DeleteObject(win32.HGDIOBJ(hBmp))

	win32.SelectObject(di.rootHdc, win32.HGDIOBJ(hBrush))

	win32.SetCapture(this.Handle)

	//
	var pt win32.POINT
	win32.ClientToScreen(this.Handle, &pt)
	win32.ScreenToClient(di.rootHwnd, &pt)

	win32.GetClientRect(this.Handle, &rc)
	di.rcStart = win32.RECT{Left: pt.X, Top: pt.Y,
		Right: pt.X + rc.Right, Bottom: pt.Y + rc.Bottom}

	//
	di.ptMove.X = -1 //
	this.dragInfo = di
	this.dragMove(ptDown)
}

func (this *SplitterObject) dragMove(ptMove Point) {
	di := this.dragInfo
	if ptMove.X == di.ptMove.X {
		return
	}

	dx := int32(ptMove.X - di.ptStart.X)
	rc := di.rcStart
	rc.Left += dx
	rc.Right += dx
	rcWidth := rc.Right - rc.Left

	if rc.Left < 16 {
		rc.Left = 16
		rc.Right = rc.Left + rcWidth
	} else if rc.Right > int32(di.rootCx)-16 {
		rc.Right = int32(di.rootCx) - 16
		rc.Left = rc.Right - rcWidth
	}

	lastRc := di.rcDraw
	if lastRc.Right != 0 {
		win32.PatBlt(di.rootHdc, lastRc.Left, lastRc.Top,
			lastRc.Right-lastRc.Left, lastRc.Bottom-lastRc.Top, win32.PATINVERT)
	}
	win32.PatBlt(di.rootHdc, rc.Left, rc.Top,
		rc.Right-rc.Left, rc.Bottom-rc.Top, win32.PATINVERT)

	di.rcDraw = rc
	di.ptMove = ptMove
}

func (this *SplitterObject) dragEnd(canceled bool) {
	di := this.dragInfo
	if di == nil {
		println("???")
		return
	}
	this.dragInfo = nil

	lastRc := di.rcDraw
	win32.PatBlt(di.rootHdc, lastRc.Left, lastRc.Top,
		lastRc.Right-lastRc.Left, lastRc.Bottom-lastRc.Top, win32.PATINVERT)

	dx := int(di.rcDraw.Left - di.rcStart.Left)

	hbrWindow := win32.GetSysColorBrush(win32.COLOR_WINDOW)
	hBrush := win32.SelectObject(di.rootHdc, win32.HGDIOBJ(hbrWindow))

	win32.ReleaseDC(di.rootHwnd, di.rootHdc)
	win32.DeleteObject(hBrush)
	win32.ReleaseCapture()

	hWndPrev, _ := win32.GetWindow(this.Handle, win32.GW_HWNDPREV)
	hWndNext, _ := win32.GetWindow(this.Handle, win32.GW_HWNDNEXT)

	//layout := this.GetContainer().GetLayoutForControl(this)
	layout := this.GetData(layouts.Data_Layout).(Layout)
	if layout != nil {
		var rc win32.RECT
		win32.GetWindowRect(hWndPrev, &rc)
		prevWidth := int(rc.Right-rc.Left) + dx

		prevControl := GetWindow(hWndPrev).(Control)
		prevItem := layout.FindItemByControl(prevControl)
		prevItem.SetWidth(prevWidth)

		win32.GetWindowRect(hWndNext, &rc)
		nextWidth := int(rc.Right-rc.Left) - dx
		nextControl := GetWindow(hWndNext).(Control)
		nextItem := layout.FindItemByControl(nextControl)
		nextItem.SetWidth(nextWidth)

		layout.Update()
	} else {
		//?
	}
}

func (this *SplitterObject) onMouseMove(wParam WPARAM, lParam LPARAM) {
	if this.dragInfo == nil {
		return
	}
	pt := Point{X: int(int16(win32.LOWORD(win32.DWORD(lParam)))),
		Y: int(int16(win32.HIWORD(win32.DWORD(lParam))))}
	this.dragMove(pt)
}

func (this *SplitterObject) onLButtonUp(wParam WPARAM, lParam LPARAM) {
	this.dragEnd(false)
}

func (this *SplitterObject) onLButtonDown(wParam WPARAM, lParam LPARAM) {
	pt := Point{X: int(win32.LOWORD(win32.DWORD(lParam))),
		Y: int(win32.HIWORD(win32.DWORD(lParam)))}
	this.dragStart(pt)
}

func (this *SplitterObject) Create(options WindowOptions) error {
	return this.super.Create(options)
}

func (this *SplitterObject) GetPreferredSize(maxWidth int, maxHeight int) (int, int) {
	width := this.Width
	if width == 0 {
		width = DefaultSplitterWidth
	}
	return width, maxHeight
}
