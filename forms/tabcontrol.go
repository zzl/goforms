package forms

import (
	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/framework/virtual"
	"log"
	"syscall"
	"unsafe"

	"github.com/zzl/go-win32api/v2/win32"
)

type TabItem struct {
	Title string
	Sheet Control
}

type TabControl interface {
	Control

	SetItems(items []TabItem)
	SetSelectedIndex(index int)
	GetSelectedIndex() int

	SelectFirst()
	SelectPrev()
	SelectNext()
	SelectLast()

	TabControlObj() *TabControlObject
}

type TabControlObject struct {
	ControlObject
	super *ControlObject

	items []TabItem
}

type NewTabControl struct {
	Parent   Container
	Name     string
	Pos      Point
	Size     Size
	TabItems []TabItem
}

func (me NewTabControl) Create(extraOpts ...*WindowOptions) TabControl {
	tabControl := NewTabControlObject()
	tabControl.name = me.Name

	opts := utils.OptionalArg(extraOpts)
	opts.Left = me.Pos.X
	opts.Top = me.Pos.Y

	opts.ParentHandle = resolveParentHandle(me.Parent)
	err := tabControl.Create(*opts)
	assertNoErr(err)
	configControlSize(tabControl, me.Size)

	if len(me.TabItems) != 0 {
		tabControl.SetItems(me.TabItems)
	}
	return tabControl
}

func NewTabControlObject() *TabControlObject {
	return virtual.New[TabControlObject]()
}

func (this *TabControlObject) TabControlObj() *TabControlObject {
	return this
}

func (this *TabControlObject) Init() {
	this.super.Init()
}

func (this *TabControlObject) GetWindowClass() string {
	return "SysTabControl32"
}

func (this *TabControlObject) SelectFirst() {
	this.SetSelectedIndex(0)
}

func (this *TabControlObject) SelectPrev() {
	index := this.GetSelectedIndex() - 1
	if index < 0 {
		index = len(this.items) - 1
	}
	this.SetSelectedIndex(index)
}

func (this *TabControlObject) SelectNext() {
	index := this.GetSelectedIndex() + 1
	if index >= len(this.items) {
		index = 0
	}
	this.SetSelectedIndex(index)
}

func (this *TabControlObject) SelectLast() {
	this.SetSelectedIndex(len(this.items) - 1)
}

func (this *TabControlObject) OnReflectNotify(msg *NotifyMessage) {
	if msg.GetNMHDR().Code == win32.TCN_SELCHANGE {
		this.onSelChange()
	}
}

func (this *TabControlObject) onSelChange() {
	index := this.GetSelectedIndex()
	this.SetSelectedIndex(index)
}

func (this *TabControlObject) SetItems(items []TabItem) {
	this.items = items
	var ti win32.TCITEM
	ti.Mask = win32.TCIF_TEXT

	for n, item := range this.items {
		title := item.Title
		if title == "" {
			if titleAware, ok := item.Sheet.(TitleAware); ok {
				title = titleAware.GetTitle()
			}
		}
		ti.PszText, _ = syscall.UTF16PtrFromString(title)
		ret, errno := SendMessage(this.Handle, win32.TCM_INSERTITEM,
			n, unsafe.Pointer(&ti))
		if ret == NegativeOne {
			log.Fatal(errno)
		}
	}
	this.updateItemsBounds()
	if len(this.items) > 0 {
		this.SetSelectedIndex(0)
	}
}

func (this *TabControlObject) GetSelectedIndex() int {
	ret, errno := SendMessage(this.Handle, win32.TCM_GETCURSEL, 0, 0)
	_ = errno
	return int(ret)
}

func (this *TabControlObject) SetSelectedIndex(index int) {
	SendMessage(this.Handle, win32.TCM_SETCURSEL, index, 0)
	for n, item := range this.items {
		sheet := item.Sheet
		if n == index {
			sheet.Show()
		} else {
			sheet.Hide()
		}
	}
}

func (this *TabControlObject) SetBounds(left, top, width, height int) {
	this.super.SetBounds(left, top, width, height)
	this.updateItemsBounds()
}

func (this *TabControlObject) updateItemsBounds() {
	b := this.GetBounds()
	var rc = b.ToRECT()
	_, errno := SendMessage(this.Handle, win32.TCM_ADJUSTRECT,
		0, unsafe.Pointer(&rc))
	_ = errno
	rect := RectFromRECT(rc)
	for _, item := range this.items {
		item.Sheet.SetBounds(rect.Left, rect.Top,
			rect.Width(), rect.Height())
	}
}

func (this *TabControlObject) GetPreferredSize(maxWidth int, maxHeight int) (int, int) {
	var rc win32.RECT
	rc.Right = int32(maxWidth)
	rc.Bottom = int32(maxHeight)
	SendMessage(this.Handle, win32.TCM_ADJUSTRECT,
		0, unsafe.Pointer(&rc))

	maxCx, maxCy := 0, 0
	sheetMaxWidth, sheetMaxHeight := int(rc.Right), int(rc.Bottom)
	for _, item := range this.items {
		cx, cy := item.Sheet.GetPreferredSize(sheetMaxWidth, sheetMaxHeight)
		if cx == 0 {
			cx = 64
		}
		if cy == 0 {
			cy = 16
		}
		if cx > maxCx {
			maxCx = cx
		}
		if cy > maxCy {
			maxCy = cy
		}
	}
	rc.Right = int32(maxCx)
	rc.Bottom = int32(maxCy)
	SendMessage(this.Handle, win32.TCM_ADJUSTRECT,
		1, unsafe.Pointer(&rc))
	cx, cy := int(rc.Right), int(rc.Bottom)
	return cx, cy
}

func (this *TabControlObject) MeasureSheetRect(controlRect Rect) Rect {
	rect32 := controlRect.ToRECT()
	SendMessage(this.Handle, win32.TCM_ADJUSTRECT, 0,
		unsafe.Pointer(&rect32))
	return RectFromRECT(rect32)
}

func (this *TabControlObject) MeasureControlRect(sheetRect Rect) Rect {
	rect32 := sheetRect.ToRECT()
	SendMessage(this.Handle, win32.TCM_ADJUSTRECT, 1,
		unsafe.Pointer(&rect32))
	return RectFromRECT(rect32)
}
