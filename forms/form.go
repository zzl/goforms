package forms

import (
	"github.com/zzl/goforms/drawing"
	"github.com/zzl/goforms/drawing/colors"
	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/framework/virtual"
	"log"

	"github.com/zzl/go-win32api/v2/win32"
)

// Form is an interface that represents a form window,
// which are designed specifically for creating application windows.
type Form interface {
	TopWindow // the parent interface

	AsFormObject() *FormObject // returns the underlying FormObject
}

// FormSpi is an interface that provides additional methods
// specific to implementing a Form.
type FormSpi interface {
	TopWindowSpi
}

// FormInterface is a composition of Form and FormSpi
type FormInterface interface {
	Form
	FormSpi
}

// FormObject implements the FormInterface
// It extends TopWindowObject.
type FormObject struct {

	// TopWindowObject is the parent struct.
	TopWindowObject

	// super is the special pointer to the parent struct.
	super *TopWindowObject

	sizeGrip bool // whether shows a size grip at the bottom right corner
}

var _form_class_registerd = false

func (this *FormObject) Init() {
	this.super.Init()
}

func (this *FormObject) AsFormObject() *FormObject {
	return this
}

func (this *FormObject) EnsureClassRegistered() {
	if _form_class_registerd {
		return
	}
	_, err := RegisterClass("goforms.form", nil, ClassOptions{
		BackgroundBrush: ToSysColorBrush(byte(win32.COLOR_3DFACE)),
		Style:           win32.CS_HREDRAW | win32.CS_VREDRAW,
	})
	if err != nil {
		log.Fatal(err)
	}
	_form_class_registerd = true
}

func (this *FormObject) GetWindowClass() string {
	return "goforms.form"
}

func (this *FormObject) GetDefaultStyle() WINDOW_STYLE {
	return win32.WS_OVERLAPPEDWINDOW
}

func (this *FormObject) GetDefaultExStyle() WINDOW_EX_STYLE {
	return win32.WS_EX_CONTROLPARENT
}

func (this *FormObject) PreCreate(opts *WindowOptions) {
	this.super.PreCreate(opts)
}

func (this *FormObject) OnHandleCreated() {
	this.super.OnHandleCreated()
	if this.sizeGrip {
		this.GetEvent(win32.WM_NCHITTEST).AddListener(func(msg *Message) {
			hitTestSizeGrip(this.Handle, msg)
		})
	}
}

func (this *FormObject) OnEraseBkgnd(hdc win32.HDC) bool {
	result := this.super.OnEraseBkgnd(hdc)
	if this.sizeGrip {
		drawSizeGrip(this.Handle, hdc)
	}
	return result
}

func NewFormObject() *FormObject {
	return virtual.New[FormObject]()
}

type WindowCenter int

const (
	CenterNone   WindowCenter = 0
	CenterScreen WindowCenter = 1
	CenterParent WindowCenter = 2
)

type WindowState int

const (
	StateNormal    WindowState = 0
	StateMaximized WindowState = 1
	StateMinimized WindowState = 2
)

type NewForm struct {
	Title      string
	Pos        Point
	Size       Size
	ClientSize Size
	BackColor  drawing.Color

	Icon    win32.HICON
	IconBig win32.HICON

	NoResize      bool
	NoMinimizeBox bool
	NoMaximizeBox bool
	NoIcon        bool
	HelpButton    bool

	Owner    TopWindow
	Center   WindowCenter
	State    WindowState
	TopMost  bool
	SizeGrip bool
}

func (me NewForm) Create(extraOpts ...*WindowOptions) Form {
	form := NewFormObject()

	opts := utils.OptionalArg(extraOpts)
	opts.WindowName = me.Title
	opts.Left = me.Pos.X
	opts.Top = me.Pos.Y

	if opts.Style == 0 {
		opts.Style = form.GetDefaultStyle()
	}
	if opts.ExStyle == 0 {
		opts.ExStyle = form.GetDefaultExStyle()
	}

	if me.NoResize {
		opts.StyleExclude |= win32.WS_THICKFRAME
	}
	if me.NoMinimizeBox {
		opts.StyleExclude |= win32.WS_MINIMIZEBOX
	}
	if me.NoMaximizeBox {
		opts.StyleExclude |= win32.WS_MAXIMIZEBOX
	}
	if me.State == StateMaximized {
		opts.StyleInclude |= win32.WS_MAXIMIZE
	} else if me.State == StateMinimized {
		opts.StyleInclude |= win32.WS_MINIMIZE
	}
	if me.TopMost {
		opts.ExStyleInclude |= win32.WS_EX_TOPMOST
	}
	if me.NoIcon {
		opts.ExStyleInclude |= win32.WS_EX_DLGMODALFRAME
	}
	if me.HelpButton {
		opts.ExStyleInclude |= win32.WS_EX_CONTEXTHELP
	}
	if !me.NoResize && me.SizeGrip {
		form.sizeGrip = true
	}

	if me.Size.Width != 0 {
		opts.Width = me.Size.Width
	}
	if me.Size.Height != 0 {
		opts.Height = me.Size.Height
	}
	if me.ClientSize.Width != 0 {
		var rc win32.RECT
		rc.Right = int32(me.ClientSize.Width)
		rc.Bottom = int32(me.ClientSize.Height)
		style := (opts.Style | opts.StyleInclude) & ^(opts.StyleExclude)
		exStyle := (opts.ExStyle | opts.ExStyleInclude) & ^(opts.ExStyleExclude)
		win32.AdjustWindowRectEx(&rc, style, win32.FALSE, exStyle)
		opts.Width = int(rc.Right - rc.Left)
		opts.Height = int(rc.Bottom - rc.Top)
	}
	if me.Owner != nil {
		opts.ParentHandle = me.Owner.GetHandle()
	}
	err := form.Create(*opts)
	assertNoErr(err)

	//
	if me.BackColor != colors.Null {
		form.SetBackColor(me.BackColor)
	}

	if me.Center == CenterScreen {
		form.CenterOnScreen()
	} else if me.Center == CenterParent {
		if opts.ParentHandle == 0 {
			form.CenterToWindow(opts.ParentHandle)
		} else {
			form.CenterOnScreen()
		}
	}
	if me.Icon != 0 {
		SendMessage(form.Handle, win32.WM_SETICON, win32.ICON_SMALL, me.Icon)
	}
	if me.IconBig != 0 {
		SendMessage(form.Handle, win32.WM_SETICON, win32.ICON_BIG, me.IconBig)
	}

	return form
}
