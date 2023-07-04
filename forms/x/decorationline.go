package forms

import (
	"github.com/zzl/go-win32api/v2/win32"
	. "github.com/zzl/goforms/forms"
	"github.com/zzl/goforms/framework/virtual"
)

type DecorationLine interface {
	DecorationControl
}

type DecorationLineObject struct {
	DecorationControlObject
	super *DecorationControlObject

	Vertical bool
	Length   int
}

func (this *DecorationLineObject) Init() {
	this.super.Init()
}

func (this *DecorationLineObject) GetWindowClass() string {
	return "Static"
}

func NewDecorationLineObject() *DecorationLineObject {
	return virtual.New[DecorationLineObject]()
}

func (this *DecorationLineObject) GetControlSpecStyle() (WINDOW_STYLE, WINDOW_STYLE) {
	include, exclude := this.super.GetControlSpecStyle()
	if this.Vertical {
		include |= WINDOW_STYLE(win32.SS_ETCHEDVERT)
	} else {
		include |= WINDOW_STYLE(win32.SS_ETCHEDHORZ)
	}
	return include, exclude
}

func (this *DecorationLineObject) Create(options WindowOptions) error {
	//WS_EX_STATICEDGE?
	return this.super.Create(options)
}

func (this *DecorationLineObject) GetPreferredSize(cxMax int, cyMax int) (int, int) {
	if this.Vertical {
		return 2, min(cyMax, 32)
	} else {
		return min(cxMax, 32), 2
	}
}
