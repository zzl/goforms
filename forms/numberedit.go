package forms

import (
	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/framework/virtual"
	"log"
	"math"
)

type NumberEdit interface {
	Edit
}

type NumberEditObject struct {
	EditObject
	super *EditObject

	//
	Min   int
	Max   int
	udWin *WindowObject
}

type NewNumberEdit struct {
	Parent Container
	Name   string
	Pos    Point
	Size   Size
	Min    int
	Max    int
	Value  int
}

func (me NewNumberEdit) Create(extraOpts ...*WindowOptions) NumberEdit {
	edit := NewNumberEditObject()
	edit.name = me.Name
	edit.Min = me.Min
	edit.Max = me.Max

	opts := utils.OptionalArg(extraOpts)
	opts.WindowName = utils.ToString(me.Value)
	opts.Left = me.Pos.X
	opts.Top = me.Pos.Y

	opts.ParentHandle = resolveParentHandle(me.Parent)
	err := edit.Create(*opts)
	assertNoErr(err)
	configControlSize(edit, me.Size)

	return edit
}

func NewNumberEditObject() *NumberEditObject {
	return virtual.New[NumberEditObject]()
}

func (this *NumberEditObject) GetWindowClass() string {
	return "Edit"
}

func (this *NumberEditObject) Init() {
	this.super.Init()

	this.Min = 0
	this.Max = math.MaxInt16
}

func (this *NumberEditObject) Create(options WindowOptions) error {
	options.Style |= WINDOW_STYLE(win32.ES_NUMBER)
	err := this.super.Create(options)
	if err != nil {
		return err
	}

	win := NewWindowObject()
	//win.Impl = this //todo:??
	if err := win.Create(WindowOptions{
		ClassName: "msctls_updown32",
		Style: win32.WS_CHILDWINDOW | win32.WS_VISIBLE |
			WINDOW_STYLE(win32.UDS_AUTOBUDDY|win32.UDS_ALIGNRIGHT|win32.UDS_SETBUDDYINT|
				win32.UDS_ARROWKEYS|
				win32.UDS_HOTTRACK),
		ExStyle:      win32.WS_EX_LEFT | win32.WS_EX_LTRREADING,
		Width:        16,
		Height:       16,
		ParentHandle: this.GetContainer().GetHandle(),
	}); err != nil {
		log.Fatal(err)
	}
	this.udWin = win
	this.updateUdPos()
	SendMessage(win.Handle, win32.UDM_SETRANGE32, this.Min, this.Max)
	//win.SetBoundsRect(Rect{380, 56, 396, 80})
	return err
}

// ?
func (this *NumberEditObject) updateUdPos() {
	vsWidth, _ := win32.GetSystemMetrics(win32.SM_CXVSCROLL)
	bounds := this.GetBounds()
	bounds.Left = bounds.Right - 1
	bounds.Right = bounds.Left + int(vsWidth+1)
	this.udWin.SetBounds(bounds.Left, bounds.Top, bounds.Width(), bounds.Height())
}

func (this *NumberEditObject) SetBounds(left, top, width, height int) {
	vsWidth, _ := win32.GetSystemMetrics(win32.SM_CXVSCROLL)
	width -= int(vsWidth)
	this.super.SetBounds(left, top, width, height)
	this.updateUdPos()
}

func (this *NumberEditObject) GetValue() any {
	return this.GetText()
}

func (this *NumberEditObject) SetValue(value any) {
	text := utils.ToString(value)
	this.SetText(text)
}
