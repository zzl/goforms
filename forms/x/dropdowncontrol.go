package forms

import (
	"github.com/zzl/go-win32api/v2/win32"
	. "github.com/zzl/goforms/forms"
	"syscall"
)

type DropdownControl interface {
	CustomControl

	DropDown()
	CloseUp()
}

type DropdownControlSpi interface {
	GetPopup() DropdownPopup
	OnDropDown()
	OnCloseUp()
}

type DropdownControlObject struct {
	CustomControlObject
	super *CustomControlObject

	PopupBorder bool

	mouseOnButton  bool
	droppingDown   bool
	popupContainer *DropdownPopupContainerObject

	//
	value interface{}
	text  string
}

func (this *DropdownControlObject) Init() {
	this.super.Init()
}

func (this *DropdownControlObject) OnDropDown() {
	//
}

func (this *DropdownControlObject) OnCloseUp() {
	//
}

func (this *DropdownControlObject) OnHandleCreated() {
	popup := this.getPopup()
	popup.GetOnOk().AddListener(func(ei *SimpleEventInfo) {
		this.value = popup.GetValue()
		this.text = popup.GetText()
		this.CloseUp()
	})
	popup.GetOnCancel().AddListener(func(ei *SimpleEventInfo) {
		this.CloseUp()
	})
}

func (this *DropdownControlObject) OnKeyDown(args KeyEventArgs) {
	if args.Key == win32.VK_F4 {
		if this.droppingDown {
			this.CloseUp()
		} else {
			this.DropDown()
		}
		this.Invalidate()
	}
}

func (this *DropdownControlObject) getPopup() DropdownPopup {
	spi := this.RealObject.(DropdownControlSpi)
	return spi.GetPopup()
}

func (this *DropdownControlObject) DropDown() {
	this.droppingDown = true

	popup := this.getPopup()
	popup.SetValue(this.value)

	ppc := NewDropdownPopupContainerObject()
	this.popupContainer = ppc

	ppc.Popup = popup
	popup.SetContainer(ppc)
	ppc.HasBorder = this.PopupBorder

	ppc.CreateFor(this.Handle)
	ppc.OnDeactivate.AddListener(func(ei *SimpleEventInfo) {
		//this.Close()
		//Dispatcher.Invoke(func() {
		//this.droppingDown = false
		//this.Invalidate()
		this.CloseUp()
		//})
	})
	ppc.Show()
}

func (this *DropdownControlObject) CloseUp() {
	if !this.droppingDown {
		return
	}
	this.popupContainer.Close()
	this.droppingDown = false
	//this.popupContainer.Destroy()
	this.popupContainer = nil
	go func() {
		hWndRoot := win32.GetAncestor(this.Handle, win32.GA_ROOT)
		win32.SendMessage(hWndRoot, win32.WM_ACTIVATE, 1, 0)
		Dispatcher.Invoke(func() {
			this.Focus()
			this.Invalidate()
		})
	}()

}

func (this *DropdownControlObject) GetBackgroundColor() win32.COLORREF {
	clr := win32.GetSysColor(win32.COLOR_WINDOW)
	return clr
}

type comboBoxMetrics struct {
	Height   int
	RcItem   win32.RECT
	RcButton win32.RECT
}

func (this *DropdownControlObject) GetPreferredSize(cxMax int, cyMax int) (int, int) {
	if this.Handle != 0 {
		cbm := GetComboBoxMetrics(this.Handle)
		return 32, cbm.Height
	}
	return 0, 0
}

func (this *DropdownControlObject) OnSetFocus() {
	this.Invalidate()
}

func (this *DropdownControlObject) OnKillFocus() {
	this.Invalidate()
}

func (this *DropdownControlObject) OnMouseEnter() {
	this.Invalidate()
}

func (this *DropdownControlObject) OnMouseDown(x int32, y int32, button byte) {
	win32.SetFocus(this.Handle)
	if this.droppingDown {
		//return this.CloseUp()
	} else {
		this.DropDown()
		this.Invalidate()
	}
}

func (this *DropdownControlObject) OnMouseMove(x int32, y int32, button byte) {
	var rcClient win32.RECT
	win32.GetClientRect(this.Handle, &rcClient)
	cbm := GetComboBoxMetrics(this.Handle)
	buttonLeft := rcClient.Right - (cbm.RcButton.Right - cbm.RcButton.Left) - 2
	mouseOnButton := x >= buttonLeft
	if mouseOnButton != this.mouseOnButton {
		this.mouseOnButton = mouseOnButton
		this.Invalidate()
	}
}

func (this *DropdownControlObject) OnMouseUp(x int32, y int32, button byte) {
	this.Invalidate()
}

func (this *DropdownControlObject) OnMouseLeave() {
	this.mouseOnButton = false
	this.Invalidate()
}

func (this *DropdownControlObject) OnPaint(hdc win32.HDC, prcClip *win32.RECT) {
	var rcClient win32.RECT
	win32.GetClientRect(this.Handle, &rcClient)
	cbm := GetComboBoxMetrics(this.Handle)

	focused := this.HasFocus()

	pwsz, _ := syscall.UTF16PtrFromString("Combobox")
	hTheme := win32.OpenThemeData(this.Handle, pwsz)
	var rcButton win32.RECT
	if hTheme == 0 {
		//?
	} else {
		//bg
		win32.DrawThemeBackground(hTheme, hdc, 2, 0, &rcClient, prcClip)
		var state int32
		state = 1
		if focused || this.droppingDown {
			state = 3
		} else if this.IsMouseHovering() {
			state = 2
		}
		//bdr
		win32.DrawThemeBackground(hTheme, hdc, 4, state, &rcClient, prcClip)
		//btn
		rcButton.Right = rcClient.Right
		rcButton.Left = rcClient.Right - (cbm.RcButton.Right - cbm.RcButton.Left) - 2
		rcButton.Bottom = rcClient.Bottom
		state = 1
		if this.mouseOnButton {
			state = 2
		}
		if this.droppingDown {
			state = 3
		}
		win32.DrawThemeBackground(hTheme, hdc, 6, state, &rcButton, prcClip)
	}

	//
	if hTheme != 0 {
		win32.CloseThemeData(hTheme)
	}

	hFont := GetDefaultFont()
	if this.GetFont() != nil {
		hFont = this.GetFont().Handle
	}
	hOriFont := win32.SelectObject(hdc, win32.HGDIOBJ(hFont))

	var clrFg, clrBg win32.COLORREF
	if this.HasFocus() {
		clrFg = win32.GetSysColor(win32.COLOR_HIGHLIGHTTEXT)
		clrBg = win32.GetSysColor(win32.COLOR_HIGHLIGHT)
	} else {
		clrFg = win32.GetSysColor(win32.COLOR_WINDOWTEXT)
		clrBg = win32.GetSysColor(win32.COLOR_WINDOW)
	}

	text := this.text

	wsz, _ := syscall.UTF16FromString(text)
	var rcText win32.RECT
	rcText.Left = 3
	rcText.Top = 3
	rcText.Bottom = rcClient.Bottom - 3
	rcText.Right = rcButton.Left - 2

	win32.SetBkColor(hdc, clrBg)
	win32.SetTextColor(hdc, clrFg)
	FillSolidRect(hdc, &rcText, clrBg)

	rcText.Left += 1
	win32.DrawText(hdc, &wsz[0], int32(len(wsz)-1), &rcText,
		win32.DT_VCENTER|win32.DT_SINGLELINE|win32.DT_END_ELLIPSIS)

	win32.SelectObject(hdc, hOriFont)

}

func (this *DropdownControlObject) GetControlSpecStyle() (include, exclude WINDOW_STYLE) {
	return win32.WS_TABSTOP, 0
}

func (this *DropdownControlObject) Create(options WindowOptions) error {
	cbm := GetComboBoxMetrics(options.ParentHandle)
	options.Height = cbm.Height
	return this.super.Create(options)
}

func (this *DropdownControlObject) SetBounds(left, top, width, height int) {
	if this.Handle != 0 {
		cbm := GetComboBoxMetrics(this.Handle)
		height = cbm.Height
	}
	this.super.SetBounds(left, top, width, height)
}
