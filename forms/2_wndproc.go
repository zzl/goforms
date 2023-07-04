package forms

import (
	"github.com/zzl/go-win32api/v2/win32"
	"log"
	"syscall"
)

// WndProcFunc is the type of Win32 window procedure function
type WndProcFunc func(hWnd HWND, uMsg uint32, wParam WPARAM, lParam LPARAM) win32.LRESULT

// WinProcProvider wraps the WinProc method
type WinProcProvider interface {
	// WinProc processes messages targeting the window object
	WinProc(win *WindowObject, m *Message) error
}

// WndProc is the default window procedure function in GoForms
func WndProc(hWnd HWND, uMsg uint32,
	wParam WPARAM, lParam LPARAM) win32.LRESULT {

	switch uMsg {
	case win32.WM_SETFOCUS:
		if IsChildWindow(hWnd) {
			hWndParent, _ := win32.GetParent(hWnd)
			SendMessage(hWndParent, WM_CHILD_SETFOCUS, wParam, hWnd)
		}
	case win32.WM_KILLFOCUS:
		if IsChildWindow(hWnd) {
			hWndParent, _ := win32.GetParent(hWnd)
			SendMessage(hWndParent, WM_CHILD_KILLFOCUS, wParam, hWnd)
		}
	}

	// get window object
	win, ok := windowMap[hWnd]
	var winObj *WindowObject
	if ok {
		winObj = win.AsWindowObject()
	} else {
		// return win32.DefWindowProc(hWnd, uMsg, wParam, lParam)
		winObj = creatingWindows[len(creatingWindows)-1]
		if winObj.Handle != 0 {
			log.Panic("something wrong")
		}
		winObj.Handle = hWnd
		windowMap[hWnd] = winObj
	}

	// build message object
	msg := &Message{hWnd, uMsg, wParam, lParam, false, 0}

	if winObj.messageProcessors != nil {
		// dispatching priority 0: messageProcessors
		if winObj.messageProcessors.ProcessMsg(msg) {
			return msg.Result
		}
	}

	// dispatching priority 1: msg event
	// fire msg events
	if winObj.FireMsgEvent(msg) {
		return msg.Result
	}

	// dispatching priority 2: WinProcFunc
	// call custom WinProc function
	if winObj.WinProcFunc != nil {
		_ = winObj.WinProcFunc(winObj, msg)
		if msg.Handled {
			return msg.Result
		}
	}
	winProcOwner, ok := win.(WinProcProvider)
	if ok {
		// dispatching priority 3: type.WinProc
		winProcOwner.WinProc(winObj, msg)
		if msg.Handled {
			return msg.Result
		}
	}
	_ = winObj.CallOriWndProc(msg)
	return msg.Result
}

var wndProcCallback = syscall.NewCallback(WndProc)
