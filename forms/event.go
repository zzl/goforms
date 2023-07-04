package forms

import (
	"github.com/zzl/goforms/framework/events"
	"unsafe"

	"github.com/zzl/go-win32api/v2/win32"
)

type Event[T events.EventInfo] struct {
	events.Event[T]
}

type MsgEvent = Event[*Message]

type MsgEventSource interface {
	GetEvent(uMsg uint32) *MsgEvent
	TryGetEvent(uMsg uint32) (*MsgEvent, bool)
}

type KeyMessage Message

func (this *KeyMessage) GetKey() int {
	return int(this.WParam)
}

type CharMessage Message

func (this *CharMessage) GetChar() int {
	return int(this.WParam)
}

type CommandMessage Message

func (this *CommandMessage) GetNotifyCode() uint16 {
	return win32.HIWORD(win32.DWORD(this.WParam))
}

func (this *CommandMessage) GetCmdId() uint16 {
	return win32.LOWORD(win32.DWORD(this.WParam))
}

func (this *CommandMessage) GetHwndCtrl() HWND {
	return HWND(this.LParam)
}

func (this *CommandMessage) FromAccelerator() bool {
	return this.GetNotifyCode() == 1
}

func (this *CommandMessage) FromMenu() bool {
	return this.GetNotifyCode() == 0
}

type ContextMenuMessage Message

func (this *ContextMenuMessage) GetClickedHwnd() HWND {
	return HWND(this.WParam)
}

func (this *ContextMenuMessage) GetScreenPos() (int, int) {
	return int(win32.LOWORD(win32.DWORD(this.LParam))),
		int(win32.HIWORD(win32.DWORD(this.LParam)))
}

type NotifyMessage struct {
	*Message
}

func (this *NotifyMessage) GetNMHDR() *win32.NMHDR {
	return (*win32.NMHDR)(unsafe.Pointer(this.LParam))
}

type SimpleEventListener = events.SimpleEventListener

// ?type EventListener = events.EventListener[*Message]
type MsgEventListener func(ei *Message)
