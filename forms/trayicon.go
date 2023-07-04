package forms

import (
	"crypto/md5"
	"log"
	"os"
	"syscall"
	"unsafe"

	"github.com/zzl/go-win32api/v2/win32"
)

type TrayIcon struct {
	Icon        win32.HICON
	Tooltip     string
	ContextMenu *PopupMenu
	trayWin     *WindowObject
	uid         uint32
}

func NewTrayIcon(icon win32.HICON, tooltip string) *TrayIcon {
	return &TrayIcon{Icon: icon, Tooltip: tooltip}
}

const trayWinClass = "goforms.traywin"

var trayWinClassRegistered = false

const trayCallbackMsg = win32.WM_APP + 1

var nextTrayId uint32 = 1

func trayWinClassProc(hWnd HWND, uMsg uint32, wParam WPARAM, lParam LPARAM) win32.LRESULT {
	if uMsg == trayCallbackMsg {
		win := windowMap[hWnd]
		trayIcon := win.GetData("TrayIcon").(*TrayIcon)
		trayIcon.handleTrayNotify(wParam, lParam)
		return 0
	}
	return win32.DefWindowProc(hWnd, uMsg, wParam, lParam)
}

func ensureTrayWinClassRegistered() {
	if trayWinClassRegistered {
		return
	}
	_, err := RegisterClass(trayWinClass, trayWinClassProc, ClassOptions{
		BackgroundBrush: 0,
	})
	if err != nil {
		log.Fatal(err)
	}
	trayWinClassRegistered = true
}

func (this *TrayIcon) Create(visible bool) {
	ensureTrayWinClassRegistered()

	trayWin := NewWindowObject()
	err := trayWin.Create(WindowOptions{
		ClassName:    trayWinClass,
		Style:        win32.WS_POPUP,
		ParentHandle: win32.HWND_MESSAGE,
	})

	if err != nil {
		log.Fatal(err)
	}
	trayWin.data = map[string]any{
		"TrayIcon": this,
	}
	this.trayWin = trayWin

	var nid win32.NOTIFYICONDATA
	*nid.UVersion() = 4 // UVersion
	//nid.CbSize = 956
	nid.CbSize = win32.DWORD(unsafe.Sizeof(nid))
	nid.UFlags = win32.NIF_ICON | win32.NIF_TIP | win32.NIF_MESSAGE | win32.NIF_STATE | win32.NIF_GUID
	nid.HIcon = this.Icon
	if !visible {
		nid.DwStateMask |= uint32(win32.NIS_HIDDEN)
		nid.DwState |= win32.NIS_HIDDEN
	}

	wszTip, _ := syscall.UTF16FromString(this.Tooltip)
	copy(nid.SzTip[:], wszTip)

	this.uid = nextTrayId
	nextTrayId += 1

	nid.GuidItem = this.getGuid()
	nid.HWnd = trayWin.Handle
	nid.UCallbackMessage = trayCallbackMsg

	bOk := win32.Shell_NotifyIcon(win32.NIM_ADD, &nid)
	if bOk == win32.FALSE {
		log.Fatal("?")
	}
}

func (this *TrayIcon) Show() {
	//todo..
}

func (this *TrayIcon) Hide() {
	//todo..
}

func (this *TrayIcon) Remove() {
	var nid win32.NOTIFYICONDATA
	*nid.UVersion() = 4 //UVersion
	nid.CbSize = win32.DWORD(unsafe.Sizeof(nid))
	nid.UFlags = win32.NIF_GUID
	nid.GuidItem = this.getGuid()
	nid.HWnd = this.trayWin.Handle
	ok := win32.Shell_NotifyIcon(win32.NIM_DELETE, &nid)
	if ok == win32.FALSE {
		log.Println("remove tray icon failed.. ")
	}
}

func (this *TrayIcon) handleTrayNotify(wParam WPARAM, lParam LPARAM) {
	ninMsg := win32.LOWORD(win32.DWORD(lParam))
	switch uint32(ninMsg) {
	case win32.WM_RBUTTONUP:
		this.showTrayMenu()
		break
	}
}

func (this *TrayIcon) showTrayMenu() {
	var pt win32.POINT
	_, _ = win32.GetCursorPos(&pt)
	_ = win32.SetForegroundWindow(this.trayWin.Handle)
	this.ContextMenu.Show(int(pt.X), int(pt.Y), this.trayWin.Handle)
}

func (this *TrayIcon) getGuid() syscall.GUID {
	exe, _ := os.Executable()
	md5Bytes := md5.Sum([]byte(exe))
	guid := *(*syscall.GUID)(unsafe.Pointer(&md5Bytes[0]))
	return guid
}
