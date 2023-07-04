package forms

import (
	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/framework/virtual"
	"log"
	"unsafe"

	"github.com/zzl/go-win32api/v2/win32"
)

type RadioButton interface {
	Control
	TextAware

	GetChecked() bool
	SetChecked(checked bool)

	GetOnClick() *SimpleEvent

	RadioButtonObj() *RadioButtonObject
}

type RadioButtonSpi interface {
	ControlSpi
}

type RadioButtonInterface interface {
	RadioButton
	RadioButtonSpi
}

type RadioButtonObject struct {
	ControlObject

	ForeColorAwareSupport
	AutoCheck  bool
	BeginGroup bool

	OnValueChange SimpleEvent
	OnClick       SimpleEvent
}

func NewRadioButtonObject() *RadioButtonObject {
	return virtual.New(&RadioButtonObject{
		AutoCheck: true,
	})
}

type NewRadioButton struct {
	Parent     Container
	Name       string
	Text       string
	Pos        Point
	Size       Size
	Checked    bool
	Disabled   bool
	BeginGroup bool
}

func (me NewRadioButton) Create(extraOpts ...*WindowOptions) RadioButton {
	radioButton := NewRadioButtonObject()
	radioButton.name = me.Name
	radioButton.BeginGroup = me.BeginGroup

	opts := utils.OptionalArg(extraOpts)
	opts.WindowName = me.Text
	opts.Left = me.Pos.X
	opts.Top = me.Pos.Y
	if me.Disabled {
		opts.StyleInclude |= win32.WS_DISABLED
	}

	opts.ParentHandle = resolveParentHandle(me.Parent)
	err := radioButton.Create(*opts)
	assertNoErr(err)

	configControlSize(radioButton, me.Size)
	if me.Checked {
		radioButton.SetChecked(true)
	}
	return radioButton
}

func (this *RadioButtonObject) RadioButtonObj() *RadioButtonObject {
	return this
}

func (this *RadioButtonObject) GetOnClick() *SimpleEvent {
	return &this.OnClick
}

func (this *RadioButtonObject) GetChecked() bool {
	ret, _ := SendMessage(this.Handle, win32.BM_GETCHECK, 0, 0)
	return win32.DLG_BUTTON_CHECK_STATE(ret) == win32.BST_CHECKED
}

func (this *RadioButtonObject) SetChecked(checked bool) {
	wParam := win32.BST_UNCHECKED
	if checked {
		wParam = win32.BST_CHECKED
	}
	SendMessage(this.Handle, win32.BM_SETCHECK, WPARAM(wParam), 0)
}

func (this *RadioButtonObject) GetValue() any {
	return this.GetChecked()
}

func (this *RadioButtonObject) SetValue(value any) {
	this.SetChecked(value.(bool))
}

func (this *RadioButtonObject) GetOnValueChangeEvent() *SimpleEvent {
	return &this.OnValueChange
}

func (this *RadioButtonObject) SetText(text string) {
	SetWindowText(this.Handle, text)
}

func (this *RadioButtonObject) GetText() string {
	text, _ := GetWindowText(this.Handle)
	return text
}

func (this *RadioButtonObject) GetWindowClass() string {
	return "BUTTON"
}

func (this *RadioButtonObject) GetControlSpecStyle() (include, exclude WINDOW_STYLE) {
	if this.AutoCheck {
		include = WINDOW_STYLE(win32.BS_AUTORADIOBUTTON)
	} else {
		include = WINDOW_STYLE(win32.BS_RADIOBUTTON)
	}
	if this.BeginGroup {
		include |= win32.WS_GROUP
	}
	return include, win32.WS_TABSTOP
}

func (this *RadioButtonObject) OnReflectCommand(msg *CommandMessage) {
	notifyCode := msg.GetNotifyCode()
	if notifyCode == uint16(win32.BN_CLICKED) {
		this.OnClick.Fire(this, &SimpleEventInfo{})
		//this.OnValueChange.Fire(this, EventInfo{}) //?
	}
}

func (this *RadioButtonObject) GetPreferredSize(int, int) (cx, cy int) {
	var sz win32.SIZE
	bOk, _ := SendMessage(this.Handle, win32.BCM_GETIDEALSIZE, 0,
		unsafe.Pointer(&sz))
	if bOk == 0 {
		log.Fatal("?")
	}
	return int(sz.Cx + 16), int(sz.Cy)
}
