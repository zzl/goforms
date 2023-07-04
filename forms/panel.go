package forms

import (
	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/framework/virtual"
	"log"
)

type Panel interface {
	ContainerControl

	PanelObj() *PanelObject
}

type PanelSpi interface {
	ControlSpi
}

type PanelInterface interface {
	Panel
	PanelSpi
}

type PanelObject struct {
	ContainerControlObject
	super *ContainerControlObject

	WindowBg bool
	Border   bool
}

type NewPanel struct {
	Parent   Container
	Name     string
	Pos      Point
	Size     Size
	WindowBg bool
	Border   bool
}

func (me NewPanel) Create(extraOpts ...*WindowOptions) Panel {
	panel := NewPanelObject()
	panel.name = me.Name
	panel.WindowBg = me.WindowBg
	panel.Border = me.Border

	opts := utils.OptionalArg(extraOpts)
	opts.Left = me.Pos.X
	opts.Top = me.Pos.Y

	opts.ParentHandle = resolveParentHandle(me.Parent)
	err := panel.Create(*opts)
	assertNoErr(err)

	configControlSize(panel, me.Size)
	return panel
}

func NewPanelObject() *PanelObject {
	return virtual.New[PanelObject]()
}

func (this *PanelObject) PanelObj() *PanelObject {
	return this
}

func (this *PanelObject) Init() {
	this.super.Init()
}

func (this *PanelObject) GetWindowClass() string {
	if this.WindowBg {
		return "goforms.panel_w"
	} else {
		return "goforms.panel"
	}
}

var _panelClassRegstered bool
var _panel_w_ClassRegstered bool

func (this *PanelObject) EnsureClassRegistered() {
	if this.WindowBg {
		this.ensureWClassRegistered()
		return
	}
	if _panelClassRegstered {
		return
	}
	_, err := RegisterClass("goforms.panel", nil, ClassOptions{
		BackgroundBrush: ToSysColorBrush(byte(win32.COLOR_3DFACE)),
	})
	if err != nil {
		log.Fatal(err)
	}
	_panelClassRegstered = true
}

func (this *PanelObject) ensureWClassRegistered() {
	if _panel_w_ClassRegstered {
		return
	}
	_, err := RegisterClass("goforms.panel_w", nil, ClassOptions{
		BackgroundBrush: ToSysColorBrush(byte(win32.COLOR_WINDOW)),
	})
	if err != nil {
		log.Fatal(err)
	}
	_panel_w_ClassRegstered = true
}

func (this *PanelObject) GetDefaultStyle() WINDOW_STYLE {
	style := this.super.GetDefaultStyle()
	if this.Border {
		style |= win32.WS_BORDER
	}
	return style //? | win32.WS_CLIPCHILDREN
}

func (this *PanelObject) OnCtlColorStatic(msg *Message) {
	if this.WindowBg {
		msg.SetResult(0)
		msg.SetHandled(true)
	}
}

func (this *PanelObject) WinProc(winObj *WindowObject, m *Message) error {
	if m.UMsg == win32.WM_PAINT && this.HasFlag(FlagDesignMode) {
		this.paintVirtualBorder(m)
		return nil
	}
	return this.super.WinProc(winObj, m)
}

func (this *PanelObject) paintVirtualBorder(m *Message) {
	this.CallOriWndProc(m)
	m.Handled = true

	hdc := win32.GetDC(this.Handle)
	hPen := win32.CreatePen(win32.PS_DOT, 0, win32.RGB(133, 133, 133))

	hOriPen := win32.SelectObject(hdc, win32.HGDIOBJ(hPen))
	hbr := win32.GetStockObject(win32.NULL_BRUSH)
	hOriBrush := win32.SelectObject(hdc, hbr)

	var rc win32.RECT
	win32.GetClientRect(this.Handle, &rc)
	win32.SetBkMode(hdc, win32.TRANSPARENT)
	win32.Rectangle(hdc, 0, 0, rc.Right, rc.Bottom)

	win32.SelectObject(hdc, hOriPen)
	win32.SelectObject(hdc, hOriBrush)
	win32.DeleteObject(win32.HGDIOBJ(hPen))
	win32.ReleaseDC(this.Handle, hdc)
}
