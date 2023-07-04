package forms

import (
	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/framework/virtual"
	"log"
	"syscall"
	"unsafe"

	"github.com/zzl/go-win32api/v2/win32"
)

type Tooltip interface {
	Window

	AddTool(control Control, textProvider TooltipTextProvider)
}

type tooltipTool struct {
	control      Control
	textProvider TooltipTextProvider
	tipText      []uint16
}

type TooltipObject struct {
	WindowObject
	super *WindowObject

	//c Container
	hWndContainer HWND

	toolsMap map[HWND]*tooltipTool
}

type TooltipTextProvider interface {
	GetTooltipText() string
}

type GetToolTipTextFunc func() string

func (me GetToolTipTextFunc) GetTooltipText() string {
	return me()
}

func NewTooltipObject() *TooltipObject {
	return virtual.New[TooltipObject]()
}

func (this *TooltipObject) Init() {
	this.super.Init()

	this.toolsMap = make(map[HWND]*tooltipTool)
}

func (this *TooltipObject) Dispose() {
	win32.DestroyWindow(this.Handle)
	this.super.Dispose()
}

func (this *TooltipObject) GetWindowClass() string {
	return "tooltips_class32"
}

func (this *TooltipObject) GetDefaultStyle() WINDOW_STYLE {
	return WINDOW_STYLE(win32.TTS_NOPREFIX) // | win32.TTS_ALWAYSTIP
}

func (this *TooltipObject) CreateIn(parent Window, extraOpts ...*WindowOptions) {
	this.hWndContainer = parent.GetHandle() //?

	opts := utils.OptionalArg(extraOpts)

	opts.ParentHandle = parent.GetHandle()
	err := this.RealObject.Create(*opts)
	if err != nil {
		log.Fatal(err)
	}

	var ti win32.TOOLINFO
	ti.CbSize = uint32(unsafe.Sizeof(ti))
	ti.Hwnd = parent.GetHandle()         //hWndParent
	ti.UFlags = uint32(win32.TTF_TRACK | //win32.TTF_IDISHWND
		win32.TTF_ABSOLUTE |
		win32.TTF_TRANSPARENT) //|win32.TTF_SUBCLASS
	ti.UId = uintptr(parent.GetHandle())
	tooltip := "tooltip.."
	pwszTt, _ := syscall.UTF16PtrFromString(tooltip)
	ti.LpszText = pwszTt
	ret, errno := SendMessage(this.Handle, win32.TTM_ADDTOOL,
		0, unsafe.Pointer(&ti))
	if ret == 0 {
		log.Fatal(errno.Error())
	}

}

func (this *TooltipObject) ShowAt(x, y int) {
	var ti win32.TOOLINFO
	ti.CbSize = uint32(unsafe.Sizeof(ti))
	ti.Hinst = HInstance
	ti.Hwnd = this.hWndContainer // this.c.GetHandle()
	ti.UId = uintptr(this.hWndContainer)

	_, errno := SendMessage(this.Handle, win32.TTM_TRACKACTIVATE, 1, unsafe.Pointer(&ti))
	if errno != win32.NO_ERROR {
		println("?")
	}

	pt := win32.POINT{X: int32(x), Y: int32(y)}
	win32.ClientToScreen(this.hWndContainer, &pt)

	lParam := win32.MAKELONG(uint16(pt.X), uint16(pt.Y))
	SendMessage(this.Handle, win32.TTM_TRACKPOSITION, 0, lParam)
}

func (this *TooltipObject) SetText(text string) {
	var ti win32.TOOLINFO
	ti.CbSize = uint32(unsafe.Sizeof(ti))
	ti.Hinst = HInstance
	ti.Hwnd = this.hWndContainer //this.c.GetHandle() //hWndParent
	ti.UId = uintptr(this.hWndContainer)
	pwszTt, _ := syscall.UTF16PtrFromString(text)
	ti.LpszText = pwszTt
	SendMessage(this.Handle, win32.TTM_UPDATETIPTEXT,
		0, unsafe.Pointer(&ti))
}

func (this *TooltipObject) OnReflectNotify(msg *NotifyMessage) {
	pNmhdr := msg.GetNMHDR()
	if pNmhdr.Code == win32.TTN_SHOW {
		//?
	} else if pNmhdr.Code == win32.TTN_GETDISPINFO {
		pNmttdi := (*win32.NMTTDISPINFO)(unsafe.Pointer(pNmhdr))
		hWndTool := HWND(pNmhdr.IdFrom)
		ttt, ok := this.toolsMap[hWndTool]
		if ok {
			text := ttt.textProvider.GetTooltipText()
			ttt.tipText, _ = syscall.UTF16FromString(text)
			pNmttdi.LpszText = &ttt.tipText[0]
			pNmttdi.Hinst = 0
		}
	}
}

func (this *TooltipObject) Show() {
	//?
}

func (this *TooltipObject) Hide() {
	var ti win32.TOOLINFO
	ti.CbSize = uint32(unsafe.Sizeof(ti))
	ti.Hwnd = this.hWndContainer //this.c.GetHandle()//0
	ti.UId = 0                   //uintptr(this.Handle)
	_, _ = SendMessage(this.Handle, win32.TTM_TRACKACTIVATE,
		0, unsafe.Pointer(&ti))
}

func (this *TooltipObject) AddTool(control Control, textProvider TooltipTextProvider) {
	var ti win32.TOOLINFO
	ti.CbSize = uint32(unsafe.Sizeof(ti))
	ti.Hinst = HInstance
	ti.Hwnd = this.hWndContainer
	ti.UFlags = uint32(win32.TTF_IDISHWND | win32.TTF_SUBCLASS)
	ti.UId = uintptr(control.GetHandle())
	ti.LpszText = (*uint16)(unsafe.Pointer(NegativeOne)) //win32.LPSTR_TEXTCALLBACK// pwszTt
	ret, errno := SendMessage(this.Handle, win32.TTM_ADDTOOL,
		0, unsafe.Pointer(&ti))
	if ret == 0 {
		log.Fatal(errno.Error())
	}
	if textProvider == nil {
		var ok bool
		textProvider, ok = control.(TooltipTextProvider)
		if !ok {
			log.Fatal("text provider unspecified")
		}
	}
	this.toolsMap[control.GetHandle()] = &tooltipTool{
		control: control, textProvider: textProvider}
}

func (this *TooltipObject) ShowInplace(x int, y int) {
	var rc win32.RECT
	SendMessage(this.Handle, win32.TTM_ADJUSTRECT,
		1, unsafe.Pointer(&rc))
	this.ShowAt(x+int(rc.Left), y+int(rc.Top))
}
