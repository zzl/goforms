package forms

import (
	"github.com/zzl/goforms/framework/consts"
	"github.com/zzl/goforms/framework/types"
	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/framework/virtual"
	"log"
	"syscall"
	"unsafe"

	"github.com/zzl/go-win32api/v2/win32"
)

type ListBox interface {
	Control
	Input

	ClearItems()

	AddItem(item types.ValueText)
	AddItems(items []types.ValueText)

	AddStringItem(item string)
	AddStringItems(items []string)

	GetItemCount() int
	GetItem(index int) *types.ValueText

	GetSelectedItem() *types.ValueText

	GetSelectedIndex() int
	SetSelectedIndex(index int)

	GetText() string
	DeleteItem(index int)

	IndexFromPoint(x, y int) int

	ListBoxObj() *ListBoxObject
}

type ListBoxObject struct {
	ControlObject

	IntegralHeight bool
	ForeColorAwareSupport

	OnValueChange SimpleEvent
	values        []int
}

type NewListBox struct {
	Parent         Container
	Name           string
	Pos            Point
	Size           Size
	IntegralHeight bool
	Items          []types.ValueText
	StringItems    []string
}

func (me NewListBox) Create(extraOpts ...*WindowOptions) ListBox {
	listBox := NewListBoxObject()
	listBox.name = me.Name
	listBox.IntegralHeight = me.IntegralHeight

	opts := utils.OptionalArg(extraOpts)
	opts.Left = me.Pos.X
	opts.Top = me.Pos.Y

	opts.ParentHandle = resolveParentHandle(me.Parent)
	err := listBox.Create(*opts)
	assertNoErr(err)
	configControlSize(listBox, me.Size)

	if len(me.Items) != 0 {
		listBox.AddItems(me.Items)
	}
	if len(me.StringItems) != 0 {
		listBox.AddStringItems(me.StringItems)
	}

	return listBox
}

func NewListBoxObject() *ListBoxObject {
	return virtual.New[ListBoxObject]()
}

func (this *ListBoxObject) ListBoxObj() *ListBoxObject {
	return this
}

func (this *ListBoxObject) GetWindowClass() string {
	return "ListBox"
}

func (this *ListBoxObject) GetValue() any {
	index := this.GetSelectedIndex()
	if index == -1 {
		return nil
	}
	value := this.values[index]
	if value == consts.Null {
		return this.GetItemText(index)
	}
	return value
}

func (this *ListBoxObject) GetText() string {
	item := this.GetSelectedItem()
	if item == nil {
		return ""
	}
	return item.Text
}

func (this *ListBoxObject) SetValue(value any) {
	if value == nil {
		this.SetSelectedIndex(-1)
		return
	}
	nValue, ok := value.(int)
	if ok {
		if nValue == consts.Null {
			this.SetSelectedIndex(-1)
			return
		}
		for n, v := range this.values {
			if v == value {
				this.SetSelectedIndex(n)
				return
			}
		}
		this.SetSelectedIndex(-1)
	} else {
		text := value.(string)
		this.SelectByText(text)
	}
}

func (this *ListBoxObject) SelectByText(text string) bool {
	pwsz, _ := syscall.UTF16PtrFromString(text)
	ret, _ := SendMessage(this.Handle, win32.LB_FINDSTRINGEXACT,
		NegativeOne, unsafe.Pointer(pwsz))
	this.SetSelectedIndex(int(ret))
	return int32(ret) != win32.CB_ERR
}

func (this *ListBoxObject) GetOnValueChange() *SimpleEvent {
	return &this.OnValueChange
}

func (this *ListBoxObject) ClearItems() {
	SendMessage(this.Handle, win32.LB_RESETCONTENT, 0, 0)
	this.values = nil
}

func (this *ListBoxObject) AddItem(item types.ValueText) {
	pwszText, _ := syscall.UTF16PtrFromString(item.Text)
	index, errno := SendMessage(this.Handle, win32.LB_ADDSTRING,
		0, unsafe.Pointer(pwszText))
	if int32(index) == win32.LB_ERR {
		log.Fatal(errno)
	}
	this.values = append(this.values, item.Value)
}

func (this *ListBoxObject) AddItems(items []types.ValueText) {
	for _, item := range items {
		this.AddItem(item)
	}
}

func (this *ListBoxObject) AddStringItem(item string) {
	pwszText, _ := syscall.UTF16PtrFromString(item)
	index, errno := SendMessage(this.Handle, win32.LB_ADDSTRING,
		0, unsafe.Pointer(pwszText))
	if int32(index) == win32.LB_ERR {
		log.Fatal(errno)
	}
	this.values = append(this.values, consts.Null)
}

func (this *ListBoxObject) AddStringItems(items []string) {
	this.ClearItems()
	for _, item := range items {
		this.AddStringItem(item)
	}
}

func (this *ListBoxObject) DeleteItem(index int) {
	SendMessage(this.Handle, win32.LB_DELETESTRING, index, 0)
	var newValues []int
	newValues = append(newValues, this.values[:index]...)
	newValues = append(newValues, this.values[index+1:]...)
	this.values = newValues
}

func (this *ListBoxObject) GetItemCount() int {
	retVal, errono := SendMessage(this.Handle, win32.LB_GETCOUNT, 0, 0)
	if int32(retVal) == win32.LB_ERR {
		log.Fatal(errono)
	}
	return int(retVal)
}

func (this *ListBoxObject) GetItem(index int) *types.ValueText {
	text := this.GetItemText(index)
	value := this.values[index]
	return &types.ValueText{Value: value, Text: text}
}

func (this *ListBoxObject) GetItemText(index int) string {
	cch, _ := SendMessage(this.Handle,
		win32.LB_GETTEXTLEN, index, 0)
	if int32(cch) == win32.LB_ERR {
		return ""
	}
	buf := make([]uint16, cch+1)
	cch, _ = SendMessage(this.Handle, win32.LB_GETTEXT,
		index, unsafe.Pointer(&buf[0]))
	if int32(cch) == win32.LB_ERR {
		return ""
	}
	text := syscall.UTF16ToString(buf)
	return text
}

func (this *ListBoxObject) GetSelectedItem() *types.ValueText {
	index, _ := SendMessage(this.Handle, win32.LB_GETCURSEL, 0, 0)
	if int32(index) == win32.LB_ERR {
		return nil
	}
	return this.GetItem(int(index))
}

func (this *ListBoxObject) GetSelectedIndex() int {
	index, _ := SendMessage(this.Handle, win32.LB_GETCURSEL, 0, 0)
	return int(index)
}

func (this *ListBoxObject) SetSelectedIndex(index int) {
	SendMessage(this.Handle, win32.LB_SETCURSEL, index, 0)
}

func (this *ListBoxObject) GetPreferredSize(int, int) (int, int) {
	if this.Handle == 0 {
		return 16, 16
	}
	itemCount := this.GetItemCount()
	ret, _ := SendMessage(this.Handle, win32.LB_GETITEMHEIGHT, 0, 0)
	itemHeight := int(ret)
	cy := itemHeight * itemCount
	var rc win32.RECT
	win32.GetWindowRect(this.Handle, &rc)
	cy += int(rc.Bottom-rc.Top) % itemHeight
	return 16, cy
}

func (this *ListBoxObject) OnReflectCommand(msg *CommandMessage) {
	notifyCode := msg.GetNotifyCode()
	if notifyCode == uint16(win32.LBN_SELCHANGE) {
		this.OnValueChange.Fire(this, &SimpleEventInfo{})
	}
}

func (this *ListBoxObject) GetControlSpecStyle() (WINDOW_STYLE, WINDOW_STYLE) {
	include := WINDOW_STYLE(win32.LBS_NOTIFY)
	if !this.IntegralHeight {
		include |= WINDOW_STYLE(win32.LBS_NOINTEGRALHEIGHT)
	}
	include |= win32.WS_VSCROLL
	return include, 0
}

func (this *ListBoxObject) GetDefaultExStyle() WINDOW_EX_STYLE {
	return win32.WS_EX_CLIENTEDGE
}

func (this *ListBoxObject) IndexFromPoint(x, y int) int {
	lparam := LPARAM(win32.MAKELONG(uint16(x), uint16(y)))
	ret, _ := SendMessage(this.Handle, win32.LB_ITEMFROMPOINT, 0, lparam)
	dwRet := win32.DWORD(ret)
	index := win32.LOWORD(dwRet)
	outside := win32.HIWORD(dwRet)
	if outside == 0 {
		return int(index)
	}
	return -1
}
