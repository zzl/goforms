package forms

import (
	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/framework/virtual"
	"unsafe"

	"github.com/zzl/go-win32api/v2/win32"
)

type GroupBox interface {
	Control

	TextAware

	//SetLayout(layout Layout)

	GroupBoxObj() *GroupBoxObject
}

type GroupBoxSpi interface {
	ControlSpi
}

type GroupBoxInterface interface {
	GroupBox
	GroupBoxSpi
}

type GroupBoxObject struct {
	ControlObject
	super *ControlObject

	Caption string

	Layout Layout
}

func NewGroupBoxObject() *GroupBoxObject {
	return virtual.New[GroupBoxObject]()
}

func (this *GroupBoxObject) GroupBoxObj() *GroupBoxObject {
	return this
}

type NewGroupBox struct {
	Parent Container
	Name   string
	Text   string
	Pos    Point
	Size   Size
}

func (me NewGroupBox) Create(extraOpts ...*WindowOptions) GroupBox {
	groupBox := NewGroupBoxObject()
	groupBox.name = me.Name

	opts := utils.OptionalArg(extraOpts)
	opts.WindowName = me.Text
	opts.Left = me.Pos.X
	opts.Top = me.Pos.Y
	opts.ParentHandle = resolveParentHandle(me.Parent)

	err := groupBox.Create(*opts)
	assertNoErr(err)
	configControlSize(groupBox, me.Size)
	return groupBox
}

func (this *GroupBoxObject) Init() {
	this.super.Init()
}

func (this *GroupBoxObject) PreCreate(opts *WindowOptions) {
	this.super.PreCreate(opts)
	if opts.WindowName == "" {
		opts.WindowName = this.Caption
	}
}

//func (this *GroupBoxObject) SetLayout(layout Layout) {
//	this.Layout = layout
//	layout.SetContainer(this.GetContainer())
//}

func (this *GroupBoxObject) GetWindowClass() string {
	return "Button"
}

func (this *GroupBoxObject) SetText(text string) {
	SetWindowText(this.Handle, text)
}

func (this *GroupBoxObject) GetText() string {
	text, _ := GetWindowText(this.Handle)
	return text
}

func (this *GroupBoxObject) Create(options WindowOptions) error {
	options.Style = win32.WS_CHILDWINDOW | win32.WS_VISIBLE |
		win32.WS_GROUP | WINDOW_STYLE(win32.BS_GROUPBOX) // | win32.WS_CLIPSIBLINGS
	return this.super.Create(options)
}

func (this *GroupBoxObject) GetPreferredSize(maxWidth int, maxHeight int) (cx, cy int) {
	if this.Layout != nil {
		cx, cy := this.Layout.GetPreferredSize(maxWidth, maxHeight)
		return cx, cy
	}
	var sz win32.SIZE
	SendMessage(this.Handle, win32.BCM_GETIDEALSIZE, 0, unsafe.Pointer(&sz))
	return int(sz.Cx + 16), int(sz.Cy)
}

func (this *GroupBoxObject) SetBounds(left, top, width, height int) {
	this.super.SetBounds(left, top, width, height)
	if this.Layout == nil {
		return
	}
	this.Layout.SetBounds(left, top, width, height)
}

func (this *GroupBoxObject) OnReflectMessage(msg *Message) {
	//if msg.UMsg == win32.WM_CTLCOLORSTATIC {
	//	hdc := msg.WParam
	//	hbr := handleCtlColor(this.RealObject, hdc, colors.Transparent)
	//	if hbr != 0 {
	//		msg.Result = hbr
	//		msg.Handled = true
	//	}
	//}
}
