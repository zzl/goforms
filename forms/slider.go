package forms

import (
	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/framework/virtual"
	"unsafe"

	"github.com/zzl/go-win32api/v2/win32"
)

type Slider interface {
	Control

	SetRange(min int, max int)
	GetRange() (min int, max int)

	SetValue(value int)
	GetValue() int
}

type SliderObject struct {
	ControlObject
	super *ControlObject

	Vertical bool
}

type NewSlider struct {
	Parent   Container
	Name     string
	Pos      Point
	Size     Size
	Vertical bool
}

func (me NewSlider) Create(extraOpts ...*WindowOptions) Slider {
	slider := NewSliderObject()
	slider.name = me.Name
	slider.Vertical = me.Vertical

	opts := utils.OptionalArg(extraOpts)
	opts.Left = me.Pos.X
	opts.Top = me.Pos.Y

	opts.ParentHandle = resolveParentHandle(me.Parent)
	err := slider.Create(*opts)
	assertNoErr(err)
	configControlSize(slider, me.Size)

	return slider
}

func NewSliderObject() *SliderObject {
	return virtual.New[SliderObject]()
}

func (this *SliderObject) Init() {
	this.super.Init()
}

func (this *SliderObject) SetRange(min int, max int) {
	SendMessage(this.Handle, win32.TBM_SETRANGEMIN, 0, min)
	SendMessage(this.Handle, win32.TBM_SETRANGEMAX, 1, max)
	this.updateTickFreq(max)
}

func (this *SliderObject) GetRange() (int, int) {
	min, _ := SendMessage(this.Handle, win32.TBM_GETRANGEMIN, 0, 0)
	max, _ := SendMessage(this.Handle, win32.TBM_GETRANGEMAX, 0, 0)
	return int(min), int(max)
}

func (this *SliderObject) SetValue(value int) {
	SendMessage(this.Handle, win32.TBM_SETPOS, 1, value)
}

func (this *SliderObject) GetValue() int {
	ret, _ := SendMessage(this.Handle, win32.TBM_GETPOS, 0, 0)
	return int(ret)
}

func (this *SliderObject) GetWindowClass() string {
	return "msctls_trackbar32"
}

func (this *SliderObject) OnHandleCreated() {
	this.super.OnHandleCreated()

	min, _ := SendMessage(this.Handle, win32.TBM_GETRANGEMIN, 0, 0)
	max, _ := SendMessage(this.Handle, win32.TBM_GETRANGEMAX, 0, 0)
	_, _ = min, max

	if this.Vertical {
		SendMessage(this.Handle, win32.TBM_SETPOS, 1, max)
	}
	this.updateTickFreq(int(max))
}

func (this *SliderObject) updateTickFreq(max int) {
	freq := max / 10
	if freq > 1 {
		SendMessage(this.Handle, win32.TBM_SETTICFREQ, freq, 0)
	}
}

func (this *SliderObject) GetControlSpecStyle() (include, exclude WINDOW_STYLE) {
	var style WINDOW_STYLE
	style |= WINDOW_STYLE(win32.TBS_AUTOTICKS | win32.TBS_TRANSPARENTBKGND)
	if this.Vertical {
		style |= WINDOW_STYLE(win32.TBS_VERT | win32.TBS_DOWNISLEFT) //?
	} else {
		style |= WINDOW_STYLE(win32.TBS_HORZ)
	}
	return style, 0
}

func (this *SliderObject) GetPreferredSize(cxMax int, cyMax int) (int, int) {
	var cx, cy int32
	if this.Vertical {
		if this.Handle == 0 {
			cx, _ := win32.GetSystemMetrics(win32.SM_CXVSCROLL)
			cx = cx * 2
		} else {
			var rc win32.RECT
			SendMessage(this.Handle, win32.TBM_GETTHUMBRECT,
				0, unsafe.Pointer(&rc))
			cx = rc.Right + 6
		}
		return int(cx), min(cyMax, 104)
	} else {
		if this.Handle == 0 {
			cy, _ := win32.GetSystemMetrics(win32.SM_CYHSCROLL)
			cy = cy * 2
		} else {
			var rc win32.RECT
			SendMessage(this.Handle, win32.TBM_GETTHUMBRECT,
				0, unsafe.Pointer(&rc))
			cy = rc.Bottom + 6
		}
		return min(cxMax, 104), int(cy)
	}
}
