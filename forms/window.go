package forms

import (
	"github.com/zzl/goforms/framework/consts"
	"github.com/zzl/goforms/framework/types"
	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/framework/virtual"
	"log"
	"math"
	"unsafe"

	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/drawing"
	"github.com/zzl/goforms/drawing/colors"
	"github.com/zzl/goforms/drawing/gdi"
)

// Window is an interface that represents a graphical window,
// it exposes the common functionalities for Win32 windows.
type Window interface {
	virtual.Virtual
	types.LifecycleAware
	types.DataAware
	EnabledAware
	VisibleAware
	FocusAware
	FontAware
	BackColorAware
	ContextMenuAware
	BoundsAware
	PosAware
	SizeAware
	MsgEventSource

	// GetHandle returns the win32 handle of the window.
	GetHandle() HWND

	// Show displays the window.
	Show()

	// Hide hides the window.
	Hide()

	// IsVisible checks whether the window is visible
	IsVisible() bool

	// Invalidate adds all the window's client area to the update region
	Invalidate()

	// Update causes an immediate redraw on the update region
	Update()

	// GetPreferredSize returns the preferred size of the window
	// given the maximum width and height.
	GetPreferredSize(maxWidth int, maxHeight int) (int, int)

	// AutoSize adjusts the size of the window to fit its contents,
	// using the maximum size provided if available.
	AutoSize(maxSize ...Size)

	// GetClientSize returns the width and height of the client area of the window.
	GetClientSize() (cx, cy int)

	// GetDpiClientSize returns the width and height of
	// the client area of the window, adjusted for DPI scaling.
	GetDpiClientSize() (cx, cy int)

	// SetDluBounds sets the size and position of the control using dialog units.
	SetDluBounds(left, top, width, height int)

	// GetParentHandle returns the handle of the parent window.
	GetParentHandle() HWND

	// Contains checks if the hWnd is the handle of a window contained in this window.
	// The contained window is either a direct child or a deeper descendant window.
	Contains(hWnd win32.HWND) bool

	// GetParent returns the parent window of this window.
	GetParent() Window

	// GetRootWindow returns the root window that contains the current window.
	GetRootWindow() Window

	// GetStyle returns the window style.
	GetStyle() WINDOW_STYLE

	// GetExStyle returns the extended window style.
	GetExStyle() WINDOW_EX_STYLE

	// ModifyStyle modifies the window style using the provided style flags.
	// A style with value 0 is equivalent to the current style.
	ModifyStyle(style, styleInclude, styleExclude win32.WINDOW_STYLE)

	// IsChild returns true if the window is a child window.
	IsChild() bool

	// Create creates the window with the specified options.
	Create(options WindowOptions) error

	// GetOnCreate returns the event that is triggered when the window is created.
	GetOnCreate() *SimpleEvent

	// Refresh redraws the window.
	Refresh()

	// AddDisposeAction adds an action to be executed when the window is disposed.
	AddDisposeAction(action Action)

	// AsWindowObject returns the underlying WindowObject.
	AsWindowObject() *WindowObject

	// AddMessageProcessor adds a message processor to handle incoming messages.
	AddMessageProcessor(processor MessageProcessor)

	// RemoveMessageProcessor removes a message processor.
	RemoveMessageProcessor(processor MessageProcessor)
}

// WindowSpi is an interface that provides additional methods
// specific to implementing a window.
type WindowSpi interface {

	// GetWindowClass returns the class name of the window.
	GetWindowClass() string

	// GetDefaultStyle returns the default window style.
	GetDefaultStyle() WINDOW_STYLE

	// GetDefaultExStyle returns the default extended window style.
	GetDefaultExStyle() WINDOW_EX_STYLE

	// EnsureClassRegistered ensures that the window class is registered.
	EnsureClassRegistered()

	// EnsureCustomWndProc ensures that the custom window procedure is set up.
	EnsureCustomWndProc()

	// OnHandleCreated is called when the window handle is created.
	OnHandleCreated()

	// CreateHandle creates a win32 window with the specified class name and options.
	CreateHandle(className string, options WindowOptions) (HWND, error)

	// OnReflectMessage processes notification messages of this window
	// that are sent to its parent window,
	OnReflectMessage(msg *Message)

	// OnReflectCommand processes the WM_COMMAND reflection message.
	OnReflectCommand(msg *CommandMessage)

	// OnReflectNotify processes the WM_NOTIFY reflection message.
	OnReflectNotify(msg *NotifyMessage)

	// PreCreate is called before the win32 window is created.
	PreCreate(opts *WindowOptions)

	// PreDispose is called before the win32 window is destroyed.
	PreDispose()

	// PostCreate is called after the win32 window is created.
	PostCreate(opts *WindowOptions)

	// PreShow is called before showing the window.
	PreShow() //?first

	// PostShow is called after the window is shown.
	PostShow()

	// OnSize is called after calling SetSize or SetBounds
	OnSize(width, height int)

	// PreProcessMsg is called before processing a window message.
	// Return true to mark the message as handled.
	PreProcessMsg(msg *win32.MSG) bool

	// CallDefaultWndProc calls the default window procedure for processing a message.
	CallDefaultWndProc(msg *Message) error
}

// WindowInterface is the interface that combines Window and WindowSpi
type WindowInterface interface {
	Window
	WindowSpi
}

// WindowObject implements the WindowInterface.
// It extends VirtualObject.
type WindowObject struct {

	// VirtualObject is the parent struct.
	virtual.VirtualObject[WindowInterface]

	// Handle is the underlying win32 window handle.
	Handle HWND

	// OnCreateEvent is fired after win32 window handle is created.
	OnCreateEvent SimpleEvent

	// WinProcFunc optionally specifies a function that
	// provide a custom WinProc implementation.
	WinProcFunc WinProcFunc

	// oriWndProc stores the original window procedure if employing subclassing
	oriWndProc uintptr

	// sysEventMap is an event registry for win32 messages lower than WM_USER
	sysEventMap []*MsgEvent

	// eventMap is an event registry for win32 messages not lower than WM_USER
	eventMap map[uint32]*MsgEvent

	// data is a map used to attach any data with this window by name
	data map[string]any

	// flags is a bit field used to store boolean states of the window
	flags WindowFlag

	// disposeActions stores actions to execute on dispose
	disposeActions []Action

	// messageProcessors is a MessageProcessors registry associated with this window
	messageProcessors *MessageProcessors
}

type WindowFlag uint32

// Window flags
const (
	FlagPreDisposed WindowFlag = 0x02 // PreDispose has been called
	FlagDisposed    WindowFlag = 0x04 // Dispose has been called
	FlagCollapsed   WindowFlag = 0x08 // collapsed in layout
	FlagDesignMode  WindowFlag = 0x10 // in design mode
)

// NewWindowObject creates a new WindowObject.
func NewWindowObject() *WindowObject {
	return virtual.New[WindowObject]()
}

// WindowObject implements Window.WindowObject
func (this *WindowObject) AsWindowObject() *WindowObject {
	return this
}

// HasFlag checks if the window has the specified flag
func (this *WindowObject) HasFlag(flag WindowFlag) bool {
	return this.flags&flag != 0
}

// SetFlag sets or unsets the specified WindowFlag of the window.
func (this *WindowObject) SetFlag(flag WindowFlag, set bool) {
	if set {
		this.flags |= flag
	} else {
		this.flags &^= flag
	}
}

// AddDisposeAction implements Window.AddDisposeAction
func (this *WindowObject) AddDisposeAction(action Action) {
	this.disposeActions = append(this.disposeActions, action)
}

// PreProcessMsg implements WindowSpi.PreProcessMsg
func (this *WindowObject) PreProcessMsg(msg *win32.MSG) bool {
	if this.messageProcessors != nil {
		message := NewMessageFromMSG(msg)
		for _, processor := range this.messageProcessors.processors {
			processor.ProcessMessage(message)
			if message.Handled {
				return true
			}
		}
	}
	return false
}

// GetParent implements Window.GetParent
func (this *WindowObject) GetParent() Window {
	return GetWindow(this.GetParentHandle())
}

// GetRootWindow implements Window.GetRootWindow
func (this *WindowObject) GetRootWindow() Window {
	return GetWindow(this.GetRootHandle())
}

// GetRootHandle returns the handle of the root ancestor windows of this window
func (this *WindowObject) GetRootHandle() HWND {
	hWnd := win32.GetAncestor(this.Handle, win32.GA_ROOT)
	return hWnd
}

// callPreDispose calls PreDispose if not already called
func (this *WindowObject) callPreDispose() {
	if this.flags&FlagPreDisposed == 0 {
		this.flags |= FlagPreDisposed
		this.RealObject.PreDispose()
	}
}

// callDispose calls Dispose if not already called
func (this *WindowObject) callDispose() {
	if this.flags&FlagDisposed == 0 {
		this.flags |= FlagDisposed
		this.RealObject.Dispose()
	}
}

// Init implements Window.Init
func (this *WindowObject) Init() {
	//
}

// PreDispose implements WindowSpi.PreDispose
func (this *WindowObject) PreDispose() {
	//
}

// Dispose implements Window.Dispose
func (this *WindowObject) Dispose() {
	font := this.GetFont()
	if font != nil {
		font.Dispose()
	}

	if this.Handle == 0 {
		return
	}
	delete(windowMap, this.Handle)
	this.Handle = 0
	for _, action := range this.disposeActions {
		action()
	}
	if data, ok := this.data[Data_Disposables]; ok {
		disposables := data.([]Disposable)
		for _, disposable := range disposables {
			disposable.Dispose()
		}
	}
}

// OnParentResized implements WindowSpi.OnParentResized
func (this *WindowObject) OnParentResized() {
	//nop
}

// OnReflectMessage implements WindowSpi.OnReflectMessage
func (this *WindowObject) OnReflectMessage(msg *Message) {
	if msg.UMsg == win32.WM_COMMAND {
		this.RealObject.OnReflectCommand((*CommandMessage)(msg))
	} else if msg.UMsg == win32.WM_NOTIFY {
		this.RealObject.OnReflectNotify(&NotifyMessage{msg})
	}
}

// OnReflectCommand implements WindowSpi.OnReflectCommand
func (this *WindowObject) OnReflectCommand(info *CommandMessage) {
	//
}

// OnReflectNotify implements WindowSpi.OnReflectNotify
func (this *WindowObject) OnReflectNotify(info *NotifyMessage) {
	//
}

// OnHandleCreated implements WindowSpi.OnHandleCreated
func (this *WindowObject) OnHandleCreated() {
	//
}

// GetData implements Window.GetData
func (this *WindowObject) GetData(key string) any {
	if this.data == nil {
		return nil
	}
	return this.data[key]
}

// SetData implements Window.SetData
func (this *WindowObject) SetData(key string, value any) {
	if value == nil {
		if this.data != nil {
			delete(this.data, key)
		}
	} else {
		if this.data == nil {
			this.data = make(map[string]any)
		}
		this.data[key] = value
	}
}

// TryGetEvent implements Window.TryGetEvent
func (this *WindowObject) TryGetEvent(uMsg uint32) (*MsgEvent, bool) {
	if uMsg < win32.WM_USER {
		if this.sysEventMap == nil {
			return nil, false
		}
		e := this.sysEventMap[uMsg]
		return e, e != nil
	}
	event, ok := this.eventMap[uMsg]
	if !ok {
		return nil, false
	}
	return event, true
}

// GetEvent implements Window.GetEvent
func (this *WindowObject) GetEvent(uMsg uint32) *MsgEvent {
	if uMsg < win32.WM_USER {
		if this.sysEventMap == nil {
			this.sysEventMap = make([]*MsgEvent, win32.WM_USER)
			if this.Handle != 0 {
				this.RealObject.(WindowSpi).EnsureCustomWndProc()
			}
		}
		e := this.sysEventMap[uMsg]
		if e == nil {
			e = &MsgEvent{}
			this.sysEventMap[uMsg] = e
		}
		return e
	}
	e, ok := this.eventMap[uMsg]
	if !ok {
		if this.eventMap == nil {
			this.eventMap = make(map[uint32]*MsgEvent)
			if this.Handle != 0 {
				this.RealObject.(WindowSpi).EnsureCustomWndProc()
			}
		}
		e = &MsgEvent{}
		this.eventMap[uMsg] = e
	}
	return e
}

// EnsureCustomWndProc implements WindowSpi.EnsureCustomWndProc
func (this *WindowObject) EnsureCustomWndProc() {
	this.ensureSubclassed()
}

// ensureSubclassed subclasses this Win32 window to use the custom WndProc
func (this *WindowObject) ensureSubclassed() {
	if this.oriWndProc != 0 {
		return
	}
	wndProc, errno := win32.SetWindowLongPtr(this.Handle,
		win32.GWLP_WNDPROC, wndProcCallback)
	if wndProc == 0 {
		log.Fatal(errno)
	}
	this.oriWndProc = wndProc
}

// unSubclass restores the Win32 window with the original window procedure
func (this *WindowObject) unSubclass() {
	if this.oriWndProc == 0 {
		return
	}
	win32.SetWindowLongPtr(this.Handle, win32.GWLP_WNDPROC, this.oriWndProc)
	this.oriWndProc = 0
}

// Destroy destroys the Win32 window and disposes this WindowObject
func (this *WindowObject) Destroy() {
	win32.DestroyWindow(this.Handle)
	this.callDispose()
}

// EnsureClassRegistered implements WindowSpi.EnsureClassRegistered
func (this *WindowObject) EnsureClassRegistered() {
	//nop
}

// PreCreate implements WindowSpi.PreCreate
func (this *WindowObject) PreCreate(opts *WindowOptions) {
	if opts.Style == 0 {
		opts.Style = this.RealObject.GetDefaultStyle()
	}
	if opts.ExStyle == 0 {
		opts.ExStyle = this.RealObject.GetDefaultExStyle()
	} else if opts.ExStyle == WINDOW_EX_STYLE(consts.Zero) {
		opts.ExStyle = 0
	}
}

// PostCreate implements WindowSpi.PostCreate
func (this *WindowObject) PostCreate(opts *WindowOptions) {
	//
}

// GetWindowClass implements WindowSpi.GetWindowClass
func (this *WindowObject) GetWindowClass() string {
	return ""
}

// GetDefaultStyle implements WindowSpi.GetDefaultStyle
func (this *WindowObject) GetDefaultStyle() WINDOW_STYLE {
	return 0
}

// GetDefaultExStyle implements WindowSpi.GetDefaultExStyle
func (this *WindowObject) GetDefaultExStyle() WINDOW_EX_STYLE {
	return 0
}

// GetStyle implements Window.GetStyle
func (this *WindowObject) GetStyle() WINDOW_STYLE {
	if this.Handle == 0 {
		return this.RealObject.(WindowSpi).GetDefaultStyle()
	}
	ret, _ := win32.GetWindowLongPtr(this.Handle, win32.GWL_STYLE)
	return WINDOW_STYLE(ret)
}

// GetExStyle implements Window.GetExStyle
func (this *WindowObject) GetExStyle() WINDOW_EX_STYLE {
	if this.Handle == 0 {
		return this.RealObject.(WindowSpi).GetDefaultExStyle()
	}
	ret, _ := win32.GetWindowLong(this.Handle, win32.GWL_EXSTYLE)
	return WINDOW_EX_STYLE(ret)
}

// ModifyStyle implements Window.ModifyStyle
func (this *WindowObject) ModifyStyle(style, styleInclude, styleExclude win32.WINDOW_STYLE) {
	if style == 0 {
		style = this.GetStyle()
	}
	style |= styleInclude
	style &^= styleExclude
	win32.SetWindowLong(this.Handle, win32.GWL_STYLE, int32(style))
}

// GetOnCreate implements Window.GetOnCreate
func (this *WindowObject) GetOnCreate() *SimpleEvent {
	return &this.OnCreateEvent
}

// CreateHandle implements WindowSpi.CreateHandle
func (this *WindowObject) CreateHandle(
	className string, options WindowOptions) (HWND, error) {
	return CreateWindow(className, options, nil)
}

// Create implements Window.Create
func (this *WindowObject) Create(options WindowOptions) error {
	win := this.RealObject
	if win == nil {
		log.Fatal("Object not realized")
	}
	win.EnsureClassRegistered()
	win.PreCreate(&options)

	//
	className := options.ClassName
	if className == "" {
		className = win.GetWindowClass()
	}

	creatingWindows = append(creatingWindows, this)
	hWnd, err := win.CreateHandle(className, options)
	creatingWindows = creatingWindows[:len(creatingWindows)-1]

	if err != nil {
		return err
	}
	this.Handle = hWnd
	windowMap[hWnd] = this.RealObject

	if this.WinProcFunc != nil {
		win.EnsureCustomWndProc()
	} else if _, ok := win.(WinProcProvider); ok {
		win.EnsureCustomWndProc()
	}

	if this.sysEventMap != nil || this.eventMap != nil {
		win.EnsureCustomWndProc()
	}

	var hFont win32.HFONT
	font := win.GetFont()

	if font == nil && options.ParentHandle != 0 {
		if parentWin, ok := windowMap[options.ParentHandle]; ok {
			font = parentWin.GetFont()
		}
	}
	if font != nil {
		font.EnsureCreated()
		hFont = font.Handle
	}
	if hFont == 0 {
		hFont = GetDefaultFont()
	}
	SendMessage(hWnd, win32.WM_SETFONT, hFont, 0)
	this.SetData(Data_FontHandle, hFont)

	win.PostCreate(&options)
	win.OnHandleCreated()
	this.OnCreateEvent.Fire(this, &SimpleEventInfo{})
	return nil
}

// Attach attaches this WindowObject to an existing window handle
func (this *WindowObject) Attach(hWnd HWND) error {
	if this.Handle != 0 {
		log.Panic("??")
	}
	this.Handle = hWnd
	windowMap[hWnd] = this
	this.ensureSubclassed()
	return nil
}

// Detach detaches this WindowObject from its window handle
func (this *WindowObject) Detach() error {
	if this.Handle == 0 {
		log.Println("??")
		return nil
	}
	if this.oriWndProc == 0 {
		log.Panic("??")
	}
	this.unSubclass()
	delete(windowMap, this.Handle)
	this.Handle = 0
	return nil
}

// IsVisible implements Window.IsVisible
func (this *WindowObject) IsVisible() bool {
	ret := win32.IsWindowVisible(this.Handle)
	return ret != win32.FALSE
}

// Hide implements Window.Hide
func (this *WindowObject) Hide() {
	win32.ShowWindow(this.Handle, win32.SW_HIDE)
}

// Show implements Window.Show
func (this *WindowObject) Show() {
	this.RealObject.PreShow()
	win32.ShowWindow(this.Handle, win32.SW_SHOW)
	this.RealObject.PostShow()
}

// PreShow implements WindowSpi.PreShow
func (this *WindowObject) PreShow() {
	//nop
}

// PostShow implements WindowSpi.PostShow
func (this *WindowObject) PostShow() {
	//nop
}

// OnSize implements WindowSpi.OnSize
func (this *WindowObject) OnSize(width, height int) {
	//nop
}

// Invalidate implements Window.Invalidate
func (this *WindowObject) Invalidate() {
	win32.InvalidateRect(this.Handle, nil, win32.TRUE)
}

// Update implements Window.Update
func (this *WindowObject) Update() {
	win32.UpdateWindow(this.Handle)
}

// SetBounds implements Window.SetBounds
func (this *WindowObject) SetBounds(left, top, width, height int) {
	if width == 1024 && height == 0 {
		if (this.flags&FlagCollapsed) == 0 && this.IsVisible() {
			this.flags |= FlagCollapsed
			this.Hide()
		}
		return
	} else if (this.flags & FlagCollapsed) != 0 {
		this.flags &^= FlagCollapsed
		this.Show()
	}
	//
	if width == 0 || height == 0 {
		cx, cy := this.RealObject.GetPreferredSize(0, 0)
		if width == 0 {
			width = cx
		} else if width == consts.Zero {
			width = 0
		}
		if height == 0 {
			height = cy
		} else if height == consts.Zero {
			height = 0
		}
	}
	win32.MoveWindow(this.Handle, int32(left), int32(top),
		int32(width), int32(height), win32.TRUE)
	this.RealObject.OnSize(width, height)
}

// CallOriWndProc calls the original or default window procedure
func (this *WindowObject) CallOriWndProc(msg *Message) error {
	if this.oriWndProc != 0 { //subclassed
		ret := win32.CallWindowProc(this.oriWndProc,
			this.Handle, msg.UMsg, msg.WParam, msg.LParam)
		msg.Result = ret
		msg.Handled = true
		return nil
	} else {
		return this.RealObject.CallDefaultWndProc(msg)
	}
}

// CallDefaultWndProc implements WindowSpi.CallDefaultWndProc
func (this *WindowObject) CallDefaultWndProc(msg *Message) error {
	msg.Result = win32.DefWindowProc(msg.HWnd, msg.UMsg, msg.WParam, msg.LParam)
	return nil
}

// SetFont implements Window.SetFont
func (this *WindowObject) SetFont(font *gdi.Font) {
	oriFont := this.GetFont()
	if oriFont != nil {
		oriFont.Dispose()
	}
	if font == nil {
		this.SetData(Data_Font, nil)
		return
	}
	this.SetData(Data_Font, font)
	if this.Handle == 0 {
		return
	}
	font.EnsureCreated()
	SendMessage(this.Handle, win32.WM_SETFONT, font.Handle, 0)
	this.SetData(Data_FontHandle, font.Handle)
}

// GetFont implements Window.GetFont
func (this *WindowObject) GetFont() *gdi.Font {
	value := this.GetData(Data_Font)
	if value == nil {
		return nil
	}
	return value.(*gdi.Font)
}

// GetParentHandle implements Window.GetParentHandle
func (this *WindowObject) GetParentHandle() HWND {
	return win32.GetAncestor(this.Handle, win32.GA_PARENT)
}

// Focus implements Window.FocusFocus
func (this *WindowObject) Focus() {
	_, _ = win32.SetFocus(this.Handle)
}

// HasFocus implements Window.HasFocus
func (this *WindowObject) HasFocus() bool {
	hWndFocus := win32.GetFocus()
	return hWndFocus == this.Handle
}

// Contains implements Window.Contains
func (this *WindowObject) Contains(hWnd win32.HWND) bool {
	ret := win32.IsChild(this.Handle, hWnd)
	return ret == win32.TRUE
}

// ContainsFocus checks whether this window has focus
// or contains the focused window handle
func (this *WindowObject) ContainsFocus() bool {
	hWnd := win32.GetFocus()
	return this.Contains(hWnd)
}

// SetEnabled implements Window.SetEnabled
func (this *WindowObject) SetEnabled(enabled bool) {
	_ = win32.EnableWindow(this.Handle, win32.BoolToBOOL(enabled))
}

// GetEnabled implements Window.GetEnabled
func (this *WindowObject) GetEnabled() bool {
	result := win32.IsWindowEnabled(this.Handle)
	return result != win32.FALSE
}

// GetHandle implements Window.GetHandle
func (this *WindowObject) GetHandle() HWND {
	return this.Handle
}

// SetPos implements Window.SetPos
func (this *WindowObject) SetPos(x int, y int) {
	win32.SetWindowPos(this.Handle, 0, int32(x), int32(y), 0, 0,
		win32.SWP_NOZORDER|win32.SWP_NOSIZE|win32.SWP_NOACTIVATE)
}

// GetPos implements Window.GetPos
func (this *WindowObject) GetPos() (x, y int) {
	var rc win32.RECT
	win32.GetWindowRect(this.Handle, &rc)
	pt := win32.POINT{X: rc.Left, Y: rc.Top}
	win32.ScreenToClient(this.GetParentHandle(), &pt)
	return int(pt.X), int(pt.Y)
}

// SetDpiPos implements Window.SetDpiPos
func (this *WindowObject) SetDpiPos(x int, y int) {
	this.RealObject.SetPos(DpiScale(x), DpiScale(y))
}

// GetDpiPos implements Window.GetDpiPos
func (this *WindowObject) GetDpiPos() (x, y int) {
	x, y = this.RealObject.GetPos()
	x, y = DpiUnscale(x), DpiUnscale(y)
	return
}

// SetSize implements Window.SetSize
func (this *WindowObject) SetSize(cx int, cy int) {
	win32.SetWindowPos(this.Handle, 0, 0, 0, int32(cx), int32(cy),
		win32.SWP_NOZORDER|win32.SWP_NOMOVE|win32.SWP_NOACTIVATE)
	this.RealObject.OnSize(cx, cy)
}

// GetSize implements Window.GetSize
func (this *WindowObject) GetSize() (cx, cy int) {
	var rc win32.RECT
	win32.GetWindowRect(this.Handle, &rc)
	return int(rc.Right - rc.Left), int(rc.Bottom - rc.Top)
}

// SetDpiSize implements Window.SetDpiSize
func (this *WindowObject) SetDpiSize(cx int, cy int) {
	cx, cy = DpiScale(cx), DpiScale(cy)
	this.RealObject.SetSize(cx, cy)
}

// GetDpiSize implements Window.GetDpiSize
func (this *WindowObject) GetDpiSize() (cx, cy int) {
	cx, cy = this.RealObject.GetSize()
	cx, cy = DpiUnscale(cx), DpiUnscale(cy)
	return
}

// GetClientSize implements Window.GetClientSize
func (this *WindowObject) GetClientSize() (cx, cy int) {
	var rc win32.RECT
	win32.GetClientRect(this.Handle, &rc)
	return int(rc.Right), int(rc.Bottom)
}

// GetDpiClientSize implements Window.GetDpiClientSize
func (this *WindowObject) GetDpiClientSize() (cx, cy int) {
	cx, cy = this.RealObject.GetClientSize()
	cx, cy = DpiUnscale(cx), DpiUnscale(cy)
	return
}

// DluToPx converts x, y in dialog units to pixel values
func (this *WindowObject) DluToPx(x, y int) (int, int) {
	var hFont win32.HFONT
	value := this.GetData(Data_FontHandle)
	if value != nil {
		hFont = value.(win32.HFONT)
	}
	if hFont == 0 {
		hFont = GetDefaultFont()
	}
	xDbu, yDbu := gdi.MeasureDbus(hFont)
	x = (int)(math.Round(float64(x) * float64(xDbu) / 4))
	y = (int)(math.Round(float64(y) * float64(yDbu) / 8))
	return x, y
}

// PxToDlu converts x, y in pixels to dialog units
func (this *WindowObject) PxToDlu(x, y int) (int, int) {
	var hFont win32.HFONT
	value := this.GetData(Data_FontHandle)
	if value != nil {
		hFont = value.(win32.HFONT)
	}
	if hFont == 0 {
		hFont = GetDefaultFont()
	}
	xDbu, yDbu := gdi.MeasureDbus(hFont)

	x = (int)(math.Round(float64(x) * 4 / float64(xDbu)))
	y = (int)(math.Round(float64(y) * 8 / float64(yDbu)))
	return x, y
}

func (this *WindowObject) GetClientSize32() (cx, cy int32) {
	var rc win32.RECT
	win32.GetClientRect(this.Handle, &rc)
	return rc.Right, rc.Bottom
}

// GetBounds implements Window.GetBounds
func (this *WindowObject) GetBounds() Rect {
	var rc win32.RECT
	win32.GetWindowRect(this.Handle, &rc)
	hWndParent, _ := win32.GetParent(this.Handle)
	win32.MapWindowPoints(0, hWndParent, (*win32.POINT)(unsafe.Pointer(&rc)), 2)
	return RectFromRECT(rc)
}

// GetPreferredSize implements Window.GetPreferredSize
func (this *WindowObject) GetPreferredSize(cxMax int, cyMax int) (int, int) {
	return 0, 0
}

// AutoSize implements Window.AutoSize
func (this *WindowObject) AutoSize(optMaxSize ...Size) {
	maxSize := utils.OptionalArgByVal(optMaxSize)
	cx, cy := this.RealObject.GetPreferredSize(maxSize.Width, maxSize.Height)
	this.RealObject.SetSize(cx, cy)
}

// Refresh implements Window.Refresh
func (this *WindowObject) Refresh() {
	this.Invalidate()
	this.Update()
}

// FireMsgEvent fires the event for the specified window message
// and returns whether it's been handled
func (this *WindowObject) FireMsgEvent(msg *Message) bool {
	event, ok := this.TryGetEvent(msg.UMsg)
	if ok {
		event.Fire(this, msg)
		if msg.Handled {
			return true
		}
	}
	return false
}

// IsChild implements Window.IsChild
func (this *WindowObject) IsChild() bool {
	return this.GetStyle()&win32.WS_CHILD != 0
}

// SetVisible implements Window.SetVisible
func (this *WindowObject) SetVisible(visible bool) {
	if this.HasFlag(FlagDesignMode) {
		return //?
	}
	if this.IsChild() {
		var swp = win32.SWP_HIDEWINDOW
		if visible {
			swp = win32.SWP_NOSIZE | win32.SWP_NOMOVE |
				win32.SWP_NOZORDER | win32.SWP_NOACTIVATE | win32.SWP_SHOWWINDOW
		}
		win32.SetWindowPos(this.Handle, 0, 0, 0, 0, 0, swp)
	} else {
		sw := win32.SW_HIDE
		if visible {
			sw = win32.SW_SHOW
		}
		win32.ShowWindow(this.Handle, sw)
	}
}

// GetVisible implements Window.GetVisible
func (this *WindowObject) GetVisible() bool {
	ret := win32.IsWindowVisible(this.Handle)
	return ret != win32.FALSE
}

// SetBackColor implements Window.SetBackColor
func (this *WindowObject) SetBackColor(color drawing.Color) {
	this.SetData(Data_BackColor, color)
}

// GetBackColor implements Window.GetBackColor
func (this *WindowObject) GetBackColor() drawing.Color {
	if obj, ok := this.data[Data_BackColor]; ok {
		return obj.(drawing.Color)
	}
	return colors.Null
}

// SetContextMenu implements Window.SetContextMenu
func (this *WindowObject) SetContextMenu(menu *PopupMenu) {
	this.SetData(Data_ContextMenu, menu)
}

// GetContextMenu implements Window.GetContextMenu
func (this *WindowObject) GetContextMenu() *PopupMenu {
	value := this.GetData(Data_ContextMenu)
	if value == nil {
		return nil
	}
	return value.(*PopupMenu)
}

// AddMessageProcessor implements Window.AddMessageProcessor
func (this *WindowObject) AddMessageProcessor(processor MessageProcessor) {
	if this.messageProcessors == nil {
		this.messageProcessors = &MessageProcessors{}
	}
	this.messageProcessors.Add(processor)
}

// RemoveMessageProcessor implements Window.RemoveMessageProcessor
func (this *WindowObject) RemoveMessageProcessor(processor MessageProcessor) {
	if this.messageProcessors == nil {
		return
	}
	this.messageProcessors.Remove(processor)
}

// SetDluBounds implements Window.SetDluBounds.
func (this *WindowObject) SetDluBounds(left, top, width, height int) {
	left, top = this.DluToPx(left, top)
	width, height = this.DluToPx(width, height)
	this.RealObject.SetBounds(left, top, width, height)
}
