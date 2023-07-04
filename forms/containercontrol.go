package forms

import (
	"log"

	"github.com/zzl/go-win32api/v2/win32"
)

type ContainerControl interface {
	Container
	Control
}

type ContainerControlObject struct {
	ContainerObject

	NameAwareSupport
}

func (this *ContainerControlObject) GetControlId() uint16 {
	ret, errno := win32.GetDlgCtrlID(this.Handle)
	if ret == 0 {
		log.Println(errno.Error())
	}
	return uint16(ret)
}

func (this *ContainerControlObject) GetDefaultStyle() WINDOW_STYLE {
	style := win32.WS_CHILD | win32.WS_VISIBLE
	incStyle, excStyle := this.RealObject.(Control).GetControlSpecStyle()
	style |= incStyle
	style &^= excStyle
	return style
}

func (this *ContainerControlObject) GetDefaultExStyle() WINDOW_EX_STYLE {
	return win32.WS_EX_CONTROLPARENT
}

func (this *ContainerControlObject) GetControlSpecStyle() (include, exclude WINDOW_STYLE) {
	return 0, 0
}

func (this *ContainerControlObject) CreateIn(
	parent Window, extraOpts ...*WindowOptions) Control {
	return createControlIn(parent, this.RealObject, extraOpts...)
}

func (this *ContainerControlObject) GetContainer() Container {
	hWndParent, _ := win32.GetParent(this.Handle)
	parentWin := GetWindow(hWndParent)
	if parentWin == nil {
		return nil //?
	}
	return parentWin.(Container)
}
