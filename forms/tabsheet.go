package forms

import (
	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/framework/virtual"
	"syscall"
	"unsafe"
)

type TabSheet interface {
	ContainerControl
	//TextAware
	TitleAware
}

type TabSheetSpi interface {
	ControlSpi
}

type TabSheetInterface interface {
	TabSheet
	TabSheetSpi
}

type TabSheetObject struct {
	ContainerControlObject
	super *ContainerControlObject
}

func NewTabSheetObject() *TabSheetObject {
	return virtual.New[TabSheetObject]()
}

type NewTabSheet struct {
	Parent Container
	Name   string
	Title  string
}

func (me NewTabSheet) Create(extraOpts ...*WindowOptions) TabSheet {
	tabSheet := NewTabSheetObject()
	tabSheet.name = me.Name

	opts := utils.OptionalArg(extraOpts)
	opts.WindowName = me.Title

	opts.ParentHandle = resolveParentHandle(me.Parent)
	err := tabSheet.Create(*opts)
	assertNoErr(err)

	return tabSheet
}

func tabSheetProc(hWnd HWND, uMsg uint32, wParam WPARAM, lParam LPARAM) win32.LRESULT {
	switch uMsg {
	case win32.WM_INITDIALOG:
		win32.EnableThemeDialogTexture(hWnd, win32.ETDT_ENABLE|win32.ETDT_USETABTEXTURE)
		break
	case win32.WM_CLOSE:
		//win32.EndTabSheet(hWnd, uintptr(win32.IDCANCEL))
		break
	}
	return 0
}

var tabSheetProcCallback = syscall.NewCallback(tabSheetProc)

//func (this *TabSheetObject) CreateIn(parent Window) Control {
//	if parent == nil {
//		parent = ContextContainer
//	}
//
//	var templ win32.DLGTEMPLATE
//	templ.Style = uint32(win32.WS_CHILDWINDOW|win32.WS_VISIBLE) |
//		uint32(win32.DS_3DLOOK|win32.DS_CONTROL)
//
//	hWnd, errno := win32.CreateDialogIndirectParam(HInstance, &templ,
//		parent.GetHandle(), tabSheetProcCallback, 0)
//	if hWnd == 0 {
//		log.Fatal(errno)
//	}
//	this.Handle = hWnd
//	if container, ok := parent.(Container); ok {
//		container.Add(this)
//	}
//	return this
//}

//func (this *TabSheetObject) Create(options WindowOptions) error {
//	win := this.RealObject
//	win.PreCreate(&options)
//	//return this.super.Create(options)
//	var templ win32.DLGTEMPLATE
//	templ.Style = uint32(win32.WS_CHILDWINDOW|win32.WS_VISIBLE) |
//		uint32(win32.DS_3DLOOK|win32.DS_CONTROL)
//
//	hWnd, errno := win32.CreateDialogIndirectParam(HInstance, &templ,
//		options.ParentHandle, tabSheetProcCallback, 0)
//	if hWnd == 0 {
//		log.Fatal(errno)
//	}
//	this.Handle = hWnd
//	windowMap[hWnd] = this.RealObject
//
//}

func (this *TabSheetObject) SetTitle(title string) {
	win32.SetWindowText(this.Handle, win32.StrToPwstr(title))
}

func (this *TabSheetObject) GetTitle() string {
	buf := make([]uint16, 255)
	cc, _ := win32.GetWindowText(this.Handle, &buf[0], int32(len(buf)))
	return syscall.UTF16ToString(buf[:cc])
}

func (this *TabSheetObject) GetControlSpecStyle() (include, exclude WINDOW_STYLE) {
	return WINDOW_STYLE(win32.DS_3DLOOK | win32.DS_CONTROL), 0
}

type _DLGTEMPLATEEX struct {
	dlgVer       win32.WORD
	signature    win32.WORD
	helpID       win32.DWORD
	exStyle      win32.DWORD
	style        win32.DWORD
	cDlgItems    win32.WORD
	x            win32.SHORT
	y            win32.SHORT
	cx           win32.SHORT
	cy           win32.SHORT
	_menu        win32.WORD
	_windowClass win32.WORD
	title        [80]uint16
}

func (this *TabSheetObject) CreateHandle(
	className string, options WindowOptions) (HWND, error) {

	resolveWindowOptions(&options)
	x, y, cx, cy := options.Left, options.Top, options.Width, options.Height
	x, y = this.PxToDlu(x, y)
	cx, cy = this.PxToDlu(cx, cy)

	templ := _DLGTEMPLATEEX{
		dlgVer:    1,
		signature: 0xFFFF,
		style:     win32.DWORD(options.Style),
		exStyle:   win32.DWORD(options.ExStyle),
		x:         win32.SHORT(x),
		y:         win32.SHORT(y),
		cx:        win32.SHORT(cx),
		cy:        win32.SHORT(cy),
	}
	copy(templ.title[:], StrToWsz(options.WindowName))

	hWnd, errno := win32.CreateDialogIndirectParam(HInstance,
		(*win32.DLGTEMPLATE)(unsafe.Pointer(&templ)),
		options.ParentHandle, tabSheetProcCallback, 0)
	var err error
	if hWnd == 0 {
		err = errno
	}
	return hWnd, err
}
