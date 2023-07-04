package forms

import (
	"github.com/zzl/goforms/framework/consts"
	"log"
	"syscall"
	"unsafe"

	"github.com/zzl/go-win32api/v2/win32"
)

// TopWindow is an interface that represents a top level window,
// which can hold other controls and manage their layout.
type TopWindow interface {
	Container  // the parent interface
	NameAware  // name aware
	TitleAware // title aware

	Close() // closes the window

	CenterOnScreen()          // centers the window on the screen
	CenterToWindow(hWnd HWND) // centers the window relative to the specified window

	SetClientSize(width, height int)    // sets the client area size
	SetDluClientSize(width, height int) // sets the client area size using Dialog Units.

	ShowModal()                // shows the top window as a modal dialog
	SetModalResult(result int) // sets the modal result
	GetModalResult() int       // returns the modal result

	GetAccelHandle() win32.HACCEL // returns the handle to the accelerator table

	SetMenuBar(menuBar *MenuBar)         // sets the menu bar
	SetMenuHelpUI(menuHelpUI MenuHelpUI) // sets the menu help UI

	SetAccelerators(accelTable *AcceleratorTable) // sets the accelerator table

	showingModal() bool // whether the top window is showing as a modal dialog

	AsTopWindowObject() *TopWindowObject // returns the underlying TopWindowObject

}

// TopWindowSpi is an interface that provides additional methods
// specific to implementing a TopWindow.
type TopWindowSpi interface {
	WindowSpi
	OnMenuCommand(wParam WPARAM, lParam LPARAM)
}

// TopWindowInterface is a composition of TopWindow and TopWindowSpi
type TopWindowInterface interface {
	TopWindow
	TopWindowSpi
}

// TopWindowObject implements the TopWindowInterface
// It extends ContainerObject.
type TopWindowObject struct {

	// ContainerObject is the parent struct.
	ContainerObject

	// super is the special pointer to the parent struct.
	super *ContainerObject

	// NameAwareSupport is the NameAware implementation
	NameAwareSupport

	_showingModal bool
	accelTable    *AcceleratorTable
	menuBar       *MenuBar
	menuHelpUI    MenuHelpUI

	initialFocusHwnd HWND
	lastFocusedHwnd  HWND
	defId            int
}

func (this *TopWindowObject) Init() {
	this.super.Init()
	this.defId = int(win32.IDOK)
}

func (this *TopWindowObject) AsTopWindowObject() *TopWindowObject {
	return this
}

func (this *TopWindowObject) PreCreate(opts *WindowOptions) {
	this.super.PreCreate(opts)
}

func (this *TopWindowObject) PostCreate(opts *WindowOptions) {
	if HWndActive == 0 {
		HWndActive = this.Handle
	}
}

func (this *TopWindowObject) Close() {
	SendMessage(this.Handle, win32.WM_CLOSE, 0, 0)
}

func (this *TopWindowObject) showingModal() bool {
	return this._showingModal
}

func (this *TopWindowObject) GetAccelHandle() win32.HACCEL {
	if this.accelTable == nil {
		return 0
	}
	return this.accelTable.Handle
}

func (this *TopWindowObject) SetAccelerators(accelTable *AcceleratorTable) {
	if accelTable == nil {
		this.accelTable = nil
		if HWndActive == this.Handle {
			hAccelActive = 0
		}
		return
	}
	this.accelTable = accelTable
	if HWndActive == this.Handle {
		hAccelActive = accelTable.Handle
	}
	this.accelTable.OnHandleChange.AddListener(func(info *SimpleEventInfo) {
		if HWndActive == this.Handle {
			hAccelActive = accelTable.Handle
		}
	})
}

func (this *TopWindowObject) GetTitle() string {
	buf := make([]uint16, 255)
	cc, _ := win32.GetWindowText(this.Handle, &buf[0], int32(len(buf)))
	return syscall.UTF16ToString(buf[:cc])
}

func (this *TopWindowObject) SetTitle(title string) {
	win32.SetWindowText(this.Handle, win32.StrToPwstr(title))
}

func (this *TopWindowObject) GetContainer() Container {
	return nil
}

func (this *TopWindowObject) IsChild() bool {
	return false
}

func (this *TopWindowObject) GetRootContainer() Container {
	return this.RealObject.(Container)
}

var _topwindow_class_registerd = false

func (this *TopWindowObject) EnsureClassRegistered() {
	if _topwindow_class_registerd {
		return
	}
	_, err := RegisterClass("goforms.topwindow", nil, ClassOptions{
		BackgroundBrush: ToSysColorBrush(byte(win32.COLOR_WINDOW)),
	})
	if err != nil {
		log.Fatal(err)
	}
	_topwindow_class_registerd = true
}

func (this *TopWindowObject) GetWindowClass() string {
	return "goforms.topwindow"
}

func (this *TopWindowObject) Create(options WindowOptions) error {
	if options.Left == 0 {
		options.Left = int(win32.CW_USEDEFAULT)
	} else if options.Left == consts.Zero {
		options.Left = 0
	}
	if options.Width == 0 {
		options.Width = int(win32.CW_USEDEFAULT)
	} else if options.Width == consts.Zero {
		options.Width = 0
	}
	return this.super.Create(options)
}

func (this *TopWindowObject) CenterOnScreen() {
	var rc win32.RECT
	win32.SystemParametersInfo(win32.SPI_GETWORKAREA, 0, unsafe.Pointer(&rc), 0)
	this.centerToRect(rc)
}

func (this *TopWindowObject) centerToRect(rc win32.RECT) {
	rect := RectFromRECT(rc)
	var rc1 win32.RECT
	win32.GetWindowRect(this.Handle, &rc1)
	thisRect := RectFromRECT(rc1)
	x := max(rect.Left, rect.Left+(rect.Width()-thisRect.Width())/2)
	y := max(rect.Top, rect.Top+(rect.Height()-thisRect.Height())/2)

	win32.SetWindowPos(this.Handle, 0, int32(x), int32(y), 0, 0,
		win32.SWP_NOZORDER|win32.SWP_NOSIZE|win32.SWP_NOACTIVATE)
}

func (this *TopWindowObject) CenterToWindow(hWnd HWND) {
	var rc win32.RECT
	win32.GetWindowRect(hWnd, &rc)
	this.centerToRect(rc)
}

func (this *TopWindowObject) GetModalResult() int {
	data := this.GetData(Data_ModalResult)
	if data != nil {
		return data.(int)
	}
	return 0
}

func (this *TopWindowObject) SetModalResult(result int) {
	this.SetData(Data_ModalResult, result)
}

func (this *TopWindowObject) ShowModal() {

	this.SetData(Data_ModalResult, nil)

	//
	hWndOwner := HWndActive //disable all?
	this._showingModal = true
	if hWndOwner != 0 {
		win32.SetWindowLongPtr(this.Handle, win32.GWLP_HWNDPARENT, uintptr(hWndOwner))
	}
	this.Show()
	if hWndOwner != 0 {
		win32.EnableWindow(hWndOwner, win32.FALSE)
	}
	MessageLoop()
	if hWndOwner != 0 {
		win32.EnableWindow(hWndOwner, win32.TRUE)
		win32.SetActiveWindow(hWndOwner) //?
		win32.DestroyWindow(this.Handle)
	}
	this._showingModal = false
}

func (this *TopWindowObject) PostShow() {
	if this.initialFocusHwnd != 0 {
		win32.SetFocus(this.initialFocusHwnd)
		this.initialFocusHwnd = 0
	}
}

func (this *TopWindowObject) SetMenuHelpUI(menuHelpUI MenuHelpUI) {
	this.menuHelpUI = menuHelpUI
}

func (this *TopWindowObject) SetMenuBar(menuBar *MenuBar) {
	this.menuBar = menuBar
	menuBar.BindTo(this)
}

func (this *TopWindowObject) onEnterMenuLoop(wParam WPARAM, lParam LPARAM) {
	if this.menuHelpUI != nil {
		this.menuHelpUI.EnterMenuHelp()
	}
}

func (this *TopWindowObject) onExitMenuLoop(wParam WPARAM, lParam LPARAM) {
	if this.menuHelpUI != nil {
		this.menuHelpUI.ExitMenuHelp()
	}
}

func (this *TopWindowObject) onMenuSelect(wParam WPARAM, lParam LPARAM) {
	if this.menuHelpUI == nil {
		return
	}

	hw := win32.HIWORD(win32.DWORD(wParam))
	lw := win32.LOWORD(win32.DWORD(wParam))
	bPopup := win32.MENU_ITEM_FLAGS(hw)&win32.MF_POPUP == win32.MF_POPUP

	hMenu := lParam
	if hMenu == 0 {
		//?
		return
	}

	var mii win32.MENUITEMINFO
	mii.CbSize = uint32(unsafe.Sizeof(mii))
	mii.FMask = win32.MIIM_DATA
	var byPos win32.BOOL
	if bPopup {
		byPos = 1
	}
	bOk, errno := win32.GetMenuItemInfo(hMenu, uint32(lw), byPos, &mii)
	if bOk == win32.FALSE {
		log.Println(errno)
	}
	hRootMenu := mii.DwItemData

	var item *MenuItem
	menu := GetMenuObject(hRootMenu)
	if menu == nil {
		//?
		return
	}
	if bPopup {
		hSubMenu := win32.GetSubMenu(win32.HMENU(lParam), int32(lw))
		item = menu.GetSubMenuItem(hSubMenu)
	} else {
		item = menu.GetItem(lw)
	}
	if item != nil {
		helpText := item.Desc
		command := item.Command
		if command != nil {
			helpText = command.Desc
		}
		this.menuHelpUI.ShowMenuHelp(helpText)
	}
}

func (this *TopWindowObject) OnMenuCommand(wParam WPARAM, lParam LPARAM) {
	if this.menuBar == nil {
		return
	}
	if this.menuBar.onMenuCommand(wParam, lParam) {
		return
	}
}

func asTopWindowObject(winObj *WindowObject) *TopWindowObject {
	return (*TopWindowObject)(unsafe.Pointer(winObj))
}

func (this *TopWindowObject) Dispose() {
	if this.menuBar != nil {
		this.menuBar.Dispose()
	}
	this.super.Dispose()
}

func (this *TopWindowObject) OnCommand(msg *CommandMessage) {
	if msg.FromAccelerator() {
		if this.accelTable != nil { //?
			this.accelTable.onReflectCommand(msg)
		}
	}
	this.super.OnCommand(msg)
}

func (this *TopWindowObject) SetClientSize(width, height int) {
	cx, cy := width, height
	rc := win32.RECT{Left: 0, Top: 0, Right: int32(cx), Bottom: int32(cy)}
	win32.AdjustWindowRect(&rc, WINDOW_STYLE(this.GetStyle()), win32.FALSE)
	cx, cy = int(rc.Right-rc.Left), int(rc.Bottom-rc.Top)
	win32.SetWindowPos(this.Handle, 0, 0, 0, int32(cx), int32(cy),
		win32.SWP_NOZORDER|win32.SWP_NOMOVE|win32.SWP_NOACTIVATE)
}

func (this *TopWindowObject) SetDluClientSize(width, height int) {
	width, height = this.DluToPx(width, height)
	this.SetClientSize(width, height)
}

func (this *TopWindowObject) OnSize(width, height int) {
	this.super.OnSize(width, height)
	if this.Layout != nil {
		this.Layout.SetBounds(0, 0, width, height)
	}
}

func (this *TopWindowObject) WinProc(winObj *WindowObject, m *Message) error {
	win := winObj.RealObject.(TopWindowInterface)
	//
	switch m.UMsg {
	case win32.WM_ACTIVATE:
		if m.WParam == 0 {
			if HWndActive == this.Handle {
				hAccelActive = 0
			}
			win.AsTopWindowObject().lastFocusedHwnd = win32.GetFocus()
		} else {
			HWndActive = this.Handle
			win := GetWindow(this.Handle)
			if win != nil {
				topWin := win.(TopWindow)
				hAccelActive = topWin.GetAccelHandle()
			} else {
				println("??")
			}
		}

	case win32.WM_MENUCOMMAND:
		win.OnMenuCommand(m.WParam, m.LParam)

	case win32.WM_MENUSELECT:
		asTopWindowObject(winObj).onMenuSelect(m.WParam, m.LParam)

	case win32.WM_ENTERMENULOOP:
		asTopWindowObject(winObj).onEnterMenuLoop(m.WParam, m.LParam)

	case win32.WM_EXITMENULOOP:
		asTopWindowObject(winObj).onExitMenuLoop(m.WParam, m.LParam)

	case win32.WM_SETFOCUS:
		topWinObj := win.AsTopWindowObject()
		hWndToFocus := topWinObj.lastFocusedHwnd
		if !IsFocusable(hWndToFocus) {
			hWndToFocus = FindFirstFocusable(topWinObj.Handle)
		}
		if hWndToFocus != 0 {
			win32.SetFocus(hWndToFocus)
		}

	case win32.WM_SIZE:
		width, height := win32.LOWORD(uint32(m.LParam)), win32.HIWORD(uint32(m.LParam))
		win.OnSize(int(width), int(height))
		return m.SetHandledWithResult(0)
	case win32.WM_CLOSE:
		if win.showingModal() {
			win32.PostQuitMessage('E' + 'N' + 'D' + 'M' + 'O' + 'D' + 'A' + 'L') //??
		} else {
			win32.DestroyWindow(win.GetHandle())
		}
		return m.SetHandledWithResult(0)
	case win32.WM_DESTROY:
		retVal := this.super.WinProc(winObj, m)
		topWindows := GetTopWindows()
		if len(topWindows) == 1 {
			win32.PostQuitMessage(0) //?
		}
		return retVal
	case win32.WM_NCDESTROY:
		//?
	case win32.WM_MOVE:
		if event, ok := win.TryGetEvent(win32.WM_MOVE); ok {
			event.Fire(win, m)
		}
		m.Handled = true
		return nil
	case win32.DM_GETDEFID:
		id := win.AsTopWindowObject().defId
		if id == 0 {
			m.Handled = true
			return nil
		} else {
			result := win32.MAKELONG(uint16(id), uint16(win32.DC_HASDEFID))
			m.Handled = true
			m.Result = uintptr(result)
			return nil
		}
	case win32.DM_SETDEFID:
		win.AsTopWindowObject().defId = int(m.WParam)
	//
	case win32.WM_SETCURSOR:
		if curWaitCursor != nil {
			curWaitCursor.Update()
			m.Result = win32.LRESULT(win32.TRUE)
			m.Handled = true
			return nil
		}
	}

	return this.super.WinProc(winObj, m)
}
