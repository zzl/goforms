package forms

import (
	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/framework/virtual"
	"math"
	"syscall"
	"unsafe"

	"github.com/zzl/go-win32api/v2/win32"
)

type Edit interface {
	Control
	Input
	TextAware

	SetPlaceholder(placeholder string)
	GetPlaceholder() string
	SetReadonly(readonly bool)
	SelectAll()
	SelectEnd()

	IsEmpty() bool

	EditObj() *EditObject
}

type EditSpi interface {
	ControlSpi
}

type EditInterface interface {
	Edit
	EditSpi
}

type EditObject struct {
	ControlObject
	super *ControlObject

	ForeColorAwareSupport
	OnValueChange SimpleEvent
}

func (this *EditObject) EditObj() *EditObject {
	return this
}

func (this *EditObject) GetWindowClass() string {
	return "Edit"
}

func NewEditObject() *EditObject {
	return virtual.New[EditObject]()
}

type NewEdit struct {
	Parent   Container
	Name     string
	Text     string
	Pos      Point
	Size     Size
	OnChange SimpleEventListener
}

func (me NewEdit) Create(extraOpts ...*WindowOptions) Edit {
	edit := NewEditObject()
	edit.name = me.Name

	opts := utils.OptionalArg(extraOpts)
	opts.WindowName = me.Text
	opts.Left = me.Pos.X
	opts.Top = me.Pos.Y

	opts.ParentHandle = resolveParentHandle(me.Parent)
	err := edit.Create(*opts)
	assertNoErr(err)
	configControlSize(edit, me.Size)

	if me.OnChange != nil {
		edit.OnValueChange.AddListener(me.OnChange)
	}
	return edit
}

func (this *EditObject) Init() {
	this.super.Init()
}

func (this *EditObject) Dispose() {
	oldHbr := this.GetData(Data_BackColorBrush)
	if oldHbr != nil {
		win32.DeleteObject(oldHbr.(win32.HBRUSH))
	}
	this.super.Dispose()
}

func (this *EditObject) OnReflectCommand(info *CommandMessage) {
	this.super.OnReflectCommand(info)
	if info.GetNotifyCode() == uint16(win32.EN_CHANGE) {
		this.OnValueChange.Fire(this, &SimpleEventInfo{})
	}
}

func (this *EditObject) OnReflectMessage(msg *Message) {
	if msg.UMsg == win32.WM_CTLCOLOREDIT {
		backColor := this.GetBackColor()
		hbr := handleCtlColor(this, msg.WParam, backColor)
		if hbr != 0 {
			msg.SetHandledWithResult(hbr)
			return
		}
		//
		//msg.Result = hbr
		//msg.Handled = false
		//msg.SetHandledWithResult(0)
		//return
		//if backColor != colors.Null && backColor != colors.Transparent {
		//hdc := win32.HDC(msg.WParam)
		//win32.SetBkMode(hdc, win32.TRANSPARENT)
		//oldHbr := this.GetData(Data_BackColorBrush)
		//if oldHbr != nil {
		//	win32.DeleteObject(oldHbr.(win32.HBRUSH))
		//}
		//hbr := win32.CreateSolidBrush(win32.RGB(255, 255, 0))
		//this.SetData(Data_BackColorBrush, hbr)
		//msg.Result = hbr
		//msg.Handled = true
		return
		//}
		//utils
	}
	this.super.OnReflectMessage(msg)
}

func (this *EditObject) GetDefaultExStyle() WINDOW_EX_STYLE {
	return win32.WS_EX_CLIENTEDGE
}

func (this *EditObject) SetPlaceholder(placeholder string) {
	pwsz, _ := syscall.UTF16PtrFromString(placeholder)
	SendMessage(this.Handle, win32.EM_SETCUEBANNER, 1,
		unsafe.Pointer(pwsz))
}

func (this *EditObject) GetPlaceholder() string {
	buf := make([]uint16, win32.MAX_PATH)
	bOk, _ := SendMessage(this.Handle, win32.EM_GETCUEBANNER,
		win32.WPARAM(unsafe.Pointer(&buf[0])), win32.MAX_PATH)
	if bOk == 1 {
		return syscall.UTF16ToString(buf)
	}
	return ""
}

func (this *EditObject) GetValue() any {
	return this.GetText()
}

func (this *EditObject) SetValue(value any) {
	text := utils.ToString(value)
	this.SetText(text)
}

func (this *EditObject) GetOnValueChange() *SimpleEvent {
	return &this.OnValueChange
}

func (this *EditObject) GetText() string {
	text, _ := GetWindowText(this.Handle)
	return text
}

func (this *EditObject) SetText(text string) {
	SetWindowText(this.Handle, text)
}

func (this *EditObject) GetTextLength() int {
	cch, _ := win32.GetWindowTextLength(this.Handle)
	return int(cch)
}

func (this *EditObject) IsEmpty() bool {
	return this.GetTextLength() == 0
}

func (this *EditObject) GetTextByteLength() int {
	cch, _ := win32.GetWindowTextLengthA(this.Handle)
	return int(cch)
}

func (this *EditObject) GetPreferredSize(int, int) (int, int) {
	cx, cy := MeasureText2(this.Handle, "|Why", this.GetText())

	ratio := float64(cy) / 16
	cy = int(math.Round(1.625*14*ratio + 1.5))
	return cx, cy
}

func (this *EditObject) SetReadonly(readonly bool) {
	SendMessage(this.Handle, win32.EM_SETREADONLY,
		win32.BoolToBOOL(readonly), 0)
}

func (this *EditObject) SelectAll() {
	SendMessage(this.Handle, win32.EM_SETSEL, 0, NegativeOne)
}

func (this *EditObject) SelectEnd() {
	pos := len(this.GetText())
	SendMessage(this.Handle, win32.EM_SETSEL, pos, pos)
}
