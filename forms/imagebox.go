package forms

import (
	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/framework/virtual"
)

type ImageBox interface {
	Control
	SetBitmap(hBitmap win32.HBITMAP)
	GetBitmap() win32.HBITMAP

	SetIcon(hIcon win32.HICON)
	GetIcon() win32.HICON

	//autosize?
	//center align?

	ImageBoxObj() *ImageBoxObject
}

type ImageBoxSpi interface {
	ControlSpi
}

type ImageBoxInterface interface {
	ImageBox
	ImageBoxSpi
}

type ImageBoxObject struct {
	ControlObject
	super *ControlObject

	hBitmap win32.HBITMAP
	hIcon   win32.HICON
}

func (this *ImageBoxObject) ImageBoxObj() *ImageBoxObject {
	return this
}

func (this *ImageBoxObject) Init() {
	this.super.Init()
}

func NewImageBoxObject() *ImageBoxObject {
	return virtual.New[ImageBoxObject]()
}

type NewImageBox struct {
	Parent Container
	Name   string
	Pos    Point
	Size   Size
	Icon   win32.HICON
	Bitmap win32.HBITMAP
}

func (me NewImageBox) Create(extraOpts ...*WindowOptions) ImageBox {
	imageBox := NewImageBoxObject()
	imageBox.name = me.Name

	opts := utils.OptionalArg(extraOpts)
	opts.Left = me.Pos.X
	opts.Top = me.Pos.Y

	opts.ParentHandle = resolveParentHandle(me.Parent)
	err := imageBox.Create(*opts)
	assertNoErr(err)
	configControlSize(imageBox, me.Size)

	if me.Icon != 0 {
		imageBox.SetIcon(me.Icon)
	} else if me.Bitmap != 0 {
		imageBox.SetBitmap(me.Bitmap)
	}
	return imageBox
}

func (this *ImageBoxObject) GetWindowClass() string {
	return "Static"
}

func (this *ImageBoxObject) SetBitmap(hBitmap win32.HBITMAP) {
	oldHbm := this.hBitmap
	if oldHbm == hBitmap {
		return
	}
	this.hBitmap = hBitmap
	if this.Handle != 0 {
		this.UpdateImage()
	}
	if oldHbm != 0 {
		win32.DeleteObject(oldHbm) //?
	}
}

func (this *ImageBoxObject) GetBitmap() win32.HBITMAP {
	return this.hBitmap
}

func (this *ImageBoxObject) SetIcon(hIcon win32.HICON) {
	oldHIcon := this.hIcon
	if oldHIcon == hIcon {
		return
	}
	this.hIcon = hIcon
	if this.Handle != 0 {
		this.UpdateImage()
	}
	//if oldHIcon != 0 {
	//	win32.DeleteObject(oldHIcon) //?
	//}
}

func (this *ImageBoxObject) GetIcon() win32.HICON {
	return this.hIcon
}

func (this *ImageBoxObject) UpdateImage() {
	style := this.GetStyle()
	newStyle := style
	if this.hBitmap != 0 {
		newStyle &^= win32.WINDOW_STYLE(win32.SS_ICON)
		newStyle |= win32.WINDOW_STYLE(win32.SS_BITMAP)
	} else {
		newStyle &^= win32.WINDOW_STYLE(win32.SS_BITMAP)
		newStyle |= win32.WINDOW_STYLE(win32.SS_ICON)
	}
	if newStyle != style {
		this.ModifyStyle(newStyle, 0, 0)
	}

	if this.hBitmap != 0 {
		SendMessage(this.Handle, win32.STM_SETIMAGE,
			WPARAM(win32.IMAGE_BITMAP), this.hBitmap)
	} else {
		//SendMessage(this.Handle, win32.STM_SETIMAGE,
		//	WPARAM(win32.IMAGE_ICON), this.hIcon)
		hOriIcon, errno := SendMessage(this.Handle, win32.STM_SETICON, this.hIcon, 0)
		_ = hOriIcon
		if errno != win32.NO_ERROR {
			println("?")
		}
		//this.SetSize(32, 32) //?
	}
}

func (this *ImageBoxObject) GetControlSpecStyle() (WINDOW_STYLE, WINDOW_STYLE) {
	//return WINDOW_STYLE(win32.SS_BITMAP | win32.SS_REALSIZEIMAGE), 0
	var style WINDOW_STYLE
	//style |= WINDOW_STYLE(win32.SS_REALSIZECONTROL)
	if this.hIcon != 0 {
		style = WINDOW_STYLE(win32.SS_ICON)
	} else {
		style = WINDOW_STYLE(win32.SS_BITMAP)
	}
	return style, win32.WS_TABSTOP
}

func (this *ImageBoxObject) PreCreate(opts *WindowOptions) {
	this.super.PreCreate(opts)
	if this.HasFlag(FlagDesignMode) {
		this.GetEvent(win32.WM_PAINT).AddListener(this.handlePaint)
	}
}

func (this *ImageBoxObject) OnHandleCreated() {
	if this.hBitmap != 0 {
		this.UpdateImage()
	}
}

func (this *ImageBoxObject) OnReflectMessage(msg *Message) {
	if msg.UMsg == win32.WM_CTLCOLORSTATIC {
		backColor := this.GetBackColor()
		hdc := msg.WParam

		hbr1 := handleCtlColor(this.RealObject, hdc, backColor)
		if hbr1 != 0 {
			msg.Handled = true
			msg.Result = hbr1
		}
	} else {
		this.super.OnReflectMessage(msg)
	}
}

func (this *ImageBoxObject) GetText() string {
	text, _ := GetWindowText(this.Handle)
	return text
}

func (this *ImageBoxObject) GetPreferredSize(int, int) (cx, cy int) {
	return MeasureText(this.Handle, "W")
}

func (this *ImageBoxObject) handlePaint(msg *Message) {
	this.CallOriWndProc(msg)
	msg.Handled = true

	if this.HasFlag(FlagDesignMode) {
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
}
