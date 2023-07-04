package forms

import (
	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/framework/virtual"
	"unsafe"

	"github.com/zzl/go-win32api/v2/win32"
)

type LinkLabel interface {
	Control
	TextAware

	GetOnClick() *Event[*LinkLabelOnClickInfo]
}

type LinkLabelOnClickInfo struct {
	SimpleEventInfo
	Index int
	Id    string
	Url   string
}

type LinkLabelObject struct {
	ControlObject
	super *ControlObject

	OnClick Event[*LinkLabelOnClickInfo]
	Action  Action
}

type NewLinkLabel struct {
	Parent Container
	Name   string
	Text   string
	Pos    Point
	Size   Size
	Action Action
}

func (me NewLinkLabel) Create(extraOpts ...*WindowOptions) LinkLabel {
	linkLabel := NewLinkLabelObject()
	linkLabel.name = me.Name

	opts := utils.OptionalArg(extraOpts)
	opts.WindowName = me.Text
	opts.Left = me.Pos.X
	opts.Top = me.Pos.Y

	opts.ParentHandle = resolveParentHandle(me.Parent)
	err := linkLabel.Create(*opts)
	assertNoErr(err)
	configControlSize(linkLabel, me.Size)

	linkLabel.Action = me.Action
	return linkLabel
}

func NewLinkLabelObject() *LinkLabelObject {
	return virtual.New[LinkLabelObject]()
}

func (this *LinkLabelObject) Init() {
	this.super.Init()
}

func (this *LinkLabelObject) SetText(text string) {
	SetWindowText(this.Handle, text)
}

func (this *LinkLabelObject) GetText() string {
	text, _ := GetWindowText(this.Handle)
	return text
}

func (this *LinkLabelObject) GetWindowClass() string {
	return "SysLink"
}

func (this *LinkLabelObject) GetOnClick() *Event[*LinkLabelOnClickInfo] {
	return &this.OnClick
}

func (this *LinkLabelObject) GetControlSpecStyle() (WINDOW_STYLE, WINDOW_STYLE) {
	return WINDOW_STYLE(win32.LWS_TRANSPARENT), 0
}

func (this *LinkLabelObject) GetPreferredSize(cxMax int, cyMax int) (int, int) {
	var size win32.SIZE
	SendMessage(this.Handle, win32.LM_GETIDEALSIZE,
		cxMax, unsafe.Pointer(&size))
	return int(size.Cx), int(size.Cy)
}

func (this *LinkLabelObject) OnReflectNotify(msg *NotifyMessage) {
	code := msg.GetNMHDR().Code
	if code == win32.NM_CLICK || code == win32.NM_RETURN {
		pNmlink := (*win32.NMLINK)(unsafe.Pointer(msg.LParam))
		item := pNmlink.Item
		if this.Action != nil {
			this.Action()
		}
		this.OnClick.Fire(this.RealObject, &LinkLabelOnClickInfo{
			Index: int(item.ILink),
			Id:    win32.WstrToStr(item.SzID[:]),
			Url:   win32.WstrToStr(item.SzUrl[:]),
		},
		)
	}
}
