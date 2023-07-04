package forms

import (
	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/forms/internal"
	"log"
	"runtime"
	"unsafe"
)

// HInstance is the handle to the current instance of the application.
var HInstance win32.HINSTANCE

// HWndActive is the handle of the last activated window
var HWndActive HWND

// hAccelActive is the active accelerator table handle
var hAccelActive win32.HACCEL

// windowMap maps window handles to their wrapper object interfaces
var windowMap map[HWND]WindowInterface

// creatingWindows contains window objects that are being created
// but have not yet been mapped with their handles in windowMap
var creatingWindows []*WindowObject

func init() {
	runtime.LockOSThread()
	HInstance, _ = win32.GetModuleHandle(nil)

	internal.ActivateCommonControlsV6IfNeeded()

	hr := win32.CoInitializeEx(unsafe.Pointer(nil), win32.COINIT_APARTMENTTHREADED)
	win32.ASSERT_SUCCEEDED(hr)

	//
	var iccEx win32.INITCOMMONCONTROLSEX
	iccEx.DwSize = uint32(unsafe.Sizeof(iccEx))
	iccEx.DwICC = win32.ICC_STANDARD_CLASSES | win32.ICC_WIN95_CLASSES |
		win32.ICC_DATE_CLASSES | win32.ICC_LINK_CLASS

	bOk := win32.InitCommonControlsEx(&iccEx)
	if bOk == win32.FALSE {
		log.Println("Warning: Common controls init failed..")
	}

	windowMap = make(map[HWND]WindowInterface)
}
