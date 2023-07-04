package forms

import (
	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/framework/virtual"
)

type ProgressBar interface {
	Control

	ProgressBarObj() *ProgressBarObject
}

type ProgressBarObject struct {
	ControlObject
	super *ControlObject

	min      int
	max      int
	progress int
}

type NewProgressBar struct {
	Parent Container
	Name   string
	Pos    Point
	Size   Size
	Min    int
	Max    int
	Value  int
}

func (me NewProgressBar) Create(extraOpts ...*WindowOptions) ProgressBar {
	bar := NewProgressBarObject()
	bar.name = me.Name
	bar.min = me.Min
	bar.max = me.Max
	bar.progress = me.Value

	opts := utils.OptionalArg(extraOpts)
	opts.Left = me.Pos.X
	opts.Top = me.Pos.Y

	opts.ParentHandle = resolveParentHandle(me.Parent)
	err := bar.Create(*opts)
	assertNoErr(err)
	configControlSize(bar, me.Size)

	return bar
}

func NewProgressBarObject() *ProgressBarObject {
	return virtual.New[ProgressBarObject]()
}

func (this *ProgressBarObject) ProgressBarObj() *ProgressBarObject {
	return this
}

func (this *ProgressBarObject) Init() {
	this.super.Init()
	this.max = 100
}

func (this *ProgressBarObject) GetWindowClass() string {
	return "msctls_progress32"
}

func (this *ProgressBarObject) Create(options WindowOptions) error {
	err := this.super.Create(options)
	if this.min != 0 || this.max != 100 {
		this.SetRange(this.min, this.max)
	}
	if this.progress != 0 {
		this.SetProgress(this.progress)
	}
	return err
}

func (this *ProgressBarObject) SetMin(min int) {
	this.SetRange(min, this.max)
}

func (this *ProgressBarObject) SetMax(max int) {
	this.SetRange(this.min, max)
}

func (this *ProgressBarObject) SetRange(min, max int) {
	this.min = min
	this.max = max
	if this.Handle != 0 {
		SendMessage(this.Handle, win32.PBM_SETRANGE32,
			this.min, this.max)
	}
}

func (this *ProgressBarObject) SetPos_(pos int) {
	this.SetProgress(pos)
}

func (this *ProgressBarObject) SetProgress(progress int) {
	this.progress = progress
	SendMessage(this.Handle, win32.PBM_SETPOS,
		this.progress, 0)
}

func (this *ProgressBarObject) ShowMarquee(marquee bool) {
	style, _ := win32.GetWindowLong(this.Handle, win32.GWL_STYLE)
	var setMarqueeWparam uintptr
	if marquee {
		style |= int32(win32.PBS_MARQUEE)
		setMarqueeWparam = 1
	} else {
		style &^= int32(win32.PBS_MARQUEE)
		setMarqueeWparam = 0
	}
	win32.SetWindowLong(this.Handle, win32.GWL_STYLE, style)
	SendMessage(this.Handle, win32.PBM_SETMARQUEE, setMarqueeWparam, 0)
}

func (this *ProgressBarObject) GetPreferredSize(cxMax int, cyMax int) (int, int) {
	ret, _ := win32.GetSystemMetrics(win32.SM_CYVSCROLL)
	return min(64, cxMax), min(int(ret), cyMax)
}
