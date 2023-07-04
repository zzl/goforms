package forms

import (
	"github.com/zzl/go-win32api/v2/win32"
	"unsafe"
)

type ColorDialog struct {
	CustomColors [16]win32.COLORREF
	ResultColor  win32.COLORREF
}

func NewColorDialog() *ColorDialog {
	return &ColorDialog{}
}

func (this *ColorDialog) Show(hWndOwner win32.HWND) bool {
	var cc win32.CHOOSECOLOR
	cc.LStructSize = uint32(unsafe.Sizeof(cc))
	cc.HwndOwner = hWndOwner
	cc.LpCustColors = &this.CustomColors[0]
	cc.RgbResult = this.ResultColor
	cc.Flags = win32.CC_FULLOPEN | win32.CC_RGBINIT
	ok := win32.ChooseColor(&cc)
	this.ResultColor = cc.RgbResult
	return ok == win32.TRUE
}
