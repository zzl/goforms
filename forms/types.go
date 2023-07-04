package forms

import (
	"github.com/zzl/goforms/framework/events"
	"github.com/zzl/goforms/framework/types"
	"log"
	"syscall"

	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/drawing"
	"github.com/zzl/goforms/drawing/gdi"
	"github.com/zzl/goforms/layouts"
)

type HWND = win32.HWND
type WPARAM = win32.WPARAM
type LPARAM = win32.LPARAM
type WIN32_ERROR = win32.WIN32_ERROR
type WINDOW_STYLE = win32.WINDOW_STYLE
type WINDOW_EX_STYLE = win32.WINDOW_EX_STYLE

type Disposable = types.Disposable
type Rect = types.Rect
type Point = types.Point
type Size = types.Size
type BoundsAware = types.BoundsAware
type Layout = layouts.Layout
type SimpleEvent = events.SimpleEvent
type SimpleEventInfo = events.SimpleEventInfo
type ExtraEvent = events.ExtraEvent
type ExtraEventInfo = events.ExtraEventInfo

func PointFromDWORD(dw win32.DWORD) Point {
	return Point{X: int(win32.LOWORD(dw)), Y: int(win32.HIWORD(dw))}
}

type Message struct {
	HWnd    HWND
	UMsg    uint32
	WParam  WPARAM
	LParam  LPARAM
	Handled bool
	Result  win32.LRESULT
}

func NewMessageFromMSG(msg *win32.MSG) *Message {
	return &Message{
		HWnd:   msg.Hwnd,
		UMsg:   msg.Message,
		WParam: msg.WParam,
		LParam: msg.LParam,
	}
}

/*
func (this *Message) SetHandled(result win32.LRESULT) error {
	this.Handled = true
	this.Result = result
	return nil
}*/

func (this *Message) GetSender() any {
	return GetWindow(this.HWnd)
}

func (this *Message) SetSender(sender any) {
	//ignore..
}

func (this *Message) GetHandled() bool {
	return this.Handled
}

func (this *Message) SetHandled(handle bool) {
	this.Handled = handle
}

func (this *Message) GetResult() uintptr {
	return this.Result
}

func (this *Message) SetResult(result uintptr) {
	this.Result = result
}

func (this *Message) SetHandledWithResult(result uintptr) error {
	this.Result = result
	this.Handled = true
	return nil
}

type WinProcFunc func(win *WindowObject, msg *Message) error

// --
type MessageProcessor interface {
	ProcessMessage(message *Message)
}

type MessageProcessors struct {
	processors []MessageProcessor
}

func (this *MessageProcessors) Add(processor MessageProcessor) {
	this.processors = append(this.processors, processor)
}

func (this *MessageProcessors) Remove(processor MessageProcessor) {
	count := len(this.processors)
	for n := 0; n < count; n++ {
		if this.processors[n] == processor {
			this.processors = append(this.processors[:n], this.processors[n+1:]...)
			count -= 1
			break
		}
	}
	if count == 0 {
		this.processors = nil
	}
}

func (this *MessageProcessors) ProcessMsg(message *Message) bool {
	for _, processor := range this.processors {
		processor.ProcessMessage(message)
		if message.Handled {
			return true
		}
	}
	return false
}

type MessageProcessFunc func(message *Message)

//func (me MessageProcessFunc) ProcessMessage(message *Message) {
//	me(message)
//}

type MessageProcessFuncHolder struct {
	MessageProcessFunc
}

func (this *MessageProcessFuncHolder) ProcessMessage(message *Message) {
	this.MessageProcessFunc(message)
}

func MessageProcessorByFunc(processFunc MessageProcessFunc) MessageProcessor {
	return &MessageProcessFuncHolder{processFunc}
}

//type MessageProcessorByFunc struct {
//	MessageProcessFunc
//}
//
//func (this *MessageProcessorByFunc) ProcessMessage(message *Message) {
//	this.MessageProcessFunc(message)
//}

// --
type TextAware interface {
	SetText(text string)
	GetText() string
}

type TitleAware interface {
	GetTitle() string
	SetTitle(title string)
}

type PosAware interface {
	SetPos(x, y int)
	GetPos() (x, y int)
}

type DpiPosAware interface {
	SetDpiPos(x, y int)
	GetDpiPos() (x, y int)
}

type SizeAware interface {
	SetSize(cx, cy int)
	GetSize() (cx, cy int)
}

type DpiSizeAware interface {
	SetDpiSize(cx, cy int)
	GetDpiSize() (cx, cy int)
}

type BackColorAware interface {
	SetBackColor(color drawing.Color)
	GetBackColor() drawing.Color
}

type FontAware interface {
	SetFont(font *gdi.Font)
	GetFont() *gdi.Font
}

type FocusAware interface {
	Focus()
	HasFocus() bool
}

type EnabledAware interface {
	SetEnabled(enabled bool)
	GetEnabled() bool
}

type VisibleAware interface {
	SetVisible(visible bool)
	GetVisible() bool
}

type ForeColorAware interface {
	SetForeColor(color drawing.Color)
	GetForeColor() drawing.Color
}

type ContextMenuAware interface {
	SetContextMenu(menu *PopupMenu)
	GetContextMenu() *PopupMenu
}

type ReadOnlyAware interface {
	SetReadOnly(readOnly bool)
	GetReadOnly() bool
}

//	type TextAwareSupport struct {
//		text string
//	}
//
//	func (this *TextAwareSupport) SetText(text string) {
//		this.text = text
//	}
//
//	func (this *TextAwareSupport) GetText() string {
//		return this.text
//	}
type NameAware = types.NameAware

type NameAwareSupport struct {
	name string
}

func (this *NameAwareSupport) SetName(name string) {
	this.name = name
}

func (this *NameAwareSupport) GetName() string {
	return this.name
}

type KeyEventArgs struct {
	Key   win32.VIRTUAL_KEY
	Ctrl  bool
	Alt   bool
	Shift bool
}

func NewKeyEventArgs(wParam WPARAM, lParam LPARAM) KeyEventArgs {
	ret := win32.GetKeyState(int32(win32.VK_CONTROL))
	ctrl := ret < 0
	ret = win32.GetKeyState(int32(win32.VK_MENU))
	alt := ret < 0
	ret = win32.GetKeyState(int32(win32.VK_SHIFT))
	shift := ret < 0
	return KeyEventArgs{
		Key:   win32.VIRTUAL_KEY(wParam),
		Ctrl:  ctrl,
		Alt:   alt,
		Shift: shift,
	}
}

//

type Action func()

type ActionAware interface {
	SetAction(action Action)
	GetAction() Action
}

type ActionAwareSupport struct {
	action Action
}

func (this *ActionAwareSupport) SetAction(action Action) {
	this.action = action
}

func (this *ActionAwareSupport) GetAction() Action {
	return this.action
}

func ChainActions(actions ...Action) Action {
	var nonNilActions []Action
	for _, action := range actions {
		if action != nil {
			nonNilActions = append(nonNilActions, action)
		}
	}
	if nonNilActions == nil {
		return nil
	}
	if len(nonNilActions) == 1 {
		return nonNilActions[0]
	}
	return func() {
		for _, action := range nonNilActions {
			action()
		}
	}
}

//

type KeyStroke struct {
	Ctrl  bool
	Alt   bool
	Shift bool
	Key   byte
}

func isExtendedKey(vk byte) bool {
	switch win32.VIRTUAL_KEY(vk) {
	case win32.VK_INSERT, win32.VK_DELETE, win32.VK_HOME, win32.VK_END,
		win32.VK_NEXT, win32.VK_PRIOR, win32.VK_LEFT,
		win32.VK_RIGHT, win32.VK_UP, win32.VK_DOWN:
		return true
	}
	return false
}

func (me KeyStroke) IsExtended() bool {
	return isExtendedKey(me.Key)
}

func (me KeyStroke) String() string {
	var s string
	if me.Ctrl {
		s += "Ctrl+"
	}
	if me.Shift {
		s += "Shift+"
	}
	if me.Alt {
		s += "Alt+"
	}
	scan := win32.MapVirtualKey(uint32(me.Key), win32.MAPVK_VK_TO_VSC)
	if scan == 0 {
		log.Println("?")
	}
	var lParam int32
	lParam = int32(scan) << 16
	if isExtendedKey(me.Key) {
		lParam |= 0x1000000
	}
	buf := make([]uint16, 33)
	ret, errno := win32.GetKeyNameText(lParam, &buf[0], 32)
	if ret == 0 {
		log.Println(errno)
	}
	keyName := syscall.UTF16ToString(buf)
	s += keyName
	return s
}

type ForeColorAwareSupport struct {
	foreColor drawing.Color
}

func (this *ForeColorAwareSupport) SetForeColor(color drawing.Color) {
	this.foreColor = color
}

func (this *ForeColorAwareSupport) GetForeColor() drawing.Color {
	return this.foreColor
}

//
//type BackColorAwareSupport struct {
//	backColor drawing.Color
//}
//
//func (this *BackColorAwareSupport) SetBackColor(color drawing.Color) {
//	this.backColor = color
//}
//
//func (this *BackColorAwareSupport) GetBackColor() drawing.Color {
//	return this.backColor
//}
