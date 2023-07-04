package forms

import (
	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/framework/virtual"
	"log"
	"syscall"
	"unsafe"
)

type Dialog interface {
	Container
}

type DialogObject struct {
	ContainerObject
}

func NewDialogObject() *DialogObject {
	return virtual.New[DialogObject]()
}

func dialogProc(hWnd HWND, uMsg uint32, wParam WPARAM, lParam LPARAM) win32.LRESULT {
	switch uMsg {
	case win32.WM_INITDIALOG:
		win32.MoveWindow(hWnd, 800, 200, 600, 400, win32.FALSE)
		break
	case win32.WM_CLOSE:
		win32.EndDialog(hWnd, uintptr(win32.IDCANCEL))
		break
	}
	return 0
}

var dialogProcCallback = syscall.NewCallback(dialogProc)

func (this *DialogObject) Show() {
	hTemplate := this.createTemplate()
	hWnd, errno := win32.CreateDialogIndirectParam(HInstance,
		(*win32.DLGTEMPLATE)(unsafe.Pointer(hTemplate)),
		HWndActive, dialogProcCallback, 0)
	if hWnd == 0 {
		log.Fatal(errno)
	}
	win32.GlobalFree(hTemplate)
	win32.ShowWindow(hWnd, win32.SW_SHOW)
}

func (this *DialogObject) ShowModal() {
	hTemplate := this.createTemplate()
	ret, errno := win32.DialogBoxIndirectParam(HInstance,
		(*win32.DLGTEMPLATE)(unsafe.Pointer(hTemplate)),
		HWndActive, dialogProcCallback, 0)
	if ret == NegativeOne {
		log.Fatal(errno)
	}
	win32.GlobalFree(hTemplate)
}

func (this *DialogObject) createTemplate() win32.HGLOBAL {
	var templ win32.DLGTEMPLATE
	templ.Style = uint32(win32.WS_POPUP | win32.WS_CAPTION |
		win32.WS_SYSMENU | WINDOW_STYLE(win32.DS_MODALFRAME))

	templ.Cx = 200
	templ.Cy = 200

	hGlobal, errno := win32.GlobalAlloc(win32.GMEM_ZEROINIT, 1024)
	if hGlobal == 0 {
		log.Fatal(errno)
	}
	pData, errno := win32.GlobalLock(hGlobal)
	*(*win32.DLGTEMPLATE)(pData) = templ

	win32.GlobalUnlock(hGlobal)
	return hGlobal
}
