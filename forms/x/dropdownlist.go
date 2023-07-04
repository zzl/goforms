package forms

import (
	"github.com/zzl/go-win32api/v2/win32"
	. "github.com/zzl/goforms/forms"
	"github.com/zzl/goforms/framework/consts"
	"github.com/zzl/goforms/framework/virtual"
)

type DropdownList interface {
	DropdownControl
}

type DropdownListObject struct {
	DropdownControlObject
	super *DropdownControlObject

	ListBox ListBox
	popup   *DropdownListPopup
}

func NewDropdownListObject() *DropdownListObject {
	return virtual.New[DropdownListObject]()
}

func (this *DropdownListObject) Init() {
	this.super.Init()
	this.PopupBorder = true
}

func (this *DropdownListObject) OnHandleCreated() {
	lb := NewListBoxObject()
	lb.IntegralHeight = true
	this.ListBox = lb
	style := lb.GetDefaultStyle()
	style &^= win32.WS_BORDER
	exStyle := lb.GetDefaultExStyle()
	exStyle &^= win32.WS_EX_CLIENTEDGE
	if exStyle == 0 {
		exStyle = consts.Zero
	}
	this.ListBox.Create(WindowOptions{
		ParentHandle: this.Handle,
		Style:        style,
		ExStyle:      exStyle,
	})
	this.popup = NewDropdownListPopup(lb)
	this.super.OnHandleCreated()
}

func (this *DropdownListObject) Dispose() {
	this.ListBox.Dispose()
	this.popup.Dispose()
	this.super.Dispose()
}

func (this *DropdownListObject) GetPopup() DropdownPopup {
	return this.popup
}
