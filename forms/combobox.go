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

type ComboBox interface {
	Control
	Input

	TextAware

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

	ComboBoxObj() *ComboBoxObject
}

type ComboBoxSpi interface {
	ControlSpi
}

type ComboBoxInterface interface {
	ComboBox
	ComboBoxSpi
}

type ComboBoxObject struct {
	ControlObject

	ForeColorAwareSupport

	Editable bool
	Height   int

	OnValueChange SimpleEvent
	values        []int
}

var _comboBoxButtonWidth int
var _comboBoxHeight int

func NewComboBoxObject() *ComboBoxObject {
	return virtual.New[ComboBoxObject]()
}

type NewComboBox struct {
	Parent   Container
	Name     string
	Text     string
	Pos      Point
	Size     Size
	Disabled bool
	Editable bool //?
	Items    []string
}

func (me NewComboBox) Create(extraOpts ...*WindowOptions) ComboBox {
	cb := NewComboBoxObject()
	cb.name = me.Name
	cb.Editable = me.Editable

	opts := utils.OptionalArg(extraOpts)
	//opts.WindowName = me.Text
	opts.Left = me.Pos.X
	opts.Top = me.Pos.Y
	if me.Disabled {
		opts.StyleInclude |= win32.WS_DISABLED
	}
	opts.ParentHandle = resolveParentHandle(me.Parent)
	err := cb.Create(*opts)
	assertNoErr(err)
	configControlSize(cb, me.Size)
	cb.AddStringItems(me.Items)
	if me.Text != "" {
		cb.SetText(me.Text)
	}
	return cb
}

func (this *ComboBoxObject) ComboBoxObj() *ComboBoxObject {
	return this
}

func (this *ComboBoxObject) GetItemCount() int {
	retVal, errono := SendMessage(this.Handle, win32.CB_GETCOUNT, 0, 0)
	if int32(retVal) == win32.CB_ERR {
		log.Fatal(errono)
	}
	return int(retVal)
}

func (this *ComboBoxObject) ClearItems() {
	_, _ = SendMessage(this.Handle, win32.CB_RESETCONTENT, 0, 0)
	this.values = nil
}

func (this *ComboBoxObject) AddStringItem(item string) {
	pwszText, _ := syscall.UTF16PtrFromString(item)
	index, errno := SendMessage(this.Handle, win32.CB_ADDSTRING,
		0, unsafe.Pointer(pwszText))
	if int32(index) == win32.CB_ERR {
		log.Fatal(errno)
	}
	this.values = append(this.values, consts.Null)
}

func (this *ComboBoxObject) AddStringItems(items []string) {
	this.ClearItems()
	for _, item := range items {
		this.AddStringItem(item)
	}
}

func (this *ComboBoxObject) AddItem(item types.ValueText) {
	pwszText, _ := syscall.UTF16PtrFromString(item.Text)
	index, errno := SendMessage(this.Handle, win32.CB_ADDSTRING,
		0, unsafe.Pointer(pwszText))
	if int32(index) == win32.CB_ERR {
		log.Fatal(errno)
	}
	this.values = append(this.values, item.Value)
	//_, _ = win32.SendMessage(this.Handle, win32.CB_SETITEMDATA,
	//	index, uintptr(item.Value))
}

func (this *ComboBoxObject) AddItems(items []types.ValueText) {
	for _, item := range items {
		this.AddItem(item)
	}
}

func (this *ComboBoxObject) GetSelectedIndex() int {
	index, _ := SendMessage(this.Handle, win32.CB_GETCURSEL, 0, 0)
	return int(index)
}

func (this *ComboBoxObject) SetSelectedIndex(index int) {
	SendMessage(this.Handle, win32.CB_SETCURSEL, index, 0)
}

func (this *ComboBoxObject) GetItem(index int) *types.ValueText {
	cch, _ := SendMessage(this.Handle,
		win32.CB_GETLBTEXTLEN, index, 0)
	if int32(cch) == win32.CB_ERR {
		return nil
	}
	buf := make([]uint16, cch+1)
	cch, _ = SendMessage(this.Handle, win32.CB_GETLBTEXT,
		index, unsafe.Pointer(&buf[0]))
	if int32(cch) == win32.CB_ERR {
		return nil
	}
	text := syscall.UTF16ToString(buf)
	value := this.values[index]
	return &types.ValueText{Value: value, Text: text}
}

func (this *ComboBoxObject) GetSelectedItem() *types.ValueText {
	index, errno := SendMessage(this.Handle, win32.CB_GETCURSEL, 0, 0)
	_ = errno
	if int32(index) == win32.CB_ERR {
		return nil
	}
	return this.GetItem(int(index))
}

func (this *ComboBoxObject) GetWindowClass() string {
	return "COMBOBOX"
}

func (this *ComboBoxObject) GetControlSpecStyle() (WINDOW_STYLE, WINDOW_STYLE) {
	var style WINDOW_STYLE
	if this.Editable {
		style = WINDOW_STYLE(win32.CBS_DROPDOWN)
	} else {
		style = WINDOW_STYLE(win32.CBS_DROPDOWNLIST)
	}
	if this.Height != 0 {
		style |= WINDOW_STYLE(win32.CBS_NOINTEGRALHEIGHT)
	}
	return style, 0
}

func (this *ComboBoxObject) SetPlaceholder(placeholder string) {
	pwsz, _ := syscall.UTF16PtrFromString(placeholder)
	SendMessage(this.Handle, win32.CB_SETCUEBANNER, 1,
		unsafe.Pointer(pwsz))
}

func (this *ComboBoxObject) OnReflectCommand(msg *CommandMessage) {
	if msg.GetHwndCtrl() != this.Handle {
		//could be combolbox..
	}
	notifyCode := msg.GetNotifyCode()
	if notifyCode == uint16(win32.CBN_SELCHANGE) ||
		notifyCode == uint16(win32.CBN_EDITCHANGE) {
		this.OnValueChange.Fire(this, &SimpleEventInfo{})
	}
}

func (this *ComboBoxObject) GetValue() any {
	index := this.GetSelectedIndex()
	if index == -1 {
		return consts.Null
	}
	return this.values[index]
}

func (this *ComboBoxObject) SetValue(value any) {
	for n, v := range this.values {
		if v == value {
			this.SetSelectedIndex(n)
			return
		}
	}
	this.SetSelectedIndex(-1)
}

func (this *ComboBoxObject) GetOnValueChange() *SimpleEvent {
	return &this.OnValueChange
}

func (this *ComboBoxObject) GetText() string {
	cch, errno := SendMessage(this.Handle, win32.WM_GETTEXTLENGTH, 0, 0)
	_ = errno
	//if cch == 0 && errno != win32.NO_ERROR {
	//	log.Fatal(errno)
	//}
	buf := make([]uint16, cch+1)
	_, _ = SendMessage(this.Handle, win32.WM_GETTEXT,
		cch+1, unsafe.Pointer(&buf[0]))
	text := syscall.UTF16ToString(buf)
	return text
}

func (this *ComboBoxObject) SetText(text string) {
	pwsz, _ := syscall.UTF16PtrFromString(text)
	index, _ := SendMessage(this.Handle, win32.CB_FINDSTRINGEXACT,
		-1, unsafe.Pointer(pwsz))
	if int32(index) == win32.CB_ERR {
		bOk, errno := SetWindowText(this.Handle, text)
		if !bOk {
			log.Fatal(errno)
		}
	} else {
		this.SetSelectedIndex(int(index))
	}
}

func (this *ComboBoxObject) GetPreferredSize(int, int) (int, int) {
	height := this.Height
	if height == 0 {
		if _comboBoxHeight == 0 {
			var rc win32.RECT
			win32.GetWindowRect(this.Handle, &rc)
			_comboBoxHeight = int(rc.Bottom - rc.Top)
		}
		height = _comboBoxHeight
	}
	if _comboBoxButtonWidth == 0 {
		ret, _ := win32.GetSystemMetrics(win32.SM_CXVSCROLL)
		_comboBoxButtonWidth = int(ret)
	}
	cx, cy := 16+_comboBoxButtonWidth, height
	return cx, cy
}
