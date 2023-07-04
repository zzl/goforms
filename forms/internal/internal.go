package internal

import (
	"github.com/zzl/go-win32api/v2/win32"
	"log"
	"os"
	"syscall"
	"unsafe"
)

type ACTCTX struct {
	CbSize                 uint32
	DwFlags                uint32
	LpSource               win32.PWSTR
	WProcessorArchitecture uint16
	WLangID                win32.LANGID
	LpAssemblyDirectory    win32.PWSTR
	LpResourceName         win32.PWSTR
	LpApplicationName      win32.PWSTR
	HModule                win32.HINSTANCE
}

func enumResName(hModule win32.HINSTANCE, lpType win32.PWSTR,
	lpName win32.PWSTR, lParam uintptr) uintptr {
	return 1
}

func ActivateCommonControlsV6IfNeeded() {

	pEnumResName := syscall.NewCallback(enumResName)
	ok := win32.EnumResourceNames(0, win32.MAKEINTRESOURCE(uint16(win32.RT_MANIFEST)),
		pEnumResName, 0)
	if ok == win32.TRUE {
		return
	}
	exePath, _ := os.Executable()
	if _, err := os.Stat(exePath + ".manifest"); err == nil {
		return
	}

	system32 := make([]uint16, win32.MAX_PATH)
	win32.GetSystemDirectory(&system32[0], uint32(len(system32)))

	var actctx ACTCTX
	actctx.CbSize = uint32(unsafe.Sizeof(actctx))
	actctx.DwFlags = win32.ACTCTX_FLAG_RESOURCE_NAME_VALID |
		win32.ACTCTX_FLAG_SET_PROCESS_DEFAULT |
		win32.ACTCTX_FLAG_ASSEMBLY_DIRECTORY_VALID

	actctx.LpSource = win32.StrToPwstr("shell32.dll")
	actctx.LpAssemblyDirectory = &system32[0]
	actctx.LpResourceName = win32.MAKEINTRESOURCE(124)

	libKernel32 := syscall.NewLazyDLL("kernel32.dll")
	procCreateActCtx := libKernel32.NewProc("CreateActCtxW")
	procActivateActCtx := libKernel32.NewProc("ActivateActCtx")

	hActCtx, _, errno := procCreateActCtx.Call(uintptr(unsafe.Pointer(&actctx)))
	if hActCtx == win32.INVALID_HANDLE_VALUE {
		log.Panic(errno)
	}
	var cookie uintptr
	ret, _, errno := procActivateActCtx.Call(hActCtx, uintptr(unsafe.Pointer(&cookie)))
	if win32.BOOL(ret) == win32.FALSE {
		log.Panic(errno)
	}
}
