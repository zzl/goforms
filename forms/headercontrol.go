package forms

import (
	"log"
	"syscall"
	"unsafe"

	"github.com/zzl/goforms/framework/virtual"
	"github.com/zzl/goforms/layouts/aligns"

	"github.com/zzl/go-win32api/v2/win32"
)

type HeaderControl interface {
	Control
}

type HeaderControlObject struct {
	ControlObject
	super *ControlObject
}

func (this *HeaderControlObject) Init() {
	this.super.Init()
}

func (this *HeaderControlObject) GetWindowClass() string {
	return "SysHeader32"
}

func (this *HeaderControlObject) GetControlSpecStyle() (include, exclude WINDOW_STYLE) {
	return WINDOW_STYLE(win32.HDS_BUTTONS | win32.HDS_HORZ | win32.HDS_FULLDRAG), win32.WS_TABSTOP
}

func NewHeaderControlObject() *HeaderControlObject {
	return virtual.New[HeaderControlObject]()
}

func (this *HeaderControlObject) GetItemCount() int {
	ret, _ := SendMessage(this.Handle, win32.HDM_GETITEMCOUNT, 0, 0)
	return int(ret)
}

func (this *HeaderControlObject) AddItem(text string, width int) {
	this._addItem(text, width, aligns.Default)
}

func (this *HeaderControlObject) AddItemAligned(text string, width int, align int) {
	this._addItem(text, width, align)
}

func (this *HeaderControlObject) _addItem(text string, width int, align int) {
	var hdi win32.HDITEM
	hdi.Mask = win32.HDI_TEXT | win32.HDI_FORMAT | win32.HDI_WIDTH
	hdi.Cxy = int32(width)
	wsz, _ := syscall.UTF16FromString(text)
	hdi.PszText = &wsz[0]
	hdi.CchTextMax = int32(len(wsz) - 1)
	hdi.Fmt = win32.HDF_STRING
	switch align {
	case aligns.Right:
		hdi.Fmt |= win32.HDF_RIGHT
	case aligns.Center:
		hdi.Fmt |= win32.HDF_CENTER
	default:
		hdi.Fmt |= win32.HDF_LEFT
	}
	count := this.GetItemCount()

	index, errno := SendMessage(this.Handle, win32.HDM_INSERTITEM,
		count, unsafe.Pointer(&hdi))
	if index == NegativeOne {
		log.Fatal(errno)
	}
}

func (this *HeaderControlObject) GetPreferredSize(cxMax int, cyMax int) (int, int) {
	if this.Handle == 0 {
		return this.super.GetPreferredSize(cxMax, cyMax)
	}
	var hdl win32.HDLAYOUT
	rc := win32.RECT{Left: 0, Top: 0, Right: int32(cxMax), Bottom: int32(cyMax)}
	hdl.Prc = &rc
	var wpos win32.WINDOWPOS
	hdl.Pwpos = &wpos
	SendMessage(this.Handle, win32.HDM_LAYOUT,
		0, unsafe.Pointer(&hdl))
	cx, cy := wpos.Cx, wpos.Cy
	return int(cx), int(cy)
}

func (this *HeaderControlObject) GetItemWidth(index int) int {
	var hdi win32.HDITEM
	hdi.Mask = win32.HDI_WIDTH
	SendMessage(this.Handle, win32.HDM_GETITEM,
		index, unsafe.Pointer(&hdi))
	return int(hdi.Cxy)
}

func (this *HeaderControlObject) GetItemText(index int) string {
	var hdi win32.HDITEM
	hdi.Mask = win32.HDI_TEXT
	buf := make([]uint16, win32.MAX_PATH)
	hdi.PszText = &buf[0]
	hdi.CchTextMax = int32(win32.MAX_PATH)
	SendMessage(this.Handle, win32.HDM_GETITEM,
		index, unsafe.Pointer(&hdi))
	return win32.PwstrToStr(hdi.PszText)
}

func (this *HeaderControlObject) SetItemWidth(index int, width int) {
	var hdi win32.HDITEM
	hdi.Mask = win32.HDI_WIDTH
	hdi.Cxy = int32(width)
	SendMessage(this.Handle, win32.HDM_SETITEM,
		index, unsafe.Pointer(&hdi))
}
