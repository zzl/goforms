package forms

import (
	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/framework/virtual"
)

type Label interface {
	Control
	TextAware
	BackColorAware
	LabelObject() *LabelObject
}

type LabelSpi interface {
	ControlSpi
}

type LabelInterface interface {
	Label
	LabelSpi
}

type LabelObject struct {
	ControlObject
	super *ControlObject

	ForeColorAwareSupport
}

func NewLabelObject() *LabelObject {
	return virtual.New[LabelObject]()
}

type NewLabel struct {
	Parent Container
	Name   string
	Text   string
	Pos    Point
	Size   Size
}

func (me NewLabel) Create(extraOpts ...*WindowOptions) Label {
	label := NewLabelObject()
	label.name = me.Name

	opts := utils.OptionalArg(extraOpts)
	opts.WindowName = me.Text
	opts.Left = me.Pos.X
	opts.Top = me.Pos.Y

	opts.ParentHandle = resolveParentHandle(me.Parent)
	err := label.Create(*opts)
	assertNoErr(err)
	configControlSize(label, me.Size)

	return label
}

func (this *LabelObject) LabelObject() *LabelObject {
	return this
}

func (this *LabelObject) GetWindowClass() string {
	return "Static"
}

func (this *LabelObject) GetText() string {
	text, _ := GetWindowText(this.Handle)
	return text
}

func (this *LabelObject) SetText(text string) {
	SetWindowText(this.Handle, text)
}

func (this *LabelObject) GetPreferredSize(int, int) (cx, cy int) {
	return MeasureText(this.Handle, this.GetText())
}

func (this *LabelObject) OnReflectMessage(msg *Message) {
	if msg.UMsg == win32.WM_CTLCOLORSTATIC {
		backColor := this.GetBackColor()
		hdc := msg.WParam

		hbr := handleCtlColor(this, hdc, backColor)
		if hbr != 0 {
			msg.SetHandledWithResult(hbr)
		}
	}
}
