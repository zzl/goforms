package forms

import (
	"github.com/zzl/go-win32api/v2/win32"
	"log"
)

type CustomControl interface {
	Control
}

type CustomControlSpi interface {
}

type CustomControlObject struct {
	CustomWindowObject
	super *CustomWindowObject

	NameAwareSupport
}

// dup..
func (this *CustomControlObject) GetRootContainer() Container {
	return this.GetRootWindow().(Container)
}

func (this *CustomControlObject) GetDefaultStyle() WINDOW_STYLE {
	style := DefaultControlStyle
	incStyle, excStyle := this.RealObject.(Control).GetControlSpecStyle()
	style |= incStyle
	style &^= excStyle
	return style
}

func (this *CustomControlObject) CreateIn(parent Window, extraOpts ...*WindowOptions) Control {
	return createControlIn(parent, this.RealObject, extraOpts...)
}

func (this *CustomControlObject) Create(options WindowOptions) error {
	if options.ParentHandle == 0 {
		log.Println("Warning: Parent handle unassigned")
	}
	if options.ControlId == 0 {
		options.ControlId = uint16(autoControlIdGen.Gen())
	}
	creatingControlMap[int(options.ControlId)] = this.RealObject
	return this.super.Create(options)
}

func (this *CustomControlObject) GetControlSpecStyle() (include, exclude WINDOW_STYLE) {
	return 0, 0
}

func (this *CustomControlObject) GetControlId() uint16 {
	ret, errno := win32.GetDlgCtrlID(this.Handle)
	if ret == 0 {
		log.Println(errno.Error())
	}
	return uint16(ret)
}

func (this *CustomControlObject) SetDluBounds(left, top, width, height int) {
	left, top = this.DluToPx(left, top)
	width, height = this.DluToPx(width, height)
	this.RealObject.SetBounds(left, top, width, height)
}

func (this *CustomControlObject) Init() {
	this.super.Init()
}

const CONTROL_TIMER_TRACK_LEAVE = 101

func (this *CustomControlObject) GetBackgroundColor() win32.COLORREF {
	return TransparentColor
}

func (this *CustomControlObject) EnsureCustomWndProc() {
	//nop
}

var customControlClass string

func (this *CustomControlObject) EnsureClassRegistered() {
	if customControlClass != "" {
		return
	}
	customControlClass = "goforms.control"
	_, err := RegisterClass(customControlClass, nil, ClassOptions{
		BackgroundBrush: 0,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (this *CustomControlObject) WinProc(win *WindowObject, m *Message) error {
	return this.super.WinProc(win, m)
}

func (this *CustomControlObject) GetWindowClass() string {
	return customControlClass
}
