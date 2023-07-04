package forms

import (
	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/framework/virtual"
	"unsafe"

	"github.com/zzl/go-win32api/v2/win32"
)

type CheckBox interface {
	Control
	Input
	TextAware

	GetChecked() bool
	SetChecked(checked bool)

	CheckBoxObj() *CheckBoxObject
}

type CheckBoxSpi interface {
	ControlSpi
}

type CheckBoxInterface interface {
	CheckBox
	CheckBoxSpi
}

type CheckBoxObject struct {
	ControlObject
	super *ControlObject

	ForeColorAwareSupport
	AutoCheck bool

	OnValueChange SimpleEvent
}

func NewCheckBoxObject() *CheckBoxObject {
	return virtual.New[CheckBoxObject]()
}

type NewCheckBox struct {
	Parent   Container
	Name     string
	Text     string
	Pos      Point
	Size     Size
	Checked  bool
	Disabled bool
}

func (me NewCheckBox) Create(extraOpts ...*WindowOptions) CheckBox {
	chk := NewCheckBoxObject()
	chk.name = me.Name

	opts := utils.OptionalArg(extraOpts)
	opts.WindowName = me.Text
	opts.Left = me.Pos.X
	opts.Top = me.Pos.Y
	opts.StyleInclude = utils.If(me.Disabled, win32.WS_DISABLED)
	opts.ParentHandle = resolveParentHandle(me.Parent)
	err := chk.Create(*opts)
	assertNoErr(err)
	configControlSize(chk, me.Size)
	if me.Checked {
		chk.SetChecked(true)
	}
	return chk
}

func (this *CheckBoxObject) CheckBoxObj() *CheckBoxObject {
	return this
}

func (this *CheckBoxObject) SetText(text string) {
	SetWindowText(this.Handle, text)
}

func (this *CheckBoxObject) GetText() string {
	text, _ := GetWindowText(this.Handle)
	return text
}

func (this *CheckBoxObject) GetOnValueChange() *SimpleEvent {
	return &this.OnValueChange
}

func (this *CheckBoxObject) GetChecked() bool {
	ret, _ := SendMessage(this.Handle, win32.BM_GETCHECK, 0, 0)
	return ret == uintptr(win32.BST_CHECKED)
}

func (this *CheckBoxObject) SetChecked(checked bool) {
	wParam := win32.BST_UNCHECKED
	if checked {
		wParam = win32.BST_CHECKED
	}
	SendMessage(this.Handle, win32.BM_SETCHECK, WPARAM(wParam), 0)
}

func (this *CheckBoxObject) GetValue() any {
	return this.GetChecked()
}

func (this *CheckBoxObject) SetValue(value any) {
	this.SetChecked(value.(bool))
}

func (this *CheckBoxObject) GetWindowClass() string {
	return "BUTTON"
}

func (this *CheckBoxObject) GetControlSpecStyle() (include WINDOW_STYLE, exclude WINDOW_STYLE) {
	if this.AutoCheck {
		include = WINDOW_STYLE(win32.BS_AUTOCHECKBOX)
	} else {
		include = WINDOW_STYLE(win32.BS_CHECKBOX)
	}
	return include, 0
}

func (this *CheckBoxObject) Init() {
	this.super.Init()
	this.AutoCheck = true
}

func (this *CheckBoxObject) OnReflectCommand(msg *CommandMessage) {
	this.super.OnReflectCommand(msg)
	if msg.GetNotifyCode() == uint16(win32.BN_CLICKED) {
		this.OnValueChange.Fire(this, &SimpleEventInfo{})
	}
}

func (this *CheckBoxObject) GetPreferredSize(int, int) (cx, cy int) {
	var sz win32.SIZE
	SendMessage(this.Handle, win32.BCM_GETIDEALSIZE, 0,
		unsafe.Pointer(&sz))
	return int(sz.Cx + 16), int(sz.Cy)
}
