package forms

import (
	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/framework/virtual"
	"log"
	"syscall"
	"unsafe"

	"github.com/zzl/go-win32api/v2/win32"
)

type StatusPart struct {
	Width    int
	MinWidth int
	MaxWidth int
	Weight   float32
	Text     string
}

type StatusBar interface {
	Control
	IsSimpleMode() bool
	SetSimpleMode(simple bool)
	SetParts(parts []StatusPart)

	SetSimpleText(text string)
	GetSimpleText() string

	SetPartText(index int, text string)
	GetPartText(index int) string

	SetText(text string)
	GetText() string
}

type StatusBarSpi interface {
	ControlSpi
}

type StatusBarInterface interface {
	StatusBar
	StatusBarSpi
}

type StatusBarObject struct {
	ControlObject
	super *ControlObject

	parts []*StatusPart
}

type NewStatusBar struct {
	Parent Container
	Name   string
	Simple bool
	Parts  []*StatusPart
}

func (me NewStatusBar) Create(extraOpts ...*WindowOptions) StatusBar {
	statusBar := NewStatusBarObject()
	statusBar.name = me.Name
	statusBar.parts = me.Parts

	opts := utils.OptionalArg(extraOpts)
	opts.ParentHandle = resolveParentHandle(me.Parent)
	err := statusBar.Create(*opts)
	assertNoErr(err)

	if me.Simple {
		statusBar.SetSimpleMode(true)
	}
	return statusBar
}

func NewStatusBarObject() *StatusBarObject {
	return virtual.New[StatusBarObject]()
}

func (this *StatusBarObject) Init() {
	this.super.Init()
}

func (this *StatusBarObject) GetWindowClass() string {
	return "msctls_statusbar32"
}

func (this *StatusBarObject) Create(options WindowOptions) error {
	options.Style = win32.WS_CHILD | win32.WS_VISIBLE | WINDOW_STYLE(win32.SBARS_SIZEGRIP)

	err := this.super.Create(options)
	return err
}

func (this *StatusBarObject) OnParentResized() {
	this.super.OnParentResized()
	SendMessage(this.Handle, win32.WM_SIZE, 0, 0)
	this._setParts()
}

func (this *StatusBarObject) GetPreferredSize(cxMax int, cyMax int) (int, int) {
	var rc win32.RECT
	win32.GetClientRect(this.Handle, &rc)
	cy := rc.Bottom - rc.Top
	return cxMax, int(cy)
}

func (this *StatusBarObject) _setParts() {
	sumWeight := float32(0)
	fixedWidth := 0
	partCount := len(this.parts)
	if partCount == 0 {
		return
	}
	widths := make([]int, partCount)
	for n, part := range this.parts {
		if part.Width == 0 {
			sumWeight += part.Weight
		} else {
			widths[n] = part.Width
			fixedWidth += part.Width
		}
	}
	var rc win32.RECT
	win32.GetClientRect(this.Handle, &rc)
	totalWidth := int(rc.Right)
	flexWidth := totalWidth - fixedWidth
	for n, part := range this.parts {
		if part.Width == 0 {
			ratio := float64(part.Weight) / float64(sumWeight)
			width := int(float64(flexWidth) * ratio)
			widths[n] = width
		}
	}
	rights := make([]int32, partCount)
	sumWidth := 0
	for n, width := range widths {
		sumWidth += width
		rights[n] = int32(sumWidth)
	}
	rights[len(rights)-1] = -1
	ret, errno := SendMessage(this.Handle, win32.SB_SETPARTS,
		partCount, unsafe.Pointer(&rights[0]))
	if ret == 0 {
		log.Fatal(errno)
	}
}

func (this *StatusBarObject) SetParts(parts []StatusPart) {

	this.parts = nil
	for n := range parts {
		part := &parts[n]
		if part.Width == 0 && part.Weight < 0.001 {
			if part.Text != "" {
				part.Width, _ = MeasureText(this.Handle, part.Text)
			}
			part.Width += 8
		}
		this.parts = append(this.parts, part)
	}
	this._setParts()
}

func (this *StatusBarObject) IsSimpleMode() bool {
	ret, _ := SendMessage(this.Handle, win32.SB_ISSIMPLE, 0, 0)
	return ret != 0
}

func (this *StatusBarObject) SetSimpleMode(simple bool) {
	var wParam WPARAM = 0
	if simple {
		wParam = 1
	}
	SendMessage(this.Handle, win32.SB_SIMPLE, wParam, 0)
}

func (this *StatusBarObject) GetSimpleText() string {
	simpleMode := this.IsSimpleMode()
	if !simpleMode {
		this.SetSimpleMode(true)
	}
	text := this._getPartText(0)
	if !simpleMode {
		this.SetSimpleMode(false)
	}
	return text
}

func (this *StatusBarObject) GetPartText(index int) string {
	simpleMode := this.IsSimpleMode()
	if simpleMode {
		this.SetSimpleMode(false)
	}
	text := this._getPartText(index)
	if simpleMode {
		this.SetSimpleMode(true)
	}
	return text
}

func (this *StatusBarObject) GetText() string {
	return this._getPartText(0)
}

func (this *StatusBarObject) SetText(text string) {
	var index = 0
	if this.IsSimpleMode() {
		index = int(win32.SB_SIMPLEID)
	}
	this.SetPartText(index, text)
}

func (this *StatusBarObject) _getPartText(index int) string {
	ret, _ := SendMessage(this.Handle, win32.SB_GETTEXTLENGTH,
		index, 0)
	cch := win32.LOWORD(win32.DWORD(ret))
	if cch == 0 {
		return ""
	}

	buf := make([]uint16, cch+1)
	ret, _ = SendMessage(this.Handle, win32.SB_GETTEXT,
		index, unsafe.Pointer(&buf[0]))

	text := syscall.UTF16ToString(buf)
	return text
}

func (this *StatusBarObject) SetSimpleText(text string) {
	//check in simple mode?
	this.SetPartText(int(win32.SB_SIMPLEID), text)
}

func (this *StatusBarObject) SetPartText(index int, text string) {
	pwsz, _ := syscall.UTF16PtrFromString(text)

	style := 0

	ret, errno := SendMessage(this.Handle, win32.SB_SETTEXT,
		win32.MAKELONG(uint16(index), uint16(style)),
		unsafe.Pointer(pwsz))
	if ret == 0 {
		log.Fatal(errno)
	}
}
